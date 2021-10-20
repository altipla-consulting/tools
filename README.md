
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
| `ci` | Actively used | Supported | |
| `gaestage` | Actively used | Supported | |
| `gendc` | Actively used | Unsupported.<br>Very opinionated | |
| `impsort` | Actively used | Supported | |
| `jnet` | Actively used | Supported | |
| `linter` | Actively used | Unsupported.<br>Very opinionated | |
| `previewer-netlify` | Deprecated.<br>Use `wave` instead. | Unsupported | |
| `pub` | Deprecated.<br>Use `wave` instead. | Unsupported | |
| `reloader` | Actively used | Supported | |
| `runner` | Deprecated.<br>Use `wave` instead. | Unsupported | |
| `wave` | Actively used. | Unsupported | |


## Contributing

You can make pull requests or create issues in GitHub. Any code you send should be formatted using `make gofmt`.


## License

[MIT License](LICENSE)
