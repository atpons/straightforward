package ssh_usocksd

import (
	"log"
	"net"
	"strconv"

	"github.com/cybozu-go/usocksd"
	"github.com/cybozu-go/usocksd/socks"
	"github.com/cybozu-go/well"
	"github.com/nbio/contextual"
	"golang.org/x/crypto/ssh"
)

type ForwardDialer struct {
	*ssh.Client
}

func (f ForwardDialer) Dial(r *socks.Request) (net.Conn, error) {
	var addr string
	if len(r.Hostname) > 0 {
		addr = net.JoinHostPort(r.Hostname, strconv.Itoa(r.Port))
	} else {
		addr = net.JoinHostPort(r.IP.String(), strconv.Itoa(r.Port))
	}

	return contextual.Dialer{SimpleDialer: f.Client}.DialContext(r.Context(), "tcp", addr)
}

func Serve(lns []net.Listener, c *usocksd.Config, dialer *ssh.Client) {
	s := NewServer(c, dialer)
	for _, ln := range lns {
		s.Serve(ln)
	}
	err := well.Wait()
	if err != nil && !well.IsSignaled(err) {
		_ = dialer.Close()
		log.Fatal(err.Error())
	}
}

func NewServer(c *usocksd.Config, dialer *ssh.Client) *socks.Server {
	s := usocksd.NewServer(c)
	s.Dialer = ForwardDialer{dialer}
	return s
}
