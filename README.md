#Front End Masters Project: Go for Pressional Developers

This project builds an API with authentication middleware and postgres database, to allow users to login and create a workout schedule (including view, update, and delete). There is no UI component and all actions are carried out using curl requests, or as in my case, via Postman.

The project also features database build and migration scripts that show how these evolve from stage to stage, including adding a column after the table already exists.

## My Experience

I found the structure that this project evolved to very similar to the NodeJS projects I have worked on in the past in a professional capacity, and it was interesting to see the evolving bespoke code using the core GO packages, but also seeing the few third party packages that were included, such as Chi.

## Useful Commands

Run project in default localhost 8080

```
go run main.go
```

Or, use another port

```
go run main.go port 8081
```

Run docker with postgres databases

```
docker compose up
```

Command line to access database

```
psql -U postgres -h localhost -p 5432
```

Run database migrations

```
goose -dir migrations postgres "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" up
```
