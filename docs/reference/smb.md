SMB
===

| Name         | Type   | Required    | Description                                                                       |
| ------------ | ------ | ----------- | --------------------------------------------------------------------------------- |
| Host         | String | Y           | IP or FQDN for the SMB server                                                     |
| Username     | String | Y           | Username for SMB share                                                            |
| Password     | String | Y           | Password for the user                                                             |
| Share        | String | Y           | Name of the SMB share                                                             |
| Domain       | String | Y           | The domain found in front of a login \(SMB\\Administrator : SMB would be domain\) |
| File         | String | Y           | The file in SMB share                                                             |
| ContentRegex | String | N :: "\.\*" | Regex to match on                                                                 |
| Port         | String | N :: "445"  | Port of the server                                                                |