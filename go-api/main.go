package main

import (
	"k8s-hyperledger-fabric-2.2/go-api/server"
)

func main() {
	s := server.NewServer()

	if err := s.Init(3001); err != nil {
		panic(err)
	}

	s.Start()
}
