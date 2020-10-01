Changelog
=========

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

[Unreleased]
------------

[0.5.1] - 2020-10-01
--------------------

An intermediate release to support the transition of Dynamicbeat to go mod.

### Dynamicbeat

#### Changed

- Migrate to go mod

[0.5.0] - 2020-09-29
--------------------

This is the first public release of Scorestack.

### General

#### Added

- Administration documentation
- Check-writing documentation
- Binary building documentation
- Deployment documentation

#### Changed

- Kibana download URL in deployment automation
- Don't run `make testsuite` for Dynamicbeat CI

#### Fixed

- Scorestack casing

### Dynamicbeat

#### Added

- Explain required settings/permissions to run the ICMP protocol

#### Removed

- RITSEC GitLab links

#### Fixed

- GitHub import links

[0.4.0] - 2020-04-25
--------------------

This release implements features for IRSeC 2020.

### General

#### Added

- Example SMB check
- Example MySQL check

### Dynamicbeat

#### Added

- SMB check
- MySQL check

[0.3.0] - 2020-04-07
--------------------

This release focuses on some housekeeping tasks and Dynamicbeat bugfixes.

### General

#### Added

- GCP deployment automation
- Docker deployment automation

#### Changed

- Consolidate attributes into far fewer indices

#### Fixed

- Ensure deployment automation generates certificates for Dynamicbeat
- Fix Nginx firewall rules
- Use TCP proxying for Logstash instead of HTTP proxying

#### Removed

- Example checks for custom ISTS services

### Dynamicbeat

#### Changed

- Store check metadata in separate struct
- Refactor protocol code to use common helper functions for creation and running
- Ensure timeouts are strictly enforced
- Use bulk querying to update definitions from Elasticsearch

#### Fixed

- Prevent Dynamicbeat from crashing if an invalid check type is used
- Respond to the interupt signal properly

[0.2.0] - 2020-02-28
--------------------

This release is in preparation for ISTS 2020.

### General

#### Added

- Ansible playbook for deploying Dynamicbeat
- Example Active Directory LDAP check
- Example DNS check
- Example FTP check
- Example Gophish check
- Example Greenbone Security Assistant check
- Example ICMP check
- Example Kibana check
- Example Kolide check
- Example Roundcube check
- Example SSH check
- Example VNC check
- Example WinRM check
- Example XMPP check
- Example checks for custom ISTS services
- Elasticsearch coordinating-only node
- Proper Elastic Stack user roles
- Create generic, admin, and group check results

#### Changed

- Limit Dynamicbeat permissions
- Template in team name to the `add-team.sh` script

#### Fixed

- Set devcontainer environment variables

### Dynamicbeat

#### Added

- DNS protocol
- FTP protocol
- IMAP protocol
- LDAP protocol
- SMTP protocol
- SSH protocol
- VNC protocol
- WinRM protocol
- XMPP protocol
- Report check completion information
- `StoreValue` HTTP protocol parameter

#### Changed

- Allow SMTP plain authentication via unencrypted connections
- Enforce timeouts on checks
- Ensure FTP connections are closed
- Run checks asynchronously

#### Removed

- Don't display `SUCCESS` message when ICMP checks pass

#### Fixed

- Allow checks to run even if they don't have attributes
- Plug major goroutine leak
- Ensure WinRM protocol can run commands properly
- Fix HTTP regex-matching system
- Prevent Dynamicbeat from crashing if it can't reach Elasticsearch
- Ensure LDAP protocol reports the check name properly
- Close SSH connections after check finishes
- Prevent Dynamicbeat from crashing if XMPP checks fail
- Ensure Dynamicbeat loads all checks from Elasticsearch
- Properly return errors in ICMP protocol
- Don't overwrite HTTP check definitions

### Kibana Plugin

#### Added

- Dashboard for viewing a team's service uptime history
- Attribute modification page
- Organize services by group on Check Attributes page

[0.1.0] - 2020-02-13
--------------------

The initial release of Scorestack.

[Unreleased]: https://github.com/scorestack/scorestack/compare/v0.5.1...HEAD
[0.5.1]: https://github.com/scorestack/scorestack/compare/v0.5.0...v0.5.1
[0.5.0]: https://github.com/scorestack/scorestack/compare/v0.4...v0.5.0
[0.4.0]: https://github.com/scorestack/scorestack/compare/v0.3...v0.4
[0.3.0]: https://github.com/scorestack/scorestack/compare/v0.2...v0.3
[0.2.0]: https://github.com/scorestack/scorestack/compare/v0.1...v0.2
[0.1.0]: https://github.com/scorestack/scorestack/releases/tag/v0.1.0