package handler

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"prosto-delaj-api/internal/usecase"
	"prosto-delaj-api/internal/utils"
	"prosto-delaj-api/models"
	"strings"
	"time"
)

type Handler struct {
	usecase *usecase.Usecase
}

func NewHandler(usecase *usecase.Usecase) *Handler {
	return &Handler{usecase: usecase}
}

func (h *Handler) InitRoutes(config models.ServerConfig) *gin.Engine {
	router := gin.Default()

	allowOrigins := strings.Split(config.Domain, ",")
	router.Use(cors.New(cors.Config{
		AllowOrigins: allowOrigins,
		AllowMethods: []string{http.MethodPut, http.MethodPost, http.MethodGet, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{"Content-Type", "Access-Control-Allow-Headers", "Access-Control-Allow-Origin",
			utils.HeaderAuthorization, utils.HeaderClientRequestId},
		ExposeHeaders: []string{"Content-Length", utils.HeaderTimestamp,
			utils.HeaderClientRequestId, utils.HeaderRequestId},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.POST("/create", h.create)
	return router
}
