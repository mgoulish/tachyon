
package tachyon

import (
         "fmt"
         "io"
         "net"
         "os"
         "strconv"
         "time"
       )



// This structure is used by connection_control()
// to control an individual connection, and to make
// new sessions on it.
// The two channels are what it uses to control 
// the twin read and write goroutines that actually
// represent the connection.
// This structure does not contain the data channels
// that get data between the app and this connection 
// because those are owned by the twin connection routines.
type connection struct {
  cnx             net.Conn
  id              int
  port            string
  read_control    chan * Message
  write_control   chan * Message
  ssn_to_cnx      chan * frame  // This is given to every session.
  session_count   int           // Later , this will become fancy -- to permit re-use of integers.

  // the Sessions Map associates each session's number
  // with the channel on which it receives incoming frames.
  // ( Incoming to the Application. )
  session_map     map [int] chan * frame
  ringbuffer    * ringbuffer

  // The State Machine keeps track of the state of this connection.
  sm            * state_machine
}





// goroutine
// This is the central control function for all connections.
func connection_control ( tach * Tachyon ) {

  var reply_type string

  listener_replies  := make ( chan * Message, 100 )
  connector_replies := make ( chan * Message, 100 )

  connection_counter := 0
  // XXX This should change to a map. To allow closure and deletion of connections.
  connections := make ( [] * connection, 0 )

  // Perpetually wait for requests from higher level,
  // and responses from lower level.
  for {
    select {
      // Get a request from higher level code.
      case request := <- tach.to_cnx_control :
        request_type := request.Info[0]
        switch request_type {
          case "listen" :
            // Launch a listener on the given port, and give
            // it a channel to reply to this function if it ever
            // hear anything.
            go listener ( request.Info[1], listener_replies )


          case "connect" :
            timeout, err := strconv.Atoi(request.Info[3])
            if err != nil {
              fp ( os.Stdout, "connection_control : bad int conv |%s|\n", request.Info[3] );
              // XXX send back error to tachyon.
              os.Exit ( 1 )
            }
            go connect ( request.Info[1], request.Info[2], timeout, connector_replies )


          case "session" :
            cnx_id, _ := strconv.Atoi(request.Info[1])
            cnx := connections[cnx_id]
            ssn_id := cnx.session_count

            // Make the channel that will send messages from the cnx to the ssn.
            // Each ssn gets its own channel from the cnx -- because we need to 
            // be able to send messages to particular sessions.
            // But they all use one channel to talk to the connection -- that is 
            // how the multiplexing of sessions into the connection happens.
            cnx_to_ssn := make ( chan * frame, 100 )
            cnx.session_map[ssn_id] = cnx_to_ssn

            // 
            tach.to_ssn_control <- & Message { Info : []string { "new_ssn" },
                                               Data : []interface{} { cnx_id, ssn_id, cnx.ssn_to_cnx, cnx_to_ssn } }
            cnx.session_count ++


          default :
            fp ( os.Stdout, 
                 "connection_control cannot yet handle |%s| requests\n", 
                 request_type )
        }


      // One of our connectors has a result.
      case connector_reply := <- connector_replies :
        reply_type = connector_reply.Info [ 0 ]
        switch reply_type {
          case "success" :
            // The connector has made a connection.
            // XXX we will use nex_cnx later. Maybe.
            make_new_connection ( tach,
                                  connector_reply, 
                                  & connections, 
                                  & connection_counter )

          default :
            fp ( os.Stdout, "connection_control got unknown connector reply |%s|\n", reply_type )
        }



      // One of our listeners has made an incoming 
      // connection on its port.
      case listener_reply := <- listener_replies :
        reply_type = listener_reply.Info[0]
        switch reply_type {
          case "success" :
            // The listener has made a connection.
            make_new_connection ( tach,
                                  listener_reply, 
                                  & connections,
                                  & connection_counter )

          default :
            fp ( os.Stdout, "connection_control got unknown listener reply |%s|\n", reply_type )
        }
    }
  }
}





// Helper function for connection_control().
func make_new_connection ( tach * Tachyon,
                           msg * Message, 
                           connections * [] * connection,
                           id_number * int ) {

  port_number := msg.Info[1]
  net_cnx := msg.Data[0].(net.Conn)

  id := * id_number
  * id_number ++


  // The new connection will be represented by two goroutines:
  // the reader and the writer. Both of them will have their 
  // own control channel, and their own data channel.
  // Connection Control will own the control channels,
  // but the connection routines will own the data channels.
  read_control  := make ( chan * Message, 10 )
  write_control := make ( chan * Message, 10 )

  // This channel will be given to all sessions on this connection.
  // They will all use it to send messages to this session.
  // TODO -- later, this will split into high and low priority.
  ssn_to_cnx    := make ( chan * frame, 100 )

  // This data structure is owned by connection_control().
  tach_cnx := & connection { cnx           : net_cnx,
                             port          : port_number,
                             id            : id,
                             read_control  : read_control, 
                             write_control : write_control, 
                             ssn_to_cnx    : ssn_to_cnx,
                             session_map   : make ( map[int]chan * frame, 100 ),
                             ringbuffer    : new_ringbuffer ( uint64(1000000) ),
                             sm            : new_cnx_state_machine ( fmt.Sprintf("cnx_%s", port_number) ),
                           } 
  *connections = append ( *connections, tach_cnx )

  // Start the twin goroutines that represent the connection.
  go read_cnx  ( tach, net_cnx, tach_cnx,   read_control )
  go write_cnx ( tach, net_cnx, ssn_to_cnx, write_control )

  // Inform the App that connection has succeeded.
  tach.Responses <- & Message { Info : []string { "connect", "success", port_number, fmt.Sprintf("%d", id) } }
}





// Make a brand-new state machine for a connection.
func new_cnx_state_machine ( name string ) ( * state_machine ) {

  states := []string { "START", 
                       "HDR_RCVD", 
                       "HDR_SENT", 
                       "HDR_EXCH", 
                       "OPEN_PIPE", 
                       "OC_PIPE", 
                       "OPEN_RCVD", 
                       "OPEN_SENT", 
                       "CLOSE_PIPE", 
                       "OPENED", 
                       "CLOSE_RCVD", 
                       "CLOSE_SENT", 
                       "DISCARDING", 
                       "END",
                     }

  return new_state_machine ( name, states )
}





// Listen to the given port forever.
// BUGALERT These port-listeners will leak if we ever decide we no longer
// want to listen to certain ports.
// However, this does not seem like a high-probability problem.
// Nor a large-magnitude problem.
func listener ( port string, reply_to chan * Message ) {

  tcp_listener, err := net.Listen ( "tcp", ":" + port )
  if err != nil {
    reply_to <- & Message { Info : []string { "error", "listen", port, err.Error() } }
    return
  }

  for {
    cnx, err := tcp_listener.Accept ( )
    if err != nil {
      reply_to <- & Message { Info : []string { "error", "accept", port, err.Error() } }
      return
    }

    reply_to <- & Message { Info : []string  { "success", port },
                          Data : []interface{} { cnx } }
  }
}





func connect ( host, port string, timeout int, reply_to chan * Message ) {
  for t := 0; t < timeout; t ++ {
    addr := host + ":" + port
    fp ( os.Stdout, "connect : dialing |%s|\n", addr )
    cnx, err := net.Dial ( "tcp", addr )
    if err == nil {
        reply_to <- & Message { Info : []string {"success", port },
                              Data : []interface{} { cnx } }
      return
    }
    time.Sleep ( time.Second )
  }

  reply_to <- & Message { Info : []string { "fail", "timeout" } }
}





func read_cnx ( tach         * Tachyon,
                cnx            net.Conn,
                tach_cnx     * connection,
                control        chan * Message ) {

  tach_cnx.ringbuffer.set_cnx ( tach_cnx )

  buffer := make ( []byte, 10000 )

  // Read bytes from the CNX, build frames out of them,
  // and send those frames to the SSN.
  for {
    bytes_read, err := cnx.Read ( buffer )

    if err != nil {
      if err == io.EOF {
        break // quit trying -- we lost the connection.
      } else {
        tach.Errors <- fmt.Sprintf ( "Connection read error |%s|", err.Error() )
        os.Exit ( 1 )
      }
    } 

    // Write the bytes we have received into the RB.
    // The RB will send out whole frames to the appropriate SSN
    // as it receives them.
    for {
      // This will busywait if RB is full.
      if tach_cnx.ringbuffer.write ( buffer [ 0 : bytes_read ] ) { 
        break
      }
    }
  }
}





func write_cnx ( tach        * Tachyon, 
                 cnx           net.Conn, 
                 outbound_data chan * frame, 
                 control       chan * Message ) {

  for {
    f := <- outbound_data
    var bytes_written uint32
    length := uint32 ( f.data.Len() )

    for ; bytes_written < length; {
      n, err := cnx.Write ( f.data.Bytes() )
      bytes_written += uint32(n)
      if err != nil {
        if err == io.EOF {
          goto done
        } else {
          tach.Errors <- fmt.Sprintf ( "Connection read error |%s|", err.Error() )
          break
        }
      }
    }
  }
  done :
}





