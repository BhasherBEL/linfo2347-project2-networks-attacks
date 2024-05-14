#!/usr/bin/env bash

IP_WS2=10.1.0.2
IP_WS3=10.1.0.3
IP_HTTP=10.12.0.10
IP_DNS=10.12.0.20
IP_NTP=10.12.0.30
IP_FTP=10.12.0.40
IP_WWW=10.2.0.2

HOSTS=( $IP_WS2 $IP_WS3 $IP_HTTP $IP_DNS $IP_NTP $IP_FTP $IP_WWW )

function doPing() {
  return $(ping -c 1 -W 1 $1 &> /dev/null 2>&1)
}

function doCurl() {
  return $(curl -fsSL -o /dev/null --connect-timeout 1 $1 &> /dev/null)
}

function pass() {
  echo -ne "\e[32m$1\e[0m"
}

function fail() {
  echo -ne "\e[31m$1\e[0m"
}

function printCheck() {
  if [[ $? -eq 0 ]]; then
    pass "$1"
  else
    fail "$1"
  fi
}

for HOST in ${HOSTS[@]}; do
  printf "$HOST\t: "
  
  doPing $HOST
  printCheck "P"
  doCurl $HOST
  printCheck "C"
  echo
done
