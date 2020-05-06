package main

import (
         "fmt"
         "os"
         "time"

         t "tachyon"
       )


var fp = fmt.Fprintf




func main ( ) {

  // Create a new Tachyon, and start listening for responses
  tach := t.New_Tachyon ( )

  go responses ( tach )

  // Add it's first topic: image.
  // This will be the topic that the basic sensor posts to.
  tach.Requests <- & t.Msg { []t.AV { { "new_topic", "image" } } }

  time.Sleep ( 1000 * time.Second )
}





func responses ( tach * t.Tachyon ) {
  for {
    msg := <- tach.Responses

    fp ( os.Stdout, "responses got a response: |%#v|\n", msg )
  }
}





