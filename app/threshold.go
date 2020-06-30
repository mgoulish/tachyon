package main

import (
         "os"
         "fmt"

         t "tachyon"
       )



func threshold ( tach * t.Tachyon, me * t.Abstractor ) ( ) {
  var id uint64
  fp ( os.Stdout, "MDEBUG thresh starting.\n" )

  message_count := 0
  for {


    input_abstraction := <- me.Input
    msg := input_abstraction.Msg
    message_count ++
    fp ( os.Stdout, "App: %s: got msg!\n", me.Name )


    logging := false
    var my_logging_root string
    fp ( os.Stdout, "MDEBUG thresh: log: |%s|\n", me.Log )
    if me.Log != "" {
      my_logging_root = me.Log + "/threshold"
      if ! t.Path_Exists ( my_logging_root ) {
        os.Mkdir ( my_logging_root, 0700 )
      }
      my_logging_root += fmt.Sprintf ( "/%d", message_count )
      os.Mkdir ( my_logging_root, 0700 )
      logging = true
    }


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
    a := & t.Abstraction { ID  : t.Abstraction_ID { Abstractor_Name : me.Name, ID : id },
                           Msg : t.Message { "request" : "post",
                                             "topic"   : "threshold",
                                             "data"    : thresh } }
    a.Timestamp()

    // My genealogy is the entire genealogy of my input
    // abstraction, plus its own ID.
    for _, id := range input_abstraction.Genealogy {
      a.Add_To_Genealogy ( id )
    }
    a.Add_To_Genealogy ( & input_abstraction.ID )

    me.Output <- a

    // TODO -- this is disgusting. Give the logging stuff its own fn.
    
    if logging {
      fp ( os.Stdout, "MDEBUG threshold is logging.\n" )
      
      // Write a text message to my log directory.
      log_file := my_logging_root + "/threshold"
      f, _ := os.Create ( log_file )
      fp ( f, "threshold %d\n", thresh ) 
      f.Close()
      
      // Find the original image in my genealogy.
      var antecedent * t.Abstraction_ID
      for _, antecedent = range a.Genealogy {
        if antecedent.Abstractor_Name == "sensor" {
          break
        }
      }
      tach.Requests <- t.Message { "request"    : "bb_request",
                                   "topic"      : "image",
                                   "abstractor" : antecedent.Abstractor_Name,
                                   "ID"         : antecedent.ID,
                                   "reply_to"   : me.Input }
      
      response := <- me.Input
      fp ( os.Stdout, "App: threshold got a response! |%#v|\n", response )

      img := response.Msg["data"].(* t.Image)
      thresholded_img := t.New_Image ( t.Image_Type_RGBA, img.Width, img.Height )

      fp ( os.Stdout, "App: threshold got img %d x %d\n", img.Width, img.Height )
      // Now apply the threshold.
      for x := uint32(0); x < img.Width; x ++ {
        for y := uint32(0); y < img.Height; y ++ {
          r, g, b, _ := img.Get ( x, y )
          luminance := uint64(r + g + b)
          var result byte
          if luminance >= thresh {
            result = 255
          } else {
            result = 0
          }
          thresholded_img.Set ( x, y, result, result, result, 255 )
        }
      }
      
      // And write the image file to my log directory.
      image_file_name := fmt.Sprintf ( "%s/%04d.jpg", my_logging_root, message_count )
      thresholded_img.Write ( image_file_name )
    }
  }
  
  fp ( os.Stdout, "App: %s exiting.\n", me.Name )
}





