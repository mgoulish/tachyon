
package tachyon


import (
         "os"
         "fmt"
       )





// It is the link that handles conversion between 
// frames and app-consumable messages.
type link struct {
  direction    string
  cnx_id       string
  ssn_id       string
  addr         string
  ssn_chan     chan * frame
  app_chan     chan * Message
  control      chan * Message
}





// goroutine
func link_control ( tach * Tachyon ) {
  
  send_links := make ( [] * link, 0 )
  recv_links := make ( [] * link, 0 )

  for {
    select {
      case msg := <- tach.to_lnk_control :

        message_type := msg.Info[0]

        switch message_type {

          case "send_LNK" :

             cnx_id     := msg.Info[1]
             ssn_id     := msg.Info[2]
             addr       := msg.Info[3]

             ssn_chan   := msg.Data[0].(chan * frame)
             control    := msg.Data[1].(chan * Message)

             app_chan   := make ( chan * Message, 100 )

             lnk := & link { direction   : "send",
                             cnx_id      : cnx_id,
                             ssn_id      : ssn_id,
                             addr        : addr,
                             ssn_chan    : ssn_chan,
                             app_chan    : app_chan,
                             control     : control }

             send_links = append ( send_links, lnk )

             // XXX -- Dude. Where's my control chan?
             go send ( tach, lnk )

             // And finally, inform the App.
             tach.Responses <- & Message { Info : []string { "send_LNK", cnx_id, ssn_id, addr },
                                         Data : []interface{} { app_chan } }


          case "recv_LNK" :

             cnx_id     := msg.Info[1]
             ssn_id     := msg.Info[2]
             addr       := msg.Info[3]

             ssn_chan   := msg.Data[0].(chan * frame)
             control    := msg.Data[1].(chan * Message)

             app_chan   := make ( chan * Message, 100 )

             lnk := & link { direction   : "recv",
                             cnx_id      : cnx_id,
                             ssn_id      : ssn_id,
                             addr        : addr,
                             ssn_chan    : ssn_chan,
                             app_chan    : app_chan,
                             control     : control }

             recv_links = append ( recv_links, lnk )

             go recv ( tach, lnk )

             // And finally, inform the App.
             tach.Responses <- & Message { Info : []string { "recv_LNK", cnx_id, ssn_id, addr },
                                         Data : []interface{} { app_chan } }


          default :
            tach.Errors <- fmt.Sprintf ( "link error : unknown message type |%s|", message_type )
        }
      // No default here, or it will busy-wait.
    }
  }
}





func send ( tach * Tachyon, lnk * link ) {

  for {
    select {
      case control_msg := <- lnk.control :
        fp ( os.Stdout, "MDEBUG sender %s got ctrl msg |%s|\n", lnk.addr, control_msg.Info[0] )

      case msg_from_app := <- lnk.app_chan :
        // Simply pass it up to SSN, for now.
        f, err := enframe ( msg_from_app )
        if err != nil {
          fp ( os.Stdout, "MDEBUG LNK send got framing error |%s|\n", err.Error() )
          os.Exit ( 1 )
        }
        lnk.ssn_chan <- f
    }
  }
}





func recv ( tach * Tachyon, lnk * link ) {

  for {
    select {

      case control_msg := <- lnk.control :
        fp ( os.Stdout, "MDEBUG send LNK %s got ctrl msg |%s|\n", lnk.addr, control_msg.Info[0] )
      
      case frame_from_ssn := <- lnk.ssn_chan :
        // Simply pass it down to App, for now.
        msg, err := deframe ( frame_from_ssn )
        if err != nil {
          fp ( os.Stdout, "MDEBUG recv : deframing error : |%s|\n", err.Error() )
        } else {
          lnk.app_chan <- msg
        }
    }
  }
}





