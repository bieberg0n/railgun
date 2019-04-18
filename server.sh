#!/usr/bin/env bash

ip addr add 10.1.0.1/24 dev tun1
ip link set dev tun1 up
