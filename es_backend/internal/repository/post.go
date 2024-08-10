package repository

import (
    "context"
    "errors"
    v1 "es_backend/api/v1"
    "es_backend/internal/model/model_type"

    "gorm.io/gorm"
)

type PostRepository interface {
    Create(ctx context.Context, post *model_type.Post) error
    Update(ctx context.Context, post *model_type.Post) error
    Delete(ctx context.Context, postID uint) error
    GetPostByID(ctx context.Context, postID uint) (*model_type.Post, error)
    ListALL(ctx context.Context, posts *[]*model_type.Post) error
    List(ctx context.Context, req *v1.ListPostRequest, posts *[]*model_type.Post) error
}

func NewPostRepository(r *Repository) PostRepository {
    return &postRepository{Repository: r}
}

type postRepository struct {
    *Repository
}

func (r *postRepository) Create(ctx context.Context, post *model_type.Post) error {
    if err := r.DB(ctx).Create(post).Error; err != nil {
        return err
    }
    return nil
}

func (r *postRepository) Update(ctx context.Context, post *model_type.Post) error {
    if err := r.DB(ctx).Save(post).Error; err != nil {
        return err
    }
    return nil
}

func (r *postRepository) GetPostByID(ctx context.Context, postID uint) (*model_type.Post, error) {
    var post model_type.Post
    if err := r.DB(ctx).Where("id = ?", postID).First(&post).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, v1.ErrNotFound
        }
        return nil, err
    }
    return &post, nil
}
func (r *postRepository) Delete(ctx context.Context, postID uint) error {
    if err := r.DB(ctx).Where("id = ?", postID).Delete(&model_type.Post{}).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return v1.ErrNotFound
        }
        return err
    }
    return nil
}
func (r *postRepository) ListALL(ctx context.Context, posts *[]*model_type.Post) error {
    if err := r.DB(ctx).Table(model_type.TableNamePost).Find(posts).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return v1.ErrNotFound
        }
        return err
    }
    return nil
}
func (r *postRepository) List(ctx context.Context, req *v1.ListPostRequest, posts *[]*model_type.Post) error {
    if err := r.DB(ctx).Table(model_type.TableNamePost).
        Offset(req.Offset).
        Limit(req.Limit).
        Find(posts).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return v1.ErrNotFound
        }
        return err
    }
    return nil
}
