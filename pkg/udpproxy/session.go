package udpproxy

import (
	"log"
	"net"
	"time"
)

type Session struct {
	serverConn *net.UDPConn
	proxyConn  *net.UDPConn
	caddr      *net.UDPAddr
	updateTime time.Time
}

func createSession(caddr *net.UDPAddr, raddr *net.UDPAddr, proxyConn *net.UDPConn) (*Session, error) {
	serverConn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		return nil, err
	}

	session := &Session{
		serverConn: serverConn,
		proxyConn:  proxyConn,
		caddr:      caddr,
		updateTime: time.Now(),
	}

	go session.listen()

	return session, nil
}

func (s *Session) listen() error {
	for {
		buf := make([]byte, 2048)
		n, err := s.serverConn.Read(buf)
		if err != nil {
			log.Println(err)
			continue
		}

		go s.proxyFrom(buf[:n])
	}
}

func (s *Session) proxyFrom(buf []byte) error {
	s.updateTime = time.Now()
	_, err := s.proxyConn.WriteToUDP(buf, s.caddr)
	if err != nil {
		return err
	}

	return nil
}

func (s *Session) proxyTo(buf []byte) error {
	s.updateTime = time.Now()
	_, err := s.serverConn.Write(buf)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
