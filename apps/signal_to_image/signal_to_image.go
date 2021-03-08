package main
  
import (
         "fmt"
         "os"
         "io"

         v "vision"
       )


var fp = fmt.Fprintf



func main ( ) {

  if len(os.Args) < 3 {
    fp ( os.Stdout, "\nsignal_to_image error: need <input> <output>\n\n" )
    os.Exit ( 1 )
  }

  input_file_name  := os.Args[1]
  output_file_name := os.Args[2]

  file, err := os.Open ( input_file_name )

  if err != nil {
    fmt.Println ( err )
    os.Exit ( 1 )
  }

  var  val     int

  width  := uint32 ( 0 )
  height := uint32 ( 0 )

  // Go through once to learn the dimensions.
  for {
    _, err := fmt.Fscanf ( file, "%d", & val )
    if err != nil {
      if err == io.EOF {
        break 
      }
      fmt.Println ( err )
      os.Exit ( 1 )
    }

    width ++
    if uint32(val) > height {
      height = uint32(val);
    }
  }

  file.Seek ( 0, io.SeekStart )
  img := v.New_image ( v.Image_type_rgba, width, height )
  img.Constant_rgba ( 255, 255, 255,  255 )

  // Scan the numbers again, but this time use them to 
  // make the columns of the image.
  for x := uint32(0); x < width; x ++ {
    _, err := fmt.Fscanf ( file, "%d", & val )
    if err != nil {
      if err == io.EOF {
        break 
      }
      fmt.Println ( err )
      os.Exit ( 1 )
    }

    var y uint32
    for y = height - 1; val > 0; y-- {
      val --
      img.Set_rgba ( x, y, 0, 0, 0, 255 )
    }
  }

  img.Vertical_line_rgba ( 119, 255, 0, 0, 255 )

  img.Write ( output_file_name )
}





