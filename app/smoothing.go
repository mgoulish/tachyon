package main

import (
         "os"

         t "tachyon"
       )


func smoothing ( tach * t.Tachyon ) {
  // To subscribe to our topic, we must supply
  // the channel that the topic will use to communicate
  // to us.
  my_input_channel := make ( chan * t.Msg, 10 )

  // Send the request.
  tach.Requests <- & t.Msg { []t.AV { { "subscribe", histo_topic },
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
    
      fp ( os.Stdout, "App: smoothing: got subscription confirmation.\n" )
    } else {
      // This is a real message.
      fp ( os.Stdout, "App: smoothing: got msg!\n" )

      histo, ok := t.Get_Val_From_Msg("data", msg).([768]uint32)
      if ! ok {
        fp ( os.Stdout, "App: smoothing: did not get uint32 array.\n" )
        os.Exit ( 1 )
      }

      fp ( os.Stdout, "MDEBUG histo[13] %d\n", histo[13] )
    }
  }
}





