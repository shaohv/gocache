package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

// Cmd ...
type Cmd struct {
	Op    string
	Key   string
	Value string
	Err   error
}

// Client ...
type Client struct {
	conn net.Conn
}

// New ..
func New(prot string, host string) *Client {
	c, err := net.Dial(prot, host)
	if err != nil {
		log.Println("dail fail", prot, host)
		return nil
	}
	return &Client{
		conn: c,
	}
}

// Cmd ..
func (c *Client) Cmd(op string, key string, value string, err error) *Cmd {
	return &Cmd{
		Op:    op,
		Key:   key,
		Value: value,
		Err:   err,
	}
}

func formatCmd(cmd *Cmd) []byte {

	var ops string

	if cmd.Op == "get" {
		ops = fmt.Sprintf("%d ", len(cmd.Key))
	} else if cmd.Op == "set" {
		ops = fmt.Sprintf("%d %d ", len(cmd.Key), len(cmd.Value))
	} else if cmd.Op == "del" {
		ops = fmt.Sprintf("%d ", len(cmd.Key))
	} else {
		panic("op not right")
	}

	var msg = append([]byte(ops), []byte(cmd.Key)...)
	if cmd.Op == "set" {
		msg = append(msg, []byte(cmd.Value)...)
	}

	return msg
}

// Run ..S
func (c *Client) Run(cmd *Cmd) error {
	var op string
	if cmd.Op == "get" {
		op = "G"
	} else if cmd.Op == "set" {
		op = "S"
	} else if cmd.Op == "del" {
		op = "D"
	} else {
		panic("op not right")
	}

	_, err := c.conn.Write([]byte(op))
	if err != nil {
		fmt.Println("client send cmd 2 server fail ", err)
		return err
	}

	bCmd := formatCmd(cmd)

	_, err = c.conn.Write(bCmd)
	if err != nil {
		fmt.Println("client send cmd 2 server fail ", err)
		return err
	}

	r := bufio.NewReader(c.conn)

	vlen, err := r.ReadString(' ')
	if err != nil {
		fmt.Println("client rcv rsp fail", err)
		return err
	}

	vLen, err := strconv.Atoi(strings.TrimSpace(vlen))
	if err != nil {
		log.Println("strconv fail", err)
		return err
	}

	var p = make([]byte, vLen)
	_, err = r.Read(p)
	if err != nil {
		fmt.Println("read value failed", err)
		return err
	}

	log.Printf("result:%s", p)
	return nil
}

func main() {
	fmt.Println("Client start...")
	server := flag.String("h", "localhost:12346", "cache server address")
	op := flag.String("c", "get", "cmd: can be get/set/del")
	key := flag.String("k", "", "key")
	value := flag.String("v", "", "value")
	flag.Parse()

	client := New("tcp", *server)
	if client == nil {
		log.Println("new client fail")
		return
	}
	cmd := client.Cmd(*op, *key, *value, nil)
	client.Run(cmd)
}
