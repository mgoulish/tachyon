
package tachyon

import (
         "os"
       )



//=================================================================
//  Public
//=================================================================


type Topic struct {

  name  string

  // This is the channel that all Abstractors use 
  // that produce abstractions for this Topic.
  inputs chan * Msg

  // No storage yet. At first, the topic is just a 
  // multicast message server.
  // storage [] * Msg

  subscribers [ ] chan * Msg
}





func New_Topic ( name string ) ( * Topic ) {
  top := & Topic { name                  : name,
                   inputs                : make ( chan * Msg, 100 ),
                   subscribers           : make ( [ ] chan * Msg, 0 ),
                 }
  go top.listen ( ) 
  return top
}





func ( top * Topic ) subscribe ( channel chan * Msg ) {
  top.subscribers = append ( top.subscribers, channel )
}





func ( top * Topic ) listen ( ) {
  for {
    msg, more := <- top.inputs 

    if ! more {
      break
    }

    for _, subscriber := range top.subscribers {
      subscriber <- msg
    }
  }
  
  fp ( os.Stdout, "topic |%s| : listener quitting.\n", top.name )
}





