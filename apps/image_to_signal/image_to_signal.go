package main
  
import (
         "fmt"
         "os"

         v "vision"
       )


var fp = fmt.Fprintf





func main ( ) {

  if len(os.Args) < 2 {
    fp ( os.Stdout, "\nimage_to_signal error: need file name\n\n" )
    os.Exit ( 1 )
  }

  file_name := os.Args[1]
  img := v.Read ( file_name )

  for x := uint32(0); x < img.Width; x ++ {
    sum := 0
    for y := uint32(0); y < img.Height; y ++ {
      g := img.Get_gray16 ( x, y )
      // We count black pixels, not white.
      if g == 0 {
        sum ++
      }
    }
    fp ( os.Stdout, "%d\n", sum )
  }
}





