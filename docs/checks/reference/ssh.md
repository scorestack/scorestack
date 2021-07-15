SSH
===

| Name         | Type   | Required     | Description                                            |
| ------------ | ------ | ------------ | ------------------------------------------------------ |
| Host         | String | Y            | IP or FQDN of the host to run the SSH check against    |
| Username     | String | Y            | The user to login with over SSH                        |
| Password     | String | Y (or PrivateKey)            | The password for the user that you wish to login with  |
| PrivateKey   | String | Y (or Password)            | The absolute filepath of the private key file for the user that you wish to login with  |
| Cmd          | String | Y            | The command to execute once SSH connection established |
| MatchContent | String | N :: "false" | Whether or not to match content like checking files    |
| ContentRegex | String | N :: "\.\*"  | Tegex to match if reading a file                       |
| Port         | String | N :: "22"    | The port to attempt an SSH connection on               |

Authentication
--------------

SSH checks require either a password, private key filepath, or both. If either is supplied Dynamicbeat will use that respective method to authenticate to the server. If both are specified, it will assume the private key is encrypted with the assigned password.

Notes on FreeBSD
----------------

SSH checks will fail on standard FreeBSD installations. This is because by default, FreeBSD does not enable `password` authentication for SSH. In order to fix this issue, you must ensure that `PasswordAuthentication yes` is set in your `/etc/ssh/sshd_config` file.


