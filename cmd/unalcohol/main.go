package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/scp10011/unalcohol/internal/doc"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"golang.org/x/mod/modfile"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed template.tpl
var stringMethodTemplate string

func getValueExpr(expr ast.Expr, packageName string) string {
	switch t := expr.(type) {
	case *ast.Ident:
		switch t.Name {
		case "int", "int8", "int16", "int32", "int64",
			"uint", "uint8", "uint16", "uint32", "uint64", "uintptr",
			"float32", "float64",
			"complex64", "complex128",
			"bool", "string":
			return t.Name
		default:
			if strings.Contains(t.Name, ".") {
				return t.Name
			} else if packageName != "" {
				return fmt.Sprintf("%s.%s", packageName, t.Name)
			} else {
				return t.Name
			}
		}
	case *ast.SelectorExpr:
		return getTypeExpr(t.X, packageName) + "." + t.Sel.Name // 处理有包名的类型，如 time.Duration
	case *ast.StarExpr:
		return getTypeExpr(t.X, packageName) // 指针类型
	case *ast.ArrayType:
		return "[]" + getTypeExpr(t.Elt, packageName) // 数组或切片类型
	case *ast.IndexExpr:
		return fmt.Sprintf("%s[%s]", getTypeExpr(t.X, packageName), getTypeExpr(t.Index, packageName))
	default:
		log.Printf("%T", expr)
		return "unknown"
	}
}

func getTypeExpr(expr ast.Expr, packageName string) string {
	switch t := expr.(type) {
	case *ast.Ident:
		switch t.Name {
		case "int", "int8", "int16", "int32", "int64",
			"uint", "uint8", "uint16", "uint32", "uint64", "uintptr",
			"float32", "float64",
			"complex64", "complex128",
			"bool", "string":
			return t.Name
		default:
			if strings.Contains(t.Name, ".") {
				return t.Name
			} else if packageName != "" {
				return fmt.Sprintf("%s.%s", packageName, t.Name)
			} else {
				return t.Name
			}
		}
	case *ast.SelectorExpr:
		return getTypeExpr(t.X, "") + "." + t.Sel.Name // 处理有包名的类型，如 time.Duration
	case *ast.StarExpr:
		return "&" + getTypeExpr(t.X, packageName) // 指针类型
	case *ast.ArrayType:
		return "[]" + getTypeExpr(t.Elt, packageName) // 数组或切片类型
	case *ast.IndexExpr:
		return fmt.Sprintf("%s[%s]", getTypeExpr(t.X, packageName), getTypeExpr(t.Index, packageName))
	default:
		log.Printf("%T", expr)
		return "unknown"
	}
}

type PreGen struct {
	Package string
	Handler map[string]PrePath
	Imports []string
}

type PrePath struct {
	Package string

	Path map[string][]PreHandler
}

type PreParam struct {
	Key  string
	Type string
}

type PreHandler struct {
	Description *doc.Doc
	Path        string
	Name        string
	Method      string
	Result      string
	In          []PreParam
}

var (
	root     = flag.String("root", ".", "package root path")
	endpoint = flag.String("entry", "main.go", "entry point file")
	handler  = flag.String("handler", ".", "handler dir")
)

func ParsePath(path string) string {
	if path == "." {
		path, _ = os.Getwd()
	}
	path, _ = filepath.Abs(path)
	return path
}

func main() {
	flag.Parse()
	rootPath := ParsePath(*root)
	handlerPath := ParsePath(*handler)
	mainFile := ParsePath(*endpoint)
	output, _ := filepath.Split(mainFile)
	output = filepath.Join(output, "unalcohol_gen.go")
	os.Chdir(rootPath)
	file, err := parser.ParseFile(token.NewFileSet(), mainFile, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}
	modBuffer, err := os.ReadFile(filepath.Join(rootPath, "go.mod"))
	if err != nil {
		log.Fatalln("Did not find the go.mod file.")
	}
	packageName := modfile.ModulePath(modBuffer)
	imports := make([]string, 0)
	data := make(map[string]PrePath)
	total := 0
	err = filepath.Walk(handlerPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf("walk: %+v", err)
		}
		if !info.IsDir() {
			return nil
		}
		fs := token.NewFileSet()
		pkgs, err := parser.ParseDir(fs, path, nil, parser.ParseComments)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error parsing directory:", err)
			os.Exit(1)
		}
		// 生成并输出代码
		for name, pkg := range pkgs {
			importPath := strings.Replace(path, handlerPath, packageName, 1)
			for _, file := range pkg.Files {
				handler := make(map[string]*doc.Doc)
				for _, group := range file.Comments {
					if description := doc.ParseDoc(group); description != nil {
						handler[description.Name] = description
					}
				}
				for _, decl := range file.Decls {
					funcDecl, ok := decl.(*ast.FuncDecl)
					if !ok {
						continue
					}
					if funcDecl.Type.Params == nil {
						continue
					}
					if funcDecl.Recv == nil {
						continue
					}
					ptr := getValueExpr(funcDecl.Recv.List[0].Type, "")
					block, ok := data[ptr]
					if !ok {
						block = PrePath{
							Package: name,
							Path:    make(map[string][]PreHandler),
						}
						data[ptr] = block
					}
					functionName := funcDecl.Name.String()
					description, ok := handler[functionName]
					if !ok {
						continue
					}
					request := PreHandler{Description: description, Path: description.URL, Name: functionName}
					for _, param := range funcDecl.Type.Params.List {
						for _, ident := range param.Names {
							t := getTypeExpr(param.Type, name)
							request.In = append(request.In, PreParam{
								Key:  ident.Name,
								Type: fmt.Sprintf("%s{}", t),
							})
						}
					}
					if len(funcDecl.Type.Results.List) != 1 {
						continue
					}
					result := funcDecl.Type.Results.List[0]
					request.Result = getTypeExpr(result.Type, name)
					block.Path[request.Path] = append(block.Path[request.Path], request)
				}
			}
			if len(data) != total {
				imports = append(imports, importPath)
				total = len(data)
			}
		}
		return nil
	})
	params := PreGen{Handler: data, Imports: imports, Package: file.Name.String()}

	tmpl, err := template.New("stringMethod").Parse(stringMethodTemplate)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error parsing template:", err)
		os.Exit(1)
	}
	var sb strings.Builder
	if err := tmpl.Execute(&sb, params); err != nil {
		fmt.Fprintln(os.Stderr, "Error executing template:", err)
		os.Exit(1)
	}
	src, err := format.Source([]byte(sb.String()))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error formatting source:", err)
		os.Exit(1)
	}

	os.WriteFile(output, src, 0644)
}
