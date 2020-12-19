Configuring Checks
==================

This section provides an overview of how to write Scorestack checks for your infrastructure.

Creating a check involves writing a JSON document which will configure the arguments for a check and define any necessary metadata information about the check. This document, which is called a **check definition** (or just **definition** for short) can optionally configure **attributes**. An attribute is a named variable with a string value that will be templated into the check at runtime by Dynamicbeat.

Attribute values can be changed during a competition by users and administrators through the Kibana UI. They are mainly useful for values that are likely to change frequently, such as service account credentials.

Once a definition document has been written, any attributes referenced in the definition must be configured with a permission level and default values.

As you read, it may be useful to consult the example check definitions in the [`examples/` folder of the repository](https://github.com/scorestack/scorestack/tree/stable/examples).