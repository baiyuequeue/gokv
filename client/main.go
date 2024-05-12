package main

import (
	"bytes"
	"encoding/gob"
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

type MessageResponse struct {
	Status string
	Error  string // cant send struct in struct
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

		m := Message{"GET", "test3", "test"}
		c := &Client{&ls, m}

		err := c.sendRequest()
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func (c *Client) sendRequest() (err error) {

	buf := new(bytes.Buffer)
	gobObj := gob.NewEncoder(buf)
	err = gobObj.Encode(&c.mes)
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
	mes := new(MessageResponse)
	gobObj := gob.NewDecoder(buf)
	err = gobObj.Decode(mes)

	log.Println("STATUS: ", mes.Status, " ERROR: ", mes.Error, " DATA: ", mes.Data)
}
