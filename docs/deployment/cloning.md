Cloning the Repository
======================

If you're just testing out Scorestack, cloning the default branch (`main`) should be fine in almost all cases:
```shell
git clone https://github.com/scorestack/scorestack.git
```

However, the `main` branch is the primary branch used for development. While we do our best to never merge broken code into `main`, there's no guarantees that what's on `main` will always work as expected.

Therefore, whenever you deploy Scorestack in a production environment (such as for active competitions) you should always use the `--branch` argument of `git clone` to specify the exact version you want to clone. For example, this command will clone [Scorestack 0.8.2](https://github.com/scorestack/scorestack/releases/tag/v0.8.2):
```shell
git clone --branch v0.8.2 https://github.com/scorestack/scorestack.git
```
