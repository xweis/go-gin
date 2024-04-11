An example of gin contains many useful features

## Installation
```
$ go get github.com/xweis/go-gin
```

## How to run

### Required

- Mysql
- Redis

### Conf

You should modify `conf/app.ini`

```
[database]
Type = mysql
User = root
Password =
Host = 127.0.0.1:3306
Name = blog
TablePrefix = blog_

[redis]
Host = 127.0.0.1:6379
Password =
MaxIdle = 30
MaxActive = 30
IdleTimeout = 200
...
```

### Run
```
$ cd $GOPATH/src/go-gin

$ go run main.go 
```

Project information and existing API

```
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /auth                     --> github.com/EDDYCJY/go-gin-example/routers/api.GetAuth (3 handlers)
[GIN-debug] GET    /swagger/*any             --> github.com/EDDYCJY/go-gin-example/vendor/github.com/swaggo/gin-swagger.WrapHandler.func1 (3 handlers)
[GIN-debug] GET    /api/v1/tags              --> github.com/EDDYCJY/go-gin-example/routers/api/v1.GetTags (4 handlers)

Listening port is 8000
Actual pid is 4393
```
Swagger doc

## Features

- RESTful API
- Gorm
- Swagger
- zap
- Jwt-go
- Gin
- App configurable
- Cron
- Redis
