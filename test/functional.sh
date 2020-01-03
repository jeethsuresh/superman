#!/bin/bash

echo "beginning functional tests"

#valid call
curl -X POST http://127.0.0.1:8080/v1 -d '{"username": "bob", "unix_timestamp": 1514764800, "event_uuid": "4B4B4499-29B1-461E-87ED-B8E163A53DAC", "ip_address": "206.81.252.6"}'

#invalid: does not contain IP
curl -X POST http://127.0.0.1:8080/v1 -d '{"username": "bob", "unix_timestamp": 1514764800, "event_uuid": "4B4B4499-29B1-461E-87ED-B8E163A53DAC"}'

# Test full chain: 
curl -X POST http://127.0.0.1:8080/v1 -d '{"username": "bob", "unix_timestamp": 1514564802, "event_uuid": "4B4B4499-29B1-461E-87ED-B8E163A53DAD", "ip_address": "220.81.252.6"}'
curl -X POST http://127.0.0.1:8080/v1 -d '{"username": "bob", "unix_timestamp": 1614764803, "event_uuid": "4B4B4499-29B1-461E-87ED-B8E163A53DBC", "ip_address": "180.81.252.6"}'
curl -X POST http://127.0.0.1:8080/v1 -d '{"username": "bob", "unix_timestamp": 1614564800, "event_uuid": "4B4B4499-29B1-461E-87ED-B9E163A53DAC", "ip_address": "206.81.252.6"}'

#Test different user