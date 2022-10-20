# Mleh

Like Helm but works for any kind of file, not just YAMLs.


## Usage

Similar to Helm, with the following directory structure:

```
chart
├── values.yaml
├── templates
│   ├── app.tf
│   └── _helpers.txt
└── values.yaml
```

You can run:

```
mleh -values values.yaml -output-dir . chart/
```

And it will generate, in this case an app.tf file, while skipping files that
start with an underscore.

## Template functions

As Helm, [Sprig](https://github.com/Masterminds/sprig) functions are available,
the documentation is hosted by
[Mastermind](https://masterminds.github.io/sprig/).

### Differences with Helm

Functions `toYaml`, `toJson`, `tpl`, are not implemented.
