
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
  inputs chan * Abstraction

  // No storage yet. At first, the topic is just a 
  // multicast message server.
  // storage [] * Msg
  subscribers [ ] chan * Abstraction
}





func New_Topic ( name string ) ( * Topic ) {
  top := & Topic { name         : name,
                   inputs       : make (     chan * Abstraction, 100 ),
                   subscribers  : make ( [ ] chan * Abstraction, 0 ),
                 }
  go top.listen ( ) 
  return top
}





func ( top * Topic ) subscribe ( subscriber_channel chan * Abstraction ) {
  // Add the subscriber's channel to my list.
  top.subscribers = append ( top.subscribers, subscriber_channel )
}





// By calling 'post', the given message 
// is pushed out to all subscribers.

func ( top * Topic ) post ( abstraction * Abstraction ) {
  for _, s := range top.subscribers {
    s <- abstraction
  }
}





func ( top * Topic ) listen ( ) {
  for {
    abstraction, more := <- top.inputs 

    if ! more {
      break
    }

    for _, subscriber := range top.subscribers {
      subscriber <- abstraction
    }
  }
  
  fp ( os.Stdout, "topic |%s| : listener quitting.\n", top.name )
}





