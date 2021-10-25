
# tools

Internal tools for developers in Altipla Consulting.


## Install

Install all tools in one go:

```shell
curl https://raw.githubusercontent.com/altipla-consulting/tools/master/install/all.sh | sudo bash
```

Install only one of the tools:

```shell
curl https://raw.githubusercontent.com/altipla-consulting/tools/master/install/jnet.sh | sudo bash
```


## Tools

There are multiple tools inside this repo with different levels of activity and support. The support columns indicates if the tool is prepared to be run outside our infrastructure and we actively fix any bug that may occur in that external scenarios. Internal tools are provided as-is without any kind of support.

| Tool | State | Support | Docs |
|------|-------|---------|------|
| `ci` | Actively used | Supported | [Docs](./cmd/ci/README.md) |
| `gaestage` | Actively used | Supported | [Docs](./cmd/gaestage/README.md) |
| `gendc` | Actively used | Unsupported.<br>Very opinionated | [Docs](./cmd/gendc/README.md) |
| `impsort` | Actively used | Supported | [Docs](./cmd/impsort/README.md) |
| `jnet` | Actively used | Supported | [Docs](./cmd/jnet/README.md) |
| `linter` | Actively used | Unsupported.<br>Very opinionated | |
| `previewer-netlify` | Deprecated.<br>Use `wave` instead. | Unsupported | [Docs](./cmd/previewer/README.md) |
| `pub` | Deprecated.<br>Use `wave` instead. | Unsupported | [Docs](./cmd/pub/README.md) |
| `reloader` | Actively used | Supported | [Docs](./cmd/reloader/README.md) |
| `runner` | Deprecated.<br>Use `wave` instead. | Unsupported | [Docs](./cmd/runner/README.md) |
| `wave` | Actively used. | Unsupported | [Docs](./cmd/wave/README.md) |


## Contributing

You can make pull requests or create issues in GitHub. Any code you send should be formatted using `make gofmt`.


## License

[MIT License](LICENSE)
