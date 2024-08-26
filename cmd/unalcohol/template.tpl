package main

import (
	"net/http"{{range .Imports }}
	"{{ . }}"
	{{- end}}
	"unalcohol/pkg/unalcohol"
)

{{- range $key, $value := .Handler }}
func RegisterHandler{{$key}}(srv *unalcohol.Server,__api *{{ $value.Package }}.{{$key}}) {

	{{- range $path, $handler := $value.Path }}
	srv.Mux.HandleFunc(__api.Path + "{{ $path }}", func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
			{{- range $handler}}
			case "{{.Method}}":
				{{- range .In}}
				v{{ .Key }} := {{ .Type }}
				if err := v{{ .Key }}.ParseRequest("{{.Key}}", request, writer); err != nil {
					writer.WriteHeader(http.StatusBadRequest)
					return
				}
				{{- end}}
				result := __api.{{ .Name }}(
					{{- range .In }}
					v{{ .Key }},
					{{- end}}
				)
				result.WriteResponse(writer)
			{{- end}}
		default:
			writer.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	})
	{{- end}}
}
{{- end}}