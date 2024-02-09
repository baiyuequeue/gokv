package main

import (
	"bytes"
	gob2 "encoding/gob"
	"errors"
	"fmt"
	"log"
	"net"
)

type Client struct {
	conn *net.Conn
	mes  Message
}

type Message struct {
	Action string
	Key    string
	Value  string
}

type Message2 struct {
	Status string
	Error  error
	Data   string
}

func main() {

	ls, err := net.Dial("tcp", "127.0.0.1:8888")
	defer func() {
		err := ls.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	if err != nil {
		fmt.Println(err.Error())
	} else {

		m := Message{"get", "sanya2", "hui sosi 2"}
		c := &Client{&ls, m}

		err := c.sendRequest()
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func (c *Client) sendRequest() (err error) {

	buf := new(bytes.Buffer)
	gob := gob2.NewEncoder(buf)
	err = gob.Encode(&c.mes)
	if err != nil {
		return err
	}

	_, err = (*c.conn).Write(buf.Bytes())

	if err != nil {
		return err
	} else {
		c.printResponse()
		return nil
	}

}

func (c *Client) printResponse() {

	tmp := make([]byte, 1024)
	_, err := (*c.conn).Read(tmp)

	if err != nil {
		var netErr net.Error
		if !errors.As(err, &netErr) || !netErr.Timeout() {
			log.Println("read error:", err, &netErr)
		}
	}

	buf := bytes.NewBuffer(tmp)
	mes := new(Message2)
	gob := gob2.NewDecoder(buf)
	gob.Decode(mes)

	log.Println("STATUS: ", mes.Status, " ERROR: ", mes.Error, " DATA: ", mes.Data)
}
