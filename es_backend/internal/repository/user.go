package repository

import (
    "context"
    "errors"
    v1 "es_backend/api/v1"
    "es_backend/internal/model/model_type"

    "gorm.io/gorm"
)

type UserRepository interface {
    Create(ctx context.Context, user *model_type.User) error
    Update(ctx context.Context, req *v1.UpdateProfileRequest) error
    GetByID(ctx context.Context, uuid uint64) (*model_type.User, error)
    GetByEmail(ctx context.Context, email string) (*model_type.User, error)
    GetByEmailORUsername(ctx context.Context, name string) (*model_type.User, error)
    ExistUserByUUID(ctx context.Context, uuid uint64) (bool, error)
    ExistUserByEmail(ctx context.Context, email string) (bool, error)
    ExistUserByUsername(ctx context.Context, username string) (bool, error)
    ListALL(ctx context.Context, posts *[]*v1.UserInfo) error
    List(ctx context.Context, req *v1.ListUserRequest, posts *[]*v1.UserInfo) error
}

func NewUserRepository(r *Repository) UserRepository {
    return &userRepository{Repository: r}
}

type userRepository struct {
    *Repository
}

func (r *userRepository) Create(ctx context.Context, user *model_type.User) error {
    err := r.DB(ctx).
        Table(model_type.TableNameUser).
        Create(user).Error
    if err != nil {
        return err
    }
    return nil
}

func (r *userRepository) Update(ctx context.Context, req *v1.UpdateProfileRequest) error {
    if err := r.DB(ctx).
        Table(model_type.TableNameUser).
        Where("uuid = ?", req.UUID).
        Updates(&model_type.User{Username: req.Username, Email: req.Email, Profile: req.Profile}).
        Error; err != nil {
        return err
    }
    return nil
}

func (r *userRepository) GetByID(ctx context.Context, uuid uint64) (*model_type.User, error) {
    var user model_type.User
    if err := r.DB(ctx).Table(model_type.TableNameUser).Where("uuid = ?", uuid).First(&user).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, v1.ErrNotFound
        }
        return nil, err
    }
    return &user, nil
}
func (r *userRepository) ExistUserByEmail(ctx context.Context, email string) (bool, error) {
    total := new(int64)
    err := r.DB(ctx).Table(model_type.TableNameUser).
        Where("email = ?", email).
        Count(total).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return false, nil
        }
        return false, err
    }
    return *total > 0, nil
}
func (r *userRepository) ExistUserByUsername(ctx context.Context, username string) (bool, error) {
    total := new(int64)
    err := r.DB(ctx).
        Table(model_type.TableNameUser).
        Where("username = ?", username).
        Count(total).Error

    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return false, nil
        }
        return false, err
    }
    return *total > 0, nil
}
func (r *userRepository) ExistUserByUUID(ctx context.Context, uuid uint64) (bool, error) {
    total := new(int64)
    if err := r.DB(ctx).Table(model_type.TableNameUser).Where("uuid = ?", uuid).Count(total).Error; err != nil {
        return false, err
    }
    return *total > 0, nil
}
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model_type.User, error) {
    user := new(model_type.User)
    if err := r.DB(ctx).Table(model_type.TableNameUser).Where("email = ?", email).First(user).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return user, nil
}
func (r *userRepository) GetByEmailORUsername(ctx context.Context, name string) (*model_type.User, error) {
    user := new(model_type.User)
    err := r.DB(ctx).Table(model_type.TableNameUser).
        Where("email = ? or username", name, name).
        First(user).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return user, nil
}
func (r *userRepository) ListALL(ctx context.Context, users *[]*v1.UserInfo) error {
    if err := r.DB(ctx).Table(model_type.TableNameUser).Find(users).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return v1.ErrNotFound
        }
        return err
    }
    return nil
}
func (r *userRepository) List(ctx context.Context, req *v1.ListUserRequest, users *[]*v1.UserInfo) error {
    if err := r.DB(ctx).Table(model_type.TableNameUser).
        Offset(req.Offset).
        Limit(req.Limit).
        Find(users).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return v1.ErrNotFound
        }
        return err
    }
    return nil
}
