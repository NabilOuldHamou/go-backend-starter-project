package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"go-backend-starter-project/initializers"
	"go-backend-starter-project/models"
	"go-backend-starter-project/tokens"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"os"
	"time"
)

func Signup(c *gin.Context) {
	var body struct {
		Username string
		Email    string
		Password string
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if body.Username == "" || body.Email == "" || body.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Status(http.StatusInternalServerError)
	}

	user := models.User{
		Email:    body.Email,
		Password: string(hashedPassword),
	}
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "An account already exists with that email.",
		})
		return
	}

	tokenString, err := createToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)

	c.JSON(http.StatusAccepted, gin.H{
		"user":  user,
		"token": tokenString,
	})
}

func Login(c *gin.Context) {
	var body struct {
		Email    string
		Password string
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if body.Email == "" || body.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
		return
	}

	var user models.User
	result := initializers.DB.First(&user, "email = ?", body.Email)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email/Password is invalid."})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email/Password is invalid."})
		return
	}

	tokenString, err := createToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)

	c.JSON(http.StatusAccepted, gin.H{
		"user":  user,
		"token": tokenString,
	})
}

func createToken(userId uuid.UUID) (string, error) {
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"bearer":    userId,
		"expiresAt": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := newToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", errors.New("could not create token")
	}

	return tokenString, nil
}

func Logout(c *gin.Context) {
	c.SetCookie("Authorization", "", -1, "", "", false, true)
	c.Status(http.StatusOK)
}

func ValidateToken(c *gin.Context) {

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

	c.Status(http.StatusAccepted)
}
