package main

import (
         "fmt"
         "time"

         t "tachyon"
       )


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





