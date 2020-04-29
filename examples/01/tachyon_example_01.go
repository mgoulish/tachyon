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

  for i := 1; i <= 10; i ++ {
    go receiver ( i, tach )
  }

  time.Sleep ( time.Second )


  tach.Requests <- & t.Msg { []t.AV { { "image", 12 } } }

  time.Sleep ( time.Second * 100 )
  fp ( os.Stdout, "App exiting.\n" )
}





func receiver ( id int, tach * t.Tachyon ) {
  channel := make ( chan * t.Msg, 10 )

  fp ( os.Stdout, "App : receiver %d : subscribing,\n", id )
  tach.Requests <- & t.Msg { []t.AV { { "subscribe", "image" },
                                      { "channel", channel   } } }
  
  for {
    msg := <- channel
    fp ( os.Stdout, "App : receiver %d got a message !  |%#v\n", id, msg )
  }
}





