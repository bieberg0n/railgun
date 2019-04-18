#!/usr/bin/env bash

ip addr add 10.1.0.1/24 dev tun0
ip link set dev tun0 up
