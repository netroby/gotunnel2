package conn_reader

import (
  "net"
  ic "../infinite_chan"
  "math/rand"
  "time"
  "io"
)

const BUFFER_SIZE = 5600

func init() {
  rand.Seed(time.Now().UnixNano())
}

type ConnReader struct {
  Messages *ic.InfiniteChan
}

type Info struct {
  TCPConn *net.TCPConn
  Id int64
}

const (
  DATA = iota
  EOF
  ERROR
)

type Message struct {
  Info *Info
  Type int
  Data []byte
}

func New() *ConnReader {
  self := new(ConnReader)
  self.Messages = ic.New()
  return self
}

func (self *ConnReader) Close() {
  self.Messages.Close()
}

func (self *ConnReader) Add(tcpConn *net.TCPConn, id int64) *Info {
  if id < int64(0) {
    id = rand.Int63()
  }
  info := &Info{
    TCPConn: tcpConn,
    Id: id,
  }
  go func() {
    for {
      buf := make([]byte, 2048)
      n, err := tcpConn.Read(buf)
      if n > 0 {
        self.Messages.In <- Message{info, DATA, buf[:n]}
      }
      if err != nil {
        if err == io.EOF { //EOF
          self.Messages.In <- Message{info, EOF, nil}
        } else { //ERROR
          self.Messages.In <- Message{info, ERROR, []byte(err.Error())}
        }
        return
      }
    }
  }()
  return info
}
