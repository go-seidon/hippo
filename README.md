# local-storage

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=go-seidon_local&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=go-seidon_local)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=go-seidon_local&metric=coverage)](https://sonarcloud.io/summary/new_code?id=go-seidon_local)

## Doc
No doc right now

## Todo
1. Refactor exec DeleteFile using Exec args
2. Uploading rest implementation
3. Change `file.id` to varchar and remove `file.unique_id` (id should be geenrated by the client app)

## Nice to have
1. Separate findFile query in DeleteFile and RetrieveFile
2. File meta for storing file related data, e.g: user_id, feature, category, etc
3. Store directory checking result in memory when uploading file to reduce r/w to the disk (dirManager)

## Technical Stack
1. Transport layer
- rest
- grpc
2. Database
- mysql
- postgres
- mongo
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
  $ run-rest-app
  $ build-rest-app
```

2. GRPC App

```
  $ run-grpc-app
  $ build-grpc-app
```

3. Hybrid App

```
  TBA
```

### Development
1. Create docker compose
```
  $ docker-compose up -d
```

2. MySQL database: 
```
 $ docker-compose up mysql-database
 $ docker-compose stop mysql-database
```

3. MySQL Test database:
```
 $ docker-compose up mysql-test-database
 $ docker-compose stop mysql-test-database
```

4. MySQL Migration
```bash
  $ migrate-mysql-create [args] # args e.g: migrate-mysql-create file-table
  $ migrate-mysql [args] # args e.g: migrate-mysql up
```
