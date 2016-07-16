package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"github.com/Microsoft/go-winio"
	"net"
)

var echoPipeName = `\\.\pipe\echo_pipe`

func printUsage() {

	log.Printf("Usage: %s url", os.Args[0])
	log.Printf("   ie: %s \\\\.\\pipe\\the_pipe", os.Args[0])
	log.Println()
	log.Printf("%s version: %s (%s on %s/%s; %s)", os.Args[0], "0.5", runtime.Version(), runtime.GOOS, runtime.GOARCH, runtime.Compiler)
	log.Println()
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	path := os.Args[1]
	if path == "" {
		printUsage()
		path = echoPipeName
	}

	log.Printf("using path: %q", path)

	listener, err := winio.ListenPipe(path, nil)
	if err != nil {
		log.Fatalf("error: %v", err)
		os.Exit(1)
	}
	if listener == nil {
		log.Fatalf("listener is nil: %q", path)
		os.Exit(1)
	}
	defer listener.Close()

	err = serve(listener)
	if err != nil {
		log.Fatalf("Serve: %v", err)
		os.Exit(1)
	}
}

func serve(l net.Listener) error {
	http.HandleFunc("/", EchoServer)
	http.HandleFunc("/exit", ExitServer)
	return http.Serve(l, nil)
}

func EchoServer(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(io.LimitReader(req.Body, 1048576))
	if err != nil {
		log.Panic(err)
	}
	defer req.Body.Close()
	if body == nil {
		io.WriteString(w, "[echo] OK")
	} else {
		io.WriteString(w, string(body))
	}
}

func ExitServer(w http.ResponseWriter, req *http.Request) {
	_, err := ioutil.ReadAll(io.LimitReader(req.Body, 1048576))
	if err != nil {
		log.Panic(err)
	}
	defer req.Body.Close()
	io.WriteString(w, "exit")
	os.Exit(0)
}
