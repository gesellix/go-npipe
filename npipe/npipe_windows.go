package npipe

import (
	"log"
	"net"
	"github.com/Microsoft/go-winio"
)

func Listen(addr string) (net.Listener, error) {
	//ls := net.Listener{}
	// allow Administrators and SYSTEM, plus whatever additional users or groups were specified
	sddl := "D:P(A;;GA;;;BA)(A;;GA;;;SY)"
	//if socketGroup != "" {
	//	for _, g := range strings.Split(socketGroup, ",") {
	//		sid, err := winio.LookupSidByName(g)
	//		if err != nil {
	//			return nil, err
	//		}
	//		sddl += fmt.Sprintf("(A;;GRGW;;;%s)", sid)
	//	}
	//}
	c := winio.PipeConfig{
		SecurityDescriptor: sddl,
		MessageMode:        true, // Use message mode so that CloseWrite() is supported
		InputBufferSize:    65536, // Use 64KB buffers to improve performance
		OutputBufferSize:   65536,
	}
	l, err := winio.ListenPipe(addr, &c)
	if err != nil {
		return nil, err
	}
	//ls := append(ls, l)

	log.Printf("Listener created on npipe %s", addr)

	return l, nil
}

//func server(l net.Listener, ch chan int) {
//	c, err := l.Accept()
//	if err != nil {
//		panic(err)
//	}
//	rw := bufio.NewReadWriter(bufio.NewReader(c), bufio.NewWriter(c))
//	s, err := rw.ReadString('\n')
//	if err != nil {
//		panic(err)
//	}
//	_, err = rw.WriteString("got " + s)
//	if err != nil {
//		panic(err)
//	}
//	err = rw.Flush()
//	if err != nil {
//		panic(err)
//	}
//	c.Close()
//	ch <- 1
//}

//func TestFullListenDialReadWrite(t *testing.T) {
//	l, err := winio.ListenPipe(testPipeName, nil)
//	if err != nil {
//		t.Fatal(err)
//	}
//	defer l.Close()
//
//	ch := make(chan int)
//	go server(l, ch)
//
//	c, err := winio.DialPipe(testPipeName, nil)
//	if err != nil {
//		t.Fatal(err)
//	}
//	defer c.Close()
//
//	rw := bufio.NewReadWriter(bufio.NewReader(c), bufio.NewWriter(c))
//	_, err = rw.WriteString("hello world\n")
//	if err != nil {
//		t.Fatal(err)
//	}
//	err = rw.Flush()
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	s, err := rw.ReadString('\n')
//	if err != nil {
//		t.Fatal(err)
//	}
//	ms := "got hello world\n"
//	if s != ms {
//		t.Errorf("expected '%s', got '%s'", ms, s)
//	}
//
//	<-ch
//}
