package http

import (
	"log/slog"
	"net/http"

	"github.com/jellydator/validation"

	"github.com/allisson/psqlqueue/domain"
)

type errorResponseCode int //@name ErrorResponseCode

const (
	internalServerErrorCode errorResponseCode = iota + 1
	malformedRequest
	requestValidationFailedCode
	queueAlreadyExists
	queueNotFound
	messageNotFound
	topicAlreadyExists
	topicNotFound
	subscriptionAlreadyExists
	subscriptionNotFound
)

var errorResponses = map[string]errorResponse{
	"internal_server_error": {
		Code:       internalServerErrorCode,
		Message:    "internal server error",
		StatusCode: http.StatusInternalServerError,
	},
	"malformed_request": {
		Code:       malformedRequest,
		Message:    "malformed request body",
		StatusCode: http.StatusBadRequest,
	},
	"request_validation_failed": {
		Code:       requestValidationFailedCode,
		Message:    "request validation failed",
		StatusCode: http.StatusBadRequest,
	},
	"queue_already_exists": {
		Code:       queueAlreadyExists,
		Message:    "queue already exists",
		StatusCode: http.StatusBadRequest,
	},
	"queue_not_found": {
		Code:       queueNotFound,
		Message:    "queue not found",
		StatusCode: http.StatusNotFound,
	},
	"message_not_found": {
		Code:       messageNotFound,
		Message:    "message not found",
		StatusCode: http.StatusNotFound,
	},
	"topic_already_exists": {
		Code:       topicAlreadyExists,
		Message:    "topic already exists",
		StatusCode: http.StatusBadRequest,
	},
	"topic_not_found": {
		Code:       topicNotFound,
		Message:    "topic not found",
		StatusCode: http.StatusNotFound,
	},
	"subscription_already_exists": {
		Code:       subscriptionAlreadyExists,
		Message:    "subscription already exists",
		StatusCode: http.StatusBadRequest,
	},
	"subscription_not_found": {
		Code:       subscriptionNotFound,
		Message:    "subscription not found",
		StatusCode: http.StatusNotFound,
	},
}

type errorResponse struct {
	Code       errorResponseCode `json:"code"`
	Message    string            `json:"message"`
	Details    string            `json:"details,omitempty"`
	StatusCode int               `json:"-"`
} //@name ErrorResponse

func parseServiceError(serviceName, serviceMethod string, err error) errorResponse {
	if _, ok := err.(validation.Errors); ok {
		er := errorResponses["request_validation_failed"]
		er.Details = err.Error()
		return er
	}

	switch err {
	case domain.ErrQueueAlreadyExists:
		return errorResponses["queue_already_exists"]
	case domain.ErrQueueNotFound:
		return errorResponses["queue_not_found"]
	case domain.ErrMessageNotFound:
		return errorResponses["message_not_found"]
	case domain.ErrTopicAlreadyExists:
		return errorResponses["topic_already_exists"]
	case domain.ErrTopicNotFound:
		return errorResponses["topic_not_found"]
	case domain.ErrSubscriptionAlreadyExists:
		return errorResponses["subscription_already_exists"]
	case domain.ErrSubscriptionNotFound:
		return errorResponses["subscription_not_found"]
	default:
		slog.Error(serviceName, "method", serviceMethod, "error", err.Error())
		return errorResponses["internal_server_error"]
	}
}
