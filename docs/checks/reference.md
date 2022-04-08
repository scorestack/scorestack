Check Reference
===============

This section of the documentation contains a reference to each of the check types available in Dynamicbeat. Each page within this section covers the available parameters for a specific check type. Example check definitions for all check types can be found under [the `examples` folder in the Scorestack repository](https://github.com/scorestack/scorestack/tree/main/examples).

Please note that the _Type_ listed in the tables on these pages refers to the type that must be used in the JSON document. For example, if the _type_ is _string_, then value for that parameter in the JSON document must be a `"string"`.

Required vs. Optional
---------------------

Some parameters are required for a check type. These _must_ be defined, otherwise the check definition is invalid and will not run.

Other parameters are not required, and have default values. If the default value for an optional parameter is acceptable for your check, you can omit the parameter from your check definition for the sake of brevity.