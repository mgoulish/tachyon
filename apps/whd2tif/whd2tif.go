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
    panic ( fmt.Errorf ( "usage : <WHD_FILE> <TIFF_FILE>" ) )
  }

  whd_file_name := os.Args[1]
  tif_file_name := os.Args[2]

  img := v.Read ( whd_file_name )
  img.Write_gray16_to_tif ( tif_file_name )
}





