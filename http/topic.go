package http

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/allisson/psqlqueue/domain"
)

// nolint:unused
type topicRequest struct {
	ID string `json:"id" example:"my-new-topic" validate:"required"`
} //@name TopicRequest

// nolint:unused
type topicResponse struct {
	ID        string    `json:"id" example:"my-new-topic"`
	CreatedAt time.Time `json:"created_at" example:"2023-08-17T00:00:00Z"`
} //@name TopicResponse

// nolint:unused
type topicListResponse struct {
	Data   []*topicResponse `json:"data"`
	Offset int              `json:"offset" example:"0"`
	Limit  int              `json:"limit" example:"10"`
} //@name QueueListResponse

// Topic exposes a REST API for domain.TopicService.
type TopicHandler struct {
	topicService domain.TopicService
}

// Create a topic.
//
//	@Summary	Add a topic
//	@Tags		topics
//	@Accept		json
//	@Produce	json
//	@Param		request	body		topicRequest	true	"Add a topic"
//	@Success	201		{object}	topicResponse
//	@Failure	400		{object}	errorResponse
//	@Failure	500		{object}	errorResponse
//	@Router		/topics [post]
func (t *TopicHandler) Create(c *gin.Context) {
	topic := domain.Topic{}

	if err := c.ShouldBindJSON(&topic); err != nil {
		slog.Error("malformed request", "error", err.Error())
		er := errorResponses["malformed_request"]
		c.JSON(er.StatusCode, &er)
		return
	}

	if err := t.topicService.Create(c.Request.Context(), &topic); err != nil {
		er := parseServiceError("topicService", "Create", err)
		c.JSON(er.StatusCode, &er)
		return
	}

	c.JSON(http.StatusCreated, &topic)
}

// Get a topic.
//
//	@Summary	Show a topic
//	@Tags		topics
//	@Accept		json
//	@Produce	json
//	@Param		topic_id	path		string	true	"Topic id"
//	@Success	200			{object}	topicResponse
//	@Failure	404			{object}	errorResponse
//	@Failure	500			{object}	errorResponse
//	@Router		/topics/{topic_id} [get]
func (t *TopicHandler) Get(c *gin.Context) {
	id := c.Param("topic_id")

	topic, err := t.topicService.Get(c.Request.Context(), id)
	if err != nil {
		er := parseServiceError("topicService", "Get", err)
		c.JSON(er.StatusCode, &er)
		return
	}

	c.JSON(http.StatusOK, &topic)
}

// List topics.
//
//	@Summary	List topics
//	@Tags		topics
//	@Accept		json
//	@Produce	json
//	@Param		limit	query		int	false	"The limit indicates the maximum number of items to return"
//	@Param		offset	query		int	false	"The offset indicates the starting position of the query in relation to the complete set of unpaginated items"
//	@Success	200		{object}	topicListResponse
//	@Failure	500		{object}	errorResponse
//	@Router		/topics [get]
func (t *TopicHandler) List(c *gin.Context) {
	request := newListRequestFromGIN(c)

	topics, err := t.topicService.List(c.Request.Context(), request.Offset, request.Limit)
	if err != nil {
		er := parseServiceError("topicService", "List", err)
		c.JSON(er.StatusCode, &er)
		return
	}

	response := listResponse{Data: topics, Offset: request.Offset, Limit: request.Limit}

	c.JSON(http.StatusOK, response)
}

// Delete a topic.
//
//	@Summary	Delete a topic
//	@Tags		topics
//	@Accept		json
//	@Produce	json
//	@Param		topic_id	path	string	true	"Topic id"
//	@Success	204			"No Content"
//	@Failure	404			{object}	errorResponse
//	@Failure	500			{object}	errorResponse
//	@Router		/topics/{topic_id} [delete]
func (t *TopicHandler) Delete(c *gin.Context) {
	id := c.Param("topic_id")

	if err := t.topicService.Delete(c.Request.Context(), id); err != nil {
		er := parseServiceError("topicService", "Delete", err)
		c.JSON(er.StatusCode, &er)
		return
	}

	c.Status(http.StatusNoContent)
}

// Create a message.
//
//	@Summary	Add a message
//	@Tags		topics
//	@Accept		json
//	@Produce	json
//	@Param		request	body		messageRequest	true	"Add a message"
//	@Success	201		{object}	topicResponse
//	@Failure	400		{object}	errorResponse
//	@Failure	500		{object}	errorResponse
//	@Router		/topics/{topic_id}/messages [post]
func (t *TopicHandler) CreateMessage(c *gin.Context) {
	message := domain.Message{}
	id := c.Param("topic_id")

	if err := c.ShouldBindJSON(&message); err != nil {
		slog.Error("malformed request", "error", err.Error())
		er := errorResponses["malformed_request"]
		c.JSON(er.StatusCode, &er)
		return
	}

	if err := t.topicService.CreateMessage(c.Request.Context(), id, &message); err != nil {
		er := parseServiceError("topicService", "CreateMessage", err)
		c.JSON(er.StatusCode, &er)
		return
	}

	c.Status(http.StatusNoContent)
}

// NewTopicHandler returns a new TopicHandler.
func NewTopicHandler(topicService domain.TopicService) *TopicHandler {
	return &TopicHandler{topicService: topicService}
}
