package main

import (
         "os"

         t "tachyon"
       )


// Create histograms of the incoming images.

func histogram ( tach * t.Tachyon, me * t.Abstractor ) ( ) {

  var id uint64


  message_count := 0
  for {
    input_abstraction := <- me.Input
    msg := input_abstraction.Msg

    message_count ++

    fp ( os.Stdout, "App: histogram: got msg!  \n")
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
    a := & t.Abstraction { ID  : t.Abstraction_ID { Abstractor_Name : me.Name, ID : id },
                           Msg : t.Message { "request" : "post",
                                             "topic"   : "histogram",
                                             "data"    : histo } }
    a.Timestamp()

    // My genealogy is the entire genealogy of my input
    // abstraction, plus its own ID.
    for _, id := range input_abstraction.Genealogy {
      a.Add_To_Genealogy ( id )
    }
    a.Add_To_Genealogy ( & input_abstraction.ID )

    me.Output <- a
  }

  fp ( os.Stdout, "App: histogram exiting.\n" )
}





