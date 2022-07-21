#!/bin/sh
#
# CL
#
# +--------+ 30.30.0.1/24  +--------+            +--------+
# | client |---------------| l2fwd  |------------| server |
# +--------+               +--------+         .2 +--------+
#
# enp0s16              enp0s17
# 0000:00:10.0         0000:00:11.0
# 08:00:27:a0:a9:3e    08:00:27:af:62:0b
#
#                          enp0s9                enp0s10
#                          0000:00:09.0          0000:00:0a.0
#                          08:00:27:15:8c:70     08:00:27:dc:7d:2a
#

ip netns add CL
ip link set enp0s16 up netns CL
ip netns exec CL ip addr add 30.30.0.1/24 dev enp0s16
ip netns exec CL ip addr add 127.0.0.1/8 dev lo
ip netns exec CL ip link set lo up


ip netns add SV
ip link set enp0s10 up netns SV
ip netns exec SV ip addr add 30.30.0.2/24 dev enp0s10
ip netns exec SV ip addr add 127.0.0.1/8 dev lo
ip netns exec SV ip link set lo up

modprobe uio_pci_generic

echo "===> Bind"
ip link set enp0s17 down
driverctl set-override 0000:00:11.0 uio_pci_generic
ip link set enp0s9 down
driverctl set-override 0000:00:09.0 uio_pci_generic
driverctl list-overrides
