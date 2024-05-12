package main

import (
	"bytes"
	gob2 "encoding/gob"
	"log"
	"net"
	"reflect"
	"strings"
)

type Request struct {
	instruction string
	params      MessageRequest
}

type Server struct {
	storage Storage
}

type MessageRequest struct {
	Action string
	Key    string
	Value  string
}

type MessageResponse struct {
	Status string
	Error  string // cant send struct in struct as bytes by tcp
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

func (s *Server) HandleRequest(request Request) (resp MessageResponse) {
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

	return MessageResponse{"FAIL", "UNKNOWN_OPERATION", ""}
}

func (s *Server) HandleConn(conn net.Conn) {

	log.Print("got new conn", conn)
	req, err := s.ParseConnection(conn)
	if err != "" {
		log.Print(err)
		_, _ = conn.Write([]byte("BAD_REQUEST_FORMAT"))
		return
	}

	resp := s.HandleRequest(req)

	buf := new(bytes.Buffer)
	gob := gob2.NewEncoder(buf)
	_ = gob.Encode(resp)
	_, _ = conn.Write(buf.Bytes())
}

func (s *Server) ParseConnection(conn net.Conn) (Request, string) {

	tmp := make([]byte, 1024)

	_, err := conn.Read(tmp)
	if err != nil {
		return Request{}, ""
	}

	log.Print("readed bytes")
	log.Print("raw conn data: ", tmp)

	// convert bytes into Buffer (which implements io.Reader/io.Writer)
	buf := bytes.NewBuffer(tmp)

	mes := new(MessageRequest)

	// creates a decoder object
	gob := gob2.NewDecoder(buf)

	// decodes buffer and unmarshals it into a MessageRequest struct
	err = gob.Decode(mes)
	if err != nil {
		return Request{}, err.Error()
	}

	action := strings.ToUpper(mes.Action)
	request := Request{instruction: action, params: *mes}
	log.Print("got request: ", request.instruction, request.params)

	return request, ""
}

func (s *Server) HandleSet(params MessageRequest) (resp MessageResponse) {
	s.storage.Set(params.Key, params.Value)
	return MessageResponse{"OK", "", ""}
}

func (s *Server) HandleGet(params MessageRequest) (resp MessageResponse) {
	value := s.storage.Get(params.Key)
	log.Println(value, reflect.TypeOf(value))
	if len(value) > 0 {
		log.Print("HandleGet 1")
		return MessageResponse{"OK", "", value}
	} else {
		log.Print("HandleGet 2")
		return MessageResponse{"FAIL", "empty value", ""}
	}
}

func (s *Server) HandleDel(params MessageRequest) (resp MessageResponse) {
	s.storage.Del(params.Key)
	return MessageResponse{"OK", "", ""}
}
