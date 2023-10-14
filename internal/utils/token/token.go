package token

import (
	"fmt"
	"os"
)

func ParseAndCheckToken(token string) error {
	if token == "" {
		return fmt.Errorf("no token in body")
	}
	ADMIN_KEY, exist := os.LookupEnv("ADMIN_KEY")
	if !exist {
		return fmt.Errorf("ADMIN_KEY not found")
	}
	if token != ADMIN_KEY {
		//	c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse{Error:  "invalid token"})
		return fmt.Errorf("invalid token")

	}
	return nil
}
