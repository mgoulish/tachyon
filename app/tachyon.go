package main



import (
         "fmt"
         "os"
         
         t "tachyon"
       )


var fp = fmt.Fprintf




func main ( ) {

  // Here's a small message.
  msg := t.Message { Type : "my_type", 
                     Data : map[string]interface{} { "Hello" : 12, "AI" : 12 } }
  fp ( os.Stdout, "Here's a message: %v\n", msg )

  // Now let's make a Bulletin Board and register a channel with it.
  bb := t.New_Bulletin_Board ( )
  var my_output_channel t.Message_Channel

  bb.Register_Channel ( "my_type", my_output_channel )
  bb.Register_Channel ( "my_type", my_output_channel )
}





