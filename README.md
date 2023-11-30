# EmailN

## Description

A Rest API for email marketing using Keycloak for authentication, Postgres as database, Tests, GoLang with go routines and workers to send emails async.

## Endpoints

_All endpoints need authentication_

**POST**

```
/campaigns -> Create a new campaign
```

**GET**

```
/campaigns/{campaignId} -> Show campaign details
```

**DELETE**

```
/campaigns/delete/{campaignId} -> Delete campaign
```

**PATCH**

```
/campaigns/start/{campaignId} -> Start campaign (send emails)
```

## How to run

- run `docker compose up -d`
- go to **localhost:8080** to configure the Keycloak
- access **postgres container** to create a database
- configure yours **.env variables**
- run `go run cmd/api/main.go` and `go run cmd/worker/main.go` to start the server and worker
- explore the routes
