package main

import (
         "fmt"
         "os"
         "time"

         t "tachyon"
       )


var fp = fmt.Fprintf


var image_topic string


func main ( ) {

  // Create a new Tachyon, and start listening for responses
  tach := t.New_Tachyon ( )
  go responses ( tach )

  // Add its first topic: image.
  // This will be the topic that the basic sensor posts to.
  image_topic = "image"
  tach.Requests <- & t.Msg { []t.AV { { "new_topic", image_topic } } }

  fp ( os.Stdout, "Hit 'enter' to quit.\n" )
  var s string
  fmt.Scanf ( "%s", & s )
}





func responses ( tach * t.Tachyon ) {
  for {
    msg := <- tach.Responses

    // fp ( os.Stdout, "App got a response: |%#v|\n", msg )

    if msg.Data[0].Attr == "new_topic" && msg.Data[0].Val == image_topic {
      // Create the image sensor, and tell it 
      // which topic it should post to.
      fp ( os.Stdout, "App: image topis has been created. Starting sensor.\n" )
      go sensor ( tach, image_topic )
    }
  }
}





// Create histograms of the incoming images.

func histogram ( tach * t.Tachyon, topic_name string ) ( ) {
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
        fp ( os.Stdout, "App: histogram: error: |%s|\n", msg.Data[0].Attr )
        break
      }
      fp ( os.Stdout, "App: histogram: got subscription confirmation.\n" )
    } else
    {
      fp ( os.Stdout, "App: histogram: got msg: |%#v|\n", msg )
    }
  }
  
  // TODO How should we really log stuff?
  //      It should include the Tachyon as well as the App,
  //      all the Abstractors, timestamps, everything.
  fp ( os.Stdout, "App: histogram exiting.\n" )
}




func sensor ( tach * t.Tachyon, topic_name string ) {
  for i := 1; i < 101; i ++ {
    image_file_name := fmt.Sprintf ( "/home/annex_2/vision_data/apollo/docking_with_lem/image-%04d.jpg", i )
    fp ( os.Stdout, "MDEBUG image_file_name: |%s|\n", image_file_name )
    time.Sleep ( time.Millisecond * 10 )
  }
}





