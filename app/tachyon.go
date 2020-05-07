package main

import (
         "fmt"
         "os"

         t "tachyon"
       )


var fp = fmt.Fprintf




func main ( ) {

  // Create a new Tachyon, and start listening for responses
  tach := t.New_Tachyon ( )
  go responses ( tach )

  // Add its first topic: image.
  // This will be the topic that the basic sensor posts to.
  image_topic := "image"
  tach.Requests <- & t.Msg { []t.AV { { "new_topic", image_topic } } }

  // Create the image sensor, and tell it what topic 
  // it should subscribe itself to.
  go sensor ( tach, image_topic )

  fp ( os.Stdout, "Hit 'enter' to quit.\n" )
  var s string
  fmt.Scanf ( "%s", & s )
}





func responses ( tach * t.Tachyon ) {
  for {
    msg := <- tach.Responses

    fp ( os.Stdout, "App got a response: |%#v|\n", msg )
  }
}





// TODO -- this is nice, except, um.  It isn't a sensor. 
// A sensor would send *to* this topic.
// Change this to be the first receiver from the image
// topic, and make an actual sensor.
// How will that work?


func sensor ( tach * t.Tachyon, topic_name string ) ( ) {
  // To subscribe to our topic, we must supply
  // the channel that the topic will use to communicate
  // to us.
  from_topic := make ( chan * t.Msg, 10 ) 

  // Send the request.
  tach.Requests <- & t.Msg { []t.AV { { "subscribe", topic_name },
                                      { "channel",   from_topic }}}

  // Now read messages that the topic sends me.
  message_count := 0
  for {
    msg := <- from_topic
    message_count ++

    if message_count == 1 {
      if msg.Data[0].Attr != "subscribed" {
        // Something bad happened.
        fp ( os.Stdout, "App: snesor: error: |%s|\n", msg.Data[0].Attr )
        break
      }
      fp ( os.Stdout, "App: sensor: got subscription confirmation.\n" )
    } else
    {
      fp ( os.Stdout, "App: sensor: got msg: |%#v|\n", msg )
    }
  }
  
  // TODO How should we really log stuff?
  //      It should include the Tachyon as well as the App,
  //      all the Abstractors, timestamps, everything.
  fp ( os.Stdout, "App: sensor exiting.\n" )
}





