package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	entities "github.com/Communinst/CWDB6Sem/backend/entity"
	"github.com/Communinst/CWDB6Sem/backend/validation"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"	
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) signUp(c *gin.Context) {
	var user entities.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := validation.IsPasswordValid(user.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid password"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user.Password = string(hashedPassword)
	user.SignUpDate = time.Now()

	userCount, err := h.service.UserServiceInterface.GetAllUsers(c.Request.Context())
	if err != nil {
		if err.Error() != "no user found" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check user count"})
			return
		}
	}

	// Only proceed with sign-up if there are no existing users
	if len(userCount) == 0 {
		roleId, _ := strconv.ParseInt(os.Getenv("DEFAULT_ADMIN_ROLE_ID"), 10, 64)
		user.RoleId = int(roleId)
	} else {
		roleId, _ := strconv.ParseInt(os.Getenv("DEFAULT_USER_ROLE_ID"), 10, 64)
		user.RoleId = int(roleId)
	}

	// Call the service layer to create the user
	if err := h.service.AuthServiceInterface.PostUser(c.Request.Context(), &user); err != nil {
		var postgresErr *pq.Error
		if errors.As(err, &postgresErr) && postgresErr.Code.Name() == "unique_violation" {
			c.JSON(http.StatusConflict, gin.H{"error": "Failed to create user"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

func (h *Handler) signIn(c *gin.Context) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Retrieve the user by email
	user, err := h.service.AuthServiceInterface.GetUserByEmail(c.Request.Context(), credentials.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Verify the password

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		fmt.Printf("%s\n", err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate an authentication token
	token := os.Getenv("AUTHORIZATION_TOKEN_SECRET")
	if token == "" {
		log.Fatalf("ACCESS_TOKEN_SECRET environment variable not set")
	}
	expiry, _ := strconv.ParseInt(os.Getenv("AUTHORIZATION_EXPIRE_TIME"), 10, 64) //TODO: move
	if expiry == 0 {
		log.Fatalf("AUTHORIZATION_EXPIRE_TIME environment variable not set")
	}
	userToken, err := h.service.GenerateAuthToken(user, token, 72) // 72 hours expiration
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	defaultAdmin, _ := strconv.ParseInt(os.Getenv("DEFAULT_ADMIN_ROLE_ID"), 10, 64)
	fmt.Printf("%s\n", userToken)
	if user.RoleId == int(defaultAdmin) {
		//c.Redirect(http.StatusOK, "/admin")
		c.JSON(http.StatusOK, gin.H{"token": userToken, "role": "admin", "userId": user.UserId})
	} else {
		//c.Redirect(http.StatusOK, "/api")
		c.JSON(http.StatusOK, gin.H{"token": userToken, "role": "user", "userId": user.UserId})
	}

	//c.JSON(http.StatusOK, gin.H{"token": userToken})
}
