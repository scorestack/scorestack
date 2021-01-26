DNS
===

| Name       | Type   | Required  | Description                                    |
| ---------- | ------ | --------- | ---------------------------------------------- |
| Server     | String | Y         | IP of the DNS server to query                  |
| Fqdn       | String | Y         | The FQDN of the host you are looking up        |
| ExpectedIP | String | Y         | The expected IP of the host you are looking up |
| Port       | String | N :: "53" | The port of the DNS server                     |