package main

import "html/template"

var tmplCaddyfile = template.Must(template.New("caddyfile").Parse(`
{{range .Apps}}
{{.Name}}.dev.localhost {
	@grpc protocol grpc
	handle @grpc {
		reverse_proxy {
			to {{.Name}}:8080
			transport http {
				versions h2c
			}
		}
	}
	reverse_proxy {{.Name}}:8080
	tls /opt/tls/cert.pem /opt/tls/key.pem
}
{{end}}

{{range .JS}}
{{.Name}}.dev.localhost {
	reverse_proxy {{.Name}}:3000
	tls /opt/tls/cert.pem /opt/tls/key.pem
}
{{end}}
`))
