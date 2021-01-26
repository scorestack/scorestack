Obtaining Binaries
==================

In order to run Dynamicbeat and Scorestack's Kibana plugin, the binaries for these components must be obtained. The Dynamicbeat binary is a compiled Golang executable that runs Dynamicbeat. The Kibana plugin binary is a zipfile containing the compiled assets that are installed into the Kibana server and loaded at runtime.

Prebuilt Binaries
-----------------

Most users will want to use the prebuilt binaries that are available on the [Scorestack Releases page](https://github.com/scorestack/scorestack/releases). The Kibana plugin zipfile and a zipped Dynamicbeat executable are attached to each release. This is the recommended way of obtaining the binaries.

Building Your Own Binaries
--------------------------

If you really want to, you can build these binaries yourself. Please see the [documentation on building](../building.md) for more information on how to do this.