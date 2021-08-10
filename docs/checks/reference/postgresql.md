PostgreSQL
==========

| Name         | Type   | Required     | Description                                                          |
| ------------ | ------ | ------------ | -------------------------------------------------------------------- |
| Host         | String | Y            | IP or FQDN for the PostgreSQL server                                 |
| Username     | String | Y            | Username for the database                                            |
| Password     | String | Y            | Password for the user                                                |
| Database     | String | Y            | Name of the database to access                                       |
| Table        | String | Y            | Name of the table to access                                          |
| Column       | String | Y            | Name of the column to access                                         |
| MatchContent | String | N :: "false" | Whether to perform a regex content match on the results of the query |
| ContentRegex | String | N :: "\.\*"  | Regex to match on                                                    |
| Port         | String | N :: "3306"  | Port for the server                                                  |

Example Check
-------------

Unlike MySQL, PostgreSQL does not ship by default with a queryable database. The example check definition provided in the `examples/` folder uses the following simple schema:

```sql
CREATE TABLE testtable(
  EntryID SERIAL PRIMARY KEY,
  Name TEXT NOT NULL
);
INSERT INTO testtable (Name) VALUES ('scorestack');
```