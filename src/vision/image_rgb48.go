
package vision



import (
        _ "golang.org/x/image/tiff"
)





func ( img * Image ) Set_rgb48 ( x, y uint32, r, g, b uint16 ) {
  address := 6 * ( x + y * img.Width )

  // r ---------------------
  img.Pixels [ address     ] = byte(r >> 8)
  img.Pixels [ address + 1 ] = byte(r&0xFF)

  // g ---------------------
  img.Pixels [ address + 2 ] = byte(g >> 8)
  img.Pixels [ address + 3 ] = byte(g&0xFF)

  // b ---------------------
  img.Pixels [ address + 4 ] = byte(b >> 8)
  img.Pixels [ address + 5 ] = byte(b&0xFF)
}





func ( img * Image ) Get_rgb48 ( x, y uint32 ) ( r, g, b  uint16 ) {
  address := 6 * ( x + y * img.Width )

  // r ---------------------
  r = uint16(img.Pixels [ address ])
  r <<= 8
  r += uint16(img.Pixels [ address + 1])

  // g ---------------------
  g = uint16(img.Pixels [ address + 2])
  g <<= 8
  g += uint16(img.Pixels [ address + 3])

  // b ---------------------
  b = uint16(img.Pixels [ address + 4])
  b <<= 8
  b += uint16(img.Pixels [ address + 5])

  return r, g, b
}





