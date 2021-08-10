Check Metadata
==============

All checks have a few metadata values. These are used by Scorestack and Dynamicbeat to determine how checks appear in Kibana and how Dynamicbeat will run them.

Here is the metadata section of our example IMAP check:

```json
{
  "name": "Email Clients",
  "type": "imap",
  "score_weight": 1,
}
```

ID (Omitted)
------------

> This field will be automatically populated by Dynamicbeat based on the filename, and should be omitted from check files. Take a look at the [page on adding checks](./adding_checks.md) for more information.

The ID field is used to uniquely identify the check. It may only contain the ASCII letters A-Z (uppercase or lowercase), the digits 0-9, underscores (`_`), and hyphens (`-`).

Additionally, the ID field must end in a hyphen followed by the [group field](#group). For example, a check that has a group of `team03` may have a valid ID of `wordpress-team03`, but the IDs `team03-wordpress`, `wordpress`, and `wordpressteam03` would be invalid.

Name
----

The Name field is a human-readable string that is used when displaying the check in Kibana, such as on the Scoreboard and the Attributes application. It does not have to be unique; actually, it's recommended that all checks that are similar to each other have the same name.

Good names are generally short, identify what service is being checked, and clarify where the service is running. For example, if you have multiple SSH checks for several Linux systems, you may want to include the hostname of the system each check connects to in the check's name.

Type
----

The Type field is a keyword string that defines the check type that will be used to execute the check.

Group
-----

> This field will be automatically populated by Dynamicbeat with the team name of each configured team, and should be omitted from check files.

The Group field contains the name of the team associated with this check. It may only contain the ASCII letters A-Z (uppercase or lowercase), the digits 0-9, and underscores (`_`).

Score Weight
------------

The Score Weight field defines the number of points that will be awarded for a successful check. This is typically set to 1 for all checks, but it can be changed to make some checks worth more than others. For example, a functioning e-commerce webserver should probably be worth more points per check than SSH access to a user's workstation.