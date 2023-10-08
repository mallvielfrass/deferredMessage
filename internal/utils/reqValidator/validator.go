package reqvalidator

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

func ValidateFlatMap(c *gin.Context, requestData any, fieldsConstraints []string) (map[string]interface{}, error) {
	requestBody := c.Request.Body
	bodyBytes, err := io.ReadAll(requestBody)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(bodyBytes, &requestData); err != nil {
		return nil, err
	}
	validate := validator.New()
	if err := validate.Struct(requestData); err != nil {
		return nil, err
	}
	body := make(map[string]interface{})
	if err := json.Unmarshal(bodyBytes, &body); err != nil {
		return nil, err
	}
	filtredBody := make(map[string]interface{})
	for _, v := range fieldsConstraints {
		fmt.Printf("v: %v\n", v)
		val, ok := body[v]
		if ok {
			filtredBody[v] = val

		}
	}
	return filtredBody, nil
}
func GetTagsFromFlatStruct(v interface{}) (tags []string) {
	t := reflect.TypeOf(v)
	for i := 0; i < t.NumField(); i++ {
		//	fmt.Printf("field: %v, type: %v, tag: %v, kind: %v\n", t.Field(i).Name, t.Field(i).Type, t.Field(i).Tag, t.Field(i).Type.Kind())

		switch t.Field(i).Type.Kind() {
		case reflect.Bool, reflect.String, reflect.Int, reflect.Int8, reflect.Int16,
			reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16,
			reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
			tag := t.Field(i).Tag.Get("json")
			if tag != "" {
				tags = append(tags, tag)
			}
		default:
			continue
		}

	}
	return
}
