# Whostodo

A todo app written in go. Session expires in 1 minute.

## Getting Started

To run the app locally, go v1.22.1 assumed:
```shell
go mod download
go build
./whostodo

# in another terminal
curl -X POST localhost:8080/v1/auth

# {"result":"470edaa57be62fc5d884aa1625aa5c0c"}
```

## REST Endpoints

### `POST /v1/auth`

Authenticates the user. Can be used to check session validity if token provided in header.

#### Initiates a new session; returns 201

```shell
curl -X -H localhost:8080/v1/auth
```

```json
{ "result": "1b030c7dd4a211d5897750cb725f1cd6" }
```

#### Confirms current session as valid; returns 304

```shell
# replace `YOUR_TOKEN` to actual value
curl -X POST -H 'Authorization: Bearer YOUR_TOKEN' localhost:8080/v1/auth
```

#### Deems current session expired and returns new session token; returns 304

```shell
# replace `YOUR_TOKEN` to actual value
curl -H 'Authorization: Bearer YOUR_TOKEN' localhost:8080/v1/tasks
```

```json
{ "result": "7810b2d06543ddee7d17ef230f13d2b7" }
```

### `GET /v1/tasks`

Lists task items.

```shell
# replace `YOUR_TOKEN` to actual value
curl -H 'Authorization: Bearer YOUR_TOKEN' localhost:8080/v1/tasks
```

```json
{
    "result": [
        {
            "id": 1,
            "name": "name",
            "status": 0
        }
    ]
}
```

### `POST /v1/task`

Creates a new task item.

```shell
# replace `YOUR_TOKEN` to actual value
# replace `TASK_NAME` to actual value
curl -X POST -H 'Content-type: application/json' -H 'Authorization: Bearer YOUR_TOKEN' -d '{"name":"TASK_NAME"}' localhost:8080/v1/task
```

```json
{
    "result": {
        "name": "name",
        "status": 0
        "id": 1,
    }
}
```

### `PUT /v1/task/:id`

Updates an existing task item.

#### Updates the task item; returns 201

```shell
# replace `YOUR_TOKEN` to actual value
# replace `YOUR_TOKEN` to actual value
# replace `YOUR_STATUS` to actual value, 0 or 1
# replace `TASK_ID` to actual value
curl -X PUT -H 'Content-type: application/json' -H 'Authorization: Bearer YOUR_TOKEN' -d '{"name":"TASK_NAME","status":TASK_STATUS}' localhost:8080/v1/task/TASK_ID
```

```json
{
    "result": {
        "name": "new name",
        "status": 1
        "id": 1,
    }
}
```

#### Fails to locate the task item; returns 404

```shell
# replace `YOUR_TOKEN` to actual value
# replace `YOUR_TOKEN` to actual value
# replace `TASK_STATUS` to actual value, 0 or 1
# replace `TASK_ID` to actual value
curl -X PUT -H 'Content-type: application/json' 'Authorization: Bearer YOUR_TOKEN' -d '{"name":"TASK_NAME","status":TASK_STATUS}' localhost:8080/v1/task/TASK_ID
```

```json
{
    "result": {}
}
```

### `DELETE /v1/task/:id`

Deletes an existing task item.

#### Deletes the task item; returns 200

```shell
# replace `YOUR_TOKEN` to actual value
# replace `TASK_ID` to actual value
curl -X DELETE -H 'Authorization: Bearer YOUR_TOKEN' localhost:8080/v1/task/TASK_ID
```

#### Fails to locate the task item; returns 404

```shell
# replace `YOUR_TOKEN` to actual value
# replace `TASK_ID` to actual value
curl -X DELETE -H 'Authorization: Bearer YOUR_TOKEN' localhost:8080/v1/task/TASK_ID
```

## Development

Under project directory:
```shell
docker build --target dev -t whostodo:dev .
docker run --rm -v $PWD:/app whostodo:dev
```

To run the test suite:
```shell
go test ./...
```

## Gotchas

- **Thread safety is not assumed**
- App state is persisted in memory, i.e., all states are gone when app restarts

### Session

- Sessions are not deleted, as intended, for possible audit purposes
- Token generator implementation is not secure

### Task

- Task status can be assigned to an arbitrary integer value
