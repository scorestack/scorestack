Adding New Checks
=================

In order to configure and add checks to Scorestack, you will have to create a folder structure containing the JSON definitions for each check. The `add-team.sh` script, which is used to add new checks to Scorestack or update existing checks, expects the following folder structure.

```
myChecks
├── dns-host1
│   ├── admin-attribs.json
│   └── check.json
└── http-host2
    ├── admin-attribs.json
    ├── check.json
    └── user-attribs.json
```

The top level directory (`myChecks`) contains subfolders for all of the desired checks (`dns-host1`, `http-host2`, etc.). Each of these subfolders will contain up to three JSON files that will define the specific check. The three JSON files are as follows:

- `check.json`: The main JSON document that provides values for any of a check's required arguments, and optionally overrides the defaults for a some of a check's optional arguments.
- `admin-attribs.json`: The attributes of a check that only an administrator (such as a competition volunteer) will be allowed to change during a competition through Kibana.
- `user-attribs.json`: The attributes of a check that both users and administrators will be allowed to change during a competition through Kibana.