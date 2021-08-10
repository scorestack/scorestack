Team Overrides
==============

A lot of times, checks can be easily configured for all teams through the use of attributes and the automatic [`TeamNum` override](../checks/attributes.md). However, sometimes a check file needs to be changed quite a bit from one team to another. This is where team overrides come in.

Team overrides are intended to make it possible to write a single check file that will work for the same service across all teams. With overrides, you shouldn't have to write any team-specific check files.

> Overrides can only be used to modify attributes. Overrides **cannot be applied directly to definitions**. Overrides must be applied to an attribute that is then used in your definition.

Defining Overrides
------------------

Team overrides are defined within the `teams` section of the [Dynamicbeat configuration file](./configuration.md#configuration-reference), within the `overrides` object for each team.

The `overrides` object can contain an arbitrary amount of key-value mappings. For each key-value pair in `overrides`, the key is used as the override name, and the value is used as the override's value.

Here's an example configuration snippet that configures two teams with an override named `MyOverride`:

```yaml
teams:
  - name: team01
    overrides:
      MyOverride: "team one's override value"
  - name: team02
    overrides:
      MyOverride: "team two's override value"
```

### `TeamNum`

If an override named `TeamNum` has not been configured, Dynamicbeat will automatically add one at runtime by parsing a number from the end of the team name and removing all leading zeros.

If the team's name

To verify that Dynamicbeat will parse the `TeamNum` from your team names properly, you can evaluate your team names against the regex using [this regex101 link](https://regex101.com/r/AlUufl/1). The `TeamNum` override will be set to the value of **Group 1** in the right-hand pane in the **MATCH INFORMATION** section.

Override Methods
----------------

Overrides will be applied with both templating and replacement.

### Templating

Override templating applies overrides to attributes in the same way [attributes are applied to check definitions](../checks/attributes.md).

Given the following teams configuration:

```yaml
teams:
  - name: team01
    overrides:
      UsernamePrefix: "oneIsNone"
```

With these attributes defined in a check file:

```json
{
  "attributes": {
    "user": {
      "Username": "admin{{.UsernamePrefix}}"
    }
  }
}
```

The `Username` attribute for `team01` would have a value of `adminoneIsNone`, because the `UsernamePrefix` override would be templated into the attribute value with the `{{.UsernamePrefix}}` reference.

### Replacement

Override replacement allows you to entirely replace the value of an attribute that has the same name as your override.

Given the following teams configuration:

```yaml
teams:
  - name: team01
    overrides:
      Username: "oneIsNone"
```

With these attributes defined in a check file:

```json
{
  "attributes": {
    "user": {
      "Username": "root"
    }
  }
}
```

The `Username` attribute for `team01` would have a value of `oneIsNone`, because the `Username` override replaces the value of the `Username` attribute.

> If you don't define a replacement attribute for one of your teams, the attribute **will not be replaced** for that team.

### Both

Overrides will perform both templating and replacement over a check's attributes if possible.

Given the following teams configuration:

```yaml
teams:
  - name: team01
    overrides:
      Domain: "example.com"
      Username: "oneIsNone"
```

With these attributes defined in a check file:

```json
{
  "attributes": {
    "admin": {
      "Host": "www.{{.Domain}}"
    },
    "user": {
      "Username": "root"
    }
  }
}
```

`team01`'s `Host` attribute would have a value of `www.example.com` and their `Username` attribute would have a value of `oneIsNone`.

Example Uses
------------

### Custom Subnets

If your subnetting configuration won't work with `TeamNum`, you can use overrides to keep your check files generic.

For example, let's say you need to run a check against three teams, but their hosts are on wildly different subnets:

| Team   | Service IP    |
| ------ | ------------- |
| team01 | 10.0.0.50     |
| team02 | 172.16.3.50   |
| team03 | 192.168.22.50 |

You can write a single check file for this service using overrides.

With these attributes defined in the check file:

```json
{
  "attributes": {
    "admin": {
      "Host": "{{.SubnetPrefix}}.50"
    }
  }
}
```

With this teams configuration:

```yaml
teams:
  - name: team01
    overrides:
      SubnetPrefix: "10.0.0"
  - name: team02
    overrides:
      SubnetPrefix: "172.16.3"
  - name: team03
    overrides:
      SubnetPrefix: "192.168.22.50"
```

### Unique Passwords

In competitions where teams can attack each other ("purple team" or "attack/defend" competitions), it might be a good idea to give each team a different set of default credentials. This can help prevent the fastest (and meanest) teams from destroying the systems of other teams with default credentials early on.

For example, let's say you have three teams, each with the following default passwords:

| Team   | Default Password   |
| ------ | ------------------ |
| team01 | `Changeme123`      |
| team02 | `PasswordPassword` |
| team03 | `GoScorestack!`    |

You can configure this using overrides.

With these attributes defined in the check file:

```json
{
  "attributes": {
    "user": {
      "Password": "supersecure"
    }
  }
}
```

With this teams configuration:

```yaml
teams:
  - name: team01
    overrides:
      Password: "Changeme123"
  - name: team02
    overrides:
      Password: "PasswordPassword"
  - name: team03
    overrides:
      Password: "GoScorestack!"
```