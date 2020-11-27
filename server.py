#!/bin/python
#########################################################################
# File Name: server.py
# Author: happyhe
# mail: heguang@qiyi.com
# Created Time: Fri 31 Jul 2020 10:40:51 AM CST
# only for python2
#########################################################################
# coding=utf8
try:
    import selectors
except ImportError:
    import selectors2 as selectors
import socket
import getopt
import sys

epoll = selectors.DefaultSelector()


def create_conneciton(server):
    conn, addres = server.accept()

    epoll.register(conn, selectors.EVENT_READ, read_data)
    return conn


def read_data(conn):
    data = conn.recv(1024)
    if data:
        print("tcp:",data)
        conn.send(data)
    else:
        epoll.unregister(conn)


def read_udp_data(conn):
    data,client_addr = conn.recvfrom(1024)
    if data:
        print("udp:",data)
        conn.sendto(data,client_addr)
    else:
        epoll.unregister(conn)


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
    print("server --tcp 100-200,1000-2000,10000 --udp 100,200,300-1000")
    sys.exit()

if __name__ == '__main__':
    opts, args = getopt.getopt(sys.argv[1:], '-h-t:-u:', ['help', 'tcp=', 'udp='])
    tcp_ports = []
    udp_ports = []
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

    if len(tcp_ports) == 0 and len(udp_ports) == 0:
        print_usage()

    for i in tcp_ports:
        try:
            server = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            server.bind(('', i))
            server.listen(6)
            epoll.register(server, selectors.EVENT_READ, create_conneciton)
        except Exception, err:
            print("error:" + str(i), err)

    for i in udp_ports:
        try:
            server = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
            server.bind(('', i))
            epoll.register(server, selectors.EVENT_READ, read_udp_data)
        except Exception, err:
            print("error:" + str(i), err)

    while True:
        events = epoll.select()

        for key, mask in events:
            sock = key.fileobj
            callback = key.data

            callback(sock)
