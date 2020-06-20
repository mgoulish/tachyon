
package tachyon

import (
         "fmt"
         "os"
       )



var fp = fmt.Fprintf


//=================================================================
//  Public
//=================================================================



// These are the processors that cooperate to
// build up an understanding of the input image.
type Abstractor struct {
  Name                 string
  Output_Topic         string
  Subscribed_Topics [] string

  Run       func ( * Tachyon, * Abstractor )
  Visualize func ( * Tachyon, * Abstractor )
}



type  Message  map[string]interface{}



// These are the artifacts that each Abstractor
// posts as a result of its work.
type Abstraction struct {
  Abstractor_Name string
  Abstraction_ID  int64    // Only unique within the namespace of this Abstractor's posts.
  Msg             Message
}



type Tachyon struct {
  Requests     chan Message       // from App to Tachyon
  Responses    chan Message       // from Tachyon to App
  Abstractions chan * Abstraction // from Abstractors to Tachyon


  abstractors [] *Abstractor
}





func New_Tachyon ( ) ( * Tachyon ) {
  tach := & Tachyon { Requests  : make ( chan Message, 100 ),
                      Responses : make ( chan Message, 100 ),
                    }
  go tach_requests ( tach )

  return tach
}





//=================================================================
//  Private
//=================================================================

func tach_requests ( tach * Tachyon ) {

  topics := make ( map[string] * Topic )

  for {
    msg, more := <- tach.Requests
    if ! more {
      break
    }

    switch msg["request"] {
      
      // Topics can't have any of the following keywords as their namnes.

      case "add_abstractor" :
        abstractor, ok := msg["abstractor"].(*Abstractor)
        if ! ok {
          fp ( os.Stdout, "tachyon tach_requests error: no abstractor in new_abstractor message.\n" )
          fp ( os.Stdout, "MDEBUG type is |%T|\n", abstractor )
          os.Exit ( 1 )
        }
        tach.abstractors = append ( tach.abstractors, abstractor )
        fp ( os.Stdout, "Tachyon: added abstractor |%s|\n", abstractor.Name )



      case "new_topic" :
        name, ok := msg["name"].(string) 
        if ! ok {
          fp ( os.Stdout, "tachyon tach_requests error: new_topic: no string value for |name|\n" )
          continue
        }
        top  := New_Topic ( name )
        topics [ name ] = top
        // Tell the App that the topic has been created.
        tach.Responses <- Message { "response" : "new_topic",
                                    "name"     : name }



      case "subscribe" :
        topic_name, ok := msg["topic"].(string)
        if ! ok {
          fp ( os.Stdout, "tach_requests error: subscribe with no topic |%#v|\n", msg )
          continue
        }
        channel, ok := msg["channel"].(chan Message)
        if ! ok {
          fp ( os.Stdout, "tach_requests error: subscribe with no channel |%#v|\n", msg )
          continue
        }
        topic, ok := topics [ topic_name ]
        if ! ok {
          fp ( os.Stdout, "tach_requests error: no such topic: |%s|\n", topic_name )
        }
        topic.subscribe ( channel )
        // The topic will send a confirmation message.



      case "start abstractors" :
        fp ( os.Stdout, "tachyon: starting abstractors.\n" )
        for _, a := range tach.abstractors {
          fp ( os.Stdout, "tachyon: starting |%s|\n", a.Name )
          go a.Run ( tach, a )
        }



      case "post" :
        topic_name, ok := msg["topic"].(string)
        if ! ok {
          fp ( os.Stdout, "tach_requests error: post with no topic name |%#v|\n", msg )
          continue
        }
        top := topics [ topic_name ]

        fp ( os.Stdout, "tachyon: tach_requests got post for |%s|\n", topic_name )

        // Post the whole message to the desired topic.
        // End-users will simply ignore key-value pairs that they do not need.
        top.post ( msg )



      default :
        fp ( os.Stdout, "tachyon: tach_requests error: unknown request |%s|\n", msg["request"] )
        os.Exit ( 1 )
    }
  }

  fp ( os.Stdout, "tach_requests exiting.\n" )
}





