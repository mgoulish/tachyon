package main



import (
         "fmt"
         "os"
         "time"
         
         t "tachyon"
       )


var fp = fmt.Fprintf



// The Abstractor sends an input request informing the BB 
// of what data type it's interested in, and providing the
// channel that the BB can use to send Abstractions of that
// type to it. Then it starts listening on that channel.
func abstractor ( name string, bb_channel t.Message_Channel ) {
  input_channel := make ( t.Message_Channel, 5 )

  input_request := t.Message { Type : "input_request",
                               Data : map[string]interface{} { "type"    : "image",
                                                               "channel" : input_channel} }
  fp ( os.Stdout, "Abstractor %s sending input request.\n", name )
  bb_channel <- input_request

  fp ( os.Stdout, "Abstractor %s listening for input.\n", name )
  for {
    image := <- input_channel
    fp ( os.Stdout, "Abstractor %s received image: |%#v|\n", image )
  }
}





// Main starts the Bulletin Board and then starts one Abstractor,
// passing it the BB's channel.
func main ( ) {
  bb_channel := t.Start_Bulletin_Board ( )
  go abstractor ( "a1", bb_channel )

  for {
    time.Sleep ( 10 * time.Second )
  }
}





