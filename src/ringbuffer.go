
package tachyon


import ( 
         "os"
       )





// The head and tail numbers increase forever,
// but are only used modulo the RB length.
type Ringbuffer struct {

  capacity     uint64
  head         uint64
  tail         uint64


  data      [] byte
}





func New_Ringbuffer ( size uint64 ) ( * Ringbuffer ) {
  rb := & Ringbuffer { capacity : size,
                       data     : make ( [] byte, size ) }
  return rb
}





func (r * Ringbuffer) Capacity ( ) ( uint64 ) {
  return r.capacity
}



func (r * Ringbuffer) Available_For_Read ( ) ( uint64 ) {
  return r.head - r.tail
}



func (r * Ringbuffer) Available_For_Write ( ) ( uint64 ) {
  used := r.head - r.tail
  avail_for_write := r.capacity - used

  return avail_for_write
}



func (r * Ringbuffer) Total_Written ( ) ( uint64 ) {
  return r.head
}



func (r * Ringbuffer) Total_Read ( ) ( uint64 ) {
  return r.tail
}





func ( r * Ringbuffer ) Dump ( n int ) {
  count := 0
  for i := r.tail; i < r.head; i ++ {
    fp ( os.Stdout, "%d : %0x\n", i, r.data [ i % r.capacity ] )
    count ++
    if n > 0 && count >= n {
      break
    }
  }
}





func ( r * Ringbuffer ) Write ( src [] byte ) ( written bool ) {

  to_be_written := uint64(len(src))
  used          := r.head - r.tail
  remaining     := r.capacity - used

  if remaining < to_be_written {
    fp ( os.Stdout, "MDEBUG Write fail : remaining %d to be written %d\n", remaining, to_be_written )
    return false
  }

  var i uint64
  var endpoint = r.head + to_be_written;
  if endpoint < r.capacity {
    //fp ( os.Stdout, "MDEBUG RB end %d < cap %d\n", endpoint, r.capacity )
    copy ( r.data [ r.head : ], src )
  } else {
    for i = 0; i < to_be_written; i ++ {
      r.data [ (r.head + i) % r.capacity ] = src [ i ]
    }
  }

  r.head += to_be_written
  //fp ( os.Stdout, "MDEBUG RB.Write : wrote %d -- ends now %d  %d\n", to_be_written, r.tail, r.head,  )

  return true
}





func ( r * Ringbuffer ) Read ( dst [] byte ) ( bool ) {


  to_be_read := uint64(len(dst))
  available := r.head - r.tail

  if available == 0 || available < to_be_read {
    fp ( os.Stdout, "MDEBUG Ringbuffer.Read fail %d available, %s desired.\n", available, to_be_read )
    return false
  }

  var i uint64
  for i = 0; i < to_be_read; i ++ {
    dst [ i ] = r.data [ (i % r.capacity) ]
  }

  r.tail += to_be_read
  //fp ( os.Stdout, "MDEBUG RB.Read : read %d -- ends now %d  %d\n", to_be_read, r.tail, r.head )

  return true
}





func ( r * Ringbuffer ) Read_Frame ( ) ( * frame, bool, uint64 ) {


  f := new ( frame )

  // We are guaranteed to reading at the beginning of a frame,
  // since we only read out full frames. 
  // Do we have at least 4 bytes available?
  // If so, then they are the beginning of the frame header -- 
  // i.e. its length.
  available := r.head - r.tail
  if available < 4 {
    //fp ( os.Stdout, "MDEBUG header not here yet.\n" )
    return nil, false, 0
  }

  // r.Dump ( 0 )

  // Get the length.
  var length uint32
  length += uint32 ( r.data [ (r.tail + 3) % r.capacity ] )
  length += uint32 ( r.data [ (r.tail + 2) % r.capacity ] ) << 8
  length += uint32 ( r.data [ (r.tail + 1) % r.capacity ] ) << 16
  length += uint32 ( r.data [ (r.tail + 0) % r.capacity ] ) << 24

  // If the entire frame is not yet in the ringbuffer,
  // return failure while not changing the tail pointer.
  // We peeked at the 4 length-bytes without consuming them.
  if available < uint64(length) {
    //fp ( os.Stdout, "MDEBUG Ringbuffer.Read_Frame : frame not all here yet.\n" )
    return nil, false, 0
  }

  // The frame is all here. Copy it and send it out.
  f = new ( frame )
  var i uint64

  // Modulo every byte in case we wrap around.
  for i = 0; i < uint64(length); i ++ {
    f.data.WriteByte ( r.data [ (r.tail + i) % r.capacity ] )
  }
  r.tail += uint64(length)

  return f, true, uint64(length)
}





