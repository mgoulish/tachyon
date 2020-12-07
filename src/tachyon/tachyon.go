package tachyon


import (
         "fmt"
         "os"
       )

var fp = fmt.Fprintf


type Message map [ string ] interface{}


type Message_Channel chan Message



type Bulletin_Board struct {
  message_channels map [string] []Message_Channel
}



func New_Bulletin_Board ( ) ( * Bulletin_Board ) {
  return & Bulletin_Board { message_channels : make ( map[string] []Message_Channel, 0 ) }
}



func ( bb * Bulletin_Board ) Register_Channel ( message_type string, 
                                                channel Message_Channel ) {
  bb.message_channels[message_type] = append ( bb.message_channels[message_type], channel )

  fp ( os.Stdout, "BB now has : \n" )
  for message_type, channels := range bb.message_channels {
    fp ( os.Stdout, "    message type: %s     channels: %d\n",  message_type, len(channels) )
  }
}





