# Configuring Checks
====================

This document details the many attributes and nuances of configuring checks for ScoreStack. Creating a check involves making a JSON document which will contain metadata information for the check as well as the specific attributes needed for a successful check. If you would like to speedrun writing checks please see the example check definitions in the _examples_ folder.


**Note:** The _Type_ listed in the tables below refers to the type that must be used in the JSON document. For example, if the _type_ is _string_ then the JSON document must have that attribute value as a `"string"`.


### Uniform Check Attributes

| Name          | Type        | Required | Description |
| ------------- | ----------- | -------- | ----------- |
| id            | String      | Y        |             |
| name          | String      | Y        |             |
| type          | String      | Y        |             |
| group         | String      | Y        |             |
| score\_weight | Int         | Y        |             |
| def           | JSON Object | Y        |             |

The table above 


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
| Body      | String | N :: "Hello from ScoreStack" | Body of the email             |
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