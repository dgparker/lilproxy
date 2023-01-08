package udpproxy

import (
	"log"
	"net"
	"testing"
	"time"
)

func TestSendAndReceive(t *testing.T) {
	go runLilProxy()
	go runUDPServer()

	paddr, err := net.ResolveUDPAddr("udp", "localhost:9000")
	if err != nil {
		t.Fatal(err)
	}

	conn, err := net.DialUDP("udp", nil, paddr)
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		for {
			buf := make([]byte, 2048)
			_, _, err = conn.ReadFromUDP(buf)
			if err != nil {
				log.Fatal(err)
			}

			log.Printf("response received: %s", string(buf))
		}
	}()

	for {
		time.Sleep(1 * time.Second)
		_, err = conn.Write([]byte("hi\n"))
		if err != nil {
			log.Fatal(err)
		}
	}
}

func runLilProxy() {
	port := ":9000"
	target := "localhost:9001"

	c, err := New(port, target)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(c.ListenAndServe())
}

func runUDPServer() {
	taddr, err := net.ResolveUDPAddr("udp", ":9001")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.ListenUDP("udp", taddr)
	if err != nil {
		log.Fatal(err)
	}

	for {
		buf := make([]byte, 2048)
		_, caddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("request received: %s", string(buf))

		_, err = conn.WriteToUDP([]byte("bye\n"), caddr)
		if err != nil {
			log.Fatal(err)
		}

	}
}
