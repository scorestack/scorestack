LDAP
====

| Name     | Type   | Required     | Description                            |
| -------- | ------ | ------------ | -------------------------------------- |
| User     | String | Y            | The user written in user@domain syntax |
| Password | String | Y            | The password for the user              |
| Fqdn     | String | Y            | The FQDN of the LDAP server            |
| Ldaps    | String | N :: "false" | Whether or not to use LDAP\+TLS        |
| Port     | String | N :: "389"   | Port for LDAP server                   |