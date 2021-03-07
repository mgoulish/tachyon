package main
  
import (
         "fmt"
         "os"

         v "vision"
       )


var fp = fmt.Fprintf





func check ( err error ) {
  if err != nil {
    panic ( err )
  }
}





func main ( ) {

  if len(os.Args) < 3 {
    fp ( os.Stdout, "\ntif2whd error: usage : tif2whd <TIFF_FILE> <WHD_FILE>\n\n" )
    os.Exit ( 1 )
  }

  tif_file_name := os.Args[1]
  whd_file_name := os.Args[2]

  img := v.Read_tif_to_gray16 ( tif_file_name )
  img.Write ( whd_file_name )
}





