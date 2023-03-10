package main

import (
	"flag"
	"fmt"
	"net"
)

func main() {
	portPtr := flag.String("port", "3309", "database port")
	hostPtr := flag.String("host", "localhost", "database host")
	typePtr := flag.String("type", "tcp", "type of connection, tcp or udp")
	localPortPtr := flag.String("local", "9999", "local port")
	flag.Parse()

	listener, err := net.Listen(*typePtr, "localhost:"+*localPortPtr)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Start listening to client connections, setting up remote connection...")
	fmt.Println("Current params: host => " + *hostPtr + " port=> " + *portPtr)

	for {
		remote, _ := net.Dial(*typePtr, *hostPtr+":"+*portPtr)
		client, _ := listener.Accept()

		fmt.Println("Accepted client, starting communication")

		go proxyClientToClient(client, remote)
		go proxyClientToClient(remote, client)

		defer client.Close()
		defer remote.Close()
	}
}

func proxyClientToClient(client1 net.Conn, client2 net.Conn) {
	for {
		buffer := make([]byte, 1024)
		l, err := client1.Read(buffer)
		if err != nil && l != 0 {
			fmt.Println("Error reading:", err.Error())
		}
		client2.Write(buffer[:l])
	}
}
