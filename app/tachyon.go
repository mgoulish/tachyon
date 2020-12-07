package main



import (
         "fmt"
         "os"
         
         t "tachyon"
       )


var fp = fmt.Fprintf




func main ( ) {

  // Here's a little message.
  msg := t.Message { "Hello" : 12, "AI" : 12 }
  fp ( os.Stdout, "Hello, AI. Here's a message: %v\n", msg )

  // Now let's make a Bulletin Board and register a channel with it.
  bb := t.New_Bulletin_Board ( )
  var my_output_channel t.Message_Channel

  bb.Register_Channel ( "my_type", my_output_channel )
  bb.Register_Channel ( "my_type", my_output_channel )
}





