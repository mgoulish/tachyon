
package tachyon

import (
         "os"
       )





//==========================================================
// State Machines 
// 
//==========================================================
type state_machine struct {
  name       string
  states  [] string

  // The 'transitions' field, given a current state, returns a map. 
  // That map, given an event, returns the next state.
  // So it's a 2D map, that you access like this:
  //      next_state := transitions [ current_state ] [ event ]
  transitions map [ string ] map [ string ] string 

  // The entities that need state machines -- such as connections -- 
  // are cloud entities -- implemented as interacting networks of
  // goroutines. Any of those goroutines may cause or detect events 
  // that alter the object state.
  // Here is the channel they will use to communicate any such events.
  events chan string

  // When the state does change, each of the goroutines implementing 
  // the object must be told about the change. They each have a 
  // separate channel for these messages. ( They send me these channels
  // when they 'register' with me. )
  state_change_announcements [] chan string
}





func new_state_machine ( name string, states [] string ) ( * state_machine ) {
  sm := new ( state_machine )
  sm.name        = name
  sm.transitions = make ( map [ string ] map [ string ] string )
  sm.events      = make ( chan string, 100 )


  for _, state := range states {
    sm.add_state ( state )
  }

  sm.add_state ( "start" )
  sm.add_state ( "end" )
  
  return sm
}





func ( sm * state_machine ) add_state ( state string ) {

  sm.transitions [ state ] = make ( map [ string ] string )
  fp ( os.Stdout, "MDEBUG added state |%s|\n", state )
}





// Users of this state machine should be the goroutines that work
// together to implement the object whose state is being represented.
// For each goroutine 'user' of this state machine, I need a channel
// on which I can send state changes. In return, I give you the common
// channel that all such goroutines use to notify me of state-changing
// events.
func ( sm * state_machine ) add_user ( state_changes chan string ) ( events chan string ) {
  sm.state_change_announcements = append ( sm.state_change_announcements, 
                                           state_changes )
  return sm.events
}





func ( sm * state_machine ) add_transition ( current_state, event, new_state string ) {
  current_state_map, ok := sm.transitions [ current_state ]

  if ! ok {
    fp ( os.Stdout, 
         "state_machine.add_transition error : unknown state |%s|\n", 
         current_state )
    os.Exit ( 1 )
  }

  current_state_map [ event ] = new_state

  fp ( os.Stdout, "MDEBUG added transition |%s| --|%s|--> |%s|\n", current_state, event, new_state )
}





