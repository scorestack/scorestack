Cloning the Repository
======================

Whenever you need to deploy Scorestack, you need to make sure you're cloning the right branch. In most cases, cloning directly from the `stable` branch (which is the default) is fine.

```shell
git clone --branch stable https://github.com/scorestack/scorestack.git
```

However, if you would like to deploy a specific version of Scorestack, just pass the tag of the version you want to deploy to the `--branch` argument of `git clone`. For example, to deploy Scorestack version 0.5.0, you would run `git clone --branch v0.5.0 https://github.com/scorestack/scorestack.git`.

If you just want to use the latest release, the `stable` branch is fine - it always points at the latest stable release.

If you want to live life on the edge and try out an unstable release or unreleased changes, try cloning the `dev` branch instead! Just note that it `dev` is the unstable development branch, so there's no guarantees that things will work for you.