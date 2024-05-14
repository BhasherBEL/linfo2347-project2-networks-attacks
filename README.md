# LINFO2347 - Network attacks and defense

## Deployment requirements

In addition to the mininet virtual machine, you will need to install the following tools on your host machine:

- [Go v1.21.9](https://golang.org/dl/)
- [Paramiko v3.4.0](https://pypi.org/project/paramiko/)
- [Scapy v2.5.0](https://pypi.org/project/scapy/)
- [NFTables v1.0.6](https://nftables.org/)

## Attacks

### Network scanning

We have implemented a python script able to perform a network scan in three different ways:

- TCP SYN scan
- TCP FIN scan
- UDP scan

The script can be executed from any machine in the network, and it will scan the network to check the response of the requested ip addresses and ports.

### SSH brute force

We have implemented a python script able to perform a brute force attack on an SSH server. The script will try to connect to the server using a list of passwords until it finds the correct one.

### FTP brute force

We have implemented a Go program able to perform a brute force attack on an FTP server. The program will try to connect to the server using a list of passwords until it finds the correct one.

## Defense

### Network scanning

The basic firewall configuration is already performant against network scanning.

### SSH brute force

<!-- TODO -->
We have implemented a snort rule to detect SSH brute force attacks. The rule will detect multiple failed login attempts from the same IP address and block the connection.

### FTP brute force

<!-- TODO -->
We have implemented a snort rule to detect FTP brute force attacks. The rule will detect multiple failed login attempts from the same IP address and block the connection.
