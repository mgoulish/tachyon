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
                     }
  n_topics = len(topics)
  for _, topic := range topics {
    tach.Requests <- & t.Msg { []t.AV { { "new_topic", topic } } }
  }

  //------------------------------------------
  // Make the Abstractors.
  //------------------------------------------

  // sensor ------------------------------------------------
  sensor := & t.Abstractor { Name         : "sensor",
                             Run          : sensor,
                             Output_Topic : "image",
                           } 
  tach.Requests <- & t.Msg { []t.AV{ {"add_abstractor", sensor} } }


  // histogram ------------------------------------------------
  histo  := & t.Abstractor { Name         : "histogram",
                             Run          : histogram,
                             Subscribed_Topics : []string{ "image" },
                             Output_Topic : "histogram",
                           } 
  tach.Requests <- & t.Msg { []t.AV{ {"add_abstractor", histo} } }


  // smoothing ------------------------------------------------
  smooth := & t.Abstractor { Name         : "smoothing",
                             Run          : smoothing,
                             Subscribed_Topics : []string{ "histogram" },
                             Output_Topic : "smoothed_histogram",
                           } 
  tach.Requests <- & t.Msg { []t.AV{ {"add_abstractor", smooth} } }



  // threshold ------------------------------------------------
  thresh := & t.Abstractor { Name         : "threshold",
                             Run          : threshold,
                             Subscribed_Topics : []string{ "smoothed_histogram" },
                           } 
  tach.Requests <- & t.Msg { []t.AV{ {"add_abstractor", thresh} } }



  // Quit when the user hits 'enter'.
  var s string
  fmt.Scanf ( "%s", & s )
}





func responses ( tach * t.Tachyon ) {
  abstractors_started := false
  created_topics      := 0

  for {
    msg := <- tach.Responses

    if msg.Data[0].Attr == "new_topic" {
      created_topics ++
      fp ( os.Stdout, "App: response: %d created topics.\n" , created_topics)
    }

    if created_topics >= n_topics && abstractors_started == false {
      abstractors_started = true
      fp ( os.Stdout, "App: starting abstractors.\n" )
      tach.Requests <- & t.Msg { []t.AV{ {"start abstractors", nil} } }
    }
  }
}





