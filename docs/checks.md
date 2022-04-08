Configuring Checks
==================

This section provides an overview of how to write Scorestack checks for your infrastructure.

In Scorestack, a "check" is a method of verifying that a specific service is functioning as expected. A check could make a request to a webserver and expect the response to contain specific content. Another check could establish an SSH session as a specific user, run some commands, and verify that the commands executed properly.

Each check is configured by - you guessed it - a check **configuration**. This configuration contains three major parts: the **metadata**, the **definition**, and the **attributes**. All check configurations _must_ include the required metadata and a definition, and may optionally include some attributes as necessary.

A check's metadata uniquely identifies it and tells Scorestack how it should be handled. Check definitions set parameters that define the check's execution, and are specific to the type of check being configured. Check attributes are values that are substituted into the check at runtime, and can be viewed or modified by administrators and team members in Kibana.

The rest of this section explains the parts of a check in greater detail, as well as how to write a check configuration file, how to add your checks to Kibana, and how to use the team override system to write check configurations that can be generic across all of your teams.

As you read the documentation and write your check files, it may be useful to consult the example check definitions in the [`examples/` folder of the repository](https://github.com/scorestack/scorestack/tree/main/examples).