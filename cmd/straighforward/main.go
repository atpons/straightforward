package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net"
	"time"

	"github.com/atpons/straightforward/ssh_usocksd"

	"github.com/cybozu-go/usocksd"
	"github.com/cybozu-go/well"
	"golang.org/x/crypto/ssh"
	"golang.org/x/net/proxy"
)

func privateKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return ssh.PublicKeys(key)
}

var (
	sshUser   = flag.String("u", "root", "ssh user name")
	sshKey    = flag.String("i", "/root/.ssh/id_rsa", "ssh key file")
	sshServer = flag.String("h", "localhost:22", "ssh server")
	port      = flag.Int("p", 9090, "listen port")
	proxyHost = flag.String("ocproxy", "", "ocproxy host")
)

func main() {
	flag.Parse()

	sshConfig := &ssh.ClientConfig{
		User: *sshUser,
		Auth: []ssh.AuthMethod{
			privateKeyFile(*sshKey),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}
	if *proxyHost == "" {
		log.Println("no proxy mode")
		serverConn, err := ssh.Dial("tcp", *sshServer, sshConfig)
		if err != nil {
			log.Fatal(err)
		}
		startStraightforward(serverConn)
	} else {
		log.Println("run ocproxy mode")
		p, err := proxy.SOCKS5("tcp", *proxyHost, nil, proxy.Direct)
		if err != nil {
			log.Fatal(err)
		}
		pc, err := p.Dial("tcp", *sshServer)
		conn, ch, req, err := ssh.NewClientConn(pc, *sshServer, sshConfig)
		if err != nil {
			log.Fatal(err)
		}
		c := ssh.NewClient(conn, ch, req)
		startStraightforward(c)
	}
}

func startStraightforward(serverConn *ssh.Client) {
	c := &usocksd.Config{
		Incoming: usocksd.IncomingConfig{
			Port:      *port,
			Addresses: []net.IP{net.ParseIP("127.0.0.1")},
		},
	}

	g := &well.Graceful{
		Listen: func() ([]net.Listener, error) {
			return usocksd.Listeners(c)
		},
		Serve: func(lns []net.Listener) {
			ssh_usocksd.Serve(lns, c, serverConn)
		},
		ExitTimeout: 1,
	}
	g.Run()

	err := well.Wait()
	if err != nil && !well.IsSignaled(err) {
		_ = serverConn.Close()
		log.Fatal(err)
	}
}
