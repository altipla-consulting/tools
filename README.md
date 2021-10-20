
# tools

Internal tools for developers in Altipla Consulting.


## Install

```shell
curl https://tools.altipla.consulting/install/tools | sudo bash
```

## Tools

There are multiple tools inside this repo with different levels of activity and support. The support columns indicates if the tool is prepared to be run outside our infrastructure and we actively fix any bug that may occur in that external scenarios. Internal tools are provided as-is without any kind of support.

| Tool | State | Support | Docs |
|------|-------|---------|------|
| `ci` | Actively used | Supported | [Docs](./ci/README.md) |
| `gaestage` | Actively used | Supported | [Docs](./gaestage/README.md) |
| `gendc` | Actively used | Unsupported.<br>Very opinionated | [Docs](./gendc/README.md) |
| `impsort` | Actively used | Supported | [Docs](./impsort/README.md) |
| `jnet` | Actively used | Supported | [Docs](./jnet/README.md) |
| `linter` | Actively used | Unsupported.<br>Very opinionated | |
| `previewer-netlify` | Deprecated.<br>Use `wave` instead. | Unsupported | [Docs](./previewer/README.md) |
| `pub` | Deprecated.<br>Use `wave` instead. | Unsupported | [Docs](./pub/README.md) |
| `reloader` | Actively used | Supported | [Docs](./reloader/README.md) |
| `runner` | Deprecated.<br>Use `wave` instead. | Unsupported | [Docs](./runner/README.md) |
| `wave` | Actively used. | Unsupported | [Docs](./wave/README.md) |


## Contributing

You can make pull requests or create issues in GitHub. Any code you send should be formatted using `make gofmt`.


## License

[MIT License](LICENSE)
