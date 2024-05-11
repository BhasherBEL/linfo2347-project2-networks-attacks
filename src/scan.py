#!/usr/bin/env python3
# -*- coding: utf-8 -*-

"""
A simple script to scan ports on one or more hosts using Scapy. This 
script supports scanning using TCP and UDP protocols, and can use either
SYN or FIN packets to scan for open ports.
"""

import argparse
from scapy.all import IP, UDP, TCP, sr1

def make_udp_packet(port: int) -> UDP:
	"""Create a UDP packet with the specified port
	
	Args:
		port (int): The port number to use
	Return:
		UDP: The UDP packet
	"""
	return UDP(dport=port)

def make_tcp_packet(port: int, mode: str) -> TCP:
	"""Create a TCP packet with the specified port and mode
	
	Args:
		port (int): The port number to use
		mode (str): The mode to use (syn or rst)
	Return:
		TCP: The TCP packet
	"""
	match mode:
		case "syn":
			return TCP(dport=port, flags="S")
		case "fin":
			return TCP(dport=port, flags="F")
		case _:
			raise ValueError(f"Invalid mode: expect 'syn' or 'rst' but got {mode}")

def make_ip_packet(host: str) -> IP:
	"""Create an IP packet with the specified host
	
	Args:
		host (str): The host to use
	Return:
		IP: The IP packet
	"""
	return IP(dst=host)

def scan_port(host: str, port: int, mode: str, transport: str) -> bool:
	"""Scan the specified port on the host using the specified mode and transport
	
	Args:
		host (str): The host to scan
		port (int): The port to scan
		mode (str): The mode to use (syn or rst)
		transport (str): The transport protocol to use (tcp or udp)
	Return:
		bool: True if the port is open, False otherwise
	"""
	packet = make_ip_packet(host)
	if transport == "tcp":
		packet /= make_tcp_packet(port, mode)
		response = sr1(packet, timeout=1, verbose=0)
		match [mode, response]:
			case [_, None]:
				return False
			case ["syn", response]:
				return response.haslayer("TCP") and response.getlayer("TCP").flags == 18
			case ["fin", None]:
				return True
			case ["fin", response]:
				return response.haslayer("TCP") and response.getlayer("TCP").flags == 17
	elif transport == "udp":
		packet /= make_udp_packet(port)
		response = sr1(packet, timeout=1, verbose=0)
		if response is not None:
			return response.haslayer("ICMP") and response.getlayer("ICMP").type == 3 and response.getlayer("ICMP").code == 3
	return False

def scan_host(host: str, port_range: tuple[int, int], mode: str, transport: str):
	"""Scan the host using the specified port range, mode, and transport
	
	Args:
		host (str): The host to scan
		port_range (tuple[int, int]): The port range to scan
		mode (str): The mode to use (syn or rst)
		transport (str): The transport protocol to use (tcp or udp)
	Post:
		prints port status
	"""
	print("+--- Scan results for host", host, "---+")
	for port in range(port_range[0], port_range[1] + 1):
		print(f"{port}", end=" ")
		if scan_port(host, port, mode, transport):
			print("open")
		else:
			print("closed")
	print("+---------------------------------+")

def scan_hosts(hosts: list[str], port_range: tuple[int, int], mode: str, transport: str):
	"""Scan the hosts using the specified port range, mode, and transport
	
	Args:
		hosts (list[str]): The hosts to scan
		port_range (tuple[int, int]): The port range to scan
		mode (str): The mode to use (syn or rst)
		transport (str): The transport protocol to use (tcp or udp)
	Post:
		prints port status
	"""
	for host in hosts:
		scan_host(host, port_range, mode, transport)

def _main():
	parser = argparse.ArgumentParser(description="Simple port scanner using Scapy")
	parser.add_argument('-p', '--port', type=int, nargs=2, default=[20, 80], help='Port range to scan, in -p <min> <max> format (e.g., -p 20 80)')
	parser.add_argument('-t', '--transport', choices=['tcp', 'udp'], default='tcp', help='Transport protocol to use (tcp or udp)')
	parser.add_argument('-m', '--mode', choices=['syn', 'fin'], default="syn", help='Mode to use for scanning (syn or rst)')
	parser.add_argument('hosts', nargs='+', help='One or more hosts to scan')

	args = parser.parse_args()

	scan_hosts(args.hosts, tuple(args.port), args.mode, args.transport)


if __name__ == "__main__":
	_main()
