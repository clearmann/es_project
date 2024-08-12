package server

import (
    "es_backend/docs"
    "es_backend/internal/handler"
    "es_backend/internal/middleware"
    "es_backend/pkg/jwt"
    "es_backend/pkg/log"
    "es_backend/pkg/server/http"

    "github.com/gin-gonic/gin"
    "github.com/spf13/viper"
    swaggerfiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
)

func NewHTTPServer(
    logger *log.Logger,
    conf *viper.Viper,
    jwt *jwt.JWT,
    userHandler *handler.UserHandler,
    postHandler *handler.PostHandler,
) *http.Server {
    gin.SetMode(gin.DebugMode)
    s := http.NewServer(
        gin.Default(),
        logger,
        http.WithServerHost(conf.GetString("http.host")),
        http.WithServerPort(conf.GetInt("http.port")),
    )

    // swagger doc
    docs.SwaggerInfo.BasePath = "/v1"
    s.GET("/swagger/*any", ginSwagger.WrapHandler(
        swaggerfiles.Handler,
        // ginSwagger.URL(fmt.Sprintf("http://localhost:%d/swagger/doc.json", conf.GetInt("app.http.port"))),
        ginSwagger.DefaultModelsExpandDepth(-1),
        ginSwagger.PersistAuthorization(true),
    ))

    s.Use(
        middleware.CORSMiddleware(),
        middleware.ResponseLogMiddleware(logger),
        middleware.RequestLogMiddleware(logger),
        // middleware.SignMiddleware(log),
    )

    v1 := s.Group("/v1")
    {
        // No route group has permission
        noAuthRouter := v1.Group("/")
        {
            noAuthRouter.POST("/register", userHandler.Register)
            noAuthRouter.POST("/login", userHandler.Login)
            noAuthRouter.POST("/post/list", postHandler.List)
            noAuthRouter.POST("/user/list", userHandler.List)
        }
        // Non-strict permission routing group
        noStrictAuthRouter := v1.Group("/").Use(middleware.NoStrictAuth(jwt, logger))
        {
            noStrictAuthRouter.GET("/user", userHandler.GetProfile)
            noAuthRouter.POST("/post/create", postHandler.Create)
        }

        // Strict permission routing group
        strictAuthRouter := v1.Group("/").Use(middleware.StrictAuth(jwt, logger))
        {
            strictAuthRouter.PUT("/user", userHandler.UpdateProfile)
        }
    }

    return s
}
