WinRM
=====

| Name         | Type   | Required     | Description                                              |
| ------------ | ------ | ------------ | -------------------------------------------------------- |
| Host         | String | Y            | IP or FQDN of the WinRM machine                          |
| Username     | String | Y            | User to login as. Must be a local user                   |
| Password     | String | Y            | Password for the user                                    |
| Cmd          | String | Y            | Command that will be executed                            |
| Encrypted    | String | N :: "true"  | Use TLS for connection                                   |
| MatchContent | String | N :: "false" | Turn this on to match content from the output of the cmd |
| ContentRegex | String | N :: "\.\*"  | Regexp for matching output of a command                  |
| Port         | String | N :: "5986"  | Port for WinRM                                           |

> This check will _only_ work with local users using basic authentication, because [the WinRM library used only supports that authentication method](https://github.com/masterzen/winrm#preparing-the-remote-windows-machine-for-basic-authentication).

Sometimes this check will fail even if you can connect to the system over WinRM using other tools, like Ansible or PowerShell. This is because the library used for this check only supports certain WinRM configurations, and configurations that are valid for other tools may not work with this one. To 