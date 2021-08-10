HTTP
====

| Name                 | Type                 | Required     | Description                                                       |
| -------------------- | -------------------- | ------------ | ----------------------------------------------------------------- |
| Verify               | String               | N :: "false" | Whether HTTPS certs should be validated                           |
| ReportMatchedContent | String               | N :: "false" | Whether the matched content should be returned in the CheckResult |
| Requests             | \[\]list of requests | Y            | A list of requests to make                                        |

Below are the parameters found within a single **request**.

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

An HTTP definition consists of as many _Requests_ as you would like to send for that check. See the _examples_ folder for clarification.

`StoreValue` Parameter
----------------------

When the `StoreValue` attribute is set to `true` and regex-based content matching is enabled, then the content in the response that matches the `ContentRegex` will be stored in the `SavedValue` variable. You can then use this value in later requests. This can be useful for multi-stage checks that require authentication, but do not use cookies to store the session ID.

One example of using the `StoreValue` attribute is the `http-kolide` example check. Before you can use the Kolide API, you must log in. The API route to log in returns a Bearer token within the response body. This Bearer token must be presented in the `Bearer:` header in order to authenticate to the API routes.

The saved value is made available through the same method as attributes - just insert `{{.SavedValue}}` into your check wherever you would like it to be used.

Please note that only one value can be stored using `StoreValue`. If you already have a value saved, and attempt to save another one, then the original value will be overwritten.