`check.json`
============

This file contains the meat and potatoes for each check definition. Every `check.json` file **must** contain the following parameters:

| Name          | Type        | Description                                                                                                                                                                          |
| ------------- | ----------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| id            | String      | A unique identifier for the check. The `group` attribute will be appended to this value in Elasticsearch. The value of this parameter _must_ be the same as the check's folder name. |
| name          | String      | This is the name of the check that will be displayed on the scoreboard.                                                                                                              |
| type          | String      | The type of check (dns, ftp, http, etc.)                                                                                                                                             |
| group         | String      | The team associated with this check.                                                                                                                                                 |
| score\_weight | Int         | This is the number of points awarded for a successful check.                                                                                                                         |
| definition    | JSON Object | This contains parameters that are specific to the kind of check being defined.                                                                                                       |

The content of the `definition` parameter depends on what kind of check is being defined. Please see [the check reference](../reference.md) for information on what parameters are expected within the `definition` parameter for each check type.

Example
-------

This is a minimal ICMP check definition that doesn't define any attributes.

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