package {{ .Package }}

import (
	"net/http"
	{{- range .Imports }}
	"{{ . }}"
	{{- end}}
	"github.com/scp10011/unalcohol"
	"github.com/getkin/kin-openapi/openapi3"
)

{{- range $key, $value := .Handler }}
func RegisterHandler{{$key}}(srv *unalcohol.Server,__api *{{ $value.Package }}.{{$key}}) {
	{{- range $path, $handler := $value.Path }}

	srv.Mux.HandleFunc(unalcohol.JoinPath(__api, "{{ $path }}"), func(writer http.ResponseWriter, request *http.Request) {
        defer func() {
            if v := recover(); v != nil {
                srv.StatusInternalServerError(request, writer)
            }
        }()
	    if err := __api.Middleware(request, writer); err != nil {
	        return
	    }
        if err := srv.Middleware(request, writer); err != nil {
            return
        }
		switch request.Method {
			{{- range $handler}}
			case "{{range $i, $v := .Description.Method }}{{if $i}}", "{{end}}{{$v}}{{end}}":
				{{- range .In}}
				v{{ .Key }} := {{ .Type }}
				if err := v{{ .Key }}.ParseRequest("{{.Key}}", request, writer); err != nil {
				    srv.StatusBadRequest(request, writer)
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
            srv.StatusMethodNotAllowed(request, writer)
			return
		}
	})
	{{- end}}
}
{{- end}}

{{- range $key, $value := .Handler }}
func RegisterDoc{{$key}}(paths *openapi3.Paths, __api *{{ $value.Package }}.{{$key}}) {
	{{- range $path, $handler := $value.Path }}
    paths.Set(unalcohol.JoinPath(__api, "{{ $path }}"), (func() *openapi3.PathItem {
		item := &openapi3.PathItem{}
		{{- range $h := $handler}}
		{{- range $h.Description.Method }}
            item.SetOperation("{{.}}",  (func() *openapi3.Operation {
                opt := openapi3.NewOperation()
                opt.Description = "{{ $h.Description.Description }}"
                opt.Summary = "{{ $h.Description.Summary }}"
                opt.Tags = []string{
                    {{- range $h.Description.Tags }}
                    "{{ . }}",
                    {{end }}
                }
                if len(opt.Tags) == 0 {
                    opt.Tags = __api.Tags
                }
                {{- range $in := $h.In}}
                v{{ $in.Key }} := {{ $in.Type }}
                if err := v{{ $in.Key }}.Doc("{{ $in.Key }}", opt); err != nil {
                    return nil
                }
                {{- end}}
                result := {{ $h.Result }}{}
                if err := result.Doc(opt); err != nil {
                    return nil
                }
                return opt
            })())
		{{- end}}
		{{- end}}
		return item
	})())
	{{- end}}
}
{{- end}}