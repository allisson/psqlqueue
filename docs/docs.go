// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/queue/{queue_id}/messages": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "messages"
                ],
                "summary": "Add a message",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Queue id",
                        "name": "queue_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Add a message",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/MessageRequest"
                        }
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            }
        },
        "/queues": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "queues"
                ],
                "summary": "List queues",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "The limit indicates the maximum number of items to return",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "The offset indicates the starting position of the query in relation to the complete set of unpaginated items",
                        "name": "offset",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/QueueListResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "queues"
                ],
                "summary": "Add a queue",
                "parameters": [
                    {
                        "description": "Add a queue",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/QueueRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/QueueResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            }
        },
        "/queues/{queue_id}": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "queues"
                ],
                "summary": "Show a queue",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Queue id",
                        "name": "queue_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/QueueResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            },
            "put": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "queues"
                ],
                "summary": "Update a queue",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Queue id",
                        "name": "queue_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Update a queue",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/QueueUpdateRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/QueueResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            },
            "delete": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "queues"
                ],
                "summary": "Delete a queue",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Queue id",
                        "name": "queue_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            }
        },
        "/queues/{queue_id}/cleanup": {
            "put": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "queues"
                ],
                "summary": "Cleanup a queue removing expired and acked messages",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Queue id",
                        "name": "queue_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            }
        },
        "/queues/{queue_id}/messages": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "messages"
                ],
                "summary": "List messages",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Queue id",
                        "name": "queue_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Filter by label",
                        "name": "label",
                        "in": "path"
                    },
                    {
                        "type": "integer",
                        "description": "The limit indicates the maximum number of items to return",
                        "name": "limit",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/MessageListResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            }
        },
        "/queues/{queue_id}/messages/{message_id}/ack": {
            "put": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "messages"
                ],
                "summary": "Ack a message",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Queue id",
                        "name": "queue_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Message id",
                        "name": "message_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            }
        },
        "/queues/{queue_id}/messages/{message_id}/nack": {
            "put": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "messages"
                ],
                "summary": "Nack a message",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Queue id",
                        "name": "queue_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Message id",
                        "name": "message_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Nack a message",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/MessageNackRequest"
                        }
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            }
        },
        "/queues/{queue_id}/purge": {
            "put": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "queues"
                ],
                "summary": "Purge a queue",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Queue id",
                        "name": "queue_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            }
        },
        "/queues/{queue_id}/stats": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "queues"
                ],
                "summary": "Get the queue stats",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Queue id",
                        "name": "queue_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/QueueStatsResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "ErrorResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "$ref": "#/definitions/ErrorResponseCode"
                },
                "details": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "ErrorResponseCode": {
            "type": "integer",
            "enum": [
                1,
                2,
                3,
                4,
                5,
                6
            ],
            "x-enum-varnames": [
                "internalServerErrorCode",
                "malformedRequest",
                "requestValidationFailedCode",
                "queueAlreadyExists",
                "queueNotFound",
                "messageNotFound"
            ]
        },
        "MessageListResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/MessageResponse"
                    }
                },
                "limit": {
                    "type": "integer",
                    "example": 10
                }
            }
        },
        "MessageNackRequest": {
            "type": "object",
            "required": [
                "visibility_timeout_seconds"
            ],
            "properties": {
                "visibility_timeout_seconds": {
                    "type": "integer"
                }
            }
        },
        "MessageRequest": {
            "type": "object",
            "required": [
                "body"
            ],
            "properties": {
                "attributes": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                },
                "body": {
                    "type": "string"
                },
                "label": {
                    "type": "string"
                }
            }
        },
        "MessageResponse": {
            "type": "object",
            "properties": {
                "attributes": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                },
                "body": {
                    "type": "string"
                },
                "created_at": {
                    "type": "string",
                    "example": "2023-08-17T00:00:00Z"
                },
                "delivery_attempts": {
                    "type": "integer",
                    "example": 1
                },
                "id": {
                    "type": "string",
                    "example": "7b98fe50-affd-4685-bd7d-3ae5e41493af"
                },
                "label": {
                    "type": "string"
                },
                "queue_id": {
                    "type": "string",
                    "example": "my-new-queue"
                }
            }
        },
        "QueueListResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/QueueResponse"
                    }
                },
                "limit": {
                    "type": "integer",
                    "example": 10
                },
                "offset": {
                    "type": "integer",
                    "example": 0
                }
            }
        },
        "QueueRequest": {
            "type": "object",
            "required": [
                "ack_deadline_seconds",
                "delivery_delay_seconds",
                "id",
                "message_retention_seconds"
            ],
            "properties": {
                "ack_deadline_seconds": {
                    "type": "integer",
                    "example": 30
                },
                "delivery_delay_seconds": {
                    "type": "integer",
                    "example": 0
                },
                "id": {
                    "type": "string",
                    "example": "my-new-queue"
                },
                "message_retention_seconds": {
                    "type": "integer",
                    "example": 604800
                }
            }
        },
        "QueueResponse": {
            "type": "object",
            "properties": {
                "ack_deadline_seconds": {
                    "type": "integer",
                    "example": 30
                },
                "created_at": {
                    "type": "string",
                    "example": "2023-08-17T00:00:00Z"
                },
                "delivery_delay_seconds": {
                    "type": "integer",
                    "example": 0
                },
                "id": {
                    "type": "string",
                    "example": "my-new-queue"
                },
                "message_retention_seconds": {
                    "type": "integer",
                    "example": 604800
                },
                "updated_at": {
                    "type": "string",
                    "example": "2023-08-17T00:00:00Z"
                }
            }
        },
        "QueueStatsResponse": {
            "type": "object",
            "properties": {
                "num_undelivered_messages": {
                    "type": "integer",
                    "example": 1
                },
                "oldest_unacked_message_age_seconds": {
                    "type": "integer",
                    "example": 1
                }
            }
        },
        "QueueUpdateRequest": {
            "type": "object",
            "required": [
                "ack_deadline_seconds",
                "delivery_delay_seconds",
                "message_retention_seconds"
            ],
            "properties": {
                "ack_deadline_seconds": {
                    "type": "integer",
                    "example": 30
                },
                "delivery_delay_seconds": {
                    "type": "integer",
                    "example": 0
                },
                "message_retention_seconds": {
                    "type": "integer",
                    "example": 604800
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
