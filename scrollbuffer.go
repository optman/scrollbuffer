package scrollbuffer

/*
mapping infinite data into a finite buffer.

data must be write to buffer before it can read, data access should has locality pattern.

when writing data to buffer, we scoll the buffer(image it as a loop circle) to the write position, remap and invalidate old buffer data.

at the begining, data position 0 is map to buffer position x.

x is use to balance backward and forward writing. if most write is forward, then small forward buffer will left more space for backward reading.

we use a RangeQueue struct to record the range of data loaded in buffer, it is a convenient way to keep track of discrete ranges.


                    0                                       n                                      2n
                    |---------------------------------------|---------------------------------------|-----   a long long data
                           	        

|<-backward buffer->|<- forward buffer->|

|---------------------------------------|   buffer
0                   x                   n


*/

type ScrollBuffer struct {
	capacity    int
	map_range   *Range
	data_ranges *RangeQueue
	pos         int
	data        []byte
}

func New(capacity, forward_buf_len int) *ScrollBuffer {

	//initial map: [-(capacity-forward_buf_len), forward_buf_len], pos at zero.

	if forward_buf_len == 0 {
		forward_buf_len = capacity / 2
	}

	return &ScrollBuffer{
		capacity:    capacity,
		pos:         0,
		map_range:   &Range{Pos: -(capacity - forward_buf_len), Len: capacity},
		data_ranges: &RangeQueue{},
		data:        make([]byte, capacity),
	}
}

func (this *ScrollBuffer) Write(pos int, data []byte) {
	if len(data) > this.capacity {
		panic("data length larger than capacity")
	}

	r := Range{Pos: pos, Len: len(data)}

	//append
	this.data_ranges.AddRange(r)

	//scroll
	this.scoll(r.End())

	//wite
	begin := pos % this.capacity
	if pos/this.capacity == r.End()/this.capacity {
		copy(this.data[begin:], data)
	} else {
		first_half_len := this.capacity - begin
		copy(this.data[begin:], data[0:first_half_len])
		copy(this.data, data[first_half_len:])
	}

}

func (this *ScrollBuffer) Read(pos int, data []byte) {

	r := Range{Pos: pos, Len: len(data)}
	rq := &RangeQueue{}
	rq.AddRange(r)

	if RangeQueueIntersect(rq, this.data_ranges).Equals(rq) != true {
		panic("read data not in ranges!")
	}

	//read
	begin := pos % this.capacity
	if pos/this.capacity == r.End()/this.capacity {
		copy(data, this.data[begin:])
	} else {
		first_half_len := this.capacity - begin
		copy(data[0:first_half_len], this.data[begin:])
		copy(data[first_half_len:], this.data)
	}

}

func (this *ScrollBuffer) scoll(new_pos int) {

	//map
	dist := new_pos - this.pos
	this.pos = new_pos
	this.map_range.Pos += dist

	//slice
	map_rqs := &RangeQueue{}
	map_rqs.AddRange(*this.map_range)

	this.data_ranges = RangeQueueIntersect(map_rqs, this.data_ranges)
}
