package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"anto.pt/x/togo"
)

func main() {
	laddr := flag.String("laddr", ":8080", "Address for listen to")
	flag.Parse()

	dirPath := flag.Arg(0)
	if dirPath == "" {
		fmt.Fprintf(os.Stderr, "error: missing directory to serve\n")
		fmt.Fprintf(os.Stderr, "usage: %s <directory>\n", os.Args[0])
		os.Exit(1)
	}

	dir := os.DirFS(dirPath)
	tree, err := togo.NewRouteTreeFromFS(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: initializing the route tree\n%v", err)
	}

	mux := tree.Mux()

	log.Println("listening on", *laddr)
	panic(http.ListenAndServe(*laddr, mux))
}
