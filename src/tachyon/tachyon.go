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
    msg := <- requests
    fp ( os.Stdout, "MDEBUG BB |%#v| gets message: |%v|\n", bb, msg )
  }
}





