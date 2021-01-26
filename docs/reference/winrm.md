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

Picking a Command
-----------------

While the `Cmd` parameter can generally be used for whatever you'd like, it is recommended to avoid the `netstat` command. Some Windows systems can have spurious check failures if check runs the `netstat` command. Other systems may take 2-3 seconds to print each line of `netstat` output, impacting usability for blueteamers and possibly causing checks to fail due to timeout issues. 

Issues with the `netstat` command _may_ be infrastructure-specific. If you'd really like to check a user's ability to run the `netstat` command, just make sure to test the check thoroughly before the competition _on your competition infrastructure_. One good way of doing this would be running a Dynamicbeat instance against all your checks with a very low period, such as 5-10 seconds. If you have no issues with using `netstat` as your checked command, you should be fine!

Troubleshooting
---------------

Sometimes this check will fail even if you can connect to the system over WinRM using other tools, like Ansible or PowerShell. This is because the library used for this check only supports certain WinRM configurations, and configurations that are valid for other tools may not work with this one.

When troubleshooting this check, it may be useful to use [winrm-cli](https://github.com/masterzen/winrm-cli), a WinRM command-line client written in golang. The tool uses the same library as this check, and is written by the same author. If you are able to make a connection using winrm-cli, then this check should be passing. If you're still having issues, open an issue! In that case, there may be an issue with the Dynamicbeat code.