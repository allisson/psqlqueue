package http

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/allisson/psqlqueue/domain"
)

// nolint:unused
type messageRequest struct {
	Body       string            `json:"body" validate:"required"`
	Label      *string           `json:"label" validate:"optional"`
	Attributes map[string]string `json:"attributes" validate:"optional"`
} //@name MessageRequest

// nolint:unused
type messageResponse struct {
	ID               string            `json:"id" example:"7b98fe50-affd-4685-bd7d-3ae5e41493af"`
	QueueID          string            `json:"queue_id" example:"my-new-queue"`
	Label            *string           `json:"label"`
	Body             string            `json:"body"`
	Attributes       map[string]string `json:"attributes"`
	DeliveryAttempts int               `json:"delivery_attempts" example:"1"`
	CreatedAt        time.Time         `json:"created_at" db:"created_at" example:"2023-08-17T00:00:00Z"`
} //@name MessageResponse

// nolint:unused
type messageListRequest struct {
	Label  *string `form:"label" validate:"optional"`
	Offset uint    `form:"offset" validate:"required"`
	Limit  uint    `form:"limit" validate:"required"`
} //@name MessageListRequest

// nolint:unused
type messageListResponse struct {
	Data  []*messageResponse `json:"data"`
	Limit int                `json:"limit" example:"10"`
} //@name MessageListResponse

// nolint:unused
type messageNackRequest struct {
	VisibilityTimeoutSeconds uint `form:"visibility_timeout_seconds" validate:"required"`
} //@name MessageNackRequest

// Message exposes a REST API for domain.MessageService.
type MessageHandler struct {
	messageService domain.MessageService
	cfg            *domain.Config
}

// Create a message.
//
//	@Summary	Add a message
//	@Tags		messages
//	@Accept		json
//	@Produce	json
//	@Param		queue_id	path	string			true	"Queue id"
//	@Param		request		body	messageRequest	true	"Add a message"
//	@Success	204			"No Content"
//	@Failure	400			{object}	errorResponse
//	@Failure	500			{object}	errorResponse
//	@Router		/queue/{queue_id}/messages [post]
func (m *MessageHandler) Create(c *gin.Context) {
	message := domain.Message{}

	if err := c.ShouldBindJSON(&message); err != nil {
		slog.Error("malformed request", "error", err.Error())
		er := errorResponses["malformed_request"]
		c.JSON(er.StatusCode, &er)
		return
	}

	message.QueueID = c.Param("queue_id")

	if err := m.messageService.Create(c.Request.Context(), &message); err != nil {
		er := parseServiceError("messageService", "Create", err)
		c.JSON(er.StatusCode, &er)
		return
	}

	c.Status(http.StatusNoContent)
}

// List messages.
//
//	@Summary	List messages
//	@Tags		messages
//	@Accept		json
//	@Produce	json
//	@Param		queue_id	path		string	true	"Queue id"
//	@Param		label		path		string	false	"Label"
//	@Param		limit		query		int		false	"The limit indicates the maximum number of items to return"
//	@Success	200			{object}	messageListResponse
//	@Failure	404			{object}	errorResponse
//	@Failure	500			{object}	errorResponse
//	@Router		/queues/{queue_id}/messages [get]
func (m *MessageHandler) List(c *gin.Context) {
	queueID := c.Param("queue_id")

	request := messageListRequest{Offset: 0, Limit: 10}
	if err := c.ShouldBindQuery(&request); err != nil {
		slog.Warn("message list request error", "error", err)
	}

	request.Limit = min(request.Limit, m.cfg.QueueMaxNumberOfMessages)
	request.Offset = 0

	messages, err := m.messageService.List(c.Request.Context(), queueID, request.Label, request.Limit)
	if err != nil {
		er := parseServiceError("messageService", "List", err)
		c.JSON(er.StatusCode, &er)
		return
	}

	response := listResponse{Data: messages, Offset: request.Offset, Limit: request.Limit}

	c.JSON(http.StatusOK, response)
}

// Ack a message.
//
//	@Summary	Ack a message
//	@Tags		messages
//	@Accept		json
//	@Produce	json
//	@Param		queue_id	path	string	true	"Queue id"
//	@Param		message_id	path	string	true	"Message id"
//	@Success	204			"No Content"
//	@Failure	404			{object}	errorResponse
//	@Failure	500			{object}	errorResponse
//	@Router		/queues/{queue_id}/messages/{message_id}/ack [put]
func (m *MessageHandler) Ack(c *gin.Context) {
	messageID := c.Param("message_id")

	if err := m.messageService.Ack(c.Request.Context(), messageID); err != nil {
		er := parseServiceError("messageService", "Ack", err)
		c.JSON(er.StatusCode, &er)
		return
	}

	c.Status(http.StatusNoContent)
}

// Nack a message.
//
//	@Summary	Nack a message
//	@Tags		messages
//	@Accept		json
//	@Produce	json
//	@Param		queue_id	path	string				true	"Queue id"
//	@Param		message_id	path	string				true	"Message id"
//	@Param		request		body	messageNackRequest	true	"Nack a message"
//	@Success	204			"No Content"
//	@Failure	404			{object}	errorResponse
//	@Failure	500			{object}	errorResponse
//	@Router		/queues/{queue_id}/messages/{message_id}/nack [put]
func (m *MessageHandler) Nack(c *gin.Context) {
	messageID := c.Param("message_id")

	request := messageNackRequest{}
	if err := c.ShouldBindQuery(&request); err != nil {
		slog.Warn("message nack request error", "error", err)
	}

	if err := m.messageService.Nack(c.Request.Context(), messageID, request.VisibilityTimeoutSeconds); err != nil {
		er := parseServiceError("messageService", "Ack", err)
		c.JSON(er.StatusCode, &er)
		return
	}

	c.Status(http.StatusNoContent)
}

// NewMessageHandler returns a new MessageHandler.
func NewMessageHandler(messageService domain.MessageService) *MessageHandler {
	return &MessageHandler{
		messageService: messageService,
		cfg:            domain.NewConfig(),
	}
}
