package udpproxy

import (
	"log"
	"net"
	"sync"
	"time"
)

type Client struct {
	laddr *net.UDPAddr
	raddr *net.UDPAddr

	proxyConn *net.UDPConn

	mutex    sync.RWMutex
	sessions map[string]*Session
}

func New(port, target string) (*Client, error) {
	laddr, err := net.ResolveUDPAddr("udp", port)
	if err != nil {
		return nil, err
	}

	raddr, err := net.ResolveUDPAddr("udp", target)
	if err != nil {
		return nil, err
	}

	return &Client{
		laddr:    laddr,
		raddr:    raddr,
		mutex:    sync.RWMutex{},
		sessions: map[string]*Session{},
	}, nil
}

func (c *Client) ListenAndServe() error {
	var err error
	c.proxyConn, err = net.ListenUDP("udp", c.laddr)
	if err != nil {
		return err
	}

	go c.pruneSessions()

	for {
		buf := make([]byte, 2048)
		n, caddr, err := c.proxyConn.ReadFromUDP(buf)
		if err != nil {
			log.Println(err)
		}

		session, found := c.sessions[caddr.String()]
		if !found {
			session, err = createSession(caddr, c.raddr, c.proxyConn)
			if err != nil {
				log.Println(err)
				continue
			}

			c.sessions[caddr.String()] = session
		}

		session.proxyTo(buf[:n])
	}
}

func (c *Client) pruneSessions() {
	ticker := time.NewTicker(1 * time.Minute)

	// the locks here could be abusive and i dont even know if this is a real
	// problem but we definitely need to clean up stale sessions
	for range ticker.C {
		for _, session := range c.sessions {
			c.mutex.RLock()
			if time.Since(session.updateTime) > time.Minute*5 {
				delete(c.sessions, session.caddr.String())
			}
			c.mutex.RUnlock()
		}
	}
}
