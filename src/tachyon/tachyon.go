
package tachyon

import (
         "fmt"
         "os"
       )



var fp = fmt.Fprintf


//=================================================================
//  Public
//=================================================================

type AV struct {
  Attr string
  Val  interface{}
}

type Msg struct {
  Data [] AV
}





type Tachyon struct {
  Requests  chan * Msg
  Responses chan * Msg
}





func New_Tachyon ( ) ( * Tachyon ) {
  tach := & Tachyon { Requests  : make ( chan * Msg, 100 ),
                      Responses : make ( chan * Msg, 100 ),
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

    fp ( os.Stdout, "MDEBUG tach got req : |%s|\n", req.Data[0].Attr )
  }

  fp ( os.Stdout, "MDEBUG Tachyon requet listener exiting.\n" )
}





