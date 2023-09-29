package dto

import (
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

// 	fmt.Println("validate: ", body)
// 	if err := c.ShouldBindJSON(&body); err != nil {
// 		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.Set("body", body)
// }
