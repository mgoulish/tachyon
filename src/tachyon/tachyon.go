
package tachyon

import (
         "fmt"
         "os"
       )



var fp = fmt.Fprintf


//=================================================================
//  Public
//=================================================================

type Message struct {
  AV_pairs map [ string ] interface{}
}





type Tachyon struct {
  Requests  chan * Message
  Responses chan * Message
}





func New_Tachyon ( ) ( * Tachyon ) {
  tach := & Tachyon { Requests  : make ( chan * Message, 100 ),
                      Responses : make ( chan * Message, 100 ),
                    }
  go tach_listen_for_requests ( tach )

  return tach
}





//=================================================================
//  Private
//=================================================================

func tach_listen_for_requests ( tach * Tachyon ) {
  for {
    req, more := <- tach.Requests
    if ! more {
      break
    }
    fp ( os.Stdout, "MDEBUG tach got request: |%#v|\n", req )
  }

  fp ( os.Stdout, "MDEBUG Tachyon requet listener exiting.\n" )
}





