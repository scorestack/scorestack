Attributes
==========

Attributes are named variables referenced within a check definition that allow string values to be inserted into a definition at runtime. This allows for dynamic updates to scoring through the Kibana UI during a competition, and is typically useful for things that change a lot - like usernames and/or passwords.

Attributes work using [golang templates](https://golang.org/pkg/text/template/) to insert attribute values into the check definitions before they are parsed by Dynamicbeat.

To define an attribute, you just have to use the `{{.Name}}` syntax somewhere within a check definition, where `Name` is the name of the attribute you are defining. Then, you must give the attribute an initial value in one of two files: `admin-attribs.json` or `user-attribs.json`.

The format of these two attribute files is a flat JSON object of key-value pairs. Each key is the name of an attribute, and each value is the string that will be inserted into the check definition at the location where it is referenced.

Permissions
-----------

The only difference between the two attribute files is the permissions on the attributes defined within them. Which file an attribute is defined in determines which users are able to view and modify the attribute.

The `admin-attribs.json` file contains _administrator_ attributes. These attributes can only be viewed and changed by users with the `attribute-admin` role. Administrator attributes can be useful for updating checks to react to injects, such as changing a check from HTTP to HTTPS. They can also be a method of tweaking and debugging checks as they are being written, but this is generally not recommended.

The `user-attribs.json` file contains _user_ or _team_ attributes. These attributes can only be viewed and changed by members of the same group as the check. For example, the attributes for a check in the `blue` group is restricted to users that are also in the `blue` group. These attributes are useful for values that may be changed by blueteamers, such as service account credentials.

Examples
--------

Here's the same definition [from the `check.json` page](./check_json.md), except the `host` argument references an attribute named `Host`.

```json
{
    "id": "icmp",
    "name": "ICMP",
    "type": "icmp",
    "group": "example",
    "score_weight": 1,
    "definition": {
        "host": "{{.Host}}"
    }
}
```

An `admin-attribs.json` file for this check could look like this.

```json
{
  "Host": "localhost"
}
```

Note that the `Host` key is case-sensitive. Also, if the `Host` attribute should be a user attribute, the file would have the same format, but it would be named `user-attribs.json` instead.