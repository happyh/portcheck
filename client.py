#!/bin/python
#########################################################################
# File Name: client.py
# Author: happyhe
# mail: heguang@qiyi.com
# Created Time: Fri 31 Jul 2020 10:40:51 AM CST
#########################################################################
# coding=utf8

import socket
import getopt
import sys

def parse_ports(ports):
    list_ports = []
    numbs = ports.split(",")
    for i in numbs:
        numbs2 = i.split("-")
        if len(numbs2) == 1:
            list_ports.append(int(numbs2[0]))
        elif len(numbs2) == 2:
            for j in range(int(numbs2[0]), int(numbs2[1])):
                list_ports.append(j)

    return list_ports

def print_usage():
    print("usage:")
    print("client -t 100-200,1000-2000,10000 -u 100,200,300-1000 -s 123.125.118.77 [--timeout 5]")
    sys.exit()


if __name__ == '__main__':
    opts, args = getopt.getopt(sys.argv[1:], '-h-t:-u:-s:-m:', ['help', 'tcp=', 'udp=','server=',"timeout="])
    print(opts)
    tcp_ports = []
    udp_ports = []
    server = ""
    OK_tcp_ports = []
    OK_udp_ports = []
    fail_tcp_ports = []
    fail_udp_ports = []
    time_out = 5
    for opt_name, opt_value in opts:
        if opt_name in ('-h', '--help'):
            print_usage()
        if opt_name in ('-t', '--tcp'):
            ports = opt_value
            print("[*] tcp ports is ", ports)
            tcp_ports = parse_ports(ports)
        if opt_name in ('-u', '--udp'):
            ports = opt_value
            print("[*] udp ports is ", ports)
            udp_ports = parse_ports(ports)
        if opt_name in ('-s', '--server'):
            server = opt_value
            print("[*] server is ", server)
        if opt_name in ('-m', '--timeout'):
            time_out = int(opt_value)
            print("[*] timeout is ", time_out)


    if (len(tcp_ports) == 0 and len(udp_ports) == 0) or server == "":
        print_usage()

    for i in tcp_ports:
        client = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        client.settimeout(time_out)
        try:
            client.connect((server, i))
            client.send(str(i).encode('utf-8'))
            data = client.recv(1024)
            #print('recv:', data.decode())
            OK_tcp_ports.append(i)
        except Exception, err:
            print("error:" + str(i), err)
            fail_tcp_ports.append(i)
        client.close()

    for i in udp_ports:
        client = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
        client.settimeout(time_out)
        try:
            ip_port = (server, i)
            client.sendto(str(i).encode('utf-8'),ip_port)
            data,server_addr = client.recvfrom(1024)
            #print('recv:', data)
            OK_udp_ports.append(i)
        except Exception, err:
            print("error:" + str(i), err)
            fail_udp_ports.append(i)
        client.close()

    print("result:")
    print("tcp_OK=", OK_tcp_ports)
    print("udp_OK=", OK_udp_ports)
    print("tcp_fail=", fail_tcp_ports)
    print("udp_fail=", fail_udp_ports)