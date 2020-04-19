
package tachyon

import (
         "fmt"
       )





//=================================================================
//  Public
//=================================================================

type Message struct {
  Info [] string
  Data [] interface{}
}





type Tachyon struct {
  // Public channels ------------------
  Requests  chan * Message
  Responses chan * Message
  Errors    chan * Message

  Outgoing  chan * Message
  Incoming  chan * Message

  ID        chan int

  // Private channels ------------------
  cnx       chan * Message

  id        int
}





func New_Tachyon ( id int ) ( * Tachyon ) {

  tach := & Tachyon { id : id }

  // Public channels --------------------------
  tach.Requests  = make ( chan * Message, 100 )  // The App sends requests to Tachyon.
  tach.Responses = make ( chan * Message, 100 )  // Tachyon sends responses to the App.
  tach.Outgoing  = make ( chan * Message, 100 )  // The App sends messages out the port.
  tach.Incoming  = make ( chan * Message, 100 )  // Messagges inbound from the port to the App.
  tach.Errors    = make ( chan * Message, 100 )  // Tachyon sends errors to the App.
  tach.ID        = make ( chan   int,       1 )  // Tachyon issues IDs to senders.

  // Private channels --------------------------
  tach.cnx       = make ( chan * Message, 100 )  // Tachyon sends requests to Connection Control.

  // Start the Tachyon components.
  go issue_ids          ( tach )
  go connection_control ( tach )                 
  go requests_handler   ( tach )                 

  return tach
}





//=================================================================
//  Private
//=================================================================


func requests_handler ( tach * Tachyon ) {

  for {
    req, more := <- tach.Requests
    if ! more {
      break
    }

    switch req.Info[0] {

      // So far, the only requests we have both go to Connection Control.
      case "connect", "listen" :
        tach.cnx <- req

      default :
        tach.Errors <- & Message { Info: []string {fmt.Sprintf ( "unrecognized request |%s|", req)} }
    }
  }
}





func issue_ids ( tach * Tachyon ) {
  id := 0
  for {
    tach.ID <- id
    id ++
  }
}





