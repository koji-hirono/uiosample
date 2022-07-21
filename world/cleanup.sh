#!/bin/sh

ip netns del CL

echo "===> Unbind"
driverctl unset-override 0000:00:11.0
driverctl list-overrides
