
package vision



import (
        "fmt"
        "image"
        "os"

        tif "golang.org/x/image/tiff"
)





func ( img * Image ) Set_gray8 ( x, y uint32, g uint8 ) {
  address := x + y * img.Width 
  img.Pixels [ address ] = g
}





func ( img * Image ) Get_gray8 ( x, y uint32 ) ( g  uint8 ) {
  address := x + y * img.Width
  val := img.Pixels [ address ]
  return val
}





func (img * Image) Write_gray8_to_tif ( file_name string ) {

  output_file, err := os.Create ( file_name )
  if err != nil {
    panic ( fmt.Errorf ( "image_write error: |%s|\n", err.Error()))
  }
  defer output_file.Close()

  output_image := image.NewRGBA ( image.Rect(0, 0, int(img.Width), int(img.Height)) )

  for y := uint32(0); y < uint32(img.Height); y ++ {
    for x := uint32(0); x < uint32(img.Width); x ++ {
      g := img.Get_gray8 ( x, y )
      addr := 4 * ( (y * img.Width) + x )
      output_image.Pix [ addr     ] = g
      output_image.Pix [ addr + 1 ] = g
      output_image.Pix [ addr + 2 ] = g
      output_image.Pix [ addr + 3 ] = 255
    }
  }

  tif.Encode ( output_file, output_image, nil )
}





func ( img_1 * Image ) And_gray8 ( img_2 * Image ) ( result * Image ) {
  
  if img_1.Image_type != Image_type_gray8 || img_2.Image_type != Image_type_gray8 {
    panic ( "Images must both be gray8." )
  }

  if img_1.Height != img_2.Height || img_1.Width != img_2.Width {
    panic ( fmt.Errorf ( "Both images must be same size." ) )
  }

  result = New_image ( Image_type_gray8, img_1.Width, img_1.Height )

  var x, y uint32
  for y = 0; y < img_1.Height; y ++ {
    for x = 0; x < img_1.Width; x ++ {
      if img_1.Get_gray8(x, y) == 255 && img_2.Get_gray8(x, y) == 255 {
        result.Set_gray8 ( x, y, 255 )
      } else {
        result.Set_gray8 ( x, y, 0 )
      }
    }
  }

  return result
}





// Image is assumed to be binary -- 0 or 255.
// Result (in place) will be binary.
func ( img * Image ) Invert_gray8 ( ) {
  
  if img.Image_type != Image_type_gray8 {
    panic ( "Image must be gray8." )
  }

  var x, y uint32
  for y = 0; y < img.Height; y ++ {
    for x = 0; x < img.Width; x ++ {
      if img.Get_gray8(x, y) == 255 {
        img.Set_gray8 ( x, y, 0 )
      } else {
        img.Set_gray8 ( x, y, 255 )
      }
    }
  }
}





func ( src * Image ) Dilate ( ) ( dst * Image ) {
  if src.Image_type != Image_type_gray8 {
    panic ( "Image must be gray8." )
  }

  dst = New_image ( Image_type_gray8, src.Width, src.Height )

  var x, y uint32
  for y = 0; y < src.Height; y ++ {
    for x = 0; x < src.Width; x ++ {
      if src.Get_gray8(x, y) == 255 {
        if 255 == src.Get_gray8 ( x, y ) {
          dst.Set_neighborhood ( x, y, 255 )
        }
      }
    }
  }

  return dst
}





func ( img * Image ) Set_neighborhood ( x, y uint32, gray uint8 ) {
  var min_x, min_y, max_x, max_y uint32

  if y == 0 {
    min_y = 0
  } else {
    min_y = y - 1
  }

  if x == 0 {
    min_x = 0
  } else {
    min_x = x - 1
  }

  max_y = min_y + 3
  max_x = min_x + 3

  if max_y >= img.Height {
    max_y = img.Height
  }

  if max_x >= img.Width {
    max_x = img.Width
  }

  for y1 := min_y; y1 < max_y; y1 ++ {
    for x1 := min_x; x1 < max_x; x1 ++ {
      img.Set_gray8 ( x1, y1, gray )
    }
  }
}





