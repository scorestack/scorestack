The Check Configuration
=======================

Each check is configured via a single JSON file that contains the check metadata, definition, and any attributes. The next few pages will go through an example check file section-by-section to explain the contents of the check definition and how to write your own. We will be using a modified version of the IMAP check found in the repository's examples folder.

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