Check Attributes
================

Attributes are named variables references within a check definition that allow string values to be inserted into a definition at runtime. This allows for dynamic updates to scoring through the Kibana UI during a competition.

Attributes work using [golang templates](https://golang.org/pkg/text/template/) to insert attribute values into the check definitions before they are executed.

Here are the attributes of our example IMAP check, modified to remove the [team overrides](#team-overrides):

```json
{
  "attributes": {
    "admin": {
      "Host": "10.0.0.50",
      "Username": "admin"
    },
    "user": {
      "Password": "changeme"
    }
  }
}
```

Each attribute is a key-value mapping. To use an attribute in a definition, add the attribute's key with a dot before it between two pairs of curly braces.

To reference our IMAP example, if we wanted to use the `Username` attribute somewhere within our definition, we would write `{{.Username}}`.

When adding checks to Scorestack, Dynamicbeat will substitute the double curly braces with the value of the referenced attribute.

Here's the definition of our example IMAP check, before attributes are applied:

```json
{
  "definition": {
    "Host": "{{.Host}}",
    "Port": "143",
    "Username": "{{.Username}}@example.com",
    "Password": "{{.Password}}",
  }
}
```

Here's what the definition will look like after attributes are applied, based on the attributes defined at the top of the page:

```json
{
  "definition": {
    "Host": "10.0.0.50",
    "Port": "143",
    "Username": "admin@example.com",
    "Password": "changeme",
  }
}
```

> If you reference an attribute that hasn't been defined, Dynamicbeat will replace the curly braces with an empty string.
>
> In our IMAP example, `"{{.DoesNotExist}} - test"` would become `" - test"`, since there's no attribute defined named `DoesNotExist`.

Attribute Types
---------------

There are two kinds of attributes that can be defined, each serving a different purpose.

### User Attributes

User attributes are the more common kind of attributes. These attributes can only be viewed or changed by members of the same group as the check, or by users with the `attribute-admin` role. For example, the attributes for a check in the `blue` group is restricted to users that are also in the `blue` group. These attributes are useful for values that may be changed by blueteamers, such as service account credentials.

Another use case for user attributes besides changing credentials is service migration. If you want to let teams migrate their services at will, then you could achieve this by configuring the host and/or port parameters in the definition via user attributes. Then, teams would be able to change where their checks are run whenever they want.

### Administrator Attributes

Administrator attributes (or admin attributes for short) are less common than user attributes, but can still be useful. Users who are members of the same group as the check may view these attributes, but cannot modify them. Users with the `attribute-admin` role may both view and modify these attributes.

The primary purpose of admin attributes is to give teams some information about their checks that they can see, but cannot modify. For example, if you want to let teams modify the password for an SSH check, but not the username, you should probably configure the username via an admin attribute. That way teams will know which user the password is associated with, but they won't be able to change the user that the check is configured to use.

It could also be useful to configure the IP address or hostname via an admin attribute. Typically it's preferred to use [the name field](./metadata.md#name) to identify the system a check is running against, but sometimes it's good to make it explicit.

Another use of admin attributes is, like user attributes, to make it easy to change check configurations on-the-fly during a competition. Let's say you want to let teams migrate their services, but you want them to submit an official request to be approved first. You could configure the host and/or port parameters via an admin attribute. Then, once an official migration request has been submitted and approved, a competition organizer can change the attributes in Kibana.

Team Overrides
--------------

Attributes can be customized on a per-team basis via the use of team override variables. These are inserted into attribute values in the same way that [attributes are inserted into definition parameters](#attribute-templating).

For more possible uses of team overrides, please see [the team override guide](../dynamicbeat/overrides.md).

### `TeamNum`

The most important team override is `TeamNum`, which is automatically created for you by Dynamicbeat (unless manually overridden). `TeamNum` is an unsigned integer parsed from the team name. If a team name ends in a number, that number will be parsed out and made available via the `TeamNum` override. For example, `team10` would have a `TeamNum` of 10.

Additionally, any leading zeros will be removed from the `TeamNum`, so `team05` would have a `TeamNum` of 5.

This is particularly useful if your teams are named in a sequence and their networks are similarly sequenced. For example, if your competition has `team01` on `10.0.1.0/24`, `team02` on `10.0.2.0/24`, and `team03` on `10.0.3.0/24`, then you could use the `TeamNum` override to configure the IP addresses for checks in a team-generic way.

Our IMAP check provides an example of this with the `Host` admin attribute:

```json
{
  "attributes": {
    "admin": {
      "Host": "10.0.{{.TeamNum}}.50",
    },
  },
}
```

If this check was added for a team named `team07`, then the `Host` attribute would have a value of `10.0.9.50`.