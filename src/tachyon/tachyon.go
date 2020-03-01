
package tachyon

import (
         "fmt"
         "os"
       )



var fp = fmt.Fprintf





type Tachyon struct {

  // Public channels ------------------
  Requests  chan * Message
  Responses chan * Message
  Errors    chan string

  // Private channels ------------------
  to_cnx_control  chan * Message
  to_ssn_control  chan * Message
  to_lnk_control  chan * Message
}





type Message struct {
  Info [] string
  Data [] interface{}
}





/*=================================================
  Public Functions
===================================================*/

func New_Tachyon ( ) ( * Tachyon ) {

  tach := & Tachyon { }

  // Public channels --------------------------
  tach.Requests  = make ( chan * Message,   100 )   // Requests from the App.
  tach.Responses = make ( chan * Message,   100 )   // Responses to the App.
  tach.Errors    = make ( chan   string,    100 )   // Error messages to the App.

  // Private channels --------------------------
  tach.to_cnx_control    = make ( chan * Message, 100 )
  tach.to_ssn_control    = make ( chan * Message, 100 )
  tach.to_lnk_control    = make ( chan * Message, 100 )

  // Start the Tachyon components.
  // They will all start listening on their dedicated channels.
  go connection_control ( tach )
  go session_control    ( tach )
  go link_control       ( tach )

  // Last of all, start accepting App requests
  go app_request_listener ( tach )

  return tach
}





func (t * Tachyon ) Reify () {
  fp ( os.Stdout, "Yeah I got your reification right here, buddy.\n" )
  os.Exit ( 666 )
}





/*=================================================
  Private Functions
===================================================*/


//----------------------------------------------------------
// goroutine.
// Fields request-messages from the Application 
// and distributes them to the Tachyon components.
//----------------------------------------------------------
func app_request_listener ( tach * Tachyon ) {
  for {
    req := <- tach.Requests
    req_name := req.Info[0] 

    switch req_name {
      
      case "connect", "listen", "close" :
        tach.to_cnx_control <- req

      case "session" :
        // A session request goes to the connection first,
        // because it must assign an ID for the new session
        // before the session is created.
        tach.to_cnx_control <- req 

      case "send_LNK", "recv_LNK" :
        // A link request is sent to Session Control first, because 
        // it must add the channel that this link will use to talk
        // to its session.
        tach.to_ssn_control <- req

      default :
        tach.Errors <- fmt.Sprintf ( "tachyon : unrecognized request |%s|.", req )
    }
  }
}




