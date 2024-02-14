package middleware

import (
	"go-backend-starter-project/initializers"
	"go-backend-starter-project/models"
	"go-backend-starter-project/tokens"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequireAuth(c *gin.Context) {

	parsedToken, err := tokens.ParseToken(c)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if time.Now().Unix() > parsedToken.ExpiresAt.Unix() {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var user models.User
	initializers.DB.First(&user, "id = ?", parsedToken.Bearer)
	if user.ID == uuid.Nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set("user", user)
	c.Next()
}
