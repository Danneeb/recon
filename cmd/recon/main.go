package main

import (
	"fmt"
	"log"
	"recon/internal/scanner"
	"recon/internal/server"
)

func main() {
	s := scanner.NewScanner([]string{"node_modules"})
	repos, err := s.Scan("../")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("found %d repos\n", len(repos))
	fmt.Println("serving at \033]8;;http://localhost:8080\033\\http://localhost:8080\033]8;;\033\\")

	srv := server.New(repos)
	log.Fatal(srv.Start(":8080"))
}
