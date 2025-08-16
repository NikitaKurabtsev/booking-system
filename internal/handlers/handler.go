package handlers

import (
	"errors"
	"fmt"
	"github.com/NikitaKurabtsev/booking-system/logger"
	swaggerFiles "github.com/swaggo/files"
	"net/http"
	"strconv"

	"github.com/NikitaKurabtsev/booking-system/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	service *services.Service
	logger  logger.Logger
}

// TODO: create Logger interface with dependency inversion
func NewHandler(service *services.Service, logger logger.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

func (h *Handler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]string{"message": "PONG"})
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.GET("/ping", h.Ping)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	auth := router.Group("/auth")
	h.RegisterAuthRoutes(auth)

	api := router.Group("/api",
		gin.Logger(),
		gin.Recovery(),
		h.prometheusMiddleware(),
		h.verifyUserMiddleware(),
	)
	resources := api.Group("/resources")
	h.RegisterResourceRoutes(resources)

	bookings := api.Group("/bookings")
	h.RegisterBookingRoutes(bookings)

	return router
}

func (h *Handler) RegisterAuthRoutes(r *gin.RouterGroup) {
	r.POST("/sign-up", h.signUp)
	r.POST("/sign-in", h.signIn)
}

func (h *Handler) RegisterResourceRoutes(r *gin.RouterGroup) {
	r.GET("/", h.getAllResources)
}

func (h *Handler) RegisterBookingRoutes(r *gin.RouterGroup) {
	r.GET("/:id", h.getBookings)
	r.POST("/", h.createBooking)
	r.DELETE("/:id", h.deleteBooking)
}

func getUserID(c *gin.Context) (int, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		return 0, errors.New("user id not found")
	}

	idInt, ok := id.(int)
	if !ok {
		return 0, errors.New("user id is of invalid type")
	}

	return idInt, nil
}

func getBookingID(c *gin.Context) (int, error) {
	id := c.Param("id")
	if id == "" {
		return 0, errors.New("empty id URL parameter")
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return 0, fmt.Errorf("invalid id parameter: %w", err)
	}

	return idInt, nil
}
