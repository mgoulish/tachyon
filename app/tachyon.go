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
      fp ( os.Stdout, "App: image topic has been created. Starting Abstractors.\n" )
      go sensor    ( tach, image_topic )
      go histogram ( tach, image_topic )
    }
  }
}





// Create histograms of the incoming images.

func histogram ( tach * t.Tachyon, topic_name string ) ( ) {
  // To subscribe to our topic, we must supply
  // the channel that the topic will use to communicate
  // to us.
  my_input_channel := make ( chan * t.Msg, 10 ) 

  // Send the request.
  tach.Requests <- & t.Msg { []t.AV { { "subscribe", topic_name },
                                      { "channel",   my_input_channel }}}

  // Now read messages that the topic sends me.
  message_count := 0
  for {
    msg := <- my_input_channel
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
      fp ( os.Stdout, "App: histogram: got msg!\n" )

      image, ok := t.Get_Val_From_Msg("data", msg).(*t.Image)
      if ! ok {
        fp ( os.Stdout, "tach_input error: post data does not contain an image.\n" )
        continue
      }

      x := uint32(100)
      y := uint32(100)
      r, g, b, _ := image.Get ( x, y )
      fp ( os.Stdout, "histogram got an image!!! pixel at %d,%d is: %d,%d,%d\n", x, y, r, g, b )
    }
  }
  
  // TODO How should we really log stuff?
  //      It should include the Tachyon as well as the App,
  //      all the Abstractors, timestamps, everything.
  fp ( os.Stdout, "App: histogram exiting.\n" )
}




func sensor ( tach * t.Tachyon, topic_name string ) {
  for i := 1; i < 101; i ++ {
    time.Sleep ( time.Second )
    image_file_name := fmt.Sprintf ( "/home/annex_2/vision_data/apollo/docking_with_lem/image-%04d.jpg", i )
    image := t.Read_Image ( image_file_name )
    //x := uint32(100)
    //y := uint32(100)
    //r, g, b, _ := image.Get ( x, y )
    //fp ( os.Stdout, "MDEBUG got an image!!! pixel at %d,%d is: %d,%d,%d\n", x, y, r, g, b )
    tach.Requests <- & t.Msg { []t.AV { {"post", image_topic}, {"data", image} } }
  }
}





