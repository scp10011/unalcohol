package unalcohol

import (
	"encoding/json"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3gen"
	"io"
	"net/http"
	"strconv"
)

type Path[T int | string | float32 | float64 | int64] struct {
	Value T
}

func (p *Path[T]) ParseRequest(key string, r *http.Request, resp http.ResponseWriter) (err error) {
	value := r.PathValue(key)
	switch v := any(&p.Value).(type) {
	case *string:
		*v = value
	case *int:
		*v, err = strconv.Atoi(value)
	case *float32:
		if f, err := strconv.ParseFloat(value, 64); err != nil {
			return err
		} else {
			*v = float32(f)
		}
	case *float64:
		if f, err := strconv.ParseFloat(value, 64); err != nil {
			return err
		} else {
			*v = f
		}
	case *[]byte:
		*v = []byte(value)
	}
	return nil
}

func (p *Path[T]) Doc(key string, operation *openapi3.Operation) error {
	parameter := openapi3.NewPathParameter(key)
	schemaRef, err := openapi3gen.NewSchemaRefForValue(new(T), nil)
	if err != nil {
		return err
	}
	parameter.Schema = schemaRef
	operation.Parameters = append(operation.Parameters, &openapi3.ParameterRef{Value: parameter})
	return nil
}

type JSON[T any] struct {
	Value T
}

func (j *JSON[T]) ParseRequest(key string, r *http.Request, resp http.ResponseWriter) (err error) {
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(&j.Value)
}

func (j *JSON[T]) Doc(key string, operation *openapi3.Operation) error {
	schemaRef, err := openapi3gen.NewSchemaRefForValue(new(T), nil)
	if err != nil {
		return err
	}
	operation.RequestBody = &openapi3.RequestBodyRef{}
	body := openapi3.NewRequestBody()
	body.Content = openapi3.NewContentWithJSONSchemaRef(schemaRef)
	operation.RequestBody = &openapi3.RequestBodyRef{Value: body}
	return nil
}

type Body[T any] struct {
	Value T
}

func (b *Body[T]) ParseRequest(key string, r *http.Request, resp http.ResponseWriter) (err error) {
	r.ParseForm()
	value := r.Form.Get(key)
	switch v := any(&b.Value).(type) {
	case *string:
		*v = value
	case *int:
		*v, err = strconv.Atoi(value)
	case *float32:
		if f, err := strconv.ParseFloat(value, 64); err != nil {
			return err
		} else {
			*v = float32(f)
		}
	case *float64:
		if f, err := strconv.ParseFloat(value, 64); err != nil {
			return err
		} else {
			*v = f
		}
	case *[]byte:
		*v = []byte(value)
	}
	return nil
}

func (b *Body[T]) Doc(key string, operation *openapi3.Operation) error {
	var body *openapi3.RequestBody
	var data *openapi3.MediaType
	if operation.RequestBody != nil {
		body = operation.RequestBody.Value
	} else {
		body = openapi3.NewRequestBody().WithContent(openapi3.Content{
			"application/x-www-form-urlencoded": openapi3.NewMediaType(),
		})
		operation.RequestBody = &openapi3.RequestBodyRef{Value: body}
	}
	data = body.GetMediaType("application/x-www-form-urlencoded")
	schema := openapi3.NewObjectSchema()
	if data.Schema != nil && data.Schema.Value != nil {
		schema = data.Schema.Value
	} else {
		data.Schema = &openapi3.SchemaRef{Value: schema}
	}
	schemaRef, err := openapi3gen.NewSchemaRefForValue(new(T), nil)
	if err != nil {
		return err
	}
	schema.Properties[key] = schemaRef
	return nil
}

type Param[T any] struct {
	Value T
}

func (p *Param[T]) ParseRequest(key string, r *http.Request, resp http.ResponseWriter) (err error) {
	values := r.URL.Query()
	value := values.Get(key)
	switch v := any(&p.Value).(type) {
	case *string:
		*v = value
	case *int:
		*v, err = strconv.Atoi(value)
	case *float32:
		if f, err := strconv.ParseFloat(value, 64); err != nil {
			return err
		} else {
			*v = float32(f)
		}
	case *float64:
		if f, err := strconv.ParseFloat(value, 64); err != nil {
			return err
		} else {
			*v = f
		}
	case *[]byte:
		*v = []byte(value)
	}
	return nil
}

func (p *Param[T]) Doc(key string, operation *openapi3.Operation) error {
	parameter := openapi3.NewQueryParameter(key)
	schemaRef, err := openapi3gen.NewSchemaRefForValue(new(T), nil)
	if err != nil {
		return err
	}
	parameter.Schema = schemaRef
	operation.Parameters = append(operation.Parameters, &openapi3.ParameterRef{Value: parameter})
	return nil
}

type Header[T any] struct {
	Value T
}

func (h *Header[T]) ParseRequest(key string, r *http.Request, resp http.ResponseWriter) (err error) {
	value := r.Header.Get(key)
	switch v := any(&h.Value).(type) {
	case *string:
		*v = value
	case *int:
		*v, err = strconv.Atoi(value)
	case *float32:
		if f, err := strconv.ParseFloat(value, 64); err != nil {
			return err
		} else {
			*v = float32(f)
		}
	case *float64:
		if f, err := strconv.ParseFloat(value, 64); err != nil {
			return err
		} else {
			*v = f
		}
	case *[]byte:
		*v = []byte(value)
	}
	return nil
}

func (p *Header[T]) Doc(key string, operation *openapi3.Operation) error {
	parameter := openapi3.NewHeaderParameter(key)
	schemaRef, err := openapi3gen.NewSchemaRefForValue(new(T), nil)
	if err != nil {
		return err
	}
	parameter.Schema = schemaRef
	operation.Parameters = append(operation.Parameters, &openapi3.ParameterRef{Value: parameter})
	return nil
}

type Request struct {
	Value *http.Request
}

func (r *Request) ParseRequest(key string, req *http.Request, resp http.ResponseWriter) (err error) {
	r.Value = req
	return nil
}

func (r *Request) Doc(key string, operation *openapi3.Operation) error {
	return nil
}

type Response struct {
	Value http.ResponseWriter
}

func (r *Response) ParseRequest(key string, req *http.Request, resp http.ResponseWriter) (err error) {
	r.Value = resp
	return nil
}

func (r *Response) Doc(key string, operation *openapi3.Operation) error {
	return nil
}

type JSONResponse[T any] struct {
	Code int
	Data T
}

type ResponseInf interface {
	WriteResponse(http.ResponseWriter) error
}

func (r *JSONResponse[T]) WriteResponse(w http.ResponseWriter) error {
	w.WriteHeader(r.Code)
	return json.NewEncoder(w).Encode(r.Data)
}

func (r *JSONResponse[T]) Doc(operation *openapi3.Operation) error {
	operation.Responses = openapi3.NewResponses()
	response := openapi3.NewResponse()
	schemaRef, err := openapi3gen.NewSchemaRefForValue(new(T), nil)
	if err != nil {
		return err
	}
	response.Content = openapi3.NewContentWithJSONSchemaRef(schemaRef)
	operation.Responses.Set("200", &openapi3.ResponseRef{Value: response})
	return nil
}

type IOResponse struct {
	Code int
	Data io.ReadCloser
}

func (r IOResponse) WriteResponse(w http.ResponseWriter) error {
	_, err := io.Copy(w, r.Data)
	return err
}
