package main

import (
         "flag"
         "fmt"
         "os"
         "strings"
         "sync"
         "time"

         "tachyon"
       )


var fp = fmt.Fprintf





func main ( ) {
  var wg sync.WaitGroup

  expected_frame_count := 1000  * 1000
  port := "5801"

  // The only asymmetry between the two instances of this 
  // program is whether we are initiating or accepting
  // the connection.
  init_ptr := flag.Bool ( "initiate", false, "Initiate the connection.")

  id_ptr := flag.Int ( "id", 0, "Provide a unique ID for each Tachyon.")

  // Get the command line args.
  flag.Parse ( )
  i_am_initiating := * init_ptr
  id := * id_ptr

  // This is the only function call into the 
  // Tachyon library. AFter this call, all 
  // interaction is through the channels in the
  // returned Tachyon structure.
  tach := tachyon.New_Tachyon ( id )

  // Start listening for errors first of all, just in case.
  go print_errors ( tach.Errors )

  // Exit only after the sender and receiver are both done.
  // Func responses will create them.
  wg.Add ( 2 )
  go responses ( tach, expected_frame_count, i_am_initiating, & wg )


  if i_am_initiating {
    fp ( os.Stdout, "App %d is running. I initiate the connection.\n", id )
    tach.Requests <- & tachyon.Message { Info : []string { "connect", 
                                                           "127.0.0.1", 
                                                           port, 
                                                           "30" } }
  } else {
    fp ( os.Stdout, "App %d is running. I accept the connection.\n", id )
    tach.Requests <-  & tachyon.Message { Info : []string { "listen", port } }
  }

  // Wait for the sender and receiver to exit.
  wg.Wait ( )
  fp ( os.Stdout, "App %d exiting.\n", id )
}





func responses ( tach * tachyon.Tachyon,
                 expected_frame_count int,
                 i_am_initiating bool,
                 wg * sync.WaitGroup ) {
  for {
    response, more := <- tach.Responses
    if ! more {
      break
    }

    switch response.Info[0] {

      case "start_receiving" :
        go receive ( expected_frame_count, 
                     tach.Incoming,
                     wg )

      case "start_sending" :
        go send ( expected_frame_count, 
                  tach.Outgoing,
                  wg )

      default :
        fp ( os.Stdout, "app error : unknown response type |%s|\n", response.Info )
    }
  }
}





func send ( expected_frame_count int, 
            outgoing_message_channel chan * tachyon.Message,
            wg * sync.WaitGroup ) {

  defer wg.Done()

  // Build up a 1024-byte message body.
  // A Tachyon message is always one frame:
  // an 8 byte header for the channel, and a
  // 1 KB body.
  str_08   := "abcdefgh"
  str_1024 := strings.Repeat ( str_08, 128 )
  channel_number := "00000012"
  bytes_per_frame := len(channel_number) + len(str_1024)

  message := & tachyon.Message{ Info: []string{channel_number}, Data: []interface{}{str_1024} }

  // Send all the messages.
  start_time := time.Now()
  sent_frames := 0
  for ; sent_frames < expected_frame_count; sent_frames ++ {
    outgoing_message_channel <- message
  }
  stop_time := time.Now()

  // Tell the user how we did.
  duration          := stop_time.Sub ( start_time )
  frames_per_second := float64(sent_frames) / duration.Seconds()
  bytes_per_second  := frames_per_second * float64(bytes_per_frame)

  fp ( os.Stdout, 
       "App %d sent %d frames in %f seconds == %f frames per second == %f bytes per second.\n",
       os.Getpid(), 
       sent_frames,
       duration.Seconds(),
       frames_per_second,
       bytes_per_second )

  // Closing this outgoing message channel 
  // tells Tachyon to close the network connection.
  // That closure will then cause the receiving side
  // to shut down its incoming message channel.
  // So between the two Apps, all 4 channels get 
  // shut down.
  close ( outgoing_message_channel )
}





func receive ( expected_frame_count int, 
               incoming_message_channel chan * tachyon.Message,
               wg * sync.WaitGroup ) {

  var start_time, stop_time time.Time

  defer wg.Done()

  frame_body_size  := 1024
  frame_count := 0

  // Receive incoming messages from Tachyon.
  for {
    msg, more := <- incoming_message_channel
    
    if frame_count == 0 {
      start_time = time.Now()
    }

    if ! more {
      stop_time = time.Now()
      break
    }

    n_bytes := len(msg.Data[0].(string))
    channel := msg.Info[0]

    // Make sure we got the right channel.
    if channel != "00000012" {
      fp ( os.Stdout, "App: channel error: |%s| at frame %d\n", channel, frame_count )
      break
    }

    // Make sure we got the right size frame.
    if n_bytes != frame_body_size {
      fp ( os.Stdout, "App: channel %s received non-frame %d bytes, at frame %d.\n", channel, n_bytes, frame_count )
      break
    }
    frame_count ++
  }

  duration          := stop_time.Sub ( start_time )
  frames_per_second := float64(frame_count) / duration.Seconds()
  bytes_per_second  := frames_per_second * 1032.0

  fp ( os.Stdout, 
       "App %d received %d frames in %f seconds == %f frames per second == %f bytes per second.\n", 
       os.Getpid(), 
       frame_count,
       duration.Seconds(),
       frames_per_second,
       bytes_per_second )
}





func print_errors ( errs chan * tachyon.Message ) {
  for {
    err, more := <- errs
    if ! more {
      break
    }
    fp ( os.Stdout, "App received tachyon error: |%s|\n", err )
  }
}





