#!/usr/bin/env python3
import argparse
from string import Template


def main():
    # Parse arguments
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "protocol",
        help="The name of the protocol to generate a schema file for."
    )
    parser.add_argument(
        "-t",
        "--transport",
        required=True,
        choices=["TCP", "UDP"],
        help="The transport-layer protocol used.",
    )
    parser.add_argument(
        "-p",
        "--port",
        required=True,
        type=int,
        choices=range(65536),
        help="The default port for servers of this protocol.",
    )
    args = parser.parse_args()

    # Read the template file
    with open("template.json", "rt") as fp:
        tmpl = Template(fp.read())

    # Render and write the template
    result = tmpl.substitute({
        "protocol": args.protocol.upper(),
        "transport": args.transport,
        "port": args.port,
        "ref": "$ref",
    })
    with open(f"{args.protocol}.json", "wt") as fp:
        fp.write(result)


if __name__ == "__main__":
    main()
