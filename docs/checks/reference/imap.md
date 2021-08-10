IMAP
====

| Name      | Type   | Required     | Description                         |
| --------- | ------ | ------------ | ----------------------------------- |
| Host      | String | Y            | IP or FQDN for the IMAP server      |
| Username  | String | Y            | Username for the IMAP server        |
| Password  | String | Y            | Password for the user               |
| Encrypted | String | N :: "false" | Whether or not to use TLS \(IMAPS\) |
| Port      | String | N :: "143"   | Port for the IMAP server            |