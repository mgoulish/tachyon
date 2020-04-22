
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
  go tach_listen_for_input ( tach )

  return tach
}





//=================================================================
//  Private
//=================================================================

func tach_listen_for_input ( tach * Tachyon ) {
  for {
    input, more := <- tach.Requests
    if ! more {
      break
    }

    switch input.Data[0].Attr {
      
      case "new_topic" :
        top := New_Topic ( input.Data[0].Val.(string) )
        fp ( os.Stdout, "tach : made topic |%s|\n", top.name )

      default :
        fp ( os.Stdout, "tach error : unrecognized input : |%s|\n", input.Data[0].Attr )
    }
  }

  fp ( os.Stdout, "tach input listener exiting.\n" )
}





