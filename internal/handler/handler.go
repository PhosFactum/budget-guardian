package handler

import (
	"net/http"
	"strconv"

	"github.com/PhosFactum/budget-guardian/internal/interfaces"
	"github.com/PhosFactum/budget-guardian/internal/rabbitmq"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	budgetService interfaces.BudgetService
	rmq           *rabbitmq.Client
}

func NewHandler(
	budgetService interfaces.BudgetService,
	rmq *rabbitmq.Client,
) *Handler {
	return &Handler{
		budgetService: budgetService,
		rmq:           rmq,
	}
}

// InitRoutes:initializing gin route and basic api path
func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()

	api := router.Group("/api")
	{
		api.POST("/transaction", h.recordTransaction)
		api.GET("/daily-limit/:userID", h.getDailyLimit)
	}

	return router
}

// recordTransaction: function to record a new transaction
func (h *Handler) recordTransaction(c *gin.Context) {
	var request struct {
		UserID uint    `json:"user_id" binding:"required"`
		Amount float64 `json:"amount" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.budgetService.RecordTransaction(request.UserID, request.Amount); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Sending notification to RabbitMQ
	msg := "Transaction recorded: " + strconv.FormatFloat(request.Amount, 'f', 2, 64)
	if err := h.rmq.Publish("transactions", msg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send notification"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// getDailyLimit: function to get limit sum of a day
func (h *Handler) getDailyLimit(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	limit, debt, err := h.budgetService.CalculateDailyLimit(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"daily_limit": limit,
		"debt":        debt,
	})
}
