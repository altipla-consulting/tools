
# tools

Internal tools for developers in Altipla Consulting.


## Install

Install all tools with a single command:

```shell
curl -sL https://tools.altipla.consulting/install/all | sudo bash
```

Install only one of the tools (only for tools **with Support = External**):

```shell
curl -sL https://tools.altipla.consulting/install/jnet | sudo bash
```


## Tools

There are multiple tools inside this repo with different levels of activity and support and this table will give an overview of all of them.

| Tool | State | Support | Docs |
|------|-------|---------|------|
| `gaestage` | ![](https://img.shields.io/badge/state-active-brightgreen) | ![](https://img.shields.io/badge/usage-external-blue) | [Docs](./cmd/gaestage/README.md) |
| `impsort` | ![](https://img.shields.io/badge/state-active-brightgreen) | ![](https://img.shields.io/badge/usage-external-blue) | [Docs](./cmd/impsort/README.md) |
| `jnet` | ![](https://img.shields.io/badge/state-active-brightgreen) | ![](https://img.shields.io/badge/usage-external-blue) | [Docs](./cmd/jnet/README.md) |
| `previewer-netlify` | ![](https://img.shields.io/badge/state-deprecated-red) | | [Docs](./cmd/previewer/README.md) |
| `pub` | ![](https://img.shields.io/badge/state-deprecated-red) | | [Docs](./cmd/pub/README.md) |
| `runner` | ![](https://img.shields.io/badge/state-deprecated-red) | | [Docs](./cmd/runner/README.md) |

### Legend

| Badge | Meaning |
|-------|---------|
| ![](https://img.shields.io/badge/state-active-brightgreen) | Actively used. |
| ![](https://img.shields.io/badge/state-deprecated-red) | The tool is deprecated and being replaced or removed.<br>Do not use for new projects, it will be removed in the future. |
| ![](https://img.shields.io/badge/usage-external-blue) | Prepared to run anywhere outside our infrastructure.<br>Anyone can use it easily.<br>Breaking changes will be avoided as much as possible.<br>Any bugs found will be promptly fixed. |


## Contributing

You can make pull requests or create issues in GitHub. Any code you send should be formatted using `make gofmt`.


## License

[MIT License](LICENSE)
