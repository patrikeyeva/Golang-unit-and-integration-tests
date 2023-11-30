//go:generate mockgen -source ./repository.go -destination=./mocks/repository.go -package=mock_repository
package repository

import (
	"context"
)

type ArticlesRepo interface {
	Add(context.Context, *Article) (*Article, error)
	GetByID(context.Context, int64) (*Article, error)
	DeleteByID(context.Context, int64) error
	Update(context.Context, *Article) error
}

type CommentsRepo interface {
	AddComment(context.Context, *Comment) (*Comment, error)
	GetCommentsForArticle(context.Context, int64) ([]Comment, error)
}
