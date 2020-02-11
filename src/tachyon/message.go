
package tachyon


import (
         "os"
       )





type Message struct {
  Info [] string
  Data [] interface{}
}





func ( msg * Message ) Size ( ) ( size uint64 ) {
  for _, data := range msg.Data {

    switch v := data.(type) {

      case []   byte :
        size += uint64(len(v))

      case    string :
        size += uint64(len(v))


      default :
        fp ( os.Stdout, "MDEBUG Message.Size cannot yet handle type %T.\n", data )
    }
  }

  return size
}





