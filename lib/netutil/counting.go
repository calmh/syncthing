package netutil

import (
	"io"
	"sync/atomic"
	"time"
)

type Counted interface {
	BytesRead() int64
	LastRead() time.Time
	BytesWritten() int64
	LastWrite() time.Time
}

var rootCounter Counter

func RootCounter() Counted {
	return &rootCounter
}

type CountedStream interface {
	Counted
	Stream
}

type Counter struct {
	readBytes  atomic.Int64
	lastRead   atomic.Int64
	writeBytes atomic.Int64
	lastWrite  atomic.Int64
	parent     *Counter
}

func NewCounter() *Counter {
	return newCounterWithParent(&rootCounter)
}

func newCounterWithParent(parent *Counter) *Counter {
	return &Counter{
		parent: parent,
	}
}

func (c *Counter) BytesRead() int64 {
	return c.readBytes.Load()
}

func (c *Counter) BytesWritten() int64 {
	return c.readBytes.Load()
}

func (c *Counter) LastRead() time.Time {
	return time.Unix(0, c.lastRead.Load())
}

func (c *Counter) LastWrite() time.Time {
	return time.Unix(0, c.lastWrite.Load())
}

func (c *Counter) addRead(n int) {
	c.readBytes.Add(int64(n))
	c.lastRead.Store(time.Now().UnixNano())
	if c.parent != nil {
		c.parent.addRead(n)
	}
}

func (c *Counter) addWrite(n int) {
	c.writeBytes.Add(int64(n))
	c.lastWrite.Store(time.Now().UnixNano())
	if c.parent != nil {
		c.parent.addWrite(n)
	}
}

type CountingStream struct {
	Stream
	*Counter
}

func NewCountingStream(s Stream) *CountingStream {
	return &CountingStream{
		Stream:  s,
		Counter: NewCounter(),
	}
}

func (c *CountingStream) Read(bs []byte) (int, error) {
	n, err := c.Stream.Read(bs)
	c.Counter.addRead(n)
	return n, err
}

func (c *CountingStream) Write(bs []byte) (int, error) {
	n, err := c.Stream.Write(bs)
	c.Counter.addWrite(n)
	return n, err
}

func (c *CountingStream) CreateSubstream() (io.ReadWriteCloser, error) {
	s, err := c.Stream.CreateSubstream()
	if err != nil {
		return nil, err
	}
	return &countingReadWriteCloser{
		ReadWriter: s,
		Closer:     s,
		Counter:    c.Counter,
	}, nil
}

func (c *CountingStream) AcceptSubstream() (io.ReadWriter, error) {
	s, err := c.Stream.AcceptSubstream()
	if err != nil {
		return nil, err
	}
	return &countingReadWriteCloser{
		ReadWriter: s,
		Closer:     io.NopCloser(s),
		Counter:    c.Counter,
	}, nil
}

type countingReadWriteCloser struct {
	io.ReadWriter
	io.Closer
	*Counter
}

func (c *countingReadWriteCloser) Read(bs []byte) (int, error) {
	n, err := c.ReadWriter.Read(bs)
	c.Counter.addRead(n)
	return n, err
}

func (c *countingReadWriteCloser) Write(bs []byte) (int, error) {
	n, err := c.ReadWriter.Write(bs)
	c.Counter.addWrite(n)
	return n, err
}
