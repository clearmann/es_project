package handler

import (
    "es_backend/pkg/jwt"
    "es_backend/pkg/log"

    "github.com/gin-gonic/gin"
)

type Handler struct {
    logger *log.Logger
}

func NewHandler(logger *log.Logger) *Handler {
    return &Handler{
        logger: logger,
    }
}
func GetUUIDFromCtx(ctx *gin.Context) uint64 {
    v, exists := ctx.Get("claims")
    if !exists {
        return 0
    }
    return v.(*jwt.MyCustomClaims).UUID
}
