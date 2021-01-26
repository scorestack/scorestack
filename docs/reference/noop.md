Noop
====

| Name    | Type   | Required | Description                                                 |
| ------- | ------ | -------- | ----------------------------------------------------------- |
| Dynamic | String | Y        | Contains attributes that can be modified by admins or users |
| Static  | String | Y        | Contains attributes                                         |

The Noop check does basically nothing. Dynamicbeat will load the check definition and attributes, template them into the check, pass the check, and then return the templated check definition in the check result's details field.

Noop can be useful to verify that your Dynamicbeat and Scorestack configurations are correct. It can also be used to play around with the templating and attributes systems.