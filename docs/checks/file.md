The Check File
==============

Each Scorestack check is configured via a single JSON file, called a "check file". Each check file contains the check's metadata, the check definitions, and any attributes defined for the check.

The next few pages will go through an example check file section-by-section to explain the content of a check file and how to write your own. We will be using a modified version of the example IMAP check found in the repository's examples folder.

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

This check file would be saved to `imap-clients.json`.