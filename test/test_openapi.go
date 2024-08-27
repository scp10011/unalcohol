package main

import (
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3gen"
	"gopkg.in/yaml.v3"
	"log"
)

type PPP struct {
	Data float64 `json:"data"`
}

type User struct {
	User string `json:"user"`
	T    *PPP   `json:"t"`
}

func main() {
	generator := openapi3gen.NewGenerator()
	doc := openapi3.T{}
	doc.Info = &openapi3.Info{
		Title:   "Test",
		Version: "1.0.0",
	}
	doc.Components = &openapi3.Components{
		Schemas: map[string]*openapi3.SchemaRef{},
	}

	schemaRef, err := generator.NewSchemaRefForValue(&User{}, doc.Components.Schemas)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("%+v", schemaRef)
	doc.Paths = openapi3.NewPaths()
	item := &openapi3.PathItem{}
	opt := openapi3.NewOperation()
	item.SetOperation("GET", opt)
	body := openapi3.NewRequestBody()
	body.Content = openapi3.NewContentWithJSONSchemaRef(schemaRef)
	opt.RequestBody = &openapi3.RequestBodyRef{Value: body}
	opt.Parameters = append(opt.Parameters, &openapi3.ParameterRef{Value: openapi3.NewQueryParameter("name")})
	doc.Paths.Set("/user", item)
	marshal, err := yaml.Marshal(doc)
	if err != nil {
		return
	}
	fmt.Println(string(marshal))
}
