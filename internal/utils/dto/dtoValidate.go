package dto

import (
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetStruct[T any](c *gin.Context, body T) (T, bool) {
	if err := c.ShouldBindJSON(&body); err != nil {
		//c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return body, false
	}
	// res, exist := c.Get("body")
	// if !exist {
	// 	return body, false
	// }
	// return res.(T), true
	return body, true
}

// middleware dto json validation
// func Validate(c *gin.Context) {
// 	//get data from request
// 	var body any
// 	// using BindJson method to serialize body with struct

//		fmt.Println("validate: ", body)
//		if err := c.ShouldBindJSON(&body); err != nil {
//			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//			return
//		}
//		c.Set("body", body)
//	}
type limits struct {
	Count  int `json:"count"`
	Offset int `json:"offset"`
}

func GetLimits(c *gin.Context) (limits, error) {
	lims := limits{
		Count:  10,
		Offset: 0,
	}
	countString := c.DefaultQuery("count", "10")
	offsetString := c.DefaultQuery("offset", "0")
	count, err := strconv.Atoi(countString)
	if err != nil {
		return lims, err
	}
	lims.Count = count
	offset, err := strconv.Atoi(offsetString)
	if err != nil {
		return lims, err

	}
	lims.Offset = offset
	return lims, nil

}

type ParamTyped struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func TypedMap(paramTyped []ParamTyped, body map[string]interface{}) map[string]interface{} {
	filtred := make(map[string]interface{})
	for _, param := range paramTyped {
		if val, ok := body[param.Name]; ok {
			if reflect.TypeOf(val).String() == param.Type {
				filtred[param.Name] = val
			}
		}
	}
	return filtred
}
