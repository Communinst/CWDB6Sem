package handler

import (
	"net/http"

	authToken "github.com/Communinst/CWDB6Sem/backend/JSONWebTokens"
	"github.com/Communinst/CWDB6Sem/backend/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.Service
}

func New(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) InitRoutes(middleware ...gin.HandlerFunc) *gin.Engine {
	router := gin.Default()
	router.Use(middleware...)

	// Configure CORS to allow all origins for testing
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// open for eo
	welcome := router.Group("/welcome")
	welcome.GET("/", func(context *gin.Context) {
		context.JSON(http.StatusOK, "welcome page")
	})

	auth := welcome.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}

	//apiRouter := router.Group("/api", authToken.JwtAuthMiddleware())

	// ADMIN
	admin := router.Group("/admin", authToken.JwtAuthMiddleware())
	{
		dumps := admin.Group("/dumps")
		{
			dumps.POST("/create", h.createDump)
			dumps.POST("/restore", h.restoreDump)
			dumps.GET("/", h.getAllDumps)
		}
		users := admin.Group("/users")
		{
			users.POST("/create", h.postUser)
			users.GET("/", h.getAllUsers)
			users.DELETE("/:id", h.deleteUser)
			users.PUT("/:id/role/:role_id", h.updateUserRole)
		}
	}

	return router
}
