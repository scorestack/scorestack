Git
===

| Name            | Type    | Required   | Description                                           |
| --------------- | ------- | ---------- | ----------------------------------------------------- |
| Host            | String  | Y          | IP or FQDN the remote repository is located           |
| Repository      | String  | Y          | The path to the remote repository                     |
| Branch          | String  | Y          | The branch to clone from the repository               |
| Port            | Integer | N :: 80    | The port to connect to for cloning the repository     |
| HTTPS           | Boolean | N :: false | Whether to use HTTP or HTTPS                          |
| HTTPSValidate   | Boolean | N :: false | Whether HTTPS certificates should be validated        |
| Username        | String  | N          | Username to use for private repositories              |
| Password        | String  | N          | Password for the user                                 |
| ContentMatch    | Boolean | N :: false | Whether to check the contents of a file               |
| ContentFile     | String  | N          | The path of the file to check the contents of         |
| ContentRegex    | String  | N :: ".*"  | The regex to match against the checked file           |
| CommitHashMatch | Boolean | N :: false | Whether or not to match the hash of the latest commit |
| CommitHash      | String  | N          | The hash to check against the latest commit           |

Default Behavior
----------------
By default, this check will attempt to clone the configured repository and  pass or fail the check depending on whether the clone was successful or not.

`ContentFile` Parameter
---------
The file path is in respect to the *root* of the repository will accept absolute or relative file paths. The check while correspond "`/content.txt`", "`./content.txt`", or "`content.txt`" to a file named "`content.txt`" at the root of the repository.
