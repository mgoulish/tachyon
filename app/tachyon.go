package main

import (
         "fmt"
         "os"
         _ "os/exec"
         "time"

         t "tachyon"
       )


var fp = fmt.Fprintf


var n_topics             int


func main ( ) {

  fp ( os.Stdout, "Tachyon starting...\n" )
  fp ( os.Stdout, "Hit 'enter' to quit.\n" )
  time.Sleep ( 3 * time.Second )

  // Create a new Tachyon, and start listening for responses
  tach := t.New_Tachyon ( )
  go responses ( tach )


  //------------------------------------------
  // Make the topics.
  //------------------------------------------
  topics := []string { "image",
                       "histogram",
                       "smoothed_histogram",
                       "threshold",
                     }
  n_topics = len(topics)
  for _, topic := range topics {
    tach.Requests <- t.Message { "request" : "new_topic",
                                 "name"    : topic }
  }

  // Quit when the user hits 'enter'.
  var s string
  fmt.Scanf ( "%s", & s )
}





func responses ( tach * t.Tachyon ) {
  abstractors_started := false
  created_topics      := 0

  for {
    msg := <- tach.Responses

    if msg["response"] == "new_topic" {
      created_topics ++
      fp ( os.Stdout, "App: response: %d created topics.\n" , created_topics)
    }

    if created_topics >= n_topics && abstractors_started == false {
      
      make_abstractors ( tach )

      abstractors_started = true
      fp ( os.Stdout, "App: starting abstractors.\n" )
      tach.Requests <- t.Message { "request": "start abstractors" }
    }
  }
}




func make_abstractors ( tach * t.Tachyon ) {

  /*---------------------------------------------------------------
    For an Abstractor input channel: 
      1. Create the channel.
      2. Ask Tachyon to subscribe it to the proper Topic.
      3. Give it to the Abstractor.

    For the output channel:
      1. Ask Tachyon for the Topic you want the Abstractor 
         to post to. (It will give you an Abstraction Channel.)
      2. Give that to the Abstractor as its output channel.
    
    The Abstractor also has to be added to Tachyon's list
    so that it will be able to start them all when everything
    is finished being wired up.
  ---------------------------------------------------------------*/



  // sensor ------------------------------------------------
  sensor := & t.Abstractor { Name   : "sensor",
                             Run    : sensor,
                             Output : tach.Get_Topic("image"),
                           } 
  tach.Requests <- t.Message { "request"    : "add_abstractor",
                               "abstractor" : sensor }


  // histogram ------------------------------------------------
  input_channel := make ( t.Abstraction_Channel, 100 )
  tach.Subscribe ( input_channel, "image" )
  histo  := & t.Abstractor { Name   : "histogram",
                             Run    : histogram,
                             Input  : input_channel,
                             Output : tach.Get_Topic("histogram"),
                           } 
  tach.Requests <- t.Message { "request"    : "add_abstractor",
                               "abstractor" : histo }


  // smoothing ------------------------------------------------
  input_channel = make ( t.Abstraction_Channel, 100 )
  tach.Subscribe ( input_channel, "histogram" )
  smooth := & t.Abstractor { Name   : "smoothing",
                             Run    : smoothing,
                             Input  : input_channel,
                             Output : tach.Get_Topic("smoothed_histogram"),
                           } 
  tach.Requests <- t.Message { "request"    : "add_abstractor",
                               "abstractor" : smooth }



  // threshold ------------------------------------------------
  input_channel = make ( t.Abstraction_Channel, 100 )
  tach.Subscribe ( input_channel, "smoothed_histogram" )
  thresh := & t.Abstractor { Name   : "threshold",
                             Run    : threshold,
                             Input  : input_channel,
                             Output : tach.Get_Topic("threshold"),
                             Log    : "./log",
                           } 
  tach.Requests <- t.Message { "request"    : "add_abstractor",
                               "abstractor" : thresh }

}





