package main

import (
         "os"

         t "tachyon"
       )


// Create histograms of the incoming images.

func histogram ( tach * t.Tachyon, me * t.Abstractor ) ( ) {
  // To subscribe to our topic, we must supply
  // the channel that the topic will use to communicate
  // to us.
  my_input_channel := make ( chan * t.Msg, 10 ) 

  // Send the request.
  tach.Requests <- & t.Msg { []t.AV { { "subscribe", me.Subscribed_Topics[0] },
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


      // Post the histogram !
      tach.Requests <- & t.Msg { []t.AV {{ "post", me.Output_Topic},
                                         { "data", histo},
                                         { "length", 768}}}

    }
  }
  
  fp ( os.Stdout, "App: histogram exiting.\n" )
}





