package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"recon/internal/scanner"
	"recon/internal/server"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	if len(os.Args) < 2 {
		runServe(nil)
		return
	}

	switch os.Args[1] {
	case "scan":
		runScan(os.Args[2:])
	case "version":
		fmt.Printf("%s (commit: %s, built: %s)\n", version, commit, date)
	default:
		// treat first arg as a flag (e.g. recon --port 9090)
		runServe(os.Args[1:])
	}
}

func runServe(args []string) {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	root := fs.String("root", ".", "root directory to scan for git repos")
	port := fs.String("port", "8484", "port to listen on")
	fs.Parse(args)

	fmt.Printf("\n  recon %s\n", version)
	fmt.Println("  ────────────────────────────")
	fmt.Printf("  version  %s\n", version)
	fmt.Printf("  commit   %s\n", commit)
	fmt.Printf("  built    %s\n", date)
	fmt.Printf("  root     %s\n", *root)
	fmt.Printf("  port     %s\n", *port)
	fmt.Println()

	repos, err := scanner.NewScanner([]string{"node_modules"}).Scan(*root)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("  found %d repos\n\n", len(repos))
	fmt.Printf("  serving at \033]8;;http://localhost:%s\033\\http://localhost:%s\033]8;;\033\\\n\n", *port, *port)

	srv := server.New(repos)
	log.Fatal(srv.Start(":" + *port))
}

func runScan(args []string) {
	fs := flag.NewFlagSet("scan", flag.ExitOnError)
	root := fs.String("root", ".", "root directory to scan for git repos")
	fs.Parse(args)

	repos, err := scanner.NewScanner([]string{"node_modules"}).Scan(*root)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("found %d repos\n", len(repos))
	for _, r := range repos {
		fmt.Printf("  %s\t%s\t(%s)\n", r.Name, r.Path, r.Branch)
	}
}
