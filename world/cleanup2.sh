#!/bin/sh

ip netns del CL
ip netns del SV

echo "===> Unbind"
driverctl unset-override 0000:00:11.0
driverctl unset-override 0000:00:09.0
driverctl list-overrides
