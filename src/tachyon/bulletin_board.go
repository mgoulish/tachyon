
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
        topic_name := request["name"].(string)
        storage [ topic_name ] = make ( topic_array, 0 )
        fp ( os.Stdout, "BB added topic |%s|\n", topic_name )

      case abstraction := <- tach.abstractions_to_bb :
        topic_name := abstraction.Msg["topic"].(string)
        storage[topic_name] = append ( storage[topic_name], abstraction )
        fp ( os.Stdout, "BB stored abstraction for topic |%s| (%d)\n", topic_name, len(storage[topic_name]) )
    }
  }
}




