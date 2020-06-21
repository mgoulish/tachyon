package main

import (
         "os"

         t "tachyon"
       )



func threshold ( tach * t.Tachyon, me * t.Abstractor ) ( ) {
  var id uint64

  // To subscribe to our topic, we must supply
  // the channel that the topic will use to 
  // communicate to us.
  my_input_channel := make ( chan * t.Abstraction, 10 ) 

  // Send the request.
  tach.Requests <- t.Message { "request" : "subscribe",
                               "topic"   : me.Subscribed_Topics[0],
                               "channel" : my_input_channel }

  // Now read messages that the topic sends me.
  message_count := 0
  for {
    abstraction := <- my_input_channel
    msg := abstraction.Msg
    message_count ++
    fp ( os.Stdout, "App: %s: got msg!\n", me.Name )

    histo, ok := msg["data"].([768]uint32)
    if ! ok {
      fp ( os.Stdout, "App: %s: did not get uint32 array.\n", me.Name )
      os.Exit ( 1 )
    }
    
    fp ( os.Stdout, "App: %s got a histo!  of size %d\n", me.Name, len(histo) )

    var thresh uint64
    max_val, max_pos := find_max ( histo )
    min_val := max_val
    for pos := max_pos; pos < 768; pos ++ {
      if histo[pos] > min_val {
        fp ( os.Stdout, "App: threshold at %d\n", pos )
        thresh = uint64(pos)
        break
      }
      min_val = histo[pos]
    }
    
    id ++
    tach.Abstractions <- & t.Abstraction { Abstractor_Name : me.Name,
                                           Abstraction_ID  : id,
                                           Msg             : t.Message { "request" : "post",
                                                                         "topic"   : me.Output_Topic,
                                                                         "data"    : thresh } }
  }
  
  fp ( os.Stdout, "App: %s exiting.\n", me.Name )
}





