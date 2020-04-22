package main

import (
         "fmt"
         "os"
         "time"

         t "tachyon"
       )


var fp = fmt.Fprintf





func main ( ) {
  tach := t.New_Tachyon ( )

  tach.Requests <- & t.Msg { []t.AV { { "new_topic", "image" } } }

  time.Sleep ( time.Second * 5 )
  fp ( os.Stdout, "App exiting.\n" )
}





