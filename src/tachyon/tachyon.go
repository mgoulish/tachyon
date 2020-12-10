package tachyon


import (
         "fmt"
         "os"
       )

var fp = fmt.Fprintf





type Message struct {
  Type string  
  Data map [ string ] interface{}
}





type Message_Channel chan Message





type Bulletin_Board struct {

  // These are channels that carry output from the BB -- 
  // which means that the Abstractors call them 'input channels'.
  // the key strings in this map are the Abstraction type names.
  // Each type name is associated with an array of channels, 
  // because there may be more than one Abstractor that produces 
  // that data type.
  output_channels map [string] []Message_Channel
}





func Start_Bulletin_Board ( ) ( Message_Channel ) {

  requests := make ( Message_Channel, 5 )
  go run_bulletin_board ( requests )
  return requests
}





func run_bulletin_board ( requests Message_Channel ) {
  bb := Bulletin_Board { output_channels : make ( map[string] []Message_Channel, 0 ) }

  for {
    request := <- requests

    if request.Type ==  "input_request" {
      channel_type, ok := request.Data["type"].(string)
      if ! ok {
        fp ( os.Stdout, "BB error: request contains no channel type: |%#v|\n", request )
        continue
      }

      abstractor_channel, ok := request.Data["channel"].(Message_Channel)
      if ! ok {
        fp ( os.Stdout, "BB error: request contains no channel: |%#v|\n", request )
        continue
      }

      bb.output_channels[channel_type] = 
        append ( bb.output_channels[channel_type], abstractor_channel )
      
      fp ( os.Stdout, "MDEBUG BB output channels are now: |%#v|\n", bb.output_channels )
    }
  }
}





