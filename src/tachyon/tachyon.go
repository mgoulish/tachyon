
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
  go tach_input ( tach )

  return tach
}





func get_val_from_msg ( attr string, msg * Msg ) ( val interface{} ) {
  for _, pair := range msg.Data {
    if attr == pair.Attr {
      return pair.Val
    }
  }

  return nil
}





//=================================================================
//  Private
//=================================================================

func tach_input ( tach * Tachyon ) {
  //var topics [] * Topic

  topics := make ( map[string] * Topic )

  for {
    msg, more := <- tach.Requests
    if ! more {
      break
    }

    switch msg.Data[0].Attr {
      
      // Topics can't have any of the following keywords as their namnes.

      case "new_topic" :
        name := msg.Data[0].Val.(string) 
        top  := New_Topic ( name )
        // topics = append ( topics, top )
        topics [ name ] = top


      case "subscribe" :

        topic_name, ok := get_val_from_msg("subscribe", msg).(string)
        if ! ok {
          fp ( os.Stdout, "tach_input error: subscribe with no topic name |%#v|\n", msg )
          continue
        }

        channel, ok := get_val_from_msg("channel", msg).(chan * Msg)
        if ! ok {
          fp ( os.Stdout, "tach_input error: subscribe with no channel |%#v|\n", msg )
          continue
        }

        // subscribers = append ( subscribers, channel )
        topic, ok := topics [ topic_name ]
        if ! ok {
          fp ( os.Stdout, "tach_input error: no such topic: |%s|\n", topic_name )
        }
        topic.subscribe ( channel )




      default :
        
        // Any other word must be a topic name.

        top_name := msg.Data[0].Attr
        top, ok := topics [ top_name ]
        if ! ok {
          fp ( os.Stdout, "tach_input error: Can't find topic |%s|\n", top_name )
          continue
        }
        top.inputs <- msg
    }
  }

  fp ( os.Stdout, "tach_input exiting.\n" )
}





