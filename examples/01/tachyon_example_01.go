package main

import (
         "fmt"
         "os"
         "time"

         t "tachyon"
       )


var fp = fmt.Fprintf
var start_time time.Time
var n_messages int

var sum_channel chan int





func main ( ) {

  n_topics    := 10
  n_receivers := 10
  n_messages   = 100000

  tach := t.New_Tachyon ( )

  sum_channel = make ( chan int, 10000 )

  // Ten topics.
  topics := []string { "a", "b", "c", "d", "e", "f", "g", "h", "i", "j" }


  // For each topic, create it in Tachyon and then
  // create the receivers for it.
  // They will subscribe themselves.
  for _, topic := range topics {
    tach.Requests <- & t.Msg { []t.AV { { "new_topic", topic } } }
    for i := 1; i <= n_receivers; i ++ {
      go receiver ( i, topic, tach )
    }
  }


  // Send messages to each topics.
  start_time = time.Now()
  for _, topic := range topics {
    for i := 0; i < n_messages; i ++ {
      tach.Requests <- & t.Msg { []t.AV { { topic, i } } }
    }
  }
  // Count the messages that we just sent.
  sum_channel <- n_topics * n_messages


  // Read from the sum channel until we have 
  // heard from everybody and totalled up everything.
  total_message_transmissions := 0
  client_count := 0
  for {
    message_count := <- sum_channel
    total_message_transmissions += message_count
    client_count ++

    if client_count >= 1 + n_topics * n_receivers {
      break
    }
  }

  duration := time.Now().Sub ( start_time ) 
  seconds := duration.Seconds()
  fp ( os.Stdout, 
       "%d transmissions in %f seconds == %f TPS.\n", 
       total_message_transmissions, 
       seconds,
       float64(total_message_transmissions) / seconds )
}





func receiver ( id int, topic string, tach * t.Tachyon ) {
  channel := make ( chan * t.Msg, 10 )

  // Ask Tachyon to subscribe us to this topic.
  tach.Requests <- & t.Msg { []t.AV { { "subscribe", topic },
                                      { "channel",   channel }}}
  
  message_count := 0
  for {
    <- channel
    message_count ++
    if message_count == n_messages {
      sum_channel <- n_messages
      break
    }
  }
}





