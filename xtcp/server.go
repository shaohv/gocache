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
		log.Println("atoi failed")
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
		tmp := fmt.Sprintf("-%d", len(errStr)) + errStr
		_, e := conn.Write([]byte(tmp))
		return e
	}

	vlen := fmt.Sprintf("%d ", len(value))
	_, e := conn.Write(append([]byte(vlen), value...))
	return e
}

func (s *Server) get(conn net.Conn, r *bufio.Reader) error {
	k, e := s.readKey(r)
	if e != nil {
		return e
	}

	v, e := s.Get(k)
	return sendResponse(v, e, conn)
}

func (s *Server) del(conn net.Conn, r *bufio.Reader) error {
	k, e := s.readKey(r)
	if e != nil {
		return e
	}

	return sendResponse(nil, s.Del(k), conn)
}

func (s *Server) set(conn net.Conn, r *bufio.Reader) error {
	k, v, e := s.readKeyAndValue(r)
	if e != nil {
		return e
	}

	return sendResponse(nil, s.Set(k, v), conn)
}

func (s *Server) process(conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)

	for {
		op, e := r.ReadByte()
		if e != nil {
			if e != io.EOF {
				log.Println("The connect is invalid")
			}
			return
		}

		if op == 'S' {
			e = s.set(conn, r)
		} else if op == 'G' {
			e = s.get(conn, r)
		} else if op == 'D' {
			e = s.del(conn, r)
		} else {
			log.Println("invalid operation!")
			return
		}

		if e != nil {
			log.Println("handle failed:", e)
			return
		}
	}
}
