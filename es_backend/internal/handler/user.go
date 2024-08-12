package handler

import (
    "es_backend/api/v1"
    "es_backend/internal/service"
    "es_backend/pkg/helper"
    "log"
    "net/http"

    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

type UserHandler struct {
    *Handler
    userService service.UserService
}

func NewUserHandler(handler *Handler, userService service.UserService) *UserHandler {
    return &UserHandler{
        Handler:     handler,
        userService: userService,
    }
}

// Register godoc
// @Summary 用户注册
// @Schemes
// @Description 目前只支持邮箱注册
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param request body v1.RegisterRequest true "params"
// @Success 200 {object} v1.BaseResponse
// @Router /register [post]
func (h *UserHandler) Register(ctx *gin.Context) {
    req := new(v1.RegisterRequest)
    if err := ctx.ShouldBindJSON(req); err != nil {
        v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
        return
    }
    log.Println(req)
    if !helper.VerifyEmail(req.Email) {
        v1.HandleError(ctx, http.StatusBadRequest, v1.ErrEmailFormat, nil)
    }
    log.Println("---------------------")
    if err := h.userService.Register(ctx, req); err != nil {
        h.logger.WithContext(ctx).Error("userService.Register error", zap.Error(err))
        v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
        return
    }
    log.Println("+++++++++++++++++++")
    v1.HandleSuccess(ctx, nil)
}

// Login godoc
// @Summary 账号登录
// @Schemes
// @Description
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param request body v1.LoginRequest true "params"
// @Success 200 {object} v1.LoginResponse
// @Router /login [post]
func (h *UserHandler) Login(ctx *gin.Context) {
    req := new(v1.LoginRequest)
    resp := new(v1.LoginResponse)
    if err := ctx.ShouldBindJSON(req); err != nil {
        v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
        return
    }

    err := h.userService.Login(ctx, req, resp)
    if err != nil {
        v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
        return
    }
    v1.HandleSuccess(ctx, resp)
}

// GetProfile godoc
// @Summary 获取用户信息
// @Schemes
// @Description
// @Tags 用户模块
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} v1.GetProfileResponse
// @Router /user [get]
func (h *UserHandler) GetProfile(ctx *gin.Context) {
    uuid := GetUUIDFromCtx(ctx)
    if uuid == 0 {
        v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
        return
    }
    var req = &v1.GetProfileRequest{
        UUID: uuid,
    }
    resp := new(v1.GetProfileResponse)
    err := h.userService.GetProfile(ctx, req, resp)
    if err != nil {
        v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
        return
    }

    v1.HandleSuccess(ctx, resp)
}

// UpdateProfile godoc
// @Summary 修改用户信息
// @Schemes
// @Description
// @Tags 用户模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.UpdateProfileRequest true "params"
// @Success 200 {object} v1.BaseResponse
// @Router /user [put]
func (h *UserHandler) UpdateProfile(ctx *gin.Context) {
    uuid := GetUUIDFromCtx(ctx)
    if uuid == 0 {
        v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
        return
    }
    var req *v1.UpdateProfileRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
        return
    }
    req.UUID = uuid
    if err := h.userService.UpdateProfile(ctx, req); err != nil {
        v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
        return
    }

    v1.HandleSuccess(ctx, nil)
}

// List godoc
// @Summary 列出用户信息
// @Schemes
// @Description 列出用户信息
// @Tags 用户模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.ListUserRequest true "params"
// @Success 200 {object} v1.ListUserResponse
// @Router /v1/user/list [post]
func (h *UserHandler) List(ctx *gin.Context) {
    req := new(v1.ListUserRequest)
    resp := new(v1.ListUserResponse)
    if err := ctx.ShouldBindJSON(&req); err != nil {
        v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
        return
    }

    err := h.userService.List(ctx, req, resp)
    if err != nil {
        v1.HandleError(ctx, http.StatusOK, err, nil)
        return
    }
    ctx.JSON(http.StatusOK, resp)
}
