
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





func Get_Val_From_Msg ( attr string, msg * Msg ) ( interface{} ) {
  for _, pair := range msg.Data {
    if attr == pair.Attr {
      //fp ( os.Stdout, "MDEBUG Get_Val_From_Msg type is %T\n", pair.Val )
      return pair.Val
    }
  }


  return nil
}





//=================================================================
//  Private
//=================================================================

func tach_input ( tach * Tachyon ) {

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
        topics [ name ] = top
        // Tell the App that the topic has been created.
        tach.Responses <- & Msg { []AV { { "new_topic", name } } }



      case "subscribe" :
        topic_name, ok := Get_Val_From_Msg("subscribe", msg).(string)
        if ! ok {
          fp ( os.Stdout, "tach_input error: subscribe with no topic name |%#v|\n", msg )
          continue
        }
        channel, ok := Get_Val_From_Msg("channel", msg).(chan * Msg)
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



      case "post" :
        topic_name, ok := Get_Val_From_Msg("post", msg).(string)
        if ! ok {
          fp ( os.Stdout, "tach_input error: post with no topic name |%#v|\n", msg )
          continue
        }
        top := topics [ topic_name ]

        fp ( os.Stdout, "MDEBUG tach_input got post for |%s|\n", topic_name )

        // No need to remove the Attr-Val pair that describes the topic.
        // End-users will simply ignore it.
        top.post ( msg )




      default :
        
        // Actually I think this is archaic and unused.
        // Let's see...

        fp ( os.Stdout, "MDEBUG Oops. I guess not.\n" )
        os.Exit ( 1 )

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





