package noauth

import (
	"deferredMessage/internal/db"
	"deferredMessage/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type NoAuth struct {
	db db.DB
}

func Init(db db.DB) NoAuth {
	return NoAuth{
		db: db,
	}
}
func (n NoAuth) Router(router *gin.RouterGroup) *gin.RouterGroup {
	r := router.Group("/")
	r.POST("/check", func(c *gin.Context) {
		//fmt.Println("register")
		body, exist := getStruct[RegisterBody](c, RegisterBody{})
		if !exist {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no body"})
			return
		}
		isExist, err := n.db.Collections.User.CheckUser(body.Mail)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if isExist {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user already exist"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "not exist"})
	})

	r.POST("/register", func(c *gin.Context) {
		//fmt.Println("register")
		body, exist := getStruct[RegisterBody](c, RegisterBody{})
		if !exist {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no body"})
			return
		}

		isExist, err := n.db.Collections.User.CheckUser(body.Mail)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if isExist {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "user already exist"})
			return
		}
		hash, err := utils.HashPassword(body.Password)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		user, err := n.db.Collections.User.CreateUser(body.Name, body.Mail, hash)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "success", "user": gin.H{"name": user.Name, "mail": user.Mail}})
	})
	return r
}
