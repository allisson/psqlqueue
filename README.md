# psqlqueue
[![test](https://github.com/allisson/psqlqueue/actions/workflows/test.yml/badge.svg)](https://github.com/allisson/psqlqueue/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/allisson/psqlqueue)](https://goreportcard.com/report/github.com/allisson/psqlqueue)
[![Docker Repository on Quay](https://quay.io/repository/allisson/psqlqueue/status "Docker Repository on Quay")](https://quay.io/repository/allisson/psqlqueue)

Simple queue system powered by Golang and PostgreSQL.

## Quickstart

The idea of this service is to offer a simple queuing system using PostgreSQL as a backend.

First, we need a PostgreSQL database, for this, we will use docker:

```bash
docker run --name postgres-psqlqueue \
    --restart unless-stopped \
    -e POSTGRES_USER=psqlqueue \
    -e POSTGRES_PASSWORD=psqlqueue \
    -e POSTGRES_DB=psqlqueue \
    -p 5432:5432 \
    -d postgres:15-alpine
```

Before running anything, take a look at the possible environment variables that can be used with this service: https://github.com/allisson/psqlqueue/blob/main/env.sample

Now let's run the database migrations before starting the server:

```bash
docker run --rm \
    -e PSQLQUEUE_DATABASE_URL='postgres://psqlqueue:psqlqueue@host.docker.internal:5432/psqlqueue?sslmode=disable' \
    quay.io/allisson/psqlqueue migrate
```

```json
{"time":"2023-12-29T21:11:39.516360369Z","level":"INFO","msg":"migration process started"}
{"time":"2023-12-29T21:11:39.54908151Z","level":"INFO","msg":"migration process finished"}
```

Starting the server:

```bash
docker run --rm \
    -e PSQLQUEUE_DATABASE_URL='postgres://psqlqueue:psqlqueue@host.docker.internal:5432/psqlqueue?sslmode=disable' \
    -p 8000:8000 \
    quay.io/allisson/psqlqueue server
```

```json
{"time":"2023-12-29T21:14:30.898080659Z","level":"INFO","msg":"http server starting","host":"0.0.0.0","port":8000}
```

For creating a new queue we have these fields:
- "id": The identifier of this new queue.
- "ack_deadline_seconds": The maximum time before the consumer should acknowledge the message, after this time the message will be delivered again to consumers.
- "message_retention_seconds": The maximum time in which the message must be delivered to consumers, after this time the message will be marked as expired.
- "delivery_delay_seconds": The number of seconds to postpone the delivery of new messages to consumers.

```bash
curl --location 'http://localhost:8000/v1/queues' \
--header 'Content-Type: application/json' \
--data '{
    "id": "my-new-queue",
    "ack_deadline_seconds": 30,
    "message_retention_seconds": 1209600,
    "delivery_delay_seconds": 0
}'
```

```json
{
    "id": "my-new-queue",
    "ack_deadline_seconds": 30,
    "message_retention_seconds": 1209600,
    "delivery_delay_seconds": 0,
    "created_at": "2023-12-29T21:30:58.682194763Z",
    "updated_at": "2023-12-29T21:30:58.682194763Z"
}
```

For creating a new message we have these fields:
- "body": The body of the message.
- "label": A label that allows this message to be filtered.
- "attributes": The message attributes.

```bash
curl --location 'http://localhost:8000/v1/queues/my-new-queue/messages' \
--header 'Content-Type: application/json' \
--data '{
    "body": "message body",
    "label": "my-label",
    "attributes": {"attribute1": "attribute1", "attribute2": "attribute2"}
}'
```

For consuming the messages we have these filters:
- "label": To filter by the message label.
- "limit": To limit the number of messages.

```bash
curl --location 'http://localhost:8000/v1/queues/my-new-queue/messages?limit=1'
```

```json
{
    "data": [
        {
            "id": "01HJVRCQVAD9VBT10MCS74T0EN",
            "queue_id": "my-new-queue",
            "label": "my-label",
            "body": "message body",
            "attributes": {
                "attribute1": "attribute1",
                "attribute2": "attribute2"
            },
            "delivery_attempts": 1,
            "created_at": "2023-12-29T21:41:25.994731Z"
        }
    ],
    "limit": 1
}
```

Now you have 30 seconds to execute the ack or nack for this message, first we can do the nack:

```bash
curl --location --request PUT 'http://localhost:8000/v1/queues/my-new-queue/messages/01HJVRCQVAD9VBT10MCS74T0EN/nack' \
--header 'Content-Type: application/json' \
--data '{
    "visibility_timeout_seconds": 30 
}'
```

Now we need to wait 30 seconds before consuming this message again, after this time:

```bash
curl --location 'http://localhost:8000/v1/queues/my-new-queue/messages?limit=1'
```

```json
{
    "data": [
        {
            "id": "01HJVRCQVAD9VBT10MCS74T0EN",
            "queue_id": "my-new-queue",
            "label": "my-label",
            "body": "message body",
            "attributes": {
                "attribute1": "attribute1",
                "attribute2": "attribute2"
            },
            "delivery_attempts": 2,
            "created_at": "2023-12-29T21:41:25.994731Z"
        }
    ],
    "limit": 1
}
```

Now it's time to ack the message:

```bash
curl --location --request PUT 'http://localhost:8000/v1/queues/my-new-queue/messages/01HJVRCQVAD9VBT10MCS74T0EN/ack'
```

Let's try to consume the messages again:

```bash
curl --location 'http://localhost:8000/v1/queues/my-new-queue/messages/?limit=1'
```

```json
{
    "data": [],
    "limit": 1
}
```

After the ack, the message remains in the database marked as expired, to remove expired messages we can use the cleanup endpoint:

```bash
curl --location --request PUT 'http://localhost:8000/v1/queues/my-new-queue/cleanup'
```

This is the basics of using this service, I recommend that you check the swagger documentation at http://localhost:8000/v1/swagger/index.html to see more options.

## Pub/Sub mode

It's possible to use a Pub/Sub approach with the topics/subscriptions endpoints.

For this example let's imagine an event system that processes orders, we will have a topic called `orders` and will push some messages into this topic.

For creating a new topic we have these fields:
- "id": The identifier of this new topic.

```bash
curl --location 'http://localhost:8000/v1/topics' \
--header 'Content-Type: application/json' \
--data '{
    "id": "orders"
}'
```

```json
{
    "id": "orders",
    "created_at": "2024-01-02T22:20:43.351647Z"
}
```

Now we will create two new queues:
- "all-orders": For receiving all messages from the orders topic.
- "processed-orders": For receiving only the messages with the `status` attribute equal to `"processed"`.

```bash
curl --location 'http://localhost:8000/v1/queues' \
--header 'Content-Type: application/json' \
--data '{
    "id": "all-orders",
    "ack_deadline_seconds": 30,
    "message_retention_seconds": 1209600,
    "delivery_delay_seconds": 0
}'
```

```json
{
    "id": "all-orders",
    "ack_deadline_seconds": 30,
    "message_retention_seconds": 1209600,
    "delivery_delay_seconds": 0,
    "created_at": "2024-01-02T22:24:58.219593Z",
    "updated_at": "2024-01-02T22:24:58.219593Z"
}
```

```bash
curl --location 'http://localhost:8000/v1/queues' \
--header 'Content-Type: application/json' \
--data '{
    "id": "processed-orders",
    "ack_deadline_seconds": 30,
    "message_retention_seconds": 1209600,
    "delivery_delay_seconds": 0
}'
```

```json
{
    "id": "processed-orders",
    "ack_deadline_seconds": 30,
    "message_retention_seconds": 1209600,
    "delivery_delay_seconds": 0,
    "created_at": "2024-01-02T22:25:28.472891Z",
    "updated_at": "2024-01-02T22:25:28.472891Z"
}
```

Now we will create two subscriptions to link the topic with the queue.

For creating a new subscription we have these fields:
- "id": The identifier of this new subscription.
- "topic_id": The id of the topic.
- "queue_id": The id of the queue.
- "message_filters": The filter for use with the message attributes.

Creating the first subscription:

```bash
curl --location 'http://localhost:8000/v1/subscriptions' \
--header 'Content-Type: application/json' \
--data '{
    "id": "orders-to-all-orders",
    "topic_id": "orders",
    "queue_id": "all-orders"
}'
```

```json
{
    "id": "orders-to-all-orders",
    "topic_id": "orders",
    "queue_id": "all-orders",
    "message_filters": null,
    "created_at": "2024-01-02T22:30:12.628323Z"
}
```

Now creating the second subscription, for this one we will use the message_filters field:

```bash
curl --location 'http://localhost:8000/v1/subscriptions' \
--header 'Content-Type: application/json' \
--data '{
    "id": "orders-to-processed-orders",
    "topic_id": "orders",
    "queue_id": "processed-orders",
    "message_filters": {"status": ["processed"]}
}'
```

```json
{
    "id": "orders-to-processed-orders",
    "topic_id": "orders",
    "queue_id": "processed-orders",
    "message_filters": {
        "status": [
            "processed"
        ]
    },
    "created_at": "2024-01-02T22:31:26.156692Z"
}
```

Now it's time to publish the first message:

```bash
curl --location 'http://localhost:8000/v1/topics/orders/messages' \
--header 'Content-Type: application/json' \
--data '{
    "body": "body-of-the-order",
    "attributes": {"status": "created"}
}'
```

And the second message:

```bash
curl --location 'http://localhost:8000/v1/topics/orders/messages' \
--header 'Content-Type: application/json' \
--data '{
    "body": "body-of-the-order",
    "attributes": {"status": "processed"}
}'
```

Now we will consume the `all-orders` queue:

```bash
curl --location 'http://localhost:8000/v1/queues/all-orders/messages'
```

```json
{
    "data": [
        {
            "id": "01HK651Q52EZMPKBYZGVK0ZX8S",
            "queue_id": "all-orders",
            "label": null,
            "body": "body-of-the-order",
            "attributes": {
                "status": "created"
            },
            "delivery_attempts": 1,
            "created_at": "2024-01-02T19:35:00.635625-03:00"
        },
        {
            "id": "01HK652W2HNW53XWV4QBT5MAJY",
            "queue_id": "all-orders",
            "label": null,
            "body": "body-of-the-order",
            "attributes": {
                "status": "processed"
            },
            "delivery_attempts": 1,
            "created_at": "2024-01-02T19:35:38.446759-03:00"
        }
    ],
    "limit": 10
}
```

As expected, this queue has the two published messages.

Now we will consume the `processed-orders` queue:

```bash
curl --location 'http://localhost:8000/v1/queues/processed-orders/messages'
```

```json
{
    "data": [
        {
            "id": "01HK652W2JK8MPN3JDXY9RATS5",
            "queue_id": "processed-orders",
            "label": null,
            "body": "body-of-the-order",
            "attributes": {
                "status": "processed"
            },
            "delivery_attempts": 1,
            "created_at": "2024-01-02T19:35:38.446759-03:00"
        }
    ],
    "limit": 10
}
```

As expected, this queue has only one message that was published with the `status` attribute equal to `"processed"`.
