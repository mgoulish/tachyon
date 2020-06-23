
package tachyon

import (
         "os"
       )



//=================================================================
//  Public
//=================================================================


type topic_array [] * Abstraction 





func bulletin_board ( tach * Tachyon ) {
  
  storage := make ( map [string] topic_array )

  for {
    
    select {

      case request := <- tach.requests_to_bb :
        
        fp ( os.Stdout, "MDEBUG request to BB: |%#v|\n", request )
        request_type := request["request"]

        switch request_type {
          case "new_topic" :
            topic_name := request["name"].(string)
            storage [ topic_name ] = make ( topic_array, 0 )
            fp ( os.Stdout, "BB added topic |%s|\n", topic_name )

          case "bb_request" :
            fp ( os.Stdout, "MDEBUG BB got a request!  |%#v|\n", request )
            topic_name := request["topic"].(string)
            topic := storage[topic_name]
            target_ID := request["ID"]
            found := false
            for _, abs := range topic {
              if abs.ID.ID == target_ID {
                found = true
                reply_to := request["reply_to"].(chan * Abstraction)
                reply_to <- abs 
              }
            }
            if ! found {
              fp ( os.Stdout, "bulletin_board error: abs not found\n" )
              os.Exit ( 1 )
            }

          default:
            fp ( os.Stdout, "bulletin_board error: unknown request: |%s|\n", request_type )
            os.Exit ( 1 )
        }


      case abstraction := <- tach.abstractions_to_bb :
        topic_name := abstraction.Msg["topic"].(string)
        storage[topic_name] = append ( storage[topic_name], abstraction )
        fp ( os.Stdout, "BB stored abstraction for topic |%s| (%d)\n", topic_name, len(storage[topic_name]) )
    }
  }
}





