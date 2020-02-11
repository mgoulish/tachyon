
package tachyon


import (
         "fmt"
         "os"
         "strconv"
       )



type session_t struct {
  cnx_id          int
  ssn_id          int
  ssn_to_cnx      chan * frame
  cnx_to_ssn      chan * frame

  // When links send messages to us, they all use the same channel.
  // TODO -- later, this will be two channels -- high and normal priority.
  links_to_ssn    chan * frame

  // But when we talk down to the SSNs, we use one channel for each.
  ssn_to_links [] chan * frame

  // This SSN's control channel from the CNX, above.
  control         chan * Message
}





func session_control ( tach * Tachyon ) {

  sessions := make ( [] * session_t, 0 )

  //=====================================================
  // Perpetually wait for requests from higher level,
  // and responses from lower level.
  //=====================================================
  for {

    select {
      //---------------------------------------------------
      // Get a message from Tachyon or another component.
      //---------------------------------------------------
      case msg := <- tach.to_ssn_control :
        msg_type := msg.Info[0]

        switch msg_type {
          case "new_ssn" :
            
            cnx_id := msg.Data[0].(int)
            ssn_id := msg.Data[1].(int)

            // The CNX has sent us the two channels that 
            // we will use to talk with it.
            ssn_to_cnx := msg.Data[2].(chan * frame)
            cnx_to_ssn := msg.Data[3].(chan * frame)

            // Make the control channel that we will use from
            // this fn, to control the SSN we are making here.
            control := make ( chan * Message, 100 )

            // All of this SSN's link will talk back to it
            // using this single link.
            // TODO -- later this will be two: high and normal priority.
            links_to_ssn := make ( chan * frame, 100 )

            // But this SSN will talk to each link on its own individual 
            // channel. The SSN itself will fill up this array as its links
            // are made.
            // TODO -- later this will be a map, not an array, to permit deletion.
            ssn_to_links := make ( [] chan * frame, 0 )

            // Create and store the new ssn info for retrieval later,
            // and launch the session as a goroutine.
            new_ssn := & session_t { cnx_id       : cnx_id,
                                     ssn_id       : ssn_id,
                                     ssn_to_cnx   : ssn_to_cnx,
                                     cnx_to_ssn   : cnx_to_ssn,
                                     links_to_ssn : links_to_ssn,
                                     ssn_to_links : ssn_to_links,
                                     control      : control }
            sessions = append ( sessions, new_ssn )

            go session ( tach, new_ssn )
          
            // Send to the App what it needs to talk to this SSN
            // about LNKs.
            tach.Responses <- & Message { Info : []string { "session_created" },
                                          Data : []interface{} { cnx_id, 
                                                                 ssn_id } }


          case "send_LNK" :
            // Retrieve the desired SSN...
            ssn_id, err := strconv.Atoi(msg.Info[2])
            if err != nil {
              tach.Errors <- fmt.Sprintf ( "ssn error : bad int conversion on ssn_id |%s|", msg.Info[3] )
            }
            ssn := sessions [ ssn_id ] // XXX BUGALERT -- this will not be valid forever.
            if ssn == nil {
              fp ( os.Stdout, "MDEBUG session_control error -- no such session!\n" )
              os.Exit ( 1 )
            }

            // ...and forward this message to it.
            ssn.control <- msg



          case "recv_LNK" :
            // Retrieve the desired SSN...
            ssn_id, err := strconv.Atoi(msg.Info[2])
            if err != nil {
              tach.Errors <- fmt.Sprintf ( "ssn error : bad int conversion on ssn_id |%s|", msg.Info[3] )
            }
            ssn := sessions [ ssn_id ] // XXX BUGALERT -- this will not be valid forever.
            if ssn == nil {
              fp ( os.Stdout, "MDEBUG session_control error -- no such session!\n" )
              os.Exit ( 1 )
            }

            // ...and forward this message to it.
            ssn.control <- msg


          default:
            fp ( os.Stdout, "session_control can't handle |%s| msg.\n", msg_type )
        }
    }
  }
}





// XXX do I need a read and write here, like in connection ?
func session ( tach * Tachyon, ssn * session_t ) {

  send_lnk_controls := make ( map [ string ] chan * Message )
  recv_lnk_controls := make ( map [ string ] chan * Message )

  // The receiving LNKs, i.e. those that receive messages
  // from this SSN, each need their own individual channel 
  // to receive from.
  recv_lnk_channels := make ( map [ string ] chan * frame )


  // All LNKs will use this channel to talk to this SSN.
  // TODO -- split into high- and low-priority.
  lnks_to_ssn := make ( chan * frame, 100 )

  for {
    select {

      //=========================================================
      // A control message is coming down from the CNX.
      //=========================================================
      case msg := <- ssn.control :

        msg_type := msg.Info[0]

        switch msg_type {
          case "send_LNK" :
            // Create a new link that will send messages to this session,
            // to be put out over the CNX.
            cnx_id := msg.Info[1]
            ssn_id := msg.Info[2]
            addr   := msg.Info[3]

            // Each LNK gets its own personal control channel.
            control := make ( chan * Message, 100 )
            send_lnk_controls [ addr ] = control

            // To this message, add the channel that the link 
            // will use to communicate to this SSN.
            tach.to_lnk_control <- & Message { Info : []string { "send_LNK", cnx_id, ssn_id, addr },
                                             Data : []interface{} { lnks_to_ssn, control } }

          case "recv_LNK" :
            // Make a new LNK that will receive frames,
            // from this SSN that were received from the CNX.
            cnx_id := msg.Info[1]
            ssn_id := msg.Info[2]
            addr   := msg.Info[3]

            // Make the channel that the link will use to
            // communicate with this SSN. 
            control    := make ( chan * Message, 100 )
            ssn_to_lnk := make ( chan * frame, 100 )

            // Store them in the maps that this fn will use later
            // to talk to them.
            recv_lnk_controls [ addr ] = control
            recv_lnk_channels [ addr ] = ssn_to_lnk

            // To this message, add the channel that the link 
            // will use to communicate to this SSN.
            tach.to_lnk_control <- & Message { Info : []string { "recv_LNK", cnx_id, ssn_id, addr },
                                             Data : []interface{} { ssn_to_lnk, control } }

          
          default :
            fp ( os.Stdout, "session got unknown request |%s|\n", msg_type )
        }
      

      //=========================================================
      // A message is coming up to us from a link.
      //=========================================================
      case msg := <- lnks_to_ssn :
        
        // Send it up to the CNX.
        ssn.ssn_to_cnx <- msg
      

      //=========================================================
      // A message is coming down to us from the CNX.
      //=========================================================
      case msg := <- ssn.cnx_to_ssn :
        
        //fp ( os.Stdout, "MDEBUG ssn got msg from cnx.\n" )
        // Find LNK with this addr.
        addr := "my_addr"
        channel, ok := recv_lnk_channels [ addr ]
        if ok {
          channel <- msg
        } else {
          os.Exit ( 1 )
        }
    }
  }
}





