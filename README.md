# local-storage

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=go-seidon_local&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=go-seidon_local)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=go-seidon_local&metric=coverage)](https://sonarcloud.io/summary/new_code?id=go-seidon_local)

## Doc
No doc right now

## Todo
1. Deploy dev
2. Add deployment script
3. Add `repository-mongo` implementation
4. Add `grpc-app` implementation

## Nice to have
1. Separate findFile query in DeleteFile and RetrieveFile
2. File meta for storing file related data, e.g: user_id, feature, category, etc
3. Store directory checking result in memory when uploading file to reduce r/w to the disk (dirManager)
4. File setting: (visibility, upload location default to daily rotator)
5. Access file using custom link with certain limitation such as access duration, attribute user_id, etc
6. Change NewDailyRotate using optional param
7. Resize image capability (?width=720&height=480)
8. Add `repository-postgre` implementation
9. Inject logger to mysql instance
10. Add `logging.WithReqCtx(ctx)` to parse `correlationId`

## Technical Stack
1. Transport layer
- rest
- grpc (TBA)
2. Database
- mysql
- mongo (TBA)
- postgre (TBA)
3. Config
- system environment
- file (config/*.toml and .env)

## Run
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
  $ make run-rest-app
  $ make build-rest-app
```

2. GRPC App

```
  $ make run-grpc-app
  $ make build-grpc-app
```

3. Hybrid App

```
  $ make run-hybrid-app
  $ make build-hybrid-app
```

### Docker
1. Build docker image
```
  $ docker build -t goseidon-local .
```

2. Check build result
```
  $ docker images
```

3. Create docker container
```
  $ docker container create --name goseidon-local-app ^
    -e REST_APP_HOST="0.0.0.0" ^
    -e REST_APP_PORT=3000 ^
    -e RPC_APP_HOST="0.0.0.0" ^
    -e RPC_APP_PORT=5000 ^
    -e MYSQL_MASTER_HOST="host.docker.internal" ^
    -e MYSQL_REPLICA_HOST="host.docker.internal" ^
    -p 3000:3000 -p 5000:5000 ^
    -v D:\startup\goseidon\local\storage:/storage ^
    goseidon-local
```

4. Check container
```
  $ docker container ls -a
```

5. Start container
```
  $ docker container start goseidon-local-app
```

6. Check container status
```
  $ docker container ls
```

### Development
#### First time setup
1. Copy `.env.example` to `.env`

2. Create docker compose
```
  $ docker-compose up -d
```

3. Setup MySQL replication
```
  $ ./development/mysql/replication.sh
```

#### Running database migration
1. MySQL Migration
```bash
  $ migrate-mysql-create [args] # args e.g: migrate-mysql-create file-table
  $ migrate-mysql [args] # args e.g: migrate-mysql up
```
