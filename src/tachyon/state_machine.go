
package tachyon

import (
         "os"
       )










//==========================================================
// Clients
//
//==========================================================
type client struct {
  name       string
  action     chan Message
  transition chan Message
}



// Little helper types for clarity.
type state_t  string
type action_t string



//==========================================================
// State Machines 
// 
//==========================================================
type state_machine struct {
  name       string
  states  [] state_t
  clients [] client

  // The 'transitions' field, given a current state, returns a map. 
  // That map, given an action, returns the next state.
  // So it's a 2D map, that you access like this:
  //      next_state := transitions [ current_state ] [ action ]
  transitions map [ state_t ] map [ action_t ] state_t 
}





func new_state_machine ( name string ) ( * state_machine ) {
  sm := new ( state_machine )
  sm.name        = name
  sm.clients     = make ( [] client, 0 )
  sm.transitions = make ( map [ state_t ] map [ action_t ] state_t )
  
  return sm
}





func ( sm * state_machine ) add_state ( state state_t ) {

  sm.transitions [ state ] = make ( map [ action_t ] state_t )
  fp ( os.Stdout, "MDEBUG added state |%s|\n", state )
}



