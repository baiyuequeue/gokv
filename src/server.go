package main

import (
	"bytes"
	gob2 "encoding/gob"
	"errors"
	"log"
	"net"
	"reflect"
	"strings"
)

type Request struct {
	instruction string
	params      Message
}

type Server struct {
	storage Storage
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

func (s *Server) Serve() {
	ls, _ := net.Listen("tcp", "127.0.0.1:8888")
	defer func() {
		err := ls.Close()
		log.Fatalln(err)
	}()
	for {
		conn, err := ls.Accept()
		if err != nil {
			continue
		}
		go func() {
			defer func() {
				err := conn.Close()
				if err != nil {
					log.Print(err)
				}
			}()
			s.HandleConn(conn)
		}()
	}
}

func (s *Server) HandleRequest(request Request) (resp Message2) {
	switch request.instruction {
	case "SET":
		{
			resp := s.HandleSet(request.params)
			return resp
		}
	case "GET":
		{
			resp := s.HandleGet(request.params)
			return resp
		}
	case "DEL":
		{
			resp := s.HandleDel(request.params)
			return resp
		}
	}

	return Message2{"FAIL", errors.New("UNKNOWN_OPERATION"), ""}
}

func (s *Server) HandleConn(conn net.Conn) {

	log.Print("got new conn", conn)
	req, err := s.ParseConnection(conn)
	if err != nil {
		log.Print(err)
		conn.Write([]byte("BAD_REQUEST_FORMAT"))
		return
	}

	resp := s.HandleRequest(req)

	buf := new(bytes.Buffer)
	gob := gob2.NewEncoder(buf)
	gob.Encode(resp)
	conn.Write(buf.Bytes())
}

func (s *Server) ParseConnection(conn net.Conn) (Request, error) {

	tmp := make([]byte, 1024)

	_, err := conn.Read(tmp)
	if err != nil {
		return Request{}, err
	}

	log.Print("readed bytes")
	log.Print("raw conn data: ", tmp)

	// convert bytes into Buffer (which implements io.Reader/io.Writer)
	buf := bytes.NewBuffer(tmp)

	mes := new(Message)

	// creates a decoder object
	gob := gob2.NewDecoder(buf)

	// decodes buffer and unmarshals it into a Message struct
	err = gob.Decode(mes)
	if err != nil {
		return Request{}, err
	}

	action := strings.ToUpper(mes.Action)
	request := Request{instruction: action, params: *mes}
	log.Print("got request: ", request.instruction, request.params)

	return request, nil
}

func (s *Server) HandleSet(params Message) (resp Message2) {
	s.storage.Set(params.Key, params.Value)
	return Message2{"OK", nil, ""}
}

func (s *Server) HandleGet(params Message) (resp Message2) {
	value := s.storage.Get(params.Key)
	log.Println(value, reflect.TypeOf(value))
	if len(value) > 0 {
		log.Print("HandleGet 1")
		return Message2{"OK", nil, value}
	} else {
		log.Print("HandleGet 2")
		// TODO не выводит ошибку, Status="",а Error=<nil>
		return Message2{"FAIL", errors.New("empty value"), ""}
	}
}

func (s *Server) HandleDel(params Message) (resp Message2) {
	s.storage.Del(params.Key)
	return Message2{"OK", nil, ""}
}
