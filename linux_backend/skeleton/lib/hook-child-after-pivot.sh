#!/bin/bash

[ -n "$DEBUG" ] && set -o xtrace
set -o nounset
set -o errexit
shopt -s nullglob

cd $(dirname $0)/../

source etc/config

mount -n -t devpts -o newinstance,ptmxmode=0666 devpts /dev/pts

mount -n -t proc none /proc

mount -n -t tmpfs -o size=64k tmpfs /dev/shm

hostname $id

ip address add 127.0.0.1/8 dev lo
ip link set lo up

ip address add $network_container_ip/30 dev $network_container_iface
ip link set $network_container_iface mtu $container_iface_mtu up

ip route add default via $network_host_ip dev $network_container_iface

if [ -e /etc/seed ]; then
  . /etc/seed
fi
