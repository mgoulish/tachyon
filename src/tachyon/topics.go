
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
  Inputs chan * Message

  // No storage yet. At first, the topic is just a 
  // multicast message server.
  // storage [] * Message

  // A channel to use if you want to subscribe.
  // All incoming messages will be pushed to you 
  // immediately after being stored.
  // Your message must include your channel for results.
  Subscription_Requests chan * Message

  subscribers map [ string ] chan * Message
}





func New_Topic ( name string ) ( * Topic ) {
  top := & Topic { name                  : name,
                   Inputs                : make ( chan * Message, 100 ),
                   Subscription_Requests : make ( chan * Message, 100 ),
                   subscribers           : make ( map [ string ] chan * Message ),
                 }
  return top
}





func topic ( top * Topic ) {
  go top_listen_for_requests ( top )
  // go top_listen_for_inputs   ( top )
}





func top_listen_for_requests ( top * Topic ) {
  for {
    req, more := <- top.Subscription_Requests 

    if ! more {
      break
    }
    fp ( os.Stdout, "MDEBUG topic |%s| got request |%#v|\n", top.name, req )
  }
  
  fp ( os.Stdout, "MDEBUG topic |%s| quitting.\n", top.name )
}





