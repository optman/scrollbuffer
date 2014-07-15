package scrollbuffer

/*
mapping infinite data into a finite buffer.

data must be write to buffer before it can read, data access should has locality pattern.

when writing data to buffer, we scoll the buffer(image it as a loop circle) to the write position, remap and invalidate old buffer data.

we use a RangeQueue struct to record the range of data loaded in buffer, it is a convenient way to keep track of discrete ranges.


                    0                                       n                                      2n
                    |---------------------------------------|---------------------------------------|-----   a long long data


|<-backward buffer->|<- forward buffer->|  initial maping
-x                  0                  n-x


|<- forward buffer->|<-backward buffer->|   actual buffer
0                  n-x                  n


*/

type ScrollBuffer struct {
	capacity    int
	map_range   Range
	data_ranges *RangeQueue
	pos         int
	data        []byte
}

func New(capacity, forward_buf_len int) *ScrollBuffer {

	//initial map: [-(capacity-forward_buf_len), forward_buf_len], pos at zero.

	return &ScrollBuffer{
		capacity:    capacity,
		pos:         0,
		map_range:   Range{Pos: -(capacity - forward_buf_len), Len: capacity},
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
	if pos >= this.pos {
		this.scoll(r.End())
	} else {
		this.scoll(pos)
	}

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

	if this.data_ranges.ContainRange(r) == false {
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

func (this *ScrollBuffer) DataRanges() *RangeQueue {
	return this.data_ranges
}

func (this *ScrollBuffer) scoll(new_pos int) {

	//map
	dist := new_pos - this.pos
	this.pos = new_pos
	this.map_range.Pos += dist

	//slice
	this.data_ranges = RangeQueueIntersect(RangeQueueFromRange(this.map_range), this.data_ranges)
}
