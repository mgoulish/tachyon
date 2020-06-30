
package tachyon

import (
         "fmt"
         "os"
         "time"
       )



var fp = fmt.Fprintf


//=================================================================
//  Public
//=================================================================


type Abstraction_Channel chan * Abstraction


// These are the processors that cooperate to
// build up an understanding of the input image.
type Abstractor struct {
  Name   string
  Input  Abstraction_Channel
  Output Abstraction_Channel
  Log    string
  Run    func ( * Tachyon, * Abstractor )
}



type  Message  map[string]interface{}


type Abstraction_ID struct {
  Abstractor_Name string
  ID              uint64    // Only unique within the namespace of this Abstractor's posts.
  Creation_Time   float64
}



// These are the artifacts that each Abstractor
// posts as a result of its work.
type Abstraction struct {
  ID              Abstraction_ID
  Topic           string
  Msg             Message
  Genealogy       [] * Abstraction_ID
}



type Tachyon struct {
  Requests     chan   Message     // from App to Tachyon
  Responses    chan   Message     // from Tachyon to App

  abstractors  [] *Abstractor
  topics       map[string] * Topic

  abstractions_to_bb Abstraction_Channel
  requests_to_bb     chan   Message
}





func New_Tachyon ( ) ( * Tachyon ) {
  tach := & Tachyon { Requests     : make ( chan Message, 100 ),
                      Responses    : make ( chan Message, 100 ),

                      topics       : make ( map[string] * Topic ),

                      abstractions_to_bb : make ( chan * Abstraction, 100 ),
                      requests_to_bb     : make ( chan Message, 100 ),
                    }

  go bulletin_board ( tach )
  go requests     ( tach )

  return tach
}





func (tach * Tachyon) Get_Topic ( topic_name string ) (Abstraction_Channel) {
  topic, ok :=  tach.topics [ topic_name ]

  if ! ok {
    fp ( os.Stdout, "Tachyon error: topic not found: |%s|\n", topic_name )
    os.Exit ( 1 )
  }

  return topic.input_channel
}




func Path_Exists ( path string ) ( bool ) {
  _, err := os.Stat ( path )
  if err == nil { 
    return true
  }

  if os.IsNotExist(err) { 
    return false 
  }

  fp ( os.Stdout, "Path_Exists error |%s|\n", err.Error() )
  os.Exit ( 1 )
  return false
}





func (tach * Tachyon) Subscribe ( channel Abstraction_Channel, topic_name string ) {
  topic, ok := tach.topics [ topic_name ]
  if ! ok {
    fp ( os.Stdout, "Tachyon error: no such topic: |%s|\n", topic_name )
    os.Exit(1)
  }
  topic.subscribe ( channel )
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
      
      // Topics can't have any of the following keywords as their names.

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
        top  := New_Topic ( tach, name )
        tach.topics [ name ] = top
        // Let the Bulletin Board know about it.
        tach.requests_to_bb <- msg
        // Tell the App that the topic has been created.
        tach.Responses <- Message { "response" : "new_topic",
                                    "name"     : name }



      case "start abstractors" :
        fp ( os.Stdout, "tachyon: starting abstractors.\n" )
        for _, a := range tach.abstractors {
          fp ( os.Stdout, "tachyon: starting |%s|\n", a.Name )
          go a.Run ( tach, a )
        }



      case "bb_request" :
        fp ( os.Stdout, "Bulletin Board request: |%#v|\n", msg )
        // Forward to the Bulletin Board.
        tach.requests_to_bb <- msg



      default :
        fp ( os.Stdout, "tachyon: tach_requests error: unknown request |%#v|\n", msg )
        os.Exit ( 1 )
    }
  }

  fp ( os.Stdout, "tach_requests exiting.\n" )
}




// Stamp the Abstraction with the number of nanoseconds 
// since Bill Joy came to Ann Arbor.
func (a * Abstraction) Timestamp () {
  now := time.Now()
  a.ID.Creation_Time = float64 ( now.UnixNano() ) / 1000000000.0
}





func (a * Abstraction) Add_To_Genealogy ( id * Abstraction_ID ) {
  a.Genealogy = append ( a.Genealogy, id )
}





func (a * Abstraction) Print_Genealogy ( ) {
  for _, id := range a.Genealogy {
    fp ( os.Stdout, "%s  ", id.Abstractor_Name )
  }
}





