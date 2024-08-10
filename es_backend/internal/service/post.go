package service

import (
    "context"
    v1 "es_backend/api/v1"
    "es_backend/internal/model/model_type"
    "es_backend/internal/repository"
)

type PostService interface {
    Create(ctx context.Context, req *v1.CreatePostRequest, resp *v1.BaseResponse) error
    Delete(ctx context.Context, req *v1.DeletePostRequest, resp *v1.BaseResponse) error
    Update(ctx context.Context, req *v1.UpdatePostRequest, resp *v1.BaseResponse) error
    List(ctx context.Context, req *v1.ListPostRequest, resp *v1.ListPostResponse) error
}

func NewPostService(
    service *Service,
    postRepo repository.PostRepository,
) PostService {
    return &postService{
        postRepo: postRepo,
        Service:  service,
    }
}

type postService struct {
    postRepo repository.PostRepository
    *Service
}

func (s *postService) Create(ctx context.Context, req *v1.CreatePostRequest, resp *v1.BaseResponse) error {
    var post = &model_type.Post{
        Title:   req.Title,
        Content: req.Content,
        Tags:    req.Tags,
        UUID:    req.UUID,
    }
    if err := s.postRepo.Create(ctx, post); err != nil {
        s.logger.Error("create post error")
        return v1.ErrInternalServerError
    }
    return nil
}

func (s *postService) Delete(ctx context.Context, req *v1.DeletePostRequest, resp *v1.BaseResponse) error {
    return nil
}

func (s *postService) Update(ctx context.Context, req *v1.UpdatePostRequest, resp *v1.BaseResponse) error {
    return nil
}
func (s *postService) List(ctx context.Context, req *v1.ListPostRequest, resp *v1.ListPostResponse) error {
    var err error
    if req.ListAll {
        err = s.postRepo.ListALL(ctx, &resp.Data)
    } else {
        err = s.postRepo.List(ctx, req, &resp.Data)
    }
    if err != nil {
        return v1.ErrInternalServerError
    }
    resp.Total = int64(len(resp.Data))
    resp.Offset = req.Offset
    resp.Limit = req.Limit
    return nil
}
