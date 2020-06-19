package main

import (
         "fmt"
         "os"
         _ "os/exec"
         "time"

         t "tachyon"
       )


var fp = fmt.Fprintf


var image_topic          string
var histo_topic          string
var smoothed_histo_topic string
var n_topics             int


func main ( ) {

  fp ( os.Stdout, "Tachyon starting...\n" )
  fp ( os.Stdout, "Hit 'enter' to quit.\n" )
  time.Sleep ( 3 * time.Second )

  // Create a new Tachyon, and start listening for responses
  tach := t.New_Tachyon ( )
  go responses ( tach )

  // Add its first topic: image.
  // This will be the topic that the basic sensor posts to.
  // The listener will start up the Abstractors when this
  // topic has been created.
  image_topic = "image"
  tach.Requests <- & t.Msg { []t.AV { { "new_topic", image_topic } } }
  n_topics ++

  histo_topic = "histo"
  tach.Requests <- & t.Msg { []t.AV { { "new_topic", histo_topic } } }
  n_topics ++

  smoothed_histo_topic = "smoothed_histo"
  tach.Requests <- & t.Msg { []t.AV { { "new_topic", smoothed_histo_topic } } }
  n_topics ++

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
      fp ( os.Stdout, "MDEBUG responses: %d created topics.\n" , created_topics)
    }

    if created_topics >= n_topics && abstractors_started == false {
      fp ( os.Stdout, "MDEBUG starting abstractors.\n" )
      abstractors_started = true
      fp ( os.Stdout, "App: All topics have been created. Starting Abstractors.\n" )
      go sensor    ( tach )
      go histogram ( tach ) 
      go smoothing ( tach )
    }
  }
}





