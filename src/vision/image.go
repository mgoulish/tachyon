
package vision



import (
         "bytes"
         "encoding/binary"
         "fmt"
         "io/ioutil"
         "os"
       )


var fp = fmt.Fprintf




const (
        Image_type_none    = uint32(iota)
        Image_type_gray16
        Image_type_rgba
        Image_type_rgb48
        Image_type_gray8
)





type Image struct {
  Image_type, Width, Height, n_bytes uint32
  Pixels [] byte
}





func New_image ( image_type, width, height uint32 ) ( * Image ) {
  bpp := Bytes_per_pixel ( image_type )
  return & Image { Image_type : image_type,
                   Width      : width,
                   Height     : height,
                   Pixels     : make ( []byte, width * height * bpp ),
                   n_bytes    : width * height * bpp,
                 }
}





func Bytes_per_pixel ( image_type uint32 ) ( uint32 ) {
  switch image_type {
    case Image_type_none :
      return 0

    case Image_type_gray16 :
      return 2

    case Image_type_rgba :
      return 4

    case Image_type_rgb48 :
      return 6
    
    case Image_type_gray8 :
      return 1
  }

  panic ( fmt.Errorf ( "Bytes_Per_Pixel error: unknown image type: %d\n", image_type ) )
  return 0
}





func Image_type_name ( image_type uint32 ) ( string ) {
  switch image_type {
    case Image_type_none :
      return "none"
    
    case Image_type_gray16 :
      return "gray16"

    case Image_type_rgba :
      return "rgba"

    case Image_type_rgb48 :
      return "rgb48"
    
    case Image_type_gray8 :
      return "gray8"
    
    default :
      return "unknown"
  }
}





func Read ( file_name string ) ( * Image ) {

  fn := "Read"

  buf, err := ioutil.ReadFile ( file_name )
  check ( err, fn )

  var image_type, width, height uint32

  buf_reader := bytes.NewBuffer ( buf )
  err = binary.Read ( buf_reader, binary.BigEndian, & image_type )
  check ( err, fn )

  err = binary.Read ( buf_reader, binary.BigEndian, & width )
  check ( err, fn )

  err = binary.Read ( buf_reader, binary.BigEndian, & height )
  check ( err, fn )

  img := New_image ( image_type, width, height )

  err = binary.Read ( buf_reader, binary.BigEndian, & img.Pixels )
  check ( err, fn )

  return img
}





func ( img * Image ) Write ( file_name string ) {
  fn := "Write"

  f, err := os.Create ( file_name )
  check ( err, fn )
  defer f.Close ( )

  var buf bytes.Buffer
  err = binary.Write ( & buf,
                       binary.BigEndian,
                       []uint32{img.Image_type, img.Width, img.Height} )
  check ( err, fn )
  err = binary.Write ( & buf,
                       binary.BigEndian,
                       img.Pixels )
  _, err = f.Write ( buf.Bytes() )
  check ( err, fn )
}





func check ( err error, fn string ) {
  if err != nil {
    panic ( fmt.Errorf ( "Image %s error: |%s|\n", fn, err.Error() ) )
  }
}





