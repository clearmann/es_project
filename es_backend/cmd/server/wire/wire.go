//go:build wireinject
// +build wireinject

package wire

import (
    "es_backend/internal/handler"
    "es_backend/internal/repository"
    "es_backend/internal/server"
    "es_backend/internal/service"
    "es_backend/pkg/app"
    "es_backend/pkg/jwt"
    "es_backend/pkg/log"
    "es_backend/pkg/server/http"
    "es_backend/pkg/sid"

    "github.com/google/wire"
    "github.com/spf13/viper"
)

var repositorySet = wire.NewSet(
    repository.NewDB,
    repository.NewRedis,
    repository.NewElasticSearch,
    repository.NewRepository,
    repository.NewTransaction,
    repository.NewUserRepository,
    repository.NewPostRepository,
)

var serviceSet = wire.NewSet(
    service.NewService,
    service.NewUserService,
    service.NewPostService,
)

var handlerSet = wire.NewSet(
    handler.NewHandler,
    handler.NewUserHandler,
    handler.NewPostHandler,
)

var serverSet = wire.NewSet(
    server.NewHTTPServer,
    server.NewJob,
)

// build App
func newApp(httpServer *http.Server, job *server.Job) *app.App {
    return app.NewApp(
        app.WithServer(httpServer, job),
        app.WithName("demo-server"),
    )
}

func NewWire(*viper.Viper, *log.Logger) (*app.App, func(), error) {
    panic(wire.Build(
        repositorySet,
        serviceSet,
        handlerSet,
        serverSet,
        sid.NewSid,
        jwt.NewJwt,
        newApp,
    ))
}
