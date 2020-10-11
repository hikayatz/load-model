package helpers

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

const (
	defaultMemory = 54 << 20 // 54 MB

	HeaderContentType   = "Content-Type"
	MIMEApplicationJSON = "application/json"
	MIMEApplicationXML  = "application/xml"
	MIMETextXML         = "text/xml"
	MIMEApplicationForm = "application/x-www-form-urlencoded"
	MIMEMultipartForm   = "multipart/form-data"
)

func LoadModel(model interface{}, req *http.Request, tag string) (err error) {
	if req.ContentLength == 0 {
		return errors.New("Request body can't be empty")
	}
	ctype := req.Header.Get(HeaderContentType)
	var dataRequest map[string][]string
	var jsonMap map[string]interface{}
	switch {
	case strings.HasPrefix(ctype, MIMEApplicationJSON):
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(body, &jsonMap); err != nil {
			return err
		}
	case strings.HasPrefix(ctype, MIMEApplicationXML), strings.HasPrefix(ctype, MIMETextXML):
		if err = xml.NewDecoder(req.Body).Decode(model); err != nil {
			if decField, ok := err.(*xml.UnsupportedTypeError); ok {
				return errors.New(fmt.Sprintf("Unsupported type error: type=%v, error=%v", decField.Type, decField.Error()))
			} else if se, ok := err.(*xml.SyntaxError); ok {
				return errors.New(fmt.Sprintf("Syntax error: line=%v, error=%v", se.Line, se.Error()))
			} else {
				return err
			}
		}
	case strings.HasPrefix(ctype, MIMEApplicationForm), strings.HasPrefix(ctype, MIMEMultipartForm):
		dataRequest, err = getValueParams(req)
		if err != nil {
			return err
		}
	}
	// check type model is struct
	typ := reflect.TypeOf(model).Elem()
	val := reflect.ValueOf(model).Elem()
	if typ.Kind() != reflect.Struct {
		return errors.New("binding element must be a struct")
	}
	// loop attr struct
	for i := 0; i < typ.NumField(); i++ {
		typeField := typ.Field(i)
		structField := val.Field(i)

		// continue if can be site
		if !structField.CanSet() {
			continue
		}

		structFieldKind := structField.Kind()
		inputFieldName := typeField.Tag.Get(tag)

		if inputFieldName == "" {
			inputFieldName = typeField.Name
		}
		var val = ""
		if dataRequest != nil{
			inputValue, exists := dataRequest[inputFieldName]
			if !exists {
				continue
			}
			val = inputValue[0]
		}else if jsonMap != nil {
			inputValue, exists := jsonMap[inputFieldName]
			if !exists {
				continue
			}
			val = fmt.Sprintf("%v", inputValue)
		}else{
			return
		}


		switch structFieldKind {
		case reflect.Int:
			err = setIntField(val, 0, structField)
		case reflect.Int8:
			err = setIntField(val, 8, structField)
		case reflect.Int16:
			err = setIntField(val, 16, structField)
		case reflect.Int32:
			err = setIntField(val, 32, structField)
		case reflect.Int64:
			err = setIntField(val, 64, structField)
		case reflect.Uint:
			err = setUintField(val, 0, structField)
		case reflect.Uint8:
			err = setUintField(val, 8, structField)
		case reflect.Uint16:
			err = setUintField(val, 16, structField)
		case reflect.Uint32:
			err = setUintField(val, 32, structField)
		case reflect.Uint64:
			err = setUintField(val, 64, structField)
		case reflect.Bool:
			err = setBoolField(val, structField)
		case reflect.Float32:
			err = setFloatField(val, 32, structField)
		case reflect.Float64:
			err = setFloatField(val, 64, structField)
		case reflect.String:
			structField.SetString(val)


		default:
			//_ = errors.New("unknown type")
		}
	}

	return err
}
func getValueParams(req *http.Request) (url.Values, error) {
	if strings.HasPrefix(req.Header.Get(HeaderContentType), MIMEMultipartForm) {
		if err := req.ParseMultipartForm(defaultMemory); err != nil {
			return nil, err
		}
	} else {
		if err := req.ParseForm(); err != nil {
			return nil, err
		}
	}
	return req.Form, nil
}

func setIntField(value string, bitSize int, field reflect.Value) error {
	if value == "" {
		value = "0"
	}
	intVal, err := strconv.ParseInt(value, 10, bitSize)
	if err == nil {
		field.SetInt(intVal)
	}
	return err
}

func setUintField(value string, bitSize int, field reflect.Value) error {
	if value == "" {
		value = "0"
	}
	uintVal, err := strconv.ParseUint(value, 10, bitSize)
	if err == nil {
		field.SetUint(uintVal)
	}
	return err
}

func setBoolField(value string, field reflect.Value) error {
	if value == "" {
		value = "false"
	}
	boolVal, err := strconv.ParseBool(value)
	if err == nil {
		field.SetBool(boolVal)
	}
	return err
}

func setFloatField(value string, bitSize int, field reflect.Value) error {
	if value == "" {
		value = "0.0"
	}
	floatVal, err := strconv.ParseFloat(value, bitSize)
	if err == nil {
		field.SetFloat(floatVal)
	}
	return err
}

