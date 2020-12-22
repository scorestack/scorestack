Adding Teams
============

Once you've written your [check definitions](./check_json.md) and [defined their attributes](./attributes.md), you need to add the checks to Scorestack so they can start running. We have written an `add-team.sh` script that will parse the check folder and file structure and add the checks to Scorestack one team at a time.

> Please note that the `add-team.sh` script assumes that all checks will be run by all of the teams. If you have some checks that will only be run by specific teams, you will have to add them manually. If this is a use case that you'd like to see implemented, [open an issue](https://github.com/scorestack/scorestack/issues/new/choose) and we'd be happy to take a look.

The `add-team.sh` script performs the following actions:

1. Add the all-team scoreboard if it does not already exist
2. Configure the default index template for Scorestack indices
3. For each team name passed as an argument:
   1. Add the checks for the team
   2. Add a role for the team
   3. Add the team user
   4. Add the team's dashboard

These actions can all be performed manually by interacting with the Elasticsearch API, but the `add-team.sh` script abstracts over those operations for you.

> Please note the `add-team.sh` script depends on the [`jq`](https://stedolan.github.io/jq/) tool. If this tool is not available on your system or you cannot install it, consider running the script within a docker container.

Before running the script, you will need to configure some values. Open the script in a text editor and modify the following values at the top of the file:

```
ELASTICSEARCH_HOST=localhost:9200
KIBANA_HOST=localhost:5601
CHECK_FOLDER=examples
USERNAME=root
PASSWORD=changeme
```

- `ELASTICSEARCH_HOST`: The address and port of the Elasticsearch server.
- `KIBANA_HOST`: The address and port of the Kibana server.
- `CHECK_FOLDER`: The top level folder containing all checks.
- `USERNAME`: The username of the Scorestack admin user. This should usually be left as `root`, unless you have added another user to the Elastic cluster with the `superuser` role.
- `PASSWORD`: The password for the user referenced in `USERNAME`.

Once these values are correct for your environment, you can run `add-team.sh`. Each argument after the script name is the name of a team you would like to add.

```shell
./add-team.sh team01 team02 team03
```

In this case, three teams will be added: `team01`, `team02`, and `team03`.

Adding Multiple Teams
---------------------

Thanks to some additional templating configured in the `add-team.sh` script, each check only has to be written once for multiple teams.

When adding checks for a team, the script will make the following modifications to a definition or attribute file:

- Append `-GROUP` to the `id` parameter, where `GROUP` is the name of the group
- Replace `${TEAM}` wherever it occurs with the name of the group
- Replace `${TEAM_NUM}` wherever it occurs with the group's number

### `${TEAM_NUM}`

When a team name is passed to the `add-team.sh` script, any numbers found at the end of the team name will be parsed as the team's number. For example, `team10` would have a team number of `10`.

Also, any leading zeros in the team number will be removed. For example, `team04` would have a team number of `4`, **not** `04`.

This can be useful for templating the team's number into a domain or IP address to customize the target of a check for each team.

### Example

Let's go through an example of adding multiple teams to Scorestack to see how we can use the features of the `add-team.sh` script to our advantage.

We want to add 5 teams to Scorestack for a competition. Let's use an example DNS check definition with some associated administrator attributes.

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
    "Fqdn": "www.team01.com",
    "ExpectedIP": "192.168.1.98"
}
```

This DNS check will query the DNS server at `192.168.1.10` for the host `host1.com`, and will expect the response to report the host's IP as `192.168.1.98`. We can use the `${TEAM}` and `${TEAM_NUM}` templates to make this check a bit more generic so it'll work for all of the teams! The following is how we would do the templating:

`check.json` with Team templates
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

`admin-attribs.json` with Team templates
```json
{
    "Server": "192.168.${TEAM_NUM}.10",
    "Port": "53",
    "Fqdn": "www.${TEAM}.com",
    "ExpectedIP": "192.168.${TEAM_NUM}.98"
}
```

Finally, let's see what the checks will look like once the script performs the initial templating and modifications. In this example, we'll assume the script was run with one argument - `team05`.

`check.json` modified for `team05`
```json
{
    "id": "dns-team05",
    "name": "DNS",
    "type": "dns",
    "group": "team05",
    "score_weight": 1,
    "definition": {
        "Server": "{{.Server}}",
        "Port": "{{.Port}}",
        "Fqdn": "{{.Fqdn}}",
        "ExpectedIP": "{{.ExpectedIP}}"
    }
}
```

`admin-attribs.json` modified for `team05`
```json
{
    "Server": "192.168.5.10",
    "Port": "53",
    "Fqdn": "www.team05.com",
    "ExpectedIP": "192.168.5.98"
}
```