
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
  Abstraction_ID  uint64    // Only unique within the namespace of this Abstractor's posts.
  Topic           string
  Msg             Message
}



type Tachyon struct {
  Requests     chan   Message     // from App to Tachyon
  Responses    chan   Message     // from Tachyon to App
  Abstractions chan * Abstraction // from Abstractors to Tachyon, and then out to the Topics.

  abstractors  [] *Abstractor
  topics       map[string] * Topic

  abstractions_to_bb chan * Abstraction
  requests_to_bb     chan   Message
}





func New_Tachyon ( ) ( * Tachyon ) {
  tach := & Tachyon { Requests     : make ( chan Message, 100 ),
                      Responses    : make ( chan Message, 100 ),
                      Abstractions : make ( chan * Abstraction, 100 ),

                      topics       : make ( map[string] * Topic ),

                      abstractions_to_bb : make ( chan * Abstraction, 100 ),
                      requests_to_bb     : make ( chan Message, 100 ),
                    }

  go bulletin_board ( tach )
  go requests     ( tach )
  go abstractions ( tach )

  return tach
}





//=================================================================
//  Private
//=================================================================

func requests ( tach * Tachyon ) {

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
        tach.topics [ name ] = top
        // Let the Bulletin Board know about it.
        tach.requests_to_bb <- msg
        // Tell the App that the topic has been created.
        tach.Responses <- Message { "response" : "new_topic",
                                    "name"     : name }



      case "subscribe" :
        topic_name, ok := msg["topic"].(string)
        if ! ok {
          fp ( os.Stdout, "tach_requests error: subscribe with no topic |%#v|\n", msg )
          continue
        }
        channel, ok := msg["channel"].(chan * Abstraction)
        if ! ok {
          fp ( os.Stdout, "tach_requests error: subscribe with no channel |%#v|\n", msg )
          fp ( os.Stdout, "MDEBUG type is %T\n", channel )
          os.Exit ( 1 )
        }
        topic, ok := tach.topics [ topic_name ]
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



      default :
        fp ( os.Stdout, "tachyon: tach_requests error: unknown request |%s|\n", msg["request"] )
        os.Exit ( 1 )
    }
  }

  fp ( os.Stdout, "tach_requests exiting.\n" )
}




// All abstractions posted by the Abstractots come here, 
// before being sent out to their individual topics.
func abstractions ( tach * Tachyon ) {
  for {
    abstraction, more := <- tach.Abstractions
    if ! more {
      break
    }
    
    topic_name, ok := abstraction.Msg["topic"].(string)
    if ! ok {
      fp ( os.Stdout, "Tachyon error: abstractions: no topic name in this post: |%#v|\n", abstraction )
      os.Exit ( 1 )
    }

    topic, ok := tach.topics[topic_name]
    if ! ok {
      fp ( os.Stdout, "tachyon error: Got Abstraction with no topic: |%#v|\n", abstraction )
      continue  // Just drop it. Topicless posts may be used in development.
    }

    topic.post ( abstraction )

    // And let the Bulletin Board know about it.
    tach.abstractions_to_bb <- abstraction
  }
}





