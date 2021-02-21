MSSQL
=====

| Name         | Type   | Required     | Description                                                          |
| ------------ | ------ | ------------ | -------------------------------------------------------------------- |
| Host         | String | Y            | IP or FQDN for the MSSQL server                                      |
| Username     | String | Y            | Username for the database                                            |
| Password     | String | Y            | Password for the user                                                |
| Database     | String | Y            | Name of the database to access                                       |
| Table        | String | Y            | Name of the table to access                                          |
| Column       | String | Y            | Name of the column to access                                         |
| MatchContent | String | N :: "false" | Whether to perform a regex content match on the results of the query |
| ContentRegex | String | N :: "\.\*"  | Regex to match on                                                    |
| Port         | String | N :: "1433"  | Port for the server                                                  |

Example Check
-------------

The example check provided in the `examples/` folder uses the following schema:

```sql
CREATE DATABASE scorestack
USE scorestack
CREATE TABLE stacks (testVal NVARCHAR(50))
INSERT INTO stacks ("hello from scorestack!")
GO
```
