package main

import (
         "os"
         "os/exec"
         "fmt"

         t "tachyon"
       )


func smoothing ( tach * t.Tachyon, me * t.Abstractor ) {
  var id uint64

  // To subscribe to our topic, we must supply
  // the channel that the topic will use to 
  // communicate to us.
  my_input_channel := make ( chan * t.Abstraction, 10 )

  // Send the request.
  tach.Requests <- t.Message { "request" : "subscribe",
                               "topic"   : me.Subscribed_Topics[0],
                               "channel" : my_input_channel }

  // Now read messages that the topic sends me.
  message_count := 0
  var saved, smoothed    [768]uint32
  var saved_energy, min_energy, new_energy int64

  for {
    abstraction := <- my_input_channel
    msg := abstraction.Msg
    message_count ++

    fp ( os.Stdout, "App: %s: got msg!\n", me.Name )

    histo, ok := msg["data"].([768]uint32)
    if ! ok {
      fp ( os.Stdout, "App: smoothing: did not get uint32 array.\n" )
      os.Exit ( 1 )
    }
    
    _, min_energy = count_reversals ( histo )
    fp ( os.Stdout, "smooth: reversal energy original : %d\n", min_energy )
    display_histo ( "original.jpg", histo )

    smoothed = smooth ( histo )
    _, new_energy = count_reversals ( smoothed )
    fp ( os.Stdout, "smooth: reversal energy after smooth %d: %d\n", 1, new_energy )
    display_histo ( "smoothed_1.jpg", smoothed )

    if new_energy >= min_energy {
      goto post
    } else {
      for i := 2; i < 5; i ++ {
        min_energy = new_energy
        // Save the current smoothed histogram
        // We save it here because sometimes (I don't know why)
        // when we re-smooth, the reversal energy actually rises
        // slightly. When that happens, we want to post the *previous*
        // smoothed histogram as the smoothest.
        saved = smoothed
        saved_energy = min_energy
        // And smooth it again.
        smoothed = smooth ( smoothed )
        _, new_energy = count_reversals ( smoothed )
        fp ( os.Stdout, "smooth: reversal energy after smooth %d: %d\n", i, new_energy )
        display_histo ( fmt.Sprintf("smoothed_%d.jpg", i), smoothed )

        if new_energy >= min_energy {
          fp ( os.Stdout, "smooth: done smoothing. Break 2.\n" )
          goto post
        }
      }
    }

    post :

    id ++
    // Post the smoothed histogram !
    if saved_energy < new_energy {
      fp ( os.Stdout, "smooth: posting smoothed histogram with reversal energy %d\n", saved_energy )
      a := & t.Abstraction { ID  : t.Abstraction_ID { Abstractor_Name : me.Name, ID : id },
                             Msg : t.Message { "request" : "post",
                                               "topic"   : me.Output_Topic,
                                               "data"    : saved } }
      a.Timestamp()
      tach.Abstractions <- a
    } else {
      fp ( os.Stdout, "smooth: posting smoothed histogram with reversal energy %d\n", new_energy )
      a := & t.Abstraction { ID  : t.Abstraction_ID { Abstractor_Name : me.Name, ID : id },
                             Msg : t.Message { "request" : "post",
                                               "topic"   : me.Output_Topic,
                                               "data"    : smoothed } }

      a.Timestamp()
      tach.Abstractions <- a
    }
  }
}





func count_reversals ( histo [768] uint32 ) ( int, int64 ) {

  first_derivative := make ( []int64, 768 )

  // Take the first derivative.
  for i := 1; i < 768; i ++ {
    first_derivative[i] = int64(histo[i]) - int64(histo[i-1])
  }

  // Now count the reversals of direction.
  reversals := 0
  energy    := int64(0)

  for i := 1; i < 768; i ++ {
    if first_derivative[i] * first_derivative[i-1] < 0 {
      reversals ++
      if first_derivative[i-1] > 0 {
        energy += first_derivative[i-1]
      } else {
        energy -= first_derivative[i-1]
      }
    }
  }

  return reversals, energy
}





func smooth ( histo [768] uint32 ) ( [768] uint32 ) {
  smoothed := histo

  half_kernel := uint32(5)

  // Do the convolution sum half to the left and
  // half to the right 

  for i := uint32(half_kernel); i < 768 - (half_kernel + 1); i ++ {
    sum := uint32(0)
    for j := i - half_kernel; j <= i + half_kernel; j ++ {
      sum += histo[j]
    }
    smoothed[i] = sum / ((2 * half_kernel) + 1)
  }
  return smoothed
}





func display_histo ( file_name string, histo [768] uint32 ) {

  f, _ := os.Create("./data")
  for i := 0; i < 768; i ++ {
    fmt.Fprintf ( f, "%d %d\n", i, histo[i] )
  }
  f.Close()

  cmd  := "gnuplot"
  args := []string{"./gnuplot_script"}
  exec.Command(cmd, args...).Run()

  cmd  = "mv"
  args = []string{"./plot.jpg", file_name}
  exec.Command(cmd, args...).Run()
}





func find_max ( histo [768] uint32 ) ( uint32, uint32 ) {

  var max_val, max_pos uint32

  max_val = 0
  for i := uint32(0); i < 768; i ++ {
    if histo[i] > max_val {
      max_val = histo[i]
      max_pos = i
    }
  }
  return max_val, max_pos
}
