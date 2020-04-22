package main

import (
         "fmt"
         "os"

         "tachyon"
       )


var fp = fmt.Fprintf





func main ( ) {
  tach := tachyon.New_Tachyon ( )
  fp ( os.Stdout, "MDEBUG tach: |%#v|\n", tach )
  fp ( os.Stdout, "App exiting.\n" )
}





