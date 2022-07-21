#!/bin/sh
#
# CL
#
# +--------+ 30.30.0.1/24  +--------+
# | client |---------------| server |
# +--------+            .2 +--------+
#
# enp0s16                  enp0s17
# 0000:00:10.0             0000:00:11.0
# 08:00:27:d4:c8:f8        08:00:27:d1:29:ba
#
#

ip netns add CL
ip link set enp0s16 up netns CL
ip netns exec CL ip addr add 30.30.0.1/24 dev enp0s16
ip netns exec CL ip addr add 127.0.0.1/8 dev lo
ip netns exec CL ip link set lo up


modprobe uio_pci_generic

echo "===> Bind"
ip link set enp0s17 down
driverctl set-override 0000:00:11.0 uio_pci_generic
driverctl list-overrides
