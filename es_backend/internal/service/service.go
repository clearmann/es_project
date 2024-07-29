package service

import (
    "gin-init/internal/repository"
    "gin-init/pkg/jwt"
    "gin-init/pkg/log"
    "gin-init/pkg/sid"
)

type Service struct {
    logger *log.Logger
    sid    *sid.Sid
    jwt    *jwt.JWT
    tm     repository.Transaction
}

func NewService(
    tm repository.Transaction,
    logger *log.Logger,
    sid *sid.Sid,
    jwt *jwt.JWT,
) *Service {
    return &Service{
        logger: logger,
        sid:    sid,
        jwt:    jwt,
        tm:     tm,
    }
}
