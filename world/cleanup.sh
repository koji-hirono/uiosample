#!/bin/sh

ip netns del CL

echo "===> Unbind"
./dpdk-devbind.py --bind e1000 0000:00:11.0
./dpdk-devbind.py --status
