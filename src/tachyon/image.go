package tachyon

import (
        "image"
        "os"

        // Package image/jpeg is not used explicitly in the code below,
        // but is imported for its initialization side-effect, which allows
        // image.Decode to understand JPEG formatted images. Uncomment these
        // two lines to also understand GIF and PNG images:
        // _ "image/gif"
        // _ "image/png"
        "image/jpeg"
)





const (
        Image_Type_None    = uint32(iota)
        Image_Type_Gray_8
        Image_Type_Gray_16
        Image_Type_Gray_32
        Image_Type_Gray_64
        Image_Type_RGBA
        Image_Type_Float
)





type Image struct {
  Image_Type, Width, Height uint32
  Pixels [] byte
}





func Bytes_Per_Pixel ( image_type uint32 ) ( uint32 ) {
  switch image_type {
    case Image_Type_None :
      return 0

    case Image_Type_Gray_8 :
      return 1

    case Image_Type_Gray_16 :
      return 2

    case Image_Type_Gray_32 :
      return 4

    case Image_Type_Gray_64 :
      return 8

    case Image_Type_RGBA :
      return 4

    case Image_Type_Float :
      return 8
  }

  fp ( os.Stdout, "Bytes_Per_Pixel error: unknown image type: %d\n", image_type )
  os.Exit ( 1 )
  return 0
}





func New_Image ( image_type, width, height uint32 ) ( * Image ) {
  bpp := Bytes_Per_Pixel ( image_type )
  return & Image { Image_Type : image_type,
                   Width      : width,
                   Height     : height,
                   Pixels     : make ( []byte, width * height * bpp ),
                 }
}





func ( img * Image ) Set ( x, y uint32, r, g, b, a byte ) {
  bpp := Bytes_Per_Pixel ( img.Image_Type )
  address := bpp * ( x + y * img.Width )
  img.Pixels [ address ] = r
  address ++
  img.Pixels [ address ] = g
  address ++
  img.Pixels [ address ] = b
  address ++
  img.Pixels [ address ] = a
  address ++
}





func ( img * Image ) Get ( x, y uint32 ) ( r, g, b, a byte ) {
  bpp := Bytes_Per_Pixel ( img.Image_Type )
  address := bpp * ( x + y * img.Width )
  r = img.Pixels [ address ]
  address ++
  g = img.Pixels [ address ]
  address ++
  b = img.Pixels [ address ]
  address ++
  a = img.Pixels [ address ]

  return r, g, b, a
}





func Read_Image ( file_name string ) ( img * Image ) {
  reader, err := os.Open ( file_name )
  if err != nil {
    fp ( os.Stdout, "image_read error: |%s|\n", err.Error() )
    os.Exit ( 1 )
  }
  defer reader.Close()

  jpg, _, err := image.Decode ( reader )
  if err != nil {
    fp ( os.Stdout, "image_read error: |%s|\n", err.Error() )
    os.Exit ( 1 )
  }
  bounds := jpg.Bounds()
  width  := uint32(bounds.Max.X)
  height := uint32(bounds.Max.Y)

  img = New_Image ( Image_Type_RGBA, width, height )

  var x, y uint32
  var r, g, b, a byte

  for y = 0; y < uint32(bounds.Max.Y); y ++ {
    for x = 0; x < uint32(bounds.Max.X); x ++ {
      R, G, B, A := jpg.At ( int(x), int(y) ).RGBA ( )
      r = byte(R >> 8)
      g = byte(G >> 8)
      b = byte(B >> 8)
      a = byte(A >> 8)
      img.Set ( x, y, r, g, b, a )
    }
  }

  return img
}





func (img * Image) Write ( file_name string ) {

  output_file, err := os.Create ( file_name )
  if err != nil {
    fp ( os.Stdout, "image_write error: |%s|\n", err.Error() )
    os.Exit ( 1 )
  }
  defer output_file.Close()

  output_image := image.NewRGBA ( image.Rect(0, 0, int(img.Width), int(img.Height)) )

  for y := uint32(0); y < uint32(img.Height); y ++ {
    for x := uint32(0); x < uint32(img.Width); x ++ {
      r, g, b, a := img.Get ( x, y )
      addr := 4 * ( (y * img.Width) + x )
      output_image.Pix [ addr     ] = r
      output_image.Pix [ addr + 1 ] = g
      output_image.Pix [ addr + 2 ] = b
      output_image.Pix [ addr + 3 ] = a
    }
  }

  jpeg.Encode ( output_file, output_image, nil )
}





