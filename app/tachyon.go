package main

import (
         "fmt"
         "os"
         _ "os/exec"
         "time"

         t "tachyon"
       )


var fp = fmt.Fprintf


var image_topic string
var histo_topic string


func main ( ) {

  // Create a new Tachyon, and start listening for responses
  tach := t.New_Tachyon ( )
  go responses ( tach )

  // Add its first topic: image.
  // This will be the topic that the basic sensor posts to.
  // The listener will start up the Abstractors when this
  // topic has been created.
  image_topic = "image"
  tach.Requests <- & t.Msg { []t.AV { { "new_topic", image_topic } } }

  histo_topic = "histo"
  tach.Requests <- & t.Msg { []t.AV { { "new_topic", histo_topic } } }

  fp ( os.Stdout, "Hit 'enter' to quit.\n" )
  var s string
  fmt.Scanf ( "%s", & s )
}





func responses ( tach * t.Tachyon ) {
  abstractors_started := false
  requested_topics    := 2
  created_topics      := 0

  for {
    msg := <- tach.Responses

    if msg.Data[0].Attr == "new_topic" {
      created_topics ++
      fp ( os.Stdout, "MDEBUG responses: %d created topics.\n" , created_topics)
    }

    if created_topics >= requested_topics && abstractors_started == false {
      fp ( os.Stdout, "MDEBUG starting abstractors.\n" )
      abstractors_started = true
      // Create the image sensor, and tell it 
      // which topic it should post to.
      fp ( os.Stdout, "App: All topics have been created. Starting Abstractors.\n" )
      go sensor    ( tach )
      go histogram ( tach ) 
      go smoothing ( tach )
    }
  }
}





// Create histograms of the incoming images.

func histogram ( tach * t.Tachyon ) ( ) {
  // To subscribe to our topic, we must supply
  // the channel that the topic will use to communicate
  // to us.
  my_input_channel := make ( chan * t.Msg, 10 ) 

  // Send the request.
  tach.Requests <- & t.Msg { []t.AV { { "subscribe", image_topic },
                                      { "channel",   my_input_channel }}}

  // Now read messages that the topic sends me.
  message_count := 0
  for {
    msg := <- my_input_channel
    message_count ++

    if message_count == 1 {
      // The first message on my channel should be a confirmation
      // from Tachyon that we are subscribed to the correct channel.
      if msg.Data[0].Attr != "subscribed" {
        // Something bad happened.
        fp ( os.Stdout, "App: histogram: error: |%s|\n", msg.Data[0].Attr )
        break
      }
      fp ( os.Stdout, "App: histogram: got subscription confirmation.\n" )
    } else {
      // This is a real message.
      fp ( os.Stdout, "App: histogram: got msg!\n" )
      image, ok := t.Get_Val_From_Msg("data", msg).(*t.Image)
      if ! ok {
        fp ( os.Stdout, "tach_input error: post data does not contain an image.\n" )
        continue
      }

      var sum uint32
      var histo [768] uint32

      for x := uint32(0); x < image.Width; x ++ {
        for y := uint32(0); y < image.Height; y ++ {
          r, g, b, _ := image.Get ( x, y )
          sum = uint32(r) + uint32(g) + uint32(b)
          histo[sum] ++
        }
      }


      tach.Requests <- & t.Msg { []t.AV {{ "post", histo_topic},
                                         { "data", histo},
                                         { "length", 768}}}

      
      /*
        This code writes data and then uses the gnuplot_script
        in this directory to create a gnuplot graph of every histogram.
        I'm commenting it out for now, but leaving it here because
        I would like to systematize something like this in the near
        future for all Abstractors.
        How?

      f, _ := os.Create("./data")
      for i := 0; i < 768; i ++ {
        // fp ( os.Stdout, "%d : %d\n", i, histo[i] )
        fmt.Fprintf ( f, "%d %d\n", i, histo[i] )
      }
      f.Close()
      
      cmd  := "gnuplot"
      args := []string{"./gnuplot_script"}
      exec.Command(cmd, args...).Run()

      cmd  = "mv"
      args = []string{"./plot.jpg", fmt.Sprintf("./graphs/%04d.jpg", message_count)}
      exec.Command(cmd, args...).Run()
      */
    }
  }
  
  // TODO How should we really log stuff?
  //      It should include the Tachyon as well as the App,
  //      all the Abstractors, timestamps, everything.
  fp ( os.Stdout, "App: histogram exiting.\n" )
}




// The Sensor is what gets everything started.
// It is a lowest-level Abstractor activated (for now) by a timer.
// It is responsible for abstracting physics (photons) into images.
// In a production system, this function would interface with an actual camera.
// Currently, it simulates a real sensor by just reading image files off a disk.
func sensor ( tach * t.Tachyon ) {
  for i := 2; i < 101; i ++ {
    time.Sleep ( time.Second )
    // These images are frames that I split out of the video 
    // that is checked in as part of this project.
    image_file_name := fmt.Sprintf ( "/home/annex_2/vision_data/apollo/docking_with_lem/image-%04d.jpg", i )
    image := t.Read_Image ( image_file_name )
    tach.Requests <- & t.Msg { []t.AV { {"post", image_topic}, {"data", image} } }
  }
}





