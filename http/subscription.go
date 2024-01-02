package http

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/allisson/psqlqueue/domain"
)

// nolint:unused
type subscriptionRequest struct {
	ID             string              `json:"id" example:"my-new-subscription" validate:"required"`
	TopicID        string              `json:"topic_id" example:"my-new-topic" validate:"required"`
	QueueID        string              `json:"queue_id" example:"my-new-queue" validate:"required"`
	MessageFilters map[string][]string `json:"message_filters"`
} //@name SubscriptionRequest

// nolint:unused
type subscriptionResponse struct {
	ID             string              `json:"id" example:"my-new-subscription"`
	TopicID        string              `json:"topic_id" example:"my-new-topic"`
	QueueID        string              `json:"queue_id" example:"my-new-queue"`
	MessageFilters map[string][]string `json:"message_filters"`
	CreatedAt      time.Time           `json:"created_at" example:"2023-08-17T00:00:00Z"`
} //@name SubscriptionResponse

// nolint:unused
type subscriptionListResponse struct {
	Data   []*topicResponse `json:"data"`
	Offset int              `json:"offset" example:"0"`
	Limit  int              `json:"limit" example:"10"`
} //@name SubscriptionListResponse

// SubscriptionHandler exposes a REST API for domain.SubscriptionService.
type SubscriptionHandler struct {
	subscriptionService domain.SubscriptionService
}

// Create a subscription.
//
//	@Summary	Add a subscription
//	@Tags		subscriptions
//	@Accept		json
//	@Produce	json
//	@Param		request	body		subscriptionRequest	true	"Add a subscription"
//	@Success	201		{object}	topicResponse
//	@Failure	400		{object}	errorResponse
//	@Failure	500		{object}	errorResponse
//	@Router		/subscriptions [post]
func (s *SubscriptionHandler) Create(c *gin.Context) {
	subscription := domain.Subscription{}

	if err := c.ShouldBindJSON(&subscription); err != nil {
		slog.Error("malformed request", "error", err.Error())
		er := errorResponses["malformed_request"]
		c.JSON(er.StatusCode, &er)
		return
	}

	if err := s.subscriptionService.Create(c.Request.Context(), &subscription); err != nil {
		er := parseServiceError("subscriptionService", "Create", err)
		c.JSON(er.StatusCode, &er)
		return
	}

	c.JSON(http.StatusCreated, &subscription)
}

// Get a subscription.
//
//	@Summary	Show a subscription
//	@Tags		subscriptions
//	@Accept		json
//	@Produce	json
//	@Param		subscription_id	path		string	true	"Subscription id"
//	@Success	200				{object}	subscriptionResponse
//	@Failure	404				{object}	errorResponse
//	@Failure	500				{object}	errorResponse
//	@Router		/subscriptions/{subscription_id} [get]
func (s *SubscriptionHandler) Get(c *gin.Context) {
	id := c.Param("subscription_id")

	subscription, err := s.subscriptionService.Get(c.Request.Context(), id)
	if err != nil {
		er := parseServiceError("subscriptionService", "Get", err)
		c.JSON(er.StatusCode, &er)
		return
	}

	c.JSON(http.StatusOK, &subscription)
}

// List subscriptions.
//
//	@Summary	List subscriptions
//	@Tags		subscriptions
//	@Accept		json
//	@Produce	json
//	@Param		limit	query		int	false	"The limit indicates the maximum number of items to return"
//	@Param		offset	query		int	false	"The offset indicates the starting position of the query in relation to the complete set of unpaginated items"
//	@Success	200		{object}	subscriptionListResponse
//	@Failure	500		{object}	errorResponse
//	@Router		/subscriptions [get]
func (s *SubscriptionHandler) List(c *gin.Context) {
	request := newListRequestFromGIN(c)

	topics, err := s.subscriptionService.List(c.Request.Context(), request.Offset, request.Limit)
	if err != nil {
		er := parseServiceError("subscriptionService", "List", err)
		c.JSON(er.StatusCode, &er)
		return
	}

	response := listResponse{Data: topics, Offset: request.Offset, Limit: request.Limit}

	c.JSON(http.StatusOK, response)
}

// Delete a subscription.
//
//	@Summary	Delete a subscription
//	@Tags		subscriptions
//	@Accept		json
//	@Produce	json
//	@Param		subscription_id	path	string	true	"Subscription id"
//	@Success	204				"No Content"
//	@Failure	404				{object}	errorResponse
//	@Failure	500				{object}	errorResponse
//	@Router		/subscriptions/{subscription_id} [delete]
func (s *SubscriptionHandler) Delete(c *gin.Context) {
	id := c.Param("subscription_id")

	if err := s.subscriptionService.Delete(c.Request.Context(), id); err != nil {
		er := parseServiceError("subscriptionService", "Delete", err)
		c.JSON(er.StatusCode, &er)
		return
	}

	c.Status(http.StatusNoContent)
}

// NewSubscriptionHandler returns a new SubscriptionHandler.
func NewSubscriptionHandler(subscriptionService domain.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{subscriptionService: subscriptionService}
}
