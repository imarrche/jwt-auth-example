# JWT

This is just an example of how JWT authorization could be implemented for REST API.

1. POST `api/v1/auth/sign-up` - for signing up.
* email, username, first_name, second_name, password are required.
* email, username should be unique.

2. POST `api/v1/auth/sign-in` - to get token pair(access and refresh JWTs).
* email and password must be provided.

3. POST `api/v1/auth/refresh` - to refresh access token.
* refresh token must be provided.

4. GET `api/v1/public` - just a public endpoint that always returns `HTTP 200 OK`.
5. GET `api/v1/private` - private endpoint which returns `HTTP 401 UNAUTHORIZED` if user is not authorized.
* To be actually authorized, header `Authorization` must be set in `Bearer <token>` format.


## Run instructions

1) Create `.env` file for server configuration. For example:
```bash
SERVER_ADDR=:8080
POSTGRES_HOST=postgres
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=123
POSTGRES_DBNAME=jwt
POSTGRES_SSLMODE=disable
```

2) Spin up `postgres` container.
```bash
$ docker-compose up postgres
```

3) Create database.
```bash
$ docker exec -it postgres bash
$ psql -U postgres
$ CREATE DATABASE jwt;
```

4) Spin up `jwtserver` container.
```bash
$ docker-compose up jwtserver
```
