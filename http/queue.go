package http

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/allisson/psqlqueue/domain"
)

// nolint:unused
type queueRequest struct {
	ID                      string `json:"id" example:"my-new-queue" validate:"required"`
	AckDeadlineSeconds      int    `json:"ack_deadline_seconds" example:"30" validate:"required"`
	MessageRetentionSeconds int    `json:"message_retention_seconds" example:"604800" validate:"required"`
	DeliveryDelaySeconds    int    `json:"delivery_delay_seconds" example:"0" validate:"required"`
} //@name QueueRequest

// nolint:unused
type queueUpdateRequest struct {
	AckDeadlineSeconds      int `json:"ack_deadline_seconds" example:"30" validate:"required"`
	MessageRetentionSeconds int `json:"message_retention_seconds" example:"604800" validate:"required"`
	DeliveryDelaySeconds    int `json:"delivery_delay_seconds" example:"0" validate:"required"`
} //@name QueueUpdateRequest

// nolint:unused
type queueResponse struct {
	ID                      string    `json:"id" example:"my-new-queue"`
	AckDeadlineSeconds      int       `json:"ack_deadline_seconds" example:"30"`
	MessageRetentionSeconds int       `json:"message_retention_seconds" example:"604800"`
	DeliveryDelaySeconds    int       `json:"delivery_delay_seconds" example:"0"`
	CreatedAt               time.Time `json:"created_at" example:"2023-08-17T00:00:00Z"`
	UpdatedAt               time.Time `json:"updated_at" example:"2023-08-17T00:00:00Z"`
} //@name QueueResponse

// nolint:unused
type queueListResponse struct {
	Data   []*queueResponse `json:"data"`
	Offset int              `json:"offset" example:"0"`
	Limit  int              `json:"limit" example:"10"`
} //@name QueueListResponse

// nolint:unused
type queueStatsResponse struct {
	NumUndeliveredMessages         int `json:"num_undelivered_messages" example:"1"`
	OldestUnackedMessageAgeSeconds int `json:"oldest_unacked_message_age_seconds" example:"1"`
} //@name QueueStatsResponse

// Queue exposes a REST API for domain.QueueService.
type QueueHandler struct {
	queueService domain.QueueService
}

// Create a queue.
//
//	@Summary	Add a queue
//	@Tags		queues
//	@Accept		json
//	@Produce	json
//	@Param		request	body		queueRequest	true	"Add a queue"
//	@Success	201		{object}	queueResponse
//	@Failure	400		{object}	errorResponse
//	@Failure	500		{object}	errorResponse
//	@Router		/queues [post]
func (q *QueueHandler) Create(c *gin.Context) {
	queue := domain.Queue{}

	if err := c.ShouldBindJSON(&queue); err != nil {
		slog.Error("malformed request", "error", err.Error())
		er := errorResponses["malformed_request"]
		c.JSON(er.StatusCode, &er)
		return
	}

	if err := q.queueService.Create(c.Request.Context(), &queue); err != nil {
		er := parseServiceError("queueService", "Create", err)
		c.JSON(er.StatusCode, &er)
		return
	}

	c.JSON(http.StatusCreated, &queue)
}

// Update a queue.
//
//	@Summary	Update a queue
//	@Tags		queues
//	@Accept		json
//	@Produce	json
//	@Param		queue_id	path		string				true	"Queue id"
//	@Param		request		body		queueUpdateRequest	true	"Update a queue"
//	@Success	200			{object}	queueResponse
//	@Failure	400			{object}	errorResponse
//	@Failure	404			{object}	errorResponse
//	@Failure	500			{object}	errorResponse
//	@Router		/queues/{queue_id} [put]
func (q *QueueHandler) Update(c *gin.Context) {
	queue := domain.Queue{}

	if err := c.ShouldBindJSON(&queue); err != nil {
		slog.Error("malformed request", "error", err.Error())
		er := errorResponses["malformed_request"]
		c.JSON(er.StatusCode, &er)
		return
	}

	queue.ID = c.Param("queue_id")

	if err := q.queueService.Update(c.Request.Context(), &queue); err != nil {
		er := parseServiceError("queueService", "Update", err)
		c.JSON(er.StatusCode, &er)
		return
	}

	c.JSON(http.StatusOK, &queue)
}

// Get a queue.
//
//	@Summary	Show a queue
//	@Tags		queues
//	@Accept		json
//	@Produce	json
//	@Param		queue_id	path		string	true	"Queue id"
//	@Success	200			{object}	queueResponse
//	@Failure	404			{object}	errorResponse
//	@Failure	500			{object}	errorResponse
//	@Router		/queues/{queue_id} [get]
func (q *QueueHandler) Get(c *gin.Context) {
	id := c.Param("queue_id")

	queue, err := q.queueService.Get(c.Request.Context(), id)
	if err != nil {
		er := parseServiceError("queueService", "Get", err)
		c.JSON(er.StatusCode, &er)
		return
	}

	c.JSON(http.StatusOK, &queue)
}

// List queues.
//
//	@Summary	List queues
//	@Tags		queues
//	@Accept		json
//	@Produce	json
//	@Param		limit	query		int	false	"The limit indicates the maximum number of items to return"
//	@Param		offset	query		int	false	"The offset indicates the starting position of the query in relation to the complete set of unpaginated items"
//	@Success	200		{object}	queueListResponse
//	@Failure	500		{object}	errorResponse
//	@Router		/queues [get]
func (q *QueueHandler) List(c *gin.Context) {
	request := newListRequestFromGIN(c)

	queues, err := q.queueService.List(c.Request.Context(), request.Offset, request.Limit)
	if err != nil {
		er := parseServiceError("queueService", "List", err)
		c.JSON(er.StatusCode, &er)
		return
	}

	response := listResponse{Data: queues, Offset: request.Offset, Limit: request.Limit}

	c.JSON(http.StatusOK, response)
}

// Delete a queue.
//
//	@Summary	Delete a queue
//	@Tags		queues
//	@Accept		json
//	@Produce	json
//	@Param		queue_id	path	string	true	"Queue id"
//	@Success	204			"No Content"
//	@Failure	404			{object}	errorResponse
//	@Failure	500			{object}	errorResponse
//	@Router		/queues/{queue_id} [delete]
func (q *QueueHandler) Delete(c *gin.Context) {
	id := c.Param("queue_id")

	if err := q.queueService.Delete(c.Request.Context(), id); err != nil {
		er := parseServiceError("queueService", "Delete", err)
		c.JSON(er.StatusCode, &er)
		return
	}

	c.Status(http.StatusNoContent)
}

// Get the queue stats.
//
//	@Summary	Get the queue stats
//	@Tags		queues
//	@Accept		json
//	@Produce	json
//	@Param		queue_id	path		string	true	"Queue id"
//	@Success	200			{object}	queueStatsResponse
//	@Failure	404			{object}	errorResponse
//	@Failure	500			{object}	errorResponse
//	@Router		/queues/{queue_id}/stats [get]
func (q *QueueHandler) Stats(c *gin.Context) {
	id := c.Param("queue_id")

	stats, err := q.queueService.Stats(c.Request.Context(), id)
	if err != nil {
		er := parseServiceError("queueService", "Stats", err)
		c.JSON(er.StatusCode, &er)
		return
	}

	c.JSON(http.StatusOK, &stats)
}

// Purge a queue.
//
//	@Summary	Purge a queue
//	@Tags		queues
//	@Accept		json
//	@Produce	json
//	@Param		queue_id	path	string	true	"Queue id"
//	@Success	204			"No Content"
//	@Failure	404			{object}	errorResponse
//	@Failure	500			{object}	errorResponse
//	@Router		/queues/{queue_id}/purge [put]
func (q *QueueHandler) Purge(c *gin.Context) {
	id := c.Param("queue_id")

	if err := q.queueService.Purge(c.Request.Context(), id); err != nil {
		er := parseServiceError("queueService", "Purge", err)
		c.JSON(er.StatusCode, &er)
		return
	}

	c.Status(http.StatusNoContent)
}

// Cleanup a queue.
//
//	@Summary	Cleanup a queue removing expired and acked messages
//	@Tags		queues
//	@Accept		json
//	@Produce	json
//	@Param		queue_id	path	string	true	"Queue id"
//	@Success	204			"No Content"
//	@Failure	404			{object}	errorResponse
//	@Failure	500			{object}	errorResponse
//	@Router		/queues/{queue_id}/cleanup [put]
func (q *QueueHandler) Cleanup(c *gin.Context) {
	id := c.Param("queue_id")

	if err := q.queueService.Cleanup(c.Request.Context(), id); err != nil {
		er := parseServiceError("queueService", "Cleanup", err)
		c.JSON(er.StatusCode, &er)
		return
	}

	c.Status(http.StatusNoContent)
}

// NewQueueHandler returns a new QueueHandler.
func NewQueueHandler(queueService domain.QueueService) *QueueHandler {
	return &QueueHandler{queueService: queueService}
}
