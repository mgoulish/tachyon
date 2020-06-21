package main

import (
         "os"

         t "tachyon"
       )


// Create histograms of the incoming images.

func histogram ( tach * t.Tachyon, me * t.Abstractor ) ( ) {

  var id uint64

  // To subscribe to our topic, we must supply
  // the channel that the topic will use to communicate
  // to us.
  my_input_channel := make ( chan * t.Abstraction, 10 ) 

  // Send subscription request.
  tach.Requests <- t.Message { "request" : "subscribe", 
                               "topic"   : me.Subscribed_Topics[0],
                               "channel" : my_input_channel }

  message_count := 0
  for {
    abstraction := <- my_input_channel
    msg := abstraction.Msg

    message_count ++

    fp ( os.Stdout, "App: histogram: got msg!\n" )
    image, ok := msg["data"].(*t.Image)
    if ! ok {
      fp ( os.Stdout, "App: histogram error: data does not contain a histogram.\n" )
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
    id ++
    tach.Abstractions <- & t.Abstraction { ID  : t.Abstraction_ID { Abstractor_Name : me.Name, ID : id },
                                           Msg : t.Message { "request" : "post",
                                                             "topic"   : me.Output_Topic,
                                                             "data"    : histo } }
  }
  fp ( os.Stdout, "App: histogram exiting.\n" )
}





