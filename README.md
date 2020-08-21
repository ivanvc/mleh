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

Mleh uses the default template definitions from
[Go](https://golang.org/pkg/text/template/), so when requiring a template,
rather than using Helm's `include` just use Go's native implementation.

```
{{- define "app.name" -}}
{{- .Values | trunc 16 | trimSuffix "-" -}}
{{- end -}}

// Later, just call it with:
app.name = {{ template "app.name" . }}
```

Functions `toYaml`, `toJson`, `tpl`, are neither implemented.
