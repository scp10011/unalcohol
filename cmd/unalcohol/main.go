package main

import (
	_ "embed"
	"fmt"
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
	Path   string
	Name   string
	Method string
	In     []PreParam
}

func main() {
	root := "./.."
	log.Println(os.Args)
	if len(os.Args) > 1 {
		root = os.Args[1]
	}
	e, err := os.Executable()
	base, _ := filepath.Split(e)
	base = strings.TrimSuffix(base, "/")
	modBuffer, err := os.ReadFile(filepath.Join(base, "go.mod"))
	if err != nil {
		log.Fatalln("Did not find the go.mod file.")
	}
	packageName := modfile.ModulePath(modBuffer)
	log.Println(root, packageName, base)
	findPath := strings.Replace(root, packageName, base, 1)
	if findPath == root {
		log.Fatalln("Package name does not match the declaration.")
	}
	imports := make([]string, 0)
	data := make(map[string]PrePath)
	total := 0
	err = filepath.Walk(findPath, func(path string, info os.FileInfo, err error) error {
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
			importPath := strings.Replace(path, findPath, root, 1)
			for _, file := range pkg.Files {
				handler := make(map[token.Pos]PreHandler)
				for _, group := range file.Comments {
					for _, comment := range group.List {
						if strings.HasPrefix(comment.Text, "//go:api") {
							split := strings.Split(comment.Text, " ")
							if len(split) != 3 {
								log.Fatalf("Compiler Error: %+v", comment.Text)
							}
							handler[comment.End()+1] = PreHandler{
								Method: split[1],
								Path:   split[2],
								In:     make([]PreParam, 0),
							}
						}
					}
				}
				for _, decl := range file.Decls {
					request, ok := handler[decl.Pos()]
					if !ok {
						continue
					}
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
					request.Name = funcDecl.Name.String()
					for _, param := range funcDecl.Type.Params.List {
						for _, ident := range param.Names {
							t := getTypeExpr(param.Type, name)
							request.In = append(request.In, PreParam{
								Key:  ident.Name,
								Type: fmt.Sprintf("%s{}", t),
							})
						}
					}
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
	params := PreGen{Handler: data, Imports: imports}

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
	os.WriteFile("unalcohol_gen.go", src, 0644)
}
