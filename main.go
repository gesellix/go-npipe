package main

import (
	"fmt"
	"github.com/Microsoft/go-winio"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
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
	}
	if listener == nil {
		log.Fatalf("listener is nil: %q", path)
	}
	defer listener.Close()

	err = serve(listener)
	if err != nil {
		log.Fatalf("Serve: %v", err)
	}
}

func serve(l net.Listener) error {
	http.HandleFunc("/", HandleDefault)
	http.HandleFunc("/containers/", HandleHijacked)
	http.HandleFunc("/hijack/", HandleHijacked)
	http.HandleFunc("/exit", HandleExit)
	return http.Serve(l, nil)
}

func HandleDefault(w http.ResponseWriter, req *http.Request) {
	log.Printf("<default handler> for %s %s\n", req.Method, req.RequestURI)
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

// everything below this is taken/copied/inspired from:
// - https://github.com/moby/moby/blob/07d60bc2571ba3d680f21adc84d87803ab4959c6/daemon/attach.go
// - https://github.com/containerd/containerd/blob/1f5b84f27cd675780bc7127f9aedbfe34cc7590b/pkg/cri/io/container_io.go#L137

var chErr chan error

// ParseForm ensures the request form is parsed even with invalid content types.
// If we don't do this, POST method without Content-type (even with empty body) will fail.
func ParseForm(r *http.Request) error {
	if r == nil {
		return nil
	}
	if err := r.ParseForm(); err != nil && !strings.HasPrefix(err.Error(), "mime:") {
		return err
	}
	return nil
}

// HijackConnection interrupts the http response writer to get the
// underlying connection and operate with it.
func HijackConnection(w http.ResponseWriter) (io.ReadCloser, io.Writer, error) {
	conn, _, err := w.(http.Hijacker).Hijack()
	if err != nil {
		return nil, nil, err
	}
	// Flush the options to make sure the client sets the raw mode
	_, _ = conn.Write([]byte{})
	return conn, conn, nil
}

func HandleHijacked(w http.ResponseWriter, req *http.Request) {
	log.Printf("<hijacking handler> for %s %s\n", req.Method, req.RequestURI)

	chErr = make(chan error, 1)
	defer close(chErr)
	if err := ParseForm(req); err != nil {
		chErr <- errors.Wrap(err, "error parsing form")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	r, rw, err := HijackConnection(w)
	if err != nil {
		chErr <- errors.Wrap(err, "error hijacking connection")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Close()

	fmt.Fprint(rw, "HTTP/1.1 101 UPGRADED\r\nContent-Type: application/vnd.docker.raw-stream\r\nConnection: Upgrade\r\nUpgrade: tcp\r\n\n")

	reader, writer := io.Pipe()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		if _, err := io.Copy(writer, r); err != nil {
			log.Printf("Failed to pipe stdin for attach: %v\n", err)
		}
		log.Printf("Attach stream stdin closed\n")
		writer.Close()
		r.Close()
		//rw.Close()
		wg.Done()
	}()

	attachStream := func(key string, close <-chan struct{}) {
		io.Copy(rw, reader)
		<-close
		log.Printf("Attach stream %q closed", key)
		// Make sure stdin gets closed.
		if r != nil {
			r.Close()
		}
		wg.Done()
	}

	wg.Add(1)
	close := make(chan struct{})
	go attachStream("stdout", close)

	wg.Wait()
}

// writeCloseInformer wraps a reader with a close function.
type wrapReadCloser struct {
	reader *io.PipeReader
	writer *io.PipeWriter
}

// NewWrapReadCloser creates a wrapReadCloser from a reader.
// NOTE(random-liu): To avoid goroutine leakage, the reader passed in
// must be eventually closed by the caller.
func NewWrapReadCloser(r io.Reader) io.ReadCloser {
	pr, pw := io.Pipe()
	go func() {
		_, _ = io.Copy(pw, r)
		pr.Close()
		pw.Close()
	}()
	return &wrapReadCloser{
		reader: pr,
		writer: pw,
	}
}

// Read reads up to len(p) bytes into p.
func (w *wrapReadCloser) Read(p []byte) (int, error) {
	n, err := w.reader.Read(p)
	if err == io.ErrClosedPipe {
		return n, io.EOF
	}
	return n, err
}

// Close closes read closer.
func (w *wrapReadCloser) Close() error {
	w.reader.Close()
	w.writer.Close()
	return nil
}
