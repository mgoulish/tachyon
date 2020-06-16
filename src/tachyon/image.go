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
        _ "image/jpeg"
)





const (
        image_type_none    = uint32(iota)
        image_type_gray_8
        image_type_gray_16
        image_type_gray_32
        image_type_gray_64
        image_type_rgba
        image_type_float
)





type Image struct {
  Image_Type, Width, Height uint32
  Pixels [] byte
}





func Bytes_Per_Pixel ( image_type uint32 ) ( uint32 ) {
  switch image_type {
    case image_type_none :
      return 0

    case image_type_gray_8 :
      return 1

    case image_type_gray_16 :
      return 2

    case image_type_gray_32 :
      return 4

    case image_type_gray_64 :
      return 8

    case image_type_rgba :
      return 4

    case image_type_float :
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





func image_read ( file_name string ) ( img * Image ) {
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

  img = New_Image ( image_type_rgba, width, height )

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

