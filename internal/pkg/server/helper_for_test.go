package server

import (
	mock_repository "homework5/internal/pkg/repository/mocks"
	"testing"

	"github.com/golang/mock/gomock"
)

type serverFixture struct {
	ctrl      *gomock.Controller
	mArticle  *mock_repository.MockArticlesRepo
	mComments *mock_repository.MockCommentsRepo
	server    Server
}

func setUpArticle(t *testing.T) serverFixture {
	ctrl := gomock.NewController(t)
	mArticle := mock_repository.NewMockArticlesRepo(ctrl)
	s := Server{ArticleRepo: mArticle}

	return serverFixture{
		ctrl:     ctrl,
		mArticle: mArticle,
		server:   s,
	}
}

func setUpComment(t *testing.T) serverFixture {
	ctrl := gomock.NewController(t)
	mComments := mock_repository.NewMockCommentsRepo(ctrl)
	s := Server{CommentRepo: mComments}

	return serverFixture{
		ctrl:      ctrl,
		mComments: mComments,
		server:    s,
	}
}

func (a *serverFixture) tearDown() {
	a.ctrl.Finish()
}
