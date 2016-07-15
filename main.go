package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"runtime"
	"github.com/gesellix/go-npipe/npipe"
)

func printUsage() {
	//var testPipeName = `\\.\pipe\winiotestpipe`

	log.Printf("Usage: %s url", os.Args[0])
	log.Printf("   ie: %s npipe:////./pipe/the_pipe", os.Args[0])
	log.Println()
	log.Printf("%s version: %s (%s on %s/%s; %s)", os.Args[0], "0.5", runtime.Version(), runtime.GOOS, runtime.GOARCH, runtime.Compiler)
	log.Println()
}

func main() {
	if len(os.Args) < 2 || os.Args[1] == "" {
		printUsage()
		os.Exit(1)
	}

	npipeURL, err := url.Parse(os.Args[1])
	if err != nil {
		log.Fatalf("could not parse url %q: %v", os.Args[1], err)
	}

	listener, err := npipe.Listen(npipeURL.Path)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	if listener == nil {
		log.Fatalf("listener is nil: %q", npipeURL.Path)
	}
	defer listener.Close()

	con := clientConns(listener)
	for {
		go handleConn(<-con)
	}
}

func clientConns(listener net.Listener) chan net.Conn {
	ch := make(chan net.Conn)
	i := 0
	go func() {
		for {
			client, err := listener.Accept()
			if client == nil {
				fmt.Printf("couldn't accept: %v", err)
				continue
			}
			i++
			fmt.Printf("%d: %v <-> %v\n", i, client.LocalAddr(), client.RemoteAddr())
			ch <- client
		}
	}()
	return ch
}

func handleConn(client net.Conn) {
	b := bufio.NewReader(client)
	for {
		line, err := b.ReadBytes('\n')
		if err != nil {
			// EOF, or worse
			break
		}
		client.Write(line)
	}
}
