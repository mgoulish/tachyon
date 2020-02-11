package main

import (
         "flag"
         "fmt"
         "os"
         "time"

         "tachyon"
       )



var fp = fmt.Fprintf





func print_errors ( errs chan string ) {
  for {
    err := <- errs
    fp ( os.Stdout, "APP received tachyon error: |%s|\n", err )
  }
}





func handle_responses ( tach * tachyon.Tachyon, i_am_initiating   bool ) {
  for {
        response := <- tach.Responses 

        response_type := response.Info[0]
        switch response_type {
          // The CNX has been created...
          case "listen", "connect" :
            if response.Info[1] == "success" {
              cnx_id := response.Info[3]
              //...so request the SSN.
              tach.Requests <- & tachyon.Message { Info : []string { "session", cnx_id } }
            }


          // The SSN has been created...
          case "session_created" :
            cnx_id   := response.Data[0].(int)
            ssn_id   := response.Data[1].(int)
            //...so request the LNKs.
            if i_am_initiating {
              tach.Requests <- & tachyon.Message { Info : []string { "send_LNK", 
                                                                      fmt.Sprintf("%d", cnx_id), 
                                                                      fmt.Sprintf("%d", ssn_id), 
                                                                      "my_addr" } }
            } else {
              tach.Requests <- & tachyon.Message { Info : []string { "recv_LNK", 
                                                                      fmt.Sprintf("%d", cnx_id), 
                                                                      fmt.Sprintf("%d", ssn_id), 
                                                                      "my_addr" } }
            }


          // The send LNK has been created. Start using it.
          case "send_LNK" :
            //fp ( os.Stdout, "App : start sending!\n" )
            cnx_id  := response.Info[1] 
            ssn_id  := response.Info[2]
            addr    := response.Info[3]
            channel := response.Data[0].(chan * tachyon.Message)

            go send ( tach,
                      cnx_id,
                      ssn_id,
                      addr,
                      channel,
                      )


          // The recv LNK has been created. Start using it.
          case "recv_LNK" :
            fp ( os.Stdout, "App : start receiving!\n" )
            cnx_id  := response.Info[1] 
            ssn_id  := response.Info[2]
            addr    := response.Info[3]
            channel := response.Data[0].(chan * tachyon.Message)

            go recv ( tach,
                      cnx_id,
                      ssn_id,
                      addr,
                      channel,
                    )


          default :
            fp ( os.Stdout, "App : unknown response |%s|\n",  response_type )
        }
  }
}





func send ( tach * tachyon.Tachyon,
            cnx_id  string,
            ssn_id  string,
            addr    string,
            channel chan * tachyon.Message,
            ) {

  n_messages  := 1000000
  // payload     := "Hello, World!"
  payload     := "0123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789"
  payload_len := len ( payload )
  fp ( os.Stdout, "payload_len == %d\n", payload_len )

  //msg := & tachyon.Message { Info : []string { addr, payload },
                             //Data : []interface{} { payload } }
  
  bytes_sent := 0
  first_send := time.Now()

  for count := 0; count < n_messages; count ++ {
    channel <- & tachyon.Message { Info : []string { addr },
                                   Data : []interface{} { payload } }
    bytes_sent += payload_len
  }

  channel <- & tachyon.Message { Info : []string { addr },
                                 Data : []interface{} { "done" } }
  
  fp ( os.Stdout, "MDEBUG last send: %s\n", time.Now() )
  fp ( os.Stdout, "bytes sent : %d     first_send : %s\n",  bytes_sent, first_send )
}





func recv ( tach * tachyon.Tachyon,
            cnx_id  string,
            ssn_id  string,
            addr    string,
            channel chan * tachyon.Message,
             ) {

  var last_recv time.Time
  bytes_received := 0

  for {
    msg := <- channel 
    // TODO -- the assymettry between outgoing and incoming messages
    // will get fixed once we have framing.
    str := msg.Data[0].(string)
    last_recv = time.Now()
    bytes_received += len ( str )
    //fp ( os.Stdout, "MDEBUG %d\n", bytes_received )

    if str == "done" {
      fp ( os.Stdout, "bytes rcvd : %d     last_recv : %s\n",  bytes_received, last_recv )
    }
  }
}





func main ( ) {

  initiate_ptr := flag.Bool ( "initiate", false, "initiate the connection")
  flag.Parse ( )
  i_am_initiating := * initiate_ptr

  // Start Tachyon, and start listening for errors and responses.
  tach := tachyon.New_Tachyon ( )
  go print_errors ( tach.Errors )
  go handle_responses ( tach, i_am_initiating )


  // Make our initial requests, to listen to a port, or to connect
  // to that port, depending on which flag we got on the command line.
  // All the action is in handle_responses() after this point.
  if i_am_initiating {
    tach.Requests <- & tachyon.Message { Info : []string { "connect", 
                                                            "127.0.0.1", 
                                                            "9090", 
                                                            "10"} }
  } else {
    tach.Requests <- & tachyon.Message { Info : []string { "listen", "9090" } }
  }


  for {
    time.Sleep ( 10 * time.Second )
  }
}





