package main

import (
         "os"

         t "tachyon"
       )



func threshold ( tach * t.Tachyon, me * t.Abstractor ) ( ) {
  // To subscribe to our topic, we must supply
  // the channel that the topic will use to communicate
  // to us.
  my_input_channel := make ( chan t.Message, 10 ) 

  // Send the request.
  tach.Requests <- t.Message { "request" : "subscribe",
                               "topic"   : me.Subscribed_Topics[0],
                               "channel" : my_input_channel }

  // Now read messages that the topic sends me.
  message_count := 0
  for {
    msg := <- my_input_channel
    message_count ++

    if message_count == 1 {
      // The first message on my channel should be a confirmation
      // from Tachyon that we are subscribed to the correct channel.
      if msg["response"] != "subscribed" {
        // Something bad happened.
        fp ( os.Stdout, "App: %s: error: got this message |%#v|\n", me.Name, msg )
        break
      }
      fp ( os.Stdout, "App: %s: got subscription confirmation.\n", me.Name )
    } else {
      // This is a real message.
      fp ( os.Stdout, "App: %s: got msg!\n", me.Name )

      histo, ok := msg["data"].([768]uint32)
      if ! ok {
        fp ( os.Stdout, "App: %s: did not get uint32 array.\n", me.Name )
        os.Exit ( 1 )
      }
      
      fp ( os.Stdout, "App: %s got a histo!  of size %d\n", me.Name, len(histo) )

      max_val, max_pos := find_max ( histo )
      
      min_val := max_val
      for pos := max_pos; pos < 768; pos ++ {
        if histo[pos] > min_val {
          fp ( os.Stdout, "App: threshold at %d\n", pos )
          break
        }

        min_val = histo[pos]
      }
    }
  }
  
  fp ( os.Stdout, "App: %s exiting.\n", me.Name )
}





