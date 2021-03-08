
package vision



import (
        "image"
        "os"

        tif "golang.org/x/image/tiff"
)





func ( img * Image ) Set_rgba ( x, y uint32, r, g, b, a byte ) {
  address := 4 * ( x + y * img.Width )
  img.Pixels [ address ] = r
  address ++
  img.Pixels [ address ] = g
  address ++
  img.Pixels [ address ] = b
  address ++
  img.Pixels [ address ] = a
}





func ( img * Image ) Get_rgba ( x, y uint32 ) ( r, g, b, a byte ) {
  address := 4 * ( x + y * img.Width )
  r = img.Pixels [ address ]
  address ++
  g = img.Pixels [ address ]
  address ++
  b = img.Pixels [ address ]
  address ++
  a = img.Pixels [ address ]

  return r, g, b, a
}





func Read_rgba ( file_name string ) ( img * Image ) {
  reader, err := os.Open ( file_name )
  if err != nil {
    fp ( os.Stdout, "image_read error: |%s|\n", err.Error() )
    os.Exit ( 1 )
  }
  defer reader.Close()

  decoded_image, _, err := image.Decode ( reader )
  if err != nil {
    fp ( os.Stdout, "image_read error: |%s|\n", err.Error() )
    os.Exit ( 1 )
  }

  bounds := decoded_image.Bounds()
  width  := uint32(bounds.Max.X)
  height := uint32(bounds.Max.Y)

  img = New_image ( Image_type_rgba, width, height )

  var x, y uint32
  var r, g, b, a byte

  for y = 0; y < uint32(bounds.Max.Y); y ++ {
    for x = 0; x < uint32(bounds.Max.X); x ++ {
      R, G, B, A := decoded_image.At ( int(x), int(y) ).RGBA ( )
      r = byte(R >> 8)
      g = byte(G >> 8)
      b = byte(B >> 8)
      a = byte(A >> 8)
      img.Set_rgba ( x, y, r, g, b, a )
    }
  }

  return img
}





func (img * Image) Write_rgba_to_tif ( file_name string ) {

  output_file, err := os.Create ( file_name )
  if err != nil {
    fp ( os.Stdout, "image_write error: |%s|\n", err.Error() )
    os.Exit ( 1 )
  }
  defer output_file.Close()

  output_image := image.NewRGBA ( image.Rect(0, 0, int(img.Width), int(img.Height)) )

  for y := uint32(0); y < uint32(img.Height); y ++ {
    for x := uint32(0); x < uint32(img.Width); x ++ {
      r, g, b, a := img.Get_rgba ( x, y )
      addr := 4 * ( (y * img.Width) + x )
      output_image.Pix [ addr     ] = r
      output_image.Pix [ addr + 1 ] = g
      output_image.Pix [ addr + 2 ] = b
      output_image.Pix [ addr + 3 ] = a
    }
  }

  tif.Encode ( output_file, output_image, nil )
}





func ( img * Image ) Constant_rgba ( r, g, b, a byte ) ( ) {
  var x, y uint32
  for y = 0; y < img.Height; y ++ {
    for x = 0; x < img.Width; x ++ {
        img.Set_rgba ( x, y, r, g, b, a )
    }
  }
}





func ( img * Image ) Vertical_line_rgba ( x uint32, r, g, b, a byte ) ( ) {
  var y uint32
  for y = 0; y < img.Height; y ++ {
    img.Set_rgba ( x, y, r, g, b, a )
  }
}





