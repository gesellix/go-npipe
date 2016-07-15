// +build !windows

package npipe

import (
	"net"
)

func Listen(addr string) (net.Listener, error) {
	panic("npipe protocol only supported on Windows")
}
