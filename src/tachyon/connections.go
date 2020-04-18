
package tachyon

import (
         "bufio"
         "fmt"
         "io"
         "net"
         "strconv"
         "time"
       )



var fp = fmt.Fprintf



func connection_control ( tach * Tachyon ) {

  for {
    req, more  := <- tach.cnx

    if ! more {
      break
    }

    switch req.Info[0] {

      //-------------------------------------------------
      // Connect to the given host and port.
      //-------------------------------------------------
      case "connect" :
        host := req.Info[1]
        port := req.Info[2]

        timeout, err := strconv.Atoi ( req.Info[3] )
        if err != nil {
          tach.Errors <- & Message { Info: []string { fmt.Sprintf ( "cnx arg |%s| did not parse to int", req.Info[3] ) } }
          continue
        }

        connection, err := dialer ( tach, host, port, timeout )
        if err != nil {
          tach.Errors <- & Message { Info: []string{ fmt.Sprintf ( "cnx failed to connect with command |%s| error : |%s|", req, err.Error() )}}
          continue
        }

        // We have successfully made a network connection.
        // Start the sender and receiver.
        // Note: Whether the App is the initiator or the accepter 
        //       of this connection -- in either case it gets both a 
        //       sender and a receiver created.
        //       After the creation of the connection, the initiating
        //       App and the accepting App are prefectly symmetrical.
        go sender   ( tach, connection )
        go receiver ( tach, connection )


      //-------------------------------------------------
      // Listen for connections.
      //-------------------------------------------------
      case "listen" :
        port := req.Info [ 1 ]
        go listen_to_port ( tach, port )


      default :
        tach.Errors <- & Message { Info: []string { fmt.Sprintf ( "cnx unrecognized request |%s|.", req.Info[1]) }}
    }
  }
}





func dialer ( tach * Tachyon, host, port string, timeout_seconds int ) ( cnx net.Conn, err error ) {

  for t := 0; t < timeout_seconds; t ++ {
    cnx, err := net.Dial ( "tcp", host + ":" + port )
    if err == nil {
      // Success.
      return cnx, nil 
    }
    time.Sleep ( time.Second )
  }

  return nil, fmt.Errorf ( "Dialer for %s:%s timed out.", host, port )
}





// Listen to the given port forever.
func listen_to_port ( tach * Tachyon, port string ) {

  tcp_listener, err := net.Listen ( "tcp", ":" + port )
  if err != nil {
    tach.Errors <- & Message { Info : []string { "listener failed", port, err.Error() } }
    return
  }

  for {
    cnx, err := tcp_listener.Accept ( )
    if err != nil {
      tach.Errors <- & Message { Info : []string { "listener accept failed", port, err.Error() } }
      return
    }
    
    // We have made a connection.
    // Start the sender and receiver.
    // Note: Whether the App is the initiator or the accepter 
    //       of this connection -- in either case it gets both a 
    //       sender and a receiver created.
    //       After the creation of the connection, the initiating
    //       App and the accepting App are prefectly symmetrical.
    go sender   ( tach, cnx )
    go receiver ( tach, cnx )
  }
}




//============================================================
// Send the app a channel for this connection,
// and send the bytes on it.
//============================================================
func receiver ( tach * Tachyon, cnx net.Conn ) {
  
  tach.Responses <- & Message { Info : []string{ "start_receiving" } }
  frame_size  := 1032
  header_size := 8

  buffer := make ( []byte, frame_size )
  buffered_cnx := bufio.NewReader ( cnx )
  frame_count := 0

  // Read the connection forever.
  for {
    _, err := io.ReadFull ( buffered_cnx, buffer )
    if err != nil {
      break
    }

    // Make a message from the raw bytes...
    channel_number := string ( buffer [  : header_size ] )
    body           := string ( buffer [ header_size : ] )
    message := & Message { Info: []string{channel_number}, Data: []interface{} {body} }
    // ...and send it to the App.
    tach.Incoming <- message

    frame_count ++
  }

  close ( tach.Incoming )
}





func sender ( tach * Tachyon, cnx net.Conn ) {
  defer cnx.Close ( )

  tach.Responses <- & Message { Info : []string{ "start_sending" } }

  for {
    // Get a message that the App is sending outbound.
    outgoing_message, more := <- tach.Outgoing

    // If the App has shut down outbound messages,
    // close the connection.
    if ! more {
      cnx.Close()
      break
    }

    // Send the message out on the network connection.
    fp ( cnx, outgoing_message.Info[0] + outgoing_message.Data[0].(string) )
  }
}





