The Check Configuration
=======================

Each check is configured via a single JSON file that contains the check metadata, definition, and any attributes. This page will go through an example check file section-by-section to explain the contents of the check definition and how to write your own. We will be using a modified version of the IMAP check found in the repository's examples folder.

Here is the complete example check:

```json
{
  "name": "Email Clients",
  "type": "imap",
  "score_weight": 1,
  "definition": {
    "Host": "{{.Host}}",
    "Port": "143",
    "Username": "{{.Username}}@example.com",
    "Password": "{{.Password}}",
  },
  "attributes": {
    "admin": {
      "Host": "10.0.{{.TeamNum}}.50",
      "Username": "admin"
    },
    "user": {
      "Password": "changeme"
    }
  }
}
```

Check Metadata
--------------

All checks have a few metadata values that are required to be present.

```json
{
  "name": "Email Clients",
  "type": "imap",
  "score_weight": 1,
}
```

### ID (Omitted)

> This field will be automatically populated by Dynamicbeat based on the filename, and should be omitted from check files. Take a look at the [page on adding checks](./adding_checks.md) for more information.

The ID field is used to uniquely identify the check. It may only contain the ASCII letters A-Z (uppercase or lowercase), the digits 0-9, underscores (`_`), and hyphens (`-`).

Additionally, the ID field must end in a hyphen followed by the [group field](#group). For example, a check that has a group of `team03` may have a valid ID of `wordpress-team03`, but the IDs `team03-wordpress`, `wordpress`, and `wordpressteam03` would be invalid.

### Name

The Name field is a human-readable string that is used when displaying the check in Kibana, such as on the Scoreboard and the Attributes application. It does not have to be unique; actually, it's recommended that all checks that are similar to each other have the same name.

Good names are generally short, identify what service is being checked, and clarify where the service is running. For example, if you have multiple SSH checks for several Linux systems, you may want to include the hostname of the system each check connects to in the check's name.

### Type

The Type field is a keyword string that defines the check type that will be used to execute the check.

### Group

> This field will be automatically populated by Dynamicbeat with the team name of each configured team, and should be omitted from check files.

The Group field contains the name of the team associated with this check. It may only contain the ASCII letters A-Z (uppercase or lowercase), the digits 0-9, and underscores (`_`).

### Score Weight

The Score Weight field defines the number of points that will be awarded for a successful check. This is typically set to 1 for all checks, but it can be changed to make some checks worth more than others. For example, a functioning e-commerce webserver should probably be worth more points per check than SSH access to a user's workstation.

Check Definition
----------------

The check definition contains the parameters that define how the check will be executed. For most checks it will just be a simple key-value object, but it depends on which check type is being defined. Please see [the check reference](../reference.md) for more information on what paremeters are expected within the check definition for each check type.

Check Attributes
----------------

Attributes are named variables references within a check definition that allow string values to be inserted into a definition at runtime. This allows for dynamic updates to scoring through the Kibana UI during a competition.

Attributes work using [golang templates](https://golang.org/pkg/text/template/) to insert attribute values into the check definitions before they are executed.

There are two kinds of attributes that can be used for different purposes.

### User Attributes

User attributes are the more common kind of attributes. These attributes can only be viewed or changed by members of the same group as the check, or by users with the `attribute-admin` role. For example, the attributes for a check in the `blue` group is restricted to users that are also in the `blue` group. These attributes are useful for values that may be changed by blueteamers, such as service account credentials.

Another use case for user attributes besides changing credentials is service migration. If you want to let teams migrate their services at will, then you could achieve this by configuring the host and/or port parameters in the definition via user attributes. Then, teams would be able to change where their checks are run whenever they want.

### Administrator Attributes

Administrator attributes (or admin attributes for short) are less common than user attributes, but can still be useful. Users who are members of the same group as the check may view these attributes, but cannot modify them. Users with the `attribute-admin` role may both view and modify these attributes.

The primary purpose of admin attributes is to give teams some information about their checks that they can see, but cannot modify. For example, if you want to let teams modify the password for an SSH check, but not the username, you should probably configure the username via an admin attribute. That way teams will know which user the password is associated with, but they won't be able to change the user that the check is configured to use.

It could also be useful to configure the IP address or hostname via an admin attribute. Typically it's preferred to use [the name field](#name) to identify the system a check is running against, but sometimes it's good to make it explicit.

Another use of admin attributes is, like user attributes, to make it easy to change check configurations on-the-fly during a competition. Let's say you want to let teams migrate their services, but you want them to submit an official request to be approved first. You could configure the host and/or port parameters via an admin attribute. Then, once an official migration request has been submitted and approved, a competition organizer can change the attributes in Kibana.