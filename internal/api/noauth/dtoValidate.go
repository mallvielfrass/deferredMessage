package noauth

import (
	"github.com/gin-gonic/gin"
)

type RegisterBody struct {
	// json tag to de-serialize json body
	Name     string `json:"name" binding:"required"`
	Mail     string `json:"mail" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type LoginBody struct {
	// json tag to de-serialize json body

	Mail     string `json:"mail" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func getStruct[T any](c *gin.Context, body T) (T, bool) {
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
