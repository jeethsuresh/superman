#!/bin/bash

echo "beginning functional tests"

#valid call
curl -X POST http://127.0.0.1:5000/v1 -d '{"username": "bob", "unix_timestamp": 1514764800, "event_uuid": "4B4B4499-29B1-461E-87ED-B8E163A53DAC", "ip_address": "206.81.252.6"}'

# Test full chain: 
curl -X POST http://127.0.0.1:5000/v1 -d '{"username": "bob", "unix_timestamp": 1514564802, "event_uuid": "4B4B4499-29B1-461E-87ED-B8E163A53DAD", "ip_address": "220.81.252.6"}'
curl -X POST http://127.0.0.1:5000/v1 -d '{"username": "bob", "unix_timestamp": 1614764803, "event_uuid": "4B4B4499-29B1-461E-87ED-B8E163A53DBC", "ip_address": "180.81.252.6"}'
curl -X POST http://127.0.0.1:5000/v1 -d '{"username": "bob", "unix_timestamp": 1614564800, "event_uuid": "4B4B4499-29B1-461E-87ED-B9E163A53DAC", "ip_address": "206.81.252.6"}'

#Test different user at the same time
curl -X POST http://127.0.0.1:5000/v1 -d '{"username": "alice", "unix_timestamp": 1614564800, "event_uuid": "4B4B4499-29B1-461E-87ED-B9E163A53EBC", "ip_address": "206.81.252.6"}'

#Test bad UUID
curl -X POST http://127.0.0.1:5000/v1 -d '{"username": "alice", "unix_timestamp": 1614564800, "event_uuid": "4B4B4499-29B1-461-87ED-B9E16", "ip_address": "206.81.252.6"}'

#Test no IP
curl -X POST http://127.0.0.1:5000/v1 -d '{"username": "alice", "unix_timestamp": 1614564800, "event_uuid": "4B4B4499-29B1-461E-87ED-B9E173A53EBC"}'

#Test bad IP string
curl -X POST http://127.0.0.1:5000/v1 -d '{"username": "alice", "unix_timestamp": 1614564800, "event_uuid": "4B4B4499-29B1-461E-87DD-B9E163A53EBC", "ip_address": "206.81.6"}'