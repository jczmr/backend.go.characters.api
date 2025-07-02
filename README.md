# backend.go.characters.api
Golang Microservice to get information about Dragon Ball characters

# Dragon Ball Character Service

This Go service allows you to create and store Dragon Ball character information, leveraging a hexagonal architecture, `log/slog` for logging, Gin for the web framework, and PostgreSQL for persistence. It interacts with an external Dragon Ball API to fetch character details.

---

## 1. Architecture Diagram

Here's a sequence diagram illustrating the flow when a `POST /characters` request is made:

![sequence diagram for API](docs/sequence-diagram.png)


## 2. How to run this


- add an .env file with this environment variables, refer to .env.example file in the project

```
PORT=8080
DB_USER=user
DB_PASSWORD=password
DB_NAME=dragonballdb
DB_HOST=db
DB_PORT=5432
```

