package service

import (
    "es_backend/internal/repository"
    "es_backend/pkg/jwt"
    "es_backend/pkg/log"
    "es_backend/pkg/sid"
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
