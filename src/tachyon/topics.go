
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
  inputs chan Message

  // No storage yet. At first, the topic is just a 
  // multicast message server.
  // storage [] * Msg

  subscribers [ ] chan Message
}





func New_Topic ( name string ) ( * Topic ) {
  top := & Topic { name         : name,
                   inputs       : make (     chan Message, 100 ),
                   subscribers  : make ( [ ] chan Message, 0 ),
                 }
  go top.listen ( ) 
  return top
}





func ( top * Topic ) subscribe ( subscriber_channel chan Message ) {

  // Add the subscriber's channel to my list.
  top.subscribers = append ( top.subscribers, subscriber_channel )

  // Send a confirmation message as the first message
  // on the subscriber's channel. 
  // NOTE : all subscribers to topics must undesratnd that 
  //        the first message they will receive will be a
  //        confirmation message -- not a 'real' message.
  subscriber_channel <- Message { "response" : "subscribed",
                                  "topic"    : top.name }
}



// By calling 'post', the given message 
// is pushed out to all subscribers.

func ( top * Topic ) post ( msg Message ) {
  for _, s := range top.subscribers {
    s <- msg
  }
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





