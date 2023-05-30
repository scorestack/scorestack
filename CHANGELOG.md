# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Each release in the changelog is divided into the following sections:

- General: changes to anything other than the Dynamicbeat or Kibana plugin code
- Dynamicbeat
- Kibana plugin

Each section organizes entries into the following subsections:

- Added
- Changed
- Deprecated
- Removed
- Fixed
- Security

## [Unreleased]

- Fix check logic in Kibana client to look for appropriate fields on newer versions of Kibana (#344)

### Dynamicbeat

#### Added
- Git check type (#346)
- Added threadding to document indexing (#347)

#### Changed
- Bumped Go to 1.20 (#384)
- Bumped golangci-lint to v1.52.2 (#384)

## [0.8.2] - 2021-09-28

This release fixes a Dynamicbeat bug in the team overrides system.

### Dynamicbeat

#### Fixed

- Ensure template overrides are properly templated into attribute values as expected (#324)

## [0.8.1] - 2021-09-12

This release fixes a few bugs in Dynamicbeat related to the new deployment and check-adding processes.

### Dynamicbeat

#### Fixed

- Replace remaining `${VAR}` template strings with `{{.Var}}` templates in saved objects (#321)
- Ensure generic check definitions are added by `dynamicbeat setup checks` (#321)

## [0.8.0] - 2021-08-10

This release overhauls the structure of Dynamicbeat and improves the deployment and check-adding processes.

### General

#### Changed

- Each check is now defined using only one JSON file that includes all attributes (#312)
- Example checks now only use attributes when necessary (#312)
- Replaced `${TEAM}` and `${TEAM_NUM}` check definition variables with golang template blocks (#312)

#### Removed

- Drop Logstash from architecture (#302)
- Delete `add-team.sh` and `update.sh` scripts (#312)
- Deployments no longer configure users, indices, and dashboards (#312)
- Check definition files no longer include `id` or `group` fields (#312)

### Dynamicbeat

#### Added

- `setup` command and subcommands for initializing Scorestack and adding checks (#312)
- Support overriding attributes on a per-team basis (#312)

#### Changed

- Dynamicbeat is now a standalone program that doesn't use libbeat (#302)
- Time-separated index patterns are no longer used for check results - each `results-*-*` pattern has been replaced by a single index (#310)
- Upgrade to Golang 1.16.2 (#312)

#### Removed

- Remove `update_period` setting, update configurations after starting each round (#302)

#### Fixed

- Ensure Dynamicbeat doesn't segfault if a check fails to parse properly (#316)

## [0.7.0] - 2021-02-21

This release adds two new check types and fixes a logging bug with the SSH check.

### Dynamicbeat

#### Added

- PostgreSQL check type (#294)
- MSSQL check type (#295)

#### Fixed

- Spurious EOF errors when closing SSH connections are ignored (#292)

## [0.6.2] - 2021-02-16

This release fixes a bug with the SSH check in Dynamicbeat.

### Dynamicbeat

#### Fixed

- SSH now matches content when MatchContent is true (#288)

## [0.6.1] - 2021-01-26

This release mainly fixes deployment issues and improves the documentation.

### General

#### Changed

- Migrated documentation to `mdbook`
- Changed WinRM example check's command to `whoami`
- Adjusted default dashboard refresh time to 30 seconds (#284)
- Standardized dashboard and visualization IDs (#284)

#### Fixed

- Documented `add-team.sh`'s dependency on `jq` (#261)
- Removed outdated `jvm.options` configurations (#276)
- Fixed permissions for Kibana plugin installation (#282)
- Set default for `fqdn` Terraform variable (#282)
- Explained Kibana sorting method (#285)
- Documented FreeBSD's behavior with SSH check (#285)

#### Removed

- Stop configuring JVM heap size for Logstash (#282)

### Dynamicbeat

#### Changed

- Replaced Golang `html/template` library with `text/template` (#251)

## [0.6.0] - 2020-10-17

This release upgrades Scorestack to use Elastic Stack 7.9.2, the latest released version as of this writing. It also fixes some bugs with Dynamicbeat's check template system.

### General

#### Added

- Prebuild scorestack/kibana:7.9.2 container for CI and devcontainer

#### Changed

- Run `yarn kbn bootstrap` during Kibana plugin container build process
- Update CI to use prebuilt Kibana container
- Update Elastic Stack to 7.9.2

#### Fixed

- Updated small/docker setup script to properly parse kibana password
- Updated SSL configuration paths in example Dynamicbeat config
- Fix template syntax error in `http-kolide` example check

### Dynamicbeat

#### Added

- Report template failure errors in check results
- Report ICMP packet statistics in check result details for failed checks

#### Changed

- Swap github.com/sparrc/go-ping with github.com/go-ping/ping
- Update dependencies
- Update to libbeat 7.9.2

#### Fixed

- Re-add the check code that was accidentally removed in v0.5.1
- Don't panic on invalid templates
- Remove typo in ICMP definition struct field tag
- Don't ignore Count field in ICMP definition

### Kibana Plugin

#### Changed

- Migrate plugin to Kibana New Platform
- Update plugin to Kibana 7.9.2
- Rewrite plugin in TypeScript

#### Fixed

- Replace TinyURL plugin link with GitHub Releases link
- Build plugin bundles and include them in the plugin zipfile

## [0.5.1] - 2020-10-01

An intermediate release to support the transition of Dynamicbeat to go mod.

### Dynamicbeat

#### Changed

- Migrate to go mod

## [0.5.0] - 2020-09-29

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

## [0.4.0] - 2020-04-25

This release implements features for IRSeC 2020.

### General

#### Added

- Example SMB check
- Example MySQL check

### Dynamicbeat

#### Added

- SMB check
- MySQL check

## [0.3.0] - 2020-04-07

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

## [0.2.0] - 2020-02-28

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

## [0.1.0] - 2020-02-13

The initial release of Scorestack.

[unreleased]: https://github.com/scorestack/scorestack/compare/v0.8.2...main
[0.8.2]: https://github.com/scorestack/scorestack/compare/v0.8.1...v0.8.2
[0.8.1]: https://github.com/scorestack/scorestack/compare/v0.8.0...v0.8.1
[0.8.0]: https://github.com/scorestack/scorestack/compare/v0.7.0...v0.8.0
[0.7.0]: https://github.com/scorestack/scorestack/compare/v0.6.2...v0.7.0
[0.6.2]: https://github.com/scorestack/scorestack/compare/v0.6.1...v0.6.2
[0.6.1]: https://github.com/scorestack/scorestack/compare/v0.6.0...v0.6.1
[0.6.0]: https://github.com/scorestack/scorestack/compare/v0.5.1...v0.6.0
[0.5.1]: https://github.com/scorestack/scorestack/compare/v0.5.0...v0.5.1
[0.5.0]: https://github.com/scorestack/scorestack/compare/v0.4...v0.5.0
[0.4.0]: https://github.com/scorestack/scorestack/compare/v0.3...v0.4
[0.3.0]: https://github.com/scorestack/scorestack/compare/v0.2...v0.3
[0.2.0]: https://github.com/scorestack/scorestack/compare/v0.1...v0.2
[0.1.0]: https://github.com/scorestack/scorestack/releases/tag/v0.1.0
