package main

import (
         "fmt"
         "time"
         "os"

         t "tachyon"
       )


// The Sensor is what gets everything started.
// It is a lowest-level Abstractor activated (for now) by a timer.
// It is responsible for abstracting physics (photons) into images.
// In a production system, this function would interface with an actual camera.
// Currently, it simulates a real sensor by just reading image files off a disk.
func sensor ( tach * t.Tachyon, me * t.Abstractor ) {
  fp ( os.Stdout, "Abstractor %s starting.\n", me.Name )

  var id uint64

  for i := 2; i < 101; i ++ {
    time.Sleep ( time.Second )
    // These images are frames that I split out of the video 
    // that is checked in as part of this project.
    image_file_name := fmt.Sprintf ( "/home/annex_2/vision_data/apollo/docking_with_lem/image-%04d.jpg", i )
    image := t.Read_Image ( image_file_name )
    fp ( os.Stdout, "\n\n\nApp: sensor: %s\n", image_file_name )
    id ++
    tach.Abstractions <- & t.Abstraction { Abstractor_Name : me.Name,
                                           Abstraction_ID  : id,
                                           Msg             : t.Message { "request" : "post",
                                                             "topic"   : me.Output_Topic,
                                                             "data"    : image } }
  }
}





