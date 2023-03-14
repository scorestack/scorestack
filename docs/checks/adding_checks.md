Adding New Checks
=================

Once you've written all your checks, you can use Dynamicbeat to add them all to Scorestack. You can download a precompiled version of Dynamicbeat for your operating system from the [latest release on GitHub](https://github.com/scorestack/scorestack/releases/latest).

Save the zipfile and extract Dynamicbeat to the system that you'll be writing check configurations on. Next, you'll need to [configure Dynamicbeat](../dynamicbeat/configuration.md) so it knows where to add checks. Once you have Dynamicbeat downloaded and configured, you're ready to add your checks using Dynamicbeat's [`setup checks`](../dynamicbeat/reference/dynamicbeat_setup_checks.md) command.

> Dynamicbeat's `run` command is only supported on Linux right now, but the `setup` command and its subcommands should work fine on other operating systems.

To keep things simple, create a folder and place all your check files in it. Each file should be named `check-id.json`, where `check-id` is a unique identifier for each check that satisfies the requirements of the [`id` field](./metadata.md#id-omitted) in the check's metadata. Take a look at the directory structure of the repository's examples folder for an example:

```
examples
├── dns.json
├── ftp.json
├── git.json
├── http-gophish.json
├── http-greenbone-security.json
├── http-kibana-auth.json
├── http-kibana.json
├── http-kolide.json
├── http-roundcube.json
├── icmp.json
├── imap.json
├── ldap-ad.json
├── mssql.json
├── mysql.json
├── noop.json
├── postgresql.json
├── smb.json
├── smtp.json
├── ssh.json
├── vnc.json
├── winrm.json
└── xmpp.json
```

If you wish, you can add other files to this folder (e.g. `.gitignore`, `README.md`, `topology.png`, etc.) as long as they don't end in `.json`. Dynamicbeat expects that any file ending in `.json` is a check file. All other files will be ignored.

Dynamicbeat's `setup checks` command is idempotent; if you have to make changes to any of your checks, all you have to do is rerun the command.