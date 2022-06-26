package main

import (
	"bufio"
	"errors"
	"log"

	//"io/ioutil"
	//"log"
	"net"
	"strings"
)

type Request struct {
	instruction string
	params      []string
}

type Server struct {
	storage Storage
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
		go s.HandleConn(conn)
	}
}

func (s *Server) HandleRequest(request Request) (answer string, err error) {
	switch request.instruction {
	case "SET":
		{
			answer, err := s.HandleSet(request.params)
			return answer, err
		}
	case "GET":
		{
			answer, err := s.HandleGet(request.params)
			return answer, err
		}
	case "DEL":
		{
			answer, err := s.HandleDel(request.params)
			return answer, err
		}
	}
	return "", errors.New("UNKNOWN_OPERATION")
}

func (s *Server) HandleConn(conn net.Conn) {
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Print(err)
		}
	}()
	log.Print("got new conn", conn)
	req, err := s.ParseConnection(conn)
	if err != nil {
		log.Print(err)
		conn.Write([]byte("BAD_REQUEST_FORMAT"))
		return
	}
	answer, err := s.HandleRequest(req)
	if err != nil {
		log.Print(err)
		conn.Write([]byte(err.Error()))
	} else {
		log.Print(answer)
		conn.Write([]byte(answer))
	}
}

func (s *Server) HandleCommand() {

}

func (s *Server) ParseConnection(conn net.Conn) (Request, error) {
	data, err := bufio.NewReader(conn).ReadString(';')
	if err != nil {
		return Request{}, err
	}
	log.Print("readed bytes")
	log.Print("raw conn data: ", data)
	input := strings.Split(data[:len(data)-1], " ")
	if len(input) < 1 {
		return Request{}, errors.New("cant parse request")
	}
	params := input[1:]
	instruction := strings.ToUpper(input[0])
	request := Request{instruction: instruction, params: params}
	log.Print("got request: ", request.instruction, request.params)
	return request, nil
}

func (s *Server) HandleSet(params []string) (status string, err error) {
	key, value := params[0], params[1]
	s.storage.Set(key, value)
	return "OK", nil
}

func (s *Server) HandleGet(params []string) (status string, err error) {
	key := params[0]
	value := s.storage.Get(key)
	return value, nil
}

func (s *Server) HandleDel(params []string) (status string, err error) {
	s.storage.Del(params[0])
	return "OK", nil
}
