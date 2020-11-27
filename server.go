package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

func parse_ports(ports string) []string {
	var list_ports []string
	numbs := strings.Split(ports, ",")

	for _, i := range numbs {
		numbs2 := strings.Split(i, "-")
		if len(numbs2) == 1 {
			list_ports = append(list_ports, numbs2[0])
		} else if len(numbs2) == 2 {
			var j, e int
			var err error
			if j, err = strconv.Atoi(numbs2[0]); err != nil {
				panic(err)
			}
			if e, err = strconv.Atoi(numbs2[1]); err != nil {
				panic(err)
			}
			for ; j < e; j++ {
				list_ports = append(list_ports, strconv.Itoa(j))
			}
		}
	}

	return list_ports
}

func tcp_handle(i string, conn net.Conn) {
	var buf []byte = make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("error tcp "+i+":", err)
	} else {
		fmt.Println("tcp:", string(buf[:n]))
		conn.Write(buf[:n])
	}
}

func tcp_server(i string) {
	listener, err := net.Listen("tcp", ":"+i)
	if err != nil {
		fmt.Println("error tcp "+i+":", err)
	} else {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("tcp error "+i+":", err)
				continue
			}
			go tcp_handle(i, conn)
		}
	}
}

func udp_server(i string) {
	addr, err := net.ResolveUDPAddr("udp", ":"+i)
	if err != nil {
		fmt.Println("error udp "+i+":", err)
	} else {
		conn, err := net.ListenUDP("udp", addr)
		if err != nil {
			fmt.Println("error udp "+i+":", err)
		} else {
			data := make([]byte, 1024)
			n, remoteAddr, err := conn.ReadFromUDP(data)
			if err != nil {
				fmt.Println("error udp "+i+":", err)
			} else {
				fmt.Println("udp:", string(data[:n]))
				conn.WriteToUDP(data[:n], remoteAddr)
			}
		}
	}
}

func main() {
	var tcp_port, udp_port string
	flag.StringVar(&tcp_port, "t", "", "tcp")
	flag.StringVar(&udp_port, "u", "", "udp")
	flag.Parse()

	if tcp_port == "" && udp_port == "" {
		flag.Usage()
		return
	}

	var tcp_ports, udp_ports []string

	if tcp_port != "" {
		tcp_ports = parse_ports(tcp_port)
	}
	if udp_port != "" {
		udp_ports = parse_ports(udp_port)
	}

	for _, i := range tcp_ports {
		go tcp_server(i)
	}

	for _, i := range udp_ports {
		go udp_server(i)
	}

	time.Sleep(time.Hour)
}
