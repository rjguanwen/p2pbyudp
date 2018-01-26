package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	tag := os.Args[1]
	sip := net.ParseIP("207.148.70.129")
	srcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 9982}
	dstAddr := &net.UDPAddr{IP: sip, Port: 9981}
	conn, err := net.DialUDP("udp", srcAddr, dstAddr)
	if err != nil {
		fmt.Println(err)
	}
	if _, err = conn.Write([]byte("hello, I'm new peer " + srcAddr.String())); err != nil {
		log.Panic(err)
	}
	data := make([]byte, 1024)
	n, remoteAddr, err := conn.ReadFromUDP(data)
	if err != nil {
		fmt.Printf("error during read: %s", err)
	}
	conn.Close()
	fmt.Printf("local:%s server:%s another:%s\n", srcAddr, remoteAddr, data[:n])
	anotherPeer := parseAddr(string(data[:n]))
	log.Printf("get another peer:%s", anotherPeer.String())

	conn, err = net.DialUDP("udp", srcAddr, &anotherPeer)
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()

	if _, err = conn.Write([]byte("handshake from:" + tag)); err != nil {
		log.Println("send handshake:", err)
	} else {
		log.Println("send handshake ok")
	}

	go func() {
		for {
			time.Sleep(10 * time.Second)
			if _, err = conn.Write([]byte("from [" + tag + "]")); err != nil {
				log.Println("send msg fail", err)
			}
		}
	}()
	for {
		data = make([]byte, 1024)
		n, remoteAddr, err = conn.ReadFromUDP(data)
		if err != nil {
			log.Printf("error during read: %s\n", err)
		}
		log.Printf("<%s> %s\n", remoteAddr, data[:n])
	}
}

func parseAddr(addr string) net.UDPAddr {
	t := strings.Split(addr, ":")
	port, _ := strconv.Atoi(t[1])
	return net.UDPAddr{
		IP:   net.ParseIP(t[0]),
		Port: port,
	}
}
