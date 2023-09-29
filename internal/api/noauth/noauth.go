package noauth

import (
	"deferredMessage/internal/db"
	"deferredMessage/internal/utils"
	"deferredMessage/internal/utils/dto"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type NoAuth struct {
	db db.DB
}

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

func Init(db db.DB) NoAuth {
	return NoAuth{
		db: db,
	}
}
func (n NoAuth) Router(router *gin.RouterGroup) *gin.RouterGroup {
	r := router.Group("/")
	r.POST("/check", func(c *gin.Context) {
		//fmt.Println("register")
		body, exist := dto.GetStruct[RegisterBody](c, RegisterBody{})
		if !exist {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no body"})
			return
		}
		isExist, err := n.db.Collections.User.CheckUserByMail(body.Mail)
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
		body, exist := dto.GetStruct[RegisterBody](c, RegisterBody{})
		if !exist {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no body"})
			return
		}

		isExist, err := n.db.Collections.User.CheckUserByMail(body.Mail)
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
		userIP := c.ClientIP()
		session, err := n.db.Collections.Session.CreateSession(user.ID, time.Now().Add(time.Hour*24*31).Unix(), userIP)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "success", "user": gin.H{"name": user.Name, "mail": user.Mail}, "session": gin.H{"id": session.ID, "expire": session.Expire}})
	})
	r.POST("/login", func(c *gin.Context) {
		body, exist := dto.GetStruct[LoginBody](c, LoginBody{})
		if !exist {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no body"})
			return
		}

		user, isExist, err := n.db.Collections.User.GetUserByMail(body.Mail)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "user or password incorrect"})
			fmt.Println("/login error: ", err.Error())
			return
		}
		if !isExist {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "user or password incorrect"})
			return
		}
		if !utils.CheckPasswordHash(body.Password, user.Hash) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "user or password incorrect"})
			return
		}
		userIP := c.ClientIP()
		session, err := n.db.Collections.Session.CreateSession(user.ID, time.Now().Add(time.Hour*24*31).Unix(), userIP)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "success", "session": gin.H{"id": session.ID, "expire": session.Expire}})
	})
	return r
}
