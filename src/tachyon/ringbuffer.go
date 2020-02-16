
package tachyon


import ( 
         "os"
       )





// The head and tail numbers increase forever,
// but are only used modulo the RB length.
type ringbuffer struct {

  cnx        * connection 
  capacity     uint64
  head         uint64
  tail         uint64
  frame_end    uint64

  data      [] byte
}





func new_ringbuffer ( size uint64 ) ( * ringbuffer ) {
  rb := & ringbuffer { capacity : size,
                       data     : make ( [] byte, size ) }
  return rb
}





func (r * ringbuffer) set_cnx ( cnx * connection ) {
  r.cnx = cnx
}





func ( r * ringbuffer ) dump ( n int ) {
  count := 0
  for i := r.tail; i < r.head; i ++ {
    fp ( os.Stdout, "%d : %0x\n", i, r.data [ i % r.capacity ] )
    count ++
    if n > 0 && count >= n {
      break
    }
  }
}





func ( r * ringbuffer ) write ( src [] byte ) ( written bool ) {

  to_be_written := uint64(len(src))
  used          := r.head - r.tail
  remaining     := r.capacity - used

  if remaining < to_be_written {
    fp ( os.Stdout, "MDEBUG write fail : remaining %d to be written %d\n", remaining, to_be_written )
    return false
  }

  var i uint64
  for i = 0; i < to_be_written; i ++ {
    r.data [ (r.head + i) % r.capacity ] = src [ i ] 
  }
  r.head += to_be_written

  // Output as many frames as we have now.
  for {
    if false == r.output_frame() {
      break
    }
  }

  return true
}





func ( r * ringbuffer ) output_frame ( ) ( bool ) {

  // We are guaranteed to be reading at the beginning of a frame.

  // Do we have at least 4 bytes available?
  // If so, then they are the beginning of the frame header -- 
  // i.e. its length.
  available := r.head - r.tail
  if available < 4 {
    return false
  }

  // Get the frame_length.
  var frame_length uint64
  frame_length += uint64 ( r.data [ (r.tail + 3) % r.capacity ] )
  frame_length += uint64 ( r.data [ (r.tail + 2) % r.capacity ] ) << 8
  frame_length += uint64 ( r.data [ (r.tail + 1) % r.capacity ] ) << 16
  frame_length += uint64 ( r.data [ (r.tail + 0) % r.capacity ] ) << 24

  // If the entire frame is not yet in the ringbuffer,
  // return failure while not changing the tail pointer.
  // We peeked at the 4 frame_length-bytes without consuming them.
  if available < frame_length {
    return false
  }

  // The frame is all here. Copy it and send it out.
  f := new ( frame )
  var i uint64

  // Modulo every byte in case we wrap around.
  for i = 0; i < frame_length; i ++ {
    f.data.WriteByte ( r.data [ (r.tail + i) % r.capacity ] )
  }
  r.tail += frame_length

  // What SSN does this frame belong to ?
  ssn_id := 0 // XXX -- parse this out of the frame.
  to_ssn := r.cnx.session_map [ ssn_id ]
  to_ssn <- f

  return true
}





func ( r * ringbuffer ) read ( dst [] byte ) ( bool ) {

  to_be_read := uint64(len(dst))
  available := r.head - r.tail

  if available == 0 || available < to_be_read {
    fp ( os.Stdout, "MDEBUG ringbuffer.read fail %d available, %s desired.\n", available, to_be_read )
    return false
  }

  var i uint64
  for i = 0; i < to_be_read; i ++ {
    dst [ i ] = r.data [ (i % r.capacity) ]
  }

  r.tail += to_be_read

  return true
}








