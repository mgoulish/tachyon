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

  switch img.Image_type {
    case v.Image_type_gray16 :
      img.Write_gray16_to_tif ( tif_file_name )

    case v.Image_type_rgba :
      img.Write_rgba_to_tif ( tif_file_name )

    default :
      fp ( os.Stdout, "whd2tif error: Can't handle image type %s yet.", v.Image_type_name(img.Image_type) )

  }
}





