# Hippo

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=go-seidon_hippo&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=go-seidon_hippo)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=go-seidon_hippo&metric=coverage)](https://sonarcloud.io/summary/new_code?id=go-seidon_hippo)

Bucket like local storage implementation

## Technical Stack
1. Transport layer
- rest
- grpc
2. Database
- mysql
- mongo
3. Config
- system environment
- file (config/*.toml and .env)

## How to Run?
### Test
1. Unit test

This particular command should test individual component and run really fast without the need of involving 3rd party dependencies such as database, disk, etc.

```
  $ make test-unit
  $ make test-watch-unit
```

2. Integration test

This particular command should test the integration between component, might run slowly and sometimes need to involving 3rd party dependencies such as database, disk, etc.

```
  $ make test-integration
  $ make test-watch-integration
```

3. Coverage test

This command should run all the test available on this project.

```
  $ make test
  $ make test-coverage
```

### App
1. REST App

```
  $ make run-restapp
  $ make build-restapp
```

2. GRPC App

```
  $ make run-grpcapp
  $ make build-grpcapp
```

3. Hybrid App

```
  $ make run-hybridapp
  $ make build-hybridapp
```

### Docker
1. Build docker image
```
  $ docker build -t hippo .
```

2. Check build result
```
  $ docker images
```

3. Create docker container
```
  $ docker container create --name hippo-app ^
    -e REST_APP_HOST="0.0.0.0" ^
    -e REST_APP_PORT=3000 ^
    -e GRPC_APP_HOST="0.0.0.0" ^
    -e GRPC_APP_PORT=5000 ^
    -e MYSQL_MASTER_HOST="host.docker.internal" ^
    -e MYSQL_REPLICA_HOST="host.docker.internal" ^
    -p 3000:3000 -p 5000:5000 ^
    -v storage:/storage ^
    hippo
```

4. Check container
```
  $ docker container ls -a
```

5. Start container
```
  $ docker container start hippo-app
```

6. Check container status
```
  $ docker container ls
```

## Development
### First time setup
1. Copy `.env.example` to `.env`

2. Create docker compose
```bash
  $ docker-compose up -d
```

### Database migration
1. MySQL Migration
```bash
  $ make migrate-mysql-create [args] # args e.g: migrate-mysql-create file-table
  $ make migrate-mysql [args] # args e.g: migrate-mysql up
```

2. Mongo Migration
```bash
  $ make migrate-mongo-create [args] # args e.g: migrate-mongo-create file-table
  $ make migrate-mongo [args] # args e.g: migrate-mongo up
```

### MySQL Replication Setup
1. Run setup
```bash
  $ ./development/mysql/replication.sh
```

### MongoDB Replication Setup
1. Generate keyFile (if necessary)
```bash
  $ cd /development/mongo
  $ openssl rand -base64 741 > mongodb.key
  $ chmod 400 mongodb.key
```

2. Setting local hosts
- Window
C:\Windows\System32\drivers\etc\hosts
```md
  127.0.0.1 mongo-db-1
  127.0.0.1 mongo-db-2
  127.0.0.1 mongo-db-3
```

- Linux
\etc\hosts
```md
  127.0.0.1 mongo-db-1
  127.0.0.1 mongo-db-2
  127.0.0.1 mongo-db-3
```

3. Run setup
```bash
  $ ./development/mongo/replication.sh
```

## Todo
1. Use gorm for mysql repository
2. Replace auth by using chariot auth
3. Upload docker image to docker hub

## Nice to have
1. Simplify factory function interface not returning error
2. Upload location strategy
3. Add repo: `postgre`
4. Add tracing: `logging.WithReqCtx(ctx)` to parse `correlationId`
5. Update github workflow (cqc.yml) instead of running docker-compose prefer to use mongo docker services
6. Separate unit test and integration test workflow (cqc.yml)

## Tech Debt
1. Separate findFile query in DeleteFile and RetrieveFile repository/mysql
2. Store directory checking result in memory when uploading file to reduce r/w to the disk (dirManager)
3. Change NewDailyRotate using optional param

## Issue
1. Verify script EOL
If you're using window make sure to change the bash script from CRLF to LF (https://stackoverflow.com/questions/29140377/sh-file-not-found)

