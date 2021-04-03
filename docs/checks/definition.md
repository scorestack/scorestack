Check Definition
================

The check definition contains the parameters that define how the check will be executed. For most checks it will just be a simple key-value object, but it depends on which check type is being defined. Please see [the check reference](./reference.md) for more information on what parameters are expected within the check definition for each check type.

Here is the definition section of our example IMAP check:

```json
{
  "definition": {
    "Host": "{{.Host}}",
    "Port": "143",
    "Username": "{{.Username}}@example.com",
    "Password": "{{.Password}}",
  },
}
```

Optional Parameters
-------------------

Not all parameters must be included in your check definition. Some parameters may be optional, with a default value configured. If you omit the parameter from the definition, then the default value will be used for the parameter.

In our example IMAP check, the `Encrypted` parameter has been omitted. This is because the server we're checking doesn't use TLS, and the default for IMAP's `Encrypted` parameter is `false`.

Attribute Templating
--------------------

Three of our parameters - `Host`, `Username`, and `Password` - have attributes defined for them. These attributes are inserted into the parameter values using basic [golang text templating](https://golang.org/pkg/text/template/). Two pairs of curly braces containing a dot followed by a [PascalCase](https://techterms.com/definition/pascalcase) name will be replaced with the value of the referenced attribute.

For example, if you define an attribute named `MyAttrib`, you can reference it in the check definition with `{{.MyAttrib}}`.

Attributes can be used as the entire value of a parameter, or they can be inserted into a string to make up part of the value.

### Entire Value

The `Host` parameter in our example check is an example of using the attribute (also named `Host`) as the entire value. Before templating, it looks like this:

```json
"Host": "{{.Host}}"
```

After templating, if the `Host` attribute is set to `localhost`, it will look like this:

```json
"Host": "localhost"
```

### Part of a String

The `Username` parameter in the example check is an example of using an attribute (again, also named `Username`) as part of the parameter's value. For our check, the username is an email address for a known domain. Before templating, it looks like this:

```json
"Username": "{{.Username}}@example.com"
```

After templating, if the `Username` attribute is set to `scorestack`, it will look like this:

```json
"Username": "scorestack@example.com"
```