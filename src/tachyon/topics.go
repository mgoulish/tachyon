
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
  input_channel chan * Abstraction

  // No storage yet. At first, the topic is just a 
  // multicast message server.
  // storage [] * Msg
  subscribers [ ] chan * Abstraction

  tach * Tachyon
}





func New_Topic ( tach * Tachyon, name string ) ( * Topic ) {
  top := & Topic { name          : name,
                   input_channel : make (     chan * Abstraction, 100 ),
                   subscribers   : make ( [ ] chan * Abstraction, 0 ),
                   tach          : tach,
                 }
  go top.listen ( ) 
  return top
}





func ( top * Topic ) subscribe ( subscriber_channel chan * Abstraction ) {
  // Add the subscriber's channel to my list.
  top.subscribers = append ( top.subscribers, subscriber_channel )
}





/*--------------------------------------------
  Listen for messages coming into this Topic.
  Send each one to the Bulletin Board first, 
  and then send out to each subscriber.
--------------------------------------------*/
func ( top * Topic ) listen ( ) {
  for {
    abstraction, more := <- top.input_channel 

    //fp ( os.Stdout, "MDEBUG Topic %s got abstraction: |%#v|\n", top.name, abstraction )

    if ! more {
      break
    }

    top.tach.abstractions_to_bb <- abstraction

    for _, subscriber := range top.subscribers {
      subscriber <- abstraction
    }
  }
  
  fp ( os.Stdout, "topic |%s| : listener quitting.\n", top.name )
}





