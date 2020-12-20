ICMP
====

| Name            | Type   | Required     | Description                                                                               |
| --------------- | ------ | ------------ | ----------------------------------------------------------------------------------------- |
| Host            | String | Y            | IP or FQDN of the host to run the ICMP check against                                      |
| Count           | Int    | N :: 1       | The number of ICMP requests to send per check                                             |
| AllowPacketLoss | String | N :: "false" | Pass check based on received pings matching Count; if false, will use percent packet loss |
| Percent         | Int    | N :: 100     | Percent of packets needed to come back to pass the check                                  |