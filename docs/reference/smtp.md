SMTP
====

| Name      | Type   | Required                     | Description                   |
| --------- | ------ | ---------------------------- | ----------------------------- |
| Host      | String | Y                            | IP or FQDN of the SMTP server |
| Username  | String | Y                            | Username for the SMTP server  |
| Password  | String | Y                            | Password for the SMTP server  |
| Sender    | String | Y                            | Who is sending the email      |
| Reciever  | String | Y                            | Who is receiving the email    |
| Body      | String | N :: "Hello from Scorestack" | Body of the email             |
| Encrypted | String | N :: "false"                 | Whether or not to use TLS     |
| Port      | String | N :: "25"                    | Port of the SMTP server       |