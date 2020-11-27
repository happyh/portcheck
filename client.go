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

func main() {
	var tcp_port, udp_port string
	var server string
	var time_out int
	flag.StringVar(&server, "s", "", "server")
	flag.IntVar(&time_out, "m", 5, "timeout")
	flag.StringVar(&tcp_port, "t", "", "tcp")
	flag.StringVar(&udp_port, "u", "", "udp")
	flag.Parse()

	if server == "" || (tcp_port == "" && udp_port == "") {
		flag.Usage()
		return
	}

	var tcp_ports, udp_ports, OK_tcp_ports, OK_udp_ports, fail_tcp_ports, fail_udp_ports []string

	if tcp_port != "" {
		tcp_ports = parse_ports(tcp_port)
	}
	if udp_port != "" {
		udp_ports = parse_ports(udp_port)
	}
	for _, i := range tcp_ports {
		conn, err := net.DialTimeout("tcp", server+":"+i, time.Duration(time_out)*time.Second)
		if err != nil {
			fmt.Println("error tcp " + i + ":" + err.Error())
			fail_tcp_ports = append(fail_tcp_ports, i)
		} else {
			conn.SetDeadline(time.Now().Add(time.Duration(time_out) * time.Second))
			conn.Write([]byte(i))
			readBytes := make([]byte, 1024)
			if _, err = conn.Read(readBytes); err != nil {
				fmt.Println("error tcp " + i + ":" + err.Error())
				fail_tcp_ports = append(fail_tcp_ports, i)
			} else {
				fmt.Println("tcp ok:", i)
				OK_tcp_ports = append(OK_tcp_ports, i)
			}
			conn.Close()
		}
	}

	for _, i := range udp_ports {
		if addr, err := net.ResolveUDPAddr("udp", server+":"+i); err != nil {
			fmt.Println("error udp " + i + ":" + err.Error())
			fail_udp_ports = append(fail_udp_ports, i)
		} else {
			conn, err := net.DialUDP("udp", nil, addr)
			if err != nil {
				fmt.Println("error udp " + i + ":" + err.Error())
				fail_udp_ports = append(fail_udp_ports, i)
			} else {
				conn.SetDeadline(time.Now().Add(time.Duration(time_out) * time.Second))
				_, err = conn.Write([]byte(i))
				if err != nil {
					fmt.Println("error udp " + i + ":" + err.Error())
					fail_udp_ports = append(fail_udp_ports, i)
				} else {
					data := make([]byte, 1024)
					_, err = conn.Read(data)
					if err != nil {
						fail_udp_ports = append(fail_udp_ports, i)
						fmt.Println("error udp " + i + ":" + err.Error())
					} else {
						fmt.Println("udp ok:", i)
						OK_udp_ports = append(OK_udp_ports, i)
					}
				}
				conn.Close()
			}
		}
	}

	fmt.Println("result:")
	fmt.Println("tcp_OK=", OK_tcp_ports)
	fmt.Println("udp_OK=", OK_udp_ports)
	fmt.Println("tcp_fail=", fail_tcp_ports)
	fmt.Println("udp_fail=", fail_udp_ports)

}
