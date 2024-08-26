package unalcohol

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

type Path[T int | string | float32 | float64] struct {
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

type JSON[T any] struct {
	Value T
}

func (j *JSON[T]) ParseRequest(key string, r *http.Request, resp http.ResponseWriter) (err error) {
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(&j.Value)
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

type Request struct {
	Value *http.Request
}

func (r *Request) ParseRequest(key string, req *http.Request, resp http.ResponseWriter) (err error) {
	r.Value = req
	return nil
}

type Response struct {
	Value http.ResponseWriter
}

func (r *Response) ParseRequest(key string, req *http.Request, resp http.ResponseWriter) (err error) {
	r.Value = resp
	return nil
}

type JSONResponse[T any] struct {
	Code int
	Data T
}

type ResponseInf interface {
	WriteResponse(http.ResponseWriter) error
}

func (r JSONResponse[T]) WriteResponse(w http.ResponseWriter) error {
	w.WriteHeader(r.Code)
	return json.NewEncoder(w).Encode(r.Data)
}

type IOResponse struct {
	Code int
	Data io.ReadCloser
}

func (r IOResponse) WriteResponse(w http.ResponseWriter) error {
	_, err := io.Copy(w, r.Data)
	return err
}
