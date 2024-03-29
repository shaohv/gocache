package xtcp

import (
	"bufio"
	"fmt"
	"gocache/cache"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)

// Server ...
type Server struct {
	cache.Cacher
}

type result struct {
	v []byte
	e error
}

// NewServer ...
func NewServer(c cache.Cacher) *Server {
	return &Server{c}
}

// Listen ...
func (s *Server) Listen() {
	l, err := net.Listen("tcp", ":12346")
	if err != nil {
		panic(err)
	}

	for {
		c, err := l.Accept()
		if err != nil {
			panic(err)
		}

		// start a new goroute for handle this new connection
		go s.process(c)
	}
}

// readLen ...
func readLen(r *bufio.Reader) (int, error) {
	tmp, err := r.ReadString(' ')
	if err != nil {
		log.Println("read len failed")
		return 0, err
	}

	l, err := strconv.Atoi(strings.TrimSpace(tmp))
	if err != nil {
		log.Println("atoi failed", tmp)
		return 0, err
	}

	return l, err
}

// readKey ...
func (s *Server) readKey(r *bufio.Reader) (string, error) {
	klen, err := readLen(r)
	if err != nil {
		log.Println("readLen fail")
		return "", err
	}

	k := make([]byte, klen)
	_, err = io.ReadFull(r, k)
	if err != nil {
		log.Println("read key failed")
		return "", err
	}

	return string(k), nil
}

// readKeyAndValue ...
func (s *Server) readKeyAndValue(r *bufio.Reader) (string, []byte, error) {
	klen, err := readLen(r)
	if err != nil {
		log.Println("read klen failed")
		return "", nil, err
	}

	vlen, err := readLen(r)
	if err != nil {
		log.Println("read vlen failed")
		return "", nil, err
	}

	k := make([]byte, klen)
	v := make([]byte, vlen)

	_, err = io.ReadFull(r, k)
	if err != nil {
		log.Println("read key failed")
		return "", nil, err
	}

	_, err = io.ReadFull(r, v)
	if err != nil {
		log.Println("read value failed")
		return "", nil, err
	}

	return string(k), v, nil
}

func sendResponse(value []byte, err error, conn net.Conn) error {
	if err != nil {
		errStr := err.Error()
		var e error
		if cache.ErrKeyNotFound == errStr {
			tmp := fmt.Sprintf("%d ", 0)
			_, e = conn.Write([]byte(tmp))
		} else {
			log.Println("Response err", err)
			tmp := fmt.Sprintf("%d ", len(errStr)) + errStr
			lenTmp := fmt.Sprintf("%d ", len(tmp))
			_, e = conn.Write(append([]byte(lenTmp), []byte(tmp)...))
		}
		return e
	}

	vlen := fmt.Sprintf("%d ", len(value))
	_, e := conn.Write(append([]byte(vlen), value...))
	return e
}

func (s *Server) get(ch chan chan *result, r *bufio.Reader) {
	c := make(chan *result)
	ch <- c
	k, e := s.readKey(r)
	if e != nil {
		c <- &result{nil, e}
		return
	}

	go func() {
		v, e := s.Get(k)
		c <- &result{v, e}
	}()
}

func (s *Server) del(ch chan chan *result, r *bufio.Reader) {
	c := make(chan *result)
	ch <- c
	k, e := s.readKey(r)
	if e != nil {
		c <- &result{nil, e}
		return
	}

	go func() {
		c <- &result{nil, s.Del(k)}
	}()
}

func (s *Server) set(ch chan chan *result, r *bufio.Reader) {
	c := make(chan *result)
	ch <- c
	k, v, e := s.readKeyAndValue(r)
	if e != nil {
		c <- &result{nil, e}
		return
	}

	go func() {
		c <- &result{nil, s.Set(k, v)}
	}()
}

func (s *Server) process(conn net.Conn) {
	//defer conn.Close()
	r := bufio.NewReader(conn)

	resultChan := make(chan chan *result, 5000)
	defer close(resultChan)
	go reply(conn, resultChan)

	for {
		op, e := r.ReadByte()
		if e != nil {
			if e != io.EOF {
				log.Println("The connect is invalid")
			}
			return
		}

		if op == 'S' {
			s.set(resultChan, r)
		} else if op == 'G' {
			s.get(resultChan, r)
		} else if op == 'D' {
			s.del(resultChan, r)
		} else {
			log.Println("invalid operation!", op)
			return
		}

		if e != nil {
			log.Println("handle failed:", e)
			return
		}
	}
}

func reply(conn net.Conn, resultCh chan chan *result) {
	defer conn.Close()
	for {
		c, open := <-resultCh
		if !open {
			return
		}

		r := <-c
		e := sendResponse(r.v, r.e, conn)
		if e != nil {
			fmt.Println("send response fail")
			return
		}
	}
}
