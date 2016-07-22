package main

import (
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"github.com/Microsoft/go-winio"
)

func printUsage() {
	oldFlags := log.Flags()
	log.SetFlags(0)
	log.Println()
	log.Printf("Usage:   %s <path>", os.Args[0])
	log.Printf("Example: %s \\\\.\\pipe\\the_pipe", os.Args[0])
	log.Println()
	log.Print("Any request on '/exit' stops the server and removes the pipe")
	log.Println()
	log.SetFlags(oldFlags)
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	path := os.Args[1]
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
	http.HandleFunc("/", HandleDefault)
	http.HandleFunc("/exit", HandleExit)
	return http.Serve(l, nil)
}

func HandleDefault(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(io.LimitReader(req.Body, 1048576))
	if err != nil {
		log.Panic(err)
	}
	defer req.Body.Close()
	if body == nil {
		log.Print("got empty body")
		io.WriteString(w, "[echo] OK")
	} else {
		log.Printf("got '%q'", string(body))
		response := []string{"[echo]", string(body)}
		io.WriteString(w, strings.Join(response, " "))
	}
}

func HandleExit(w http.ResponseWriter, req *http.Request) {
	_, err := ioutil.ReadAll(io.LimitReader(req.Body, 1048576))
	if err != nil {
		log.Panic(err)
	}
	defer req.Body.Close()
	log.Print("exiting")
	io.WriteString(w, "[echo] exit\n")
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
	defer os.Exit(0)
}
