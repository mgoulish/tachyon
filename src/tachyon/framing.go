
package tachyon


import (
         "bytes"
         "encoding/binary"
         "os"
       )



// A frame contains the actual bytes as they came in 
// or as they are to be sent out on the wire.
type frame struct {
  frame_type string
  data       bytes.Buffer
}



// A frame descriptor is the structure that describes a
// a frame that is about to be encoded
type frame_descriptor struct {
  frame_type    string
  fields     [] frame_descriptor_field
}



type frame_descriptor_field struct {
  name          string
  default_value interface{}
  value         interface{}
}





const (  
        amqp_datacode_transfer        = 0x14
        amqp_datacode_uint0           = 0x43
        amqp_datacode_smalluint       = 0x52
        amqp_datacode_smallulong      = 0x53
        amqp_datacode_applicationdata = 0x75
        amqp_datacode_vbin8           = 0xA0
        amqp_datacode_list32          = 0xD0
      )





func ( f * frame ) dump () {
  fp ( os.Stdout, "\n\n\nframe dump\n" )
  for i := 0; i < f.data.Len(); i ++ {
    fp ( os.Stdout, "  %d : %2x\n", i, f.data.Bytes()[i] )
  }
  fp ( os.Stdout, "end frame dump\n\n\n" )
}





func ( f * frame ) write_transfer_fields ( ) {

  //fp ( os.Stdout, "MDEBUG what is len at start of write_transfer_fields? : %d\n", f.data.Len() )

  f.data.WriteByte ( byte(amqp_datacode_list32) )
  // Remember the buffer's current length, so we 
  // can come back to this point and write the correct 
  // field info length.
  field_info_length := uint32(18)
  binary.Write ( & f.data, binary.BigEndian, field_info_length )

  field_count := uint32(4)
  binary.Write ( & f.data, binary.BigEndian, field_count )

  // field : handle
  f.data.WriteByte ( byte(amqp_datacode_uint0) )
  field_info_length ++

  // field : delivery_id
  f.data.WriteByte ( byte(amqp_datacode_smalluint) )          
  f.data.WriteByte ( byte(1) )      // XXX TEMP -- delivery ID

  // field : delivery_tag 
  f.data.WriteByte ( byte(amqp_datacode_vbin8) ) 
  f.data.WriteByte ( byte(8) )
  var delivery_tag uint64
  binary.Write ( & f.data, binary.BigEndian, delivery_tag )

  // field : message_format
  f.data.WriteByte ( byte(amqp_datacode_uint0) )  // format 0
}




// This makes a transfer frame.
func enframe ( msg * Message ) ( f * frame, err error ) {

  f = new ( frame )
  f.frame_type = "transfer"

  // Start out with placeholder for the frame header.
  var preamble uint64
  err = binary.Write ( & f.data, binary.BigEndian, preamble )
  if err != nil {
    return nil, err
  }

  f.data.WriteByte ( byte(0) )
  f.data.WriteByte ( byte(amqp_datacode_smallulong) )
  f.data.WriteByte ( byte(amqp_datacode_transfer) )

  f.write_transfer_fields ( )

  // XXX Assume string content!
  f.data.WriteByte ( byte(0) )
  f.data.WriteByte ( byte(amqp_datacode_smallulong) )
  f.data.WriteByte ( byte(amqp_datacode_applicationdata) )

  f.data.WriteByte ( byte(amqp_datacode_vbin8) )
  f.data.WriteByte ( byte(len(msg.Data[0].(string)) ) )
  f.data.WriteString ( msg.Data[0].(string) )


  // Now we can create the frame preamble.
  preamble_buf := bytes.NewBuffer ( f.data.Bytes()[:0] )

  var data_offset uint8
  var frame_type  uint8 
  var channel     uint16

  //length      := uint32(4 + 1 + 1 + 2 + f.data.Len())
  length      := uint32(f.data.Len())
  data_offset = 2
  frame_type  = 0 // XXX
  channel     = 0 // XXX

  // Encode all fields in the preamble.
  // XXX log errors
  err = binary.Write ( preamble_buf, binary.BigEndian, length )
  if err != nil {
    return nil, err
  }

  err = binary.Write ( preamble_buf, binary.BigEndian, data_offset )
  if err != nil {
    return nil, err
  }

  err = binary.Write ( preamble_buf, binary.BigEndian,  frame_type )
  if err != nil {
    return nil, err
  }

  err = binary.Write ( preamble_buf, binary.BigEndian,     channel )
  if err != nil {
    return nil, err
  }

  return f, nil
}




// Turn a transfer frame into a Message.
func deframe ( f * frame ) ( msg * Message, err error ) {

  var length      uint32
  // var data_offset uint8
  // var frame_type  uint8
  var channel     uint16

  frame_data := f.data.Bytes()

  /*
  for foo := 0; foo < len(frame_data); foo ++ {
    fp ( os.Stdout, "MDEBUG %d : %.2x\n", foo, frame_data[foo] )
  }
  */

  // buf := bytes.NewReader ( f.data.Bytes() )
  //--------------------------------------------------------
  // Frame Header
  //--------------------------------------------------------
  length += uint32 ( frame_data [ 0 ] ) << 24
  length += uint32 ( frame_data [ 1 ] ) << 16
  length += uint32 ( frame_data [ 2 ] ) << 8
  length += uint32 ( frame_data [ 3 ] )
  //fp ( os.Stdout, "MDEBUG deframe: length == %0x\n", length )


  //data_offset := uint8 ( frame_data [ 4 ] )
  //fp ( os.Stdout, "MDEBUG deframe: data_offset == %0x\n", data_offset )

  //frame_type := uint8 ( frame_data [ 5 ] )
  //fp ( os.Stdout, "MDEBUG deframe: frame_type == %0x\n", frame_type )

  channel += uint16 ( frame_data [ 6 ] ) << 8
  channel += uint16 ( frame_data [ 7 ] )
  //fp ( os.Stdout, "MDEBUG deframe: channel == %0x\n", channel )

  var string_length byte
  string_length = uint8 ( frame_data [ 38 ] )
  //fp ( os.Stdout, "MDEBUG deframe: string_length == %0x\n", string_length )

  var str string
  current_pos := 39
  stop        := current_pos + int(string_length)
  for i := current_pos; i < stop; i ++ {
    str +=  string ( frame_data [ i ] )
  }

  // fp ( os.Stdout, "MDEBUG here's the string: |%s|\n", str )

  msg = & Message { Data : []interface{} { str } }

  return msg, nil
}










