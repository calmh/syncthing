// Copyright (C) 2014 The Protocol Authors.

package protocol

import (
	"io"
	"sync/atomic"
	"time"
)

var (
	totalIncoming atomic.Int64
	totalOutgoing atomic.Int64
)

func TotalInOut() (int64, int64) {
	return totalIncoming.Load(), totalOutgoing.Load()
}

type counter struct {
	tot  atomic.Int64 // bytes
	last atomic.Int64 // unix nanos
}

func (c *counter) Add(n int) {
	c.tot.Add(int64(n))
	c.last.Store(time.Now().UnixNano())
}

func (c *counter) Last() time.Time {
	return time.Unix(0, c.last.Load())
}

func (c *counter) Tot() int64 {
	return c.tot.Load()
}

type countingUnderlying struct {
	Underlying
	read  *counter
	write *counter
}

func (c *countingUnderlying) Read(bs []byte) (int, error) {
	n, err := c.Underlying.Read(bs)
	c.read.Add(n)
	totalIncoming.Add(int64(n))
	return n, err
}

func (c *countingUnderlying) Write(bs []byte) (int, error) {
	n, err := c.Underlying.Write(bs)
	c.write.Add(n)
	totalOutgoing.Add(int64(n))
	return n, err
}

func (c *countingUnderlying) CreateSecondaryStream() (io.ReadWriteCloser, error) {
	s, err := c.Underlying.CreateSecondaryStream()
	if err != nil {
		return nil, err
	}
	return &countingReandWriteCloser{
		ReadWriter: s,
		Closer:     s,
		read:       c.read,
		write:      c.write,
	}, nil
}

func (c *countingUnderlying) AcceptSecondaryStream() (io.ReadWriter, error) {
	s, err := c.Underlying.AcceptSecondaryStream()
	if err != nil {
		return nil, err
	}
	return &countingReandWriteCloser{
		ReadWriter: s,
		Closer:     io.NopCloser(s),
		read:       c.read,
		write:      c.write,
	}, nil
}

type countingReandWriteCloser struct {
	io.ReadWriter
	io.Closer
	read  *counter
	write *counter
}

func (c *countingReandWriteCloser) Read(bs []byte) (int, error) {
	n, err := c.ReadWriter.Read(bs)
	c.read.Add(n)
	totalIncoming.Add(int64(n))
	return n, err
}

func (c *countingReandWriteCloser) Write(bs []byte) (int, error) {
	n, err := c.ReadWriter.Write(bs)
	c.write.Add(n)
	totalOutgoing.Add(int64(n))
	return n, err
}
