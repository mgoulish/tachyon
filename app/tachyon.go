package main



import (
         "fmt"
         "time"
         
         t "tachyon"
       )


var fp = fmt.Fprintf




func main ( ) {

  msg := t.Message { Type : "my_type", 
                     Data : map[string]interface{} { "Hello" : 12, "AI" : 12 } }

  bb_channel := t.Start_Bulletin_Board ( )

  for {
    bb_channel <- msg
    time.Sleep ( 2 * time.Second )
  }
}





