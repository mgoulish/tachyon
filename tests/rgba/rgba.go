package main
  
import (
         "fmt"

         v "vision"
       )


var fp = fmt.Fprintf





func check ( err error ) {
  if err != nil {
    panic ( err )
  }
}





func main ( ) {

  rgb_img  := v.New_image ( v.Image_type_rgba, 256, 240 )

  y_step  := rgb_img.Height / 6
  y_start := uint32(0)
  y_limit := y_step


  // R ---------------------------------------------
  for y := y_start; y < y_limit; y ++ {
    for x := uint32(0); x < rgb_img.Width; x ++ {
      rgb_img.Set_rgba ( x, y, 255, 0, 0, 255 )
    }
  }
  y_start += y_step
  y_limit += y_step


  // G ---------------------------------------------
  for y := y_start; y < y_limit; y ++ {
    for x := uint32(0); x < rgb_img.Width; x ++ {
      rgb_img.Set_rgba ( x, y, 0, 255, 0, 255 )
    }
  }
  y_start += y_step
  y_limit += y_step


  // B ---------------------------------------------
  for y := y_start; y < y_limit; y ++ {
    for x := uint32(0); x < rgb_img.Width; x ++ {
      rgb_img.Set_rgba ( x, y, 0, 0, 255, 255 )
    }
  }
  y_start += y_step
  y_limit += y_step


  // C ---------------------------------------------
  for y := y_start; y < y_limit; y ++ {
    for x := uint32(0); x < rgb_img.Width; x ++ {
      rgb_img.Set_rgba ( x, y, 0, 255, 255, 255 )
    }
  }
  y_start += y_step
  y_limit += y_step


  // M ---------------------------------------------
  for y := y_start; y < y_limit; y ++ {
    for x := uint32(0); x < rgb_img.Width; x ++ {
      rgb_img.Set_rgba ( x, y, 255, 0, 255, 255 )
    }
  }
  y_start += y_step
  y_limit += y_step


  // Y ---------------------------------------------
  for y := y_start; y < y_limit; y ++ {
    for x := uint32(0); x < rgb_img.Width; x ++ {
      rgb_img.Set_rgba ( x, y, 255, 255, 0, 255 )
    }
  }


  rgb_img.Write_rgba_to_tif ( "./rgba.tif" )
}





