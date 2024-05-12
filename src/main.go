package main

func main() {
	s := Server{storage: Storage{data: map[string]string{}}}
	s.Serve()
	//s.Serve2()
}
