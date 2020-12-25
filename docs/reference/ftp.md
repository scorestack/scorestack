FTP
===

| Name             | Type   | Required     | Description                                                                        |
| ---------------- | ------ | ------------ | ---------------------------------------------------------------------------------- |
| Host             | String | Y            | IP or hostname of the host to run the FTP check against                            |
| Username         | String | Y            | The user to login with over FTP                                                    |
| Password         | String | Y            | The password for the user that you wish to login with                              |
| File             | String | Y            | The path to the file to access during the FTP check                                |
| ContentRegex     | String | N :: "\.\*"  | The regex to use to match against the file contents                                |
| HashContentMatch | String | N :: "false" | Whether to use hash-based matching to check the file contents                      |
| Hash             | String | N            | The sha3\-256 hash to use when checking the file contents with hash-based matching |
| Port             | String | N :: "21"    | The port to attempt an FTP connection on                                           |
| Simple           | String | N :: "false" | Very simple FTP check for older servers                                            |

The `simple` parameter should usually be left as the default unless you are running a check against an FTP server that is very old and supports a limited set of FTP commands. When `simple` is set to `"true"`, the check will only change to the directory specified in the `file` parameter and then query for the current working directory. If both of these operations succeed, then the check will pass. All other parameters will be ignored.

Note that if you are using the simple version of the FTP check, the `file` parameter should point to a directory, _not_ a file on the system.