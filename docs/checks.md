# Configuring Checks

This document details the many attributes and nuances of configuring checks for Scorestack. Creating a check involves making a JSON document which will contain metadata information for the check as well as the specific attributes needed for a successful check. If you would like to speedrun writing checks please see the example check definitions in the _examples_ folder.

## Adding Checks

In order to configure and add checks to Scorestack, you will have to create a folder structure containing the JSON definitions for each check. 
The following is the folder structure expected by `add-team.sh`.

```
myChecks
├── dns-host1
│   ├── admin-attribs.json
│   └── check.json
├── http-host2
    ├── admin-attribs.json
    ├── check.json
    └── user-attribs.json
```

The top level directory (_myChecks_) contains subfolders for all of the desired checks (_dns-host1, http-host2, etc._). Each of these subfolders will contain up to three JSON files that will define the specifc check. The three JSON files are as follows:

- `check.json`          --> The main JSON document that defines everything about a check
- `admin-attribs.json`  --> The attributes of a check that only an Admin will be allowed to change during competition through Kibana
- `user-attribs.json`   --> The attributes of a check that Users and Admins will be allowed to change during competition through Kibana

### check.json

This file contains the meat and potatoes for each check definition. Every `check.json` file **must** contain the following attributes:

| Name          | Type        | Required | Description                                                                                              |
| ------------- | ----------- | -------- | -------------------------------------------------------------------------------------------------------- |
| id            | String      | Y        | A unique identifier for the check. The `group` attribute will be appended to this value in Elasticsearch |
| name          | String      | Y        | This is the name of the check that will be displayed on the scoreboard                                   |
| type          | String      | Y        | The type of check (dns, ftp, http, etc.)                                                                 |
| group         | String      | Y        | The team associated with this check                                                                      |
| score\_weight | Int         | Y        | This is the number of points awarded for a successful check                                              |
| definition    | JSON Object | Y        | This contains the check specific attributes                                                              |

See the _examples_ folder for example `check.json` files.


### admin-attribs.json

This file will contain the attributes that Administrators will be able to modify through the Scorestack Kibana plugin during a competition. This allows dynamic updates to scoring such as changing an HTTP check to HTTPS after an inject to migrate a webserver to HTTPS. This is also useful for troubleshooting both during setup and during the competition. It does this my templating the JSON values in `admin-attribs.json` into `check.json`. Below are two example `check.json` files. One does not have Admin attributes and the other does.

`check.json` without Admin attributes
```json
{
    "id": "icmp",
    "name": "ICMP",
    "type": "icmp",
    "group": "example",
    "score_weight": 1,
    "definition": {
        "host": "127.0.0.1"
    }
}
```

`check.json` with Admin attributes
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

`admin-attribs.json` for the above example
```json
{
    "Host": "localhost"
}
```

As you can see, the values in `admin-attribs.json` will be used to fill in the template inside `check.json`. This can be done with any check specific attribtue (see all [check specific attributes](#check-specific-attributes)).


### user-attribs.json

Similar to `admin-attribs.json`, `user-attribs.json` will allow both Users and Admins to change attributes during a competition. The most common User attribute will be the password attribute so that Users can update passwords for scored users themselves. The use of `user-attribs.json` file is the same as the `admin-attribs.json` example above. The only difference comes in the permissions associated with the attributes in Elastic and the Kibana app.


## Adding Teams

Now that we know the folder structure for our checks as well as what attributes we want to make available to Admins and Users, we can now add a team and checks to Elastic and start scoring! To add a team with the checks you have configured, you will use the `add-team.sh` script. Before we run the script you will need to configure some values. Open the `add-team.sh` script in your editor of choice and modify the following values at the top of the file:

```
ELASTICSEARCH_HOST=localhost:9200
KIBANA_HOST=localhost:5601
CHECK_FOLDER=examples
USERNAME=root
PASSWORD=changeme
```

- `ELASTICSEARCH_HOST` --> Address and port of the Elastic host
- `KIBANA_HOST` --> Address and port of the Kibana host
- `CHECK_FOLDER` --> The top level folder containing all checks
- `USERNAME` --> Default username for Admin user
- `PASSWORD` --> Default password for Admin user

Once these values are correct for your environment you can run `add-team.sh`.

`add-team.sh TEAM_NAME`

Where `TEAM_NAME` is the name you want the team to be. The script will go through and setup a User account, dashboard, and add the checks for the team you passed to the script.

## Adding Multiple Teams at Once

Up until this point we are now ready to add our checks to Elastic and start scoring! But there is a problem. We can only add one team at a time. What if we want to add multiple teams at once without having to change all of the JSON documents to work for each team? The answer: templating!

`add-team.sh` allows you to pass multiple teams to it and it will configure the checks according to each team. This way you only have to write JSON once for the entire set of teams.

Let's go through an example. We want to add 5 teams to Scorestack for a competition. To do this we must configure our checks in a slightly different way. When specifying the `group` attribute you can use `${TEAM}` instead of the literal name of the team. The value `${TEAM}` will be replaced by `add-team.sh` with each team name you pass in as it adds teams and checks to Elastic. For example, say we have the following check definition:

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

This will only add an ICMP check for the `example` team. To make it more generic we can do this:

```json
{
    "id": "icmp",
    "name": "ICMP",
    "type": "icmp",
    "group": "${TEAM}",
    "score_weight": 1,
    "definition": {
        "host": "{{.Host}}"
    }
}
```

Now the value `${TEAM}` will be replaced with the team names you pass to the `add-team.sh` script. The `add-team.sh` also allows you to use the `${TEAM_NUM}` template. This can be used for templating IP addresses for all teams when configuring checks. For example, generally the IP address for a blue team would be in the format `192.168.X.0/24` where `X` is the team number. By using the `${TEAM_NUM}` template you can configure an IP like: `192.168.${TEAM_NUM}.0/24`

Let's look at an example for a check definition. Consider the following DNS check definition and it's associated Admin attributes:

`check.json`
```json
{
    "id": "dns",
    "name": "DNS",
    "type": "dns",
    "group": "example",
    "score_weight": 1,
    "definition": {
        "Server": "{{.Server}}",
        "Port": "{{.Port}}",
        "Fqdn": "{{.Fqdn}}",
        "ExpectedIP": "{{.ExpectedIP}}"
    }
}
```

`admin-attribs.json`
```json
{
    "Server": "192.168.1.10",
    "Port": "53",
    "Fqdn": "host1.com",
    "ExpectedIP": "192.168.1.98"
}
```

This DNS check will query Team1's DNS server `192.168.1.10` for the host `host1.com` and the expected IP for that host is `192.168.1.98`. We can template these values to work for all teams so that we only have to configure the JSON one time! The following is how we would do the templating:

`check.json` with Team Number template
```json
{
    "id": "dns",
    "name": "DNS",
    "type": "dns",
    "group": "${TEAM}",
    "score_weight": 1,
    "definition": {
        "Server": "{{.Server}}",
        "Port": "{{.Port}}",
        "Fqdn": "{{.Fqdn}}",
        "ExpectedIP": "{{.ExpectedIP}}"
    }
}
```

`admin-attribs.json` with Team Number template
```json
{
    "Server": "192.168.${TEAM_NUM}.10",
    "Port": "53",
    "Fqdn": "host${TEAM_NUM}.com",
    "ExpectedIP": "192.168.${TEAM_NUM}.98"
}
```

Now this check definition will work for all teams that are passed to `add-team.sh`. To take advantage of this templating feature you can add multiple teams like so:

`add-team.sh blue_1 blue_2 blue_3 ...`

The team name needs to be in the format of `blue_1` or `blue1` in order to take advantage of both the `${TEAM}` and `${TEAM_NUM}` templates.

## Updating All Check Definitions

In order to update all check definitions at once follow the steps below:

- Stop _Dynamicbeat_
- Delete all indexes using Kibana developer console
  - `DELETE /check*`
  - `DELETE /attrib*`
  - `DELETE /result*`
- Rerun `add-team.sh` with the same teams you had originally


## Check Specific Attributes

**Note:** The _Type_ listed in the tables below refers to the type that must be used in the JSON document. For example, if the _type_ is _string_ then the JSON document must have that attribute value as a `"string"`.

### DNS

| Name       | Type   | Required  | Description                                    |
| ---------- | ------ | --------- | ---------------------------------------------- |
| Server     | String | Y         | IP of the DNS server to query                  |
| Fqdn       | String | Y         | The FQDN of the host you are looking up        |
| ExpectedIP | String | Y         | The expected IP of the host you are looking up |
| Port       | String | N :: "53" | The port of the DNS server                     |


### FTP

| Name             | Type   | Required     | Description                                               |
| ---------------- | ------ | ------------ | --------------------------------------------------------- |
| Host             | String | Y            | IP or hostname of the host to run the FTP check against   |
| Username         | String | Y            | The user to login with over FTP                           |
| Password         | String | Y            | The password for the user that you wish to login with     |
| File             | String | Y            | The path to the file to access during the FTP check       |
| ContentRegex     | String | N :: "\.\*"  | Regex to match if reading file contents                   |
| HashContentMatch | String | N :: "false" | Whether or not to match a hash of the file contents       |
| Hash             | String | N            | The hash digest from from sha3\-256 to compare the hashed |
| Port             | String | N :: "21"    | The port to attempt an FTP connection on                  |
| Simple           | String | N :: "false" | Very simple FTP check for older servers                   |


### HTTP

| Name                 | Type                 | Required     | Description                                                       |
| -------------------- | -------------------- | ------------ | ----------------------------------------------------------------- |
| Verify               | String               | N :: "false" | Whether HTTPS certs should be validated                           |
| ReportMatchedContent | String               | N :: "false" | Whether the matched content should be returned in the CheckResult |
| Requests             | \[\]list of requests | Y            | A list of requests to make                                        |

Below are the **Request** Attribs

| Name         | Type                    | Required    | Description                                                                |
| ------------ | ----------------------- | ----------- | -------------------------------------------------------------------------- |
| Host         | String                  | Y           | IP or FQDN of the HTTP server                                              |
| Path         | String                  | Y           | Path to request \- see RFC3986, section 3\.3                               |
| HTTPS        | Bool                    | N :: false  | Whether or not HTTPS should be used                                        |
| Port         | UInt16                  | N :: 80     | TCP port number the HTTP server is listening on                            |
| Method       | String                  | N :: "GET"  | HTTP method to use                                                         |
| Headers      | map\[string\]\[string\] | N           | Name\-Value pairs of header fields to add/override                         |
| Body         | String                  | N           | The request body                                                           |
| MatchCode    | Bool                    | N :: false  | Whether the response code must match a defined value for the check to pass |
| Code         | Int                     | N :: 200    | The response status code to match                                          |
| MatchContent | Bool                    | N :: false  | Whether the response body must match a defined regex for the check to pass |
| ContentRegex | String                  | N :: "\.\*" | Regex for the response body to match                                       |
| StoreValue   | Bool                    | N :: false  | Whether the matched content should be saved for use in a later request     |

An HTTP definition consists of as many _Requests_ as you would like to send for that check. See the _examples_ folder for clarification.

### ICMP

| Name            | Type   | Required     | Description                                                                               |
| --------------- | ------ | ------------ | ----------------------------------------------------------------------------------------- |
| Host            | String | Y            | IP or FQDN of the host to run the ICMP check against                                      |
| Count           | Int    | N :: 1       | The number of ICMP requests to send per check                                             |
| AllowPacketLoss | String | N :: "false" | Pass check based on received pings matching Count; if false, will use percent packet loss |
| Percent         | Int    | N :: 100     | Percent of packets needed to come back to pass the check                                  |


## IMAP

| Name      | Type   | Required     | Description                         |
| --------- | ------ | ------------ | ----------------------------------- |
| Host      | String | Y            | IP or FQDN for the IMAP server      |
| Username  | String | Y            | Username for the IMAP server        |
| Password  | String | Y            | Password for the user               |
| Encrypted | String | N :: "false" | Whether or not to use TLS \(IMAPS\) |
| Port      | String | N :: "143"   | Port for the IMAP server            |


## LDAP

| Name     | Type   | Required     | Description                            |
| -------- | ------ | ------------ | -------------------------------------- |
| User     | String | Y            | The user written in user@domain syntax |
| Password | String | Y            | The password for the user              |
| Fqdn     | String | Y            | The FQDN of the LDAP server            |
| Ldaps    | String | N :: "false" | Whether or not to use LDAP\+TLS        |
| Port     | String | N :: "389"   | Port for LDAP server                   |


## MySQL

| Name         | Type   | Required     | Description                                                          |
| ------------ | ------ | ------------ | -------------------------------------------------------------------- |
| Host         | String | Y            | IP or FQDN for the MySQL server                                      |
| Username     | String | Y            | Username for the database                                            |
| Password     | String | Y            | Password for the user                                                |
| Database     | String | Y            | Name of the database to access                                       |
| Table        | String | Y            | Name of the table to access                                          |
| Column       | String | Y            | Name of the column to access                                         |
| MatchContent | String | N :: "false" | Whether to perform a regex content match on the results of the query |
| ContentRegex | String | N :: "\.\*"  | Regex to match on                                                    |
| Port         | String | N :: "3306"  | Port for the server                                                  |


## SMB

| Name         | Type   | Required    | Description                                                                       |
| ------------ | ------ | ----------- | --------------------------------------------------------------------------------- |
| Host         | String | Y           | IP or FQDN for the SMB server                                                     |
| Username     | String | Y           | Username for SMB share                                                            |
| Password     | String | Y           | Password for the user                                                             |
| Share        | String | Y           | Name of the SMB share                                                             |
| Domain       | String | Y           | The domain found in front of a login \(SMB\\Administrator : SMB would be domain\) |
| File         | String | Y           | The file in SMB share                                                             |
| ContentRegex | String | N :: "\.\*" | Regex to match on                                                                 |
| Port         | String | N :: "445"  | Port of the server                                                                |


## SMTP

| Name      | Type   | Required                     | Description                   |
| --------- | ------ | ---------------------------- | ----------------------------- |
| Host      | String | Y                            | IP or FQDN of the SMTP server |
| Username  | String | Y                            | Username for the SMTP server  |
| Password  | String | Y                            | Password for the SMTP server  |
| Sender    | String | Y                            | Who is sending the email      |
| Reciever  | String | Y                            | Who is receiving the email    |
| Body      | String | N :: "Hello from Scorestack" | Body of the email             |
| Encrypted | String | N :: False                   | Whether or not to use TLS     |
| Port      | String | N :: "25"                    | Port of the SMTP server       |


## SSH

| Name         | Type   | Required     | Description                                            |
| ------------ | ------ | ------------ | ------------------------------------------------------ |
| Host         | String | Y            | IP or FQDN of the host to run the SSH check against    |
| Username     | String | Y            | The user to login with over SSH                        |
| Password     | String | Y            | The password for the user that you wish to login with  |
| Cmd          | String | Y            | The command to execute once SSH connection established |
| MatchContent | String | N :: "false" | Whether or not to match content like checking files    |
| ContentRegex | String | N :: "\.\*"  | Tegex to match if reading a file                       |
| Port         | String | N :: "22"    | The port to attempt an SSH connection on               |


## VNC 

| Name     | Type   | Required | Description                                           |
| -------- | ------ | -------- | ----------------------------------------------------- |
| Host     | String | Y        | IP or FQDN of the host to run the SSH check against   |
| Port     | String | Y        | The port for the VNC server                           |
| Password | String | Y        | The password for the user that you wish to login with |


## WinRM

| Name         | Type   | Required     | Description                                              |
| ------------ | ------ | ------------ | -------------------------------------------------------- |
| Host         | String | Y            | IP or FQDN of the WinRM machine                          |
| Username     | String | Y            | User to login as                                         |
| Password     | String | Y            | Password for the user                                    |
| Cmd          | String | Y            | Command that will be executed                            |
| Encrypted    | String | N :: "true"  | Use TLS for connection                                   |
| MatchContent | String | N :: "false" | Turn this on to match content from the output of the cmd |
| ContentRegex | String | N :: "\.\*"  | Regexp for matching output of a command                  |
| Port         | String | N :: "5986"  | Port for WinRM                                           |


## XMPP

| Name      | Type   | Required    | Description                         |
| --------- | ------ | ----------- | ----------------------------------- |
| Host      | String | Y           | IP or FQDN of the XMPP Server       |
| Username  | String | Y           | Username to use for the XMPP server |
| Password  | String | Y           | Password for the user               |
| Encrypted | String | N :: "true" | Whether or not to use TLS           |
| Port      | String | N :: "5222" | The port for the XMPP server        |

## NOOP

| Name    | Type   | Required | Description                                                 |
| ------- | ------ | -------- | ----------------------------------------------------------- |
| Dynamic | String | Y        | Contains attributes that can be modified by admins or users |
| Static  | String | Y        | Contains attributes                                         |