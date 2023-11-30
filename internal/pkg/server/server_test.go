package server

import (
	"context"
	"errors"
	"homework5/internal/pkg/repository"
	"homework5/tests/fixtures"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_GetArticle(t *testing.T) {

	var (
		ctx context.Context
		id  = 1
	)

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		s := setUpArticle(t)
		defer s.tearDown()
		s.mArticle.EXPECT().GetByID(gomock.Any(), int64(id)).Return(fixtures.Article().AllData().P(), nil)

		data, status, err := s.server.GetArticle(ctx, int64(id))
		require.Equal(t, http.StatusOK, status)
		require.Nil(t, err)
		assert.Equal(
			t,
			`{"ID":1,"Name":"asd","Rating":10,"CreatedAt":"0001-01-01T00:00:00Z"}`,
			string(data))
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()
		s := setUpArticle(t)
		defer s.tearDown()
		s.mArticle.EXPECT().GetByID(gomock.Any(), int64(id)).Return(nil, repository.ErrObjectNotFound)

		data, status, err := s.server.GetArticle(ctx, int64(id))
		require.Equal(t, http.StatusNotFound, status)
		require.NotNil(t, err)
		require.Nil(t, data)
	})

	t.Run("internal error", func(t *testing.T) {
		t.Parallel()

		s := setUpArticle(t)
		defer s.tearDown()

		s.mArticle.EXPECT().GetByID(gomock.Any(), int64(id)).Return(nil, errors.New("some internal error"))

		data, status, err := s.server.GetArticle(ctx, int64(id))
		require.Equal(t, http.StatusInternalServerError, status)
		require.NotNil(t, err)
		require.Nil(t, data)
	})

}

func Test_parseGetID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		req := httptest.NewRequest("GET", "http://localhost:9000/article?id=42", nil)
		expectedID := int64(42)

		articleID, status, err := parseGetID(req)

		require.Nil(t, err)
		require.Equal(t, http.StatusOK, status)
		assert.Equal(t, expectedID, articleID)

	})

	t.Run("ivalid id", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://localhost:9000/article?id=invalid", nil)
		expectedID := int64(0)

		articleID, status, err := parseGetID(req)

		require.NotNil(t, err)
		require.Equal(t, http.StatusBadRequest, status)
		assert.Equal(t, expectedID, articleID)

	})

	t.Run("missing id parameter", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://localhost:9000/article", nil)
		expectedID := int64(0)

		articleID, status, err := parseGetID(req)

		require.NotNil(t, err)
		require.Equal(t, http.StatusBadRequest, status)
		assert.Equal(t, expectedID, articleID)

	})

}

func Test_GetComments(t *testing.T) {
	var (
		ctx       context.Context
		articleID int64 = 1
	)

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		s := setUpComment(t)
		defer s.tearDown()

		s.mComments.EXPECT().GetCommentsForArticle(gomock.Any(), articleID).Return(fixtures.Comments().SetComments(articleID, 2).AllComments(), nil)

		data, status, err := s.server.GetComments(ctx, articleID)

		expectedData := strings.Join([]string{
			`{"ID":1,"ArticleID":1,"Text":"comment1","CreatedAt":"0001-01-01T00:00:00Z"}`,
			`{"ID":2,"ArticleID":1,"Text":"comment2","CreatedAt":"0001-01-01T00:00:00Z"}`},
			"\n")

		require.Nil(t, err)
		require.Equal(t, http.StatusOK, status)
		assert.Equal(t, expectedData, string(data))

	})

	t.Run("success-no comments", func(t *testing.T) {
		t.Parallel()

		s := setUpComment(t)
		defer s.tearDown()

		s.mComments.EXPECT().GetCommentsForArticle(gomock.Any(), articleID).Return([]repository.Comment{}, nil)

		data, status, err := s.server.GetComments(ctx, articleID)

		require.Nil(t, err)
		require.Equal(t, http.StatusOK, status)
		assert.Equal(t, "", string(data))

	})

	t.Run("fail", func(t *testing.T) {
		t.Parallel()

		s := setUpComment(t)
		defer s.tearDown()

		s.mComments.EXPECT().GetCommentsForArticle(gomock.Any(), articleID).Return(nil, errors.New("some internal error"))

		data, status, err := s.server.GetComments(ctx, articleID)

		require.NotNil(t, err)
		require.Equal(t, http.StatusInternalServerError, status)
		require.Nil(t, data)

	})

}

func Test_CreateArticle(t *testing.T) {
	t.Parallel()
	var (
		ctx context.Context
	)

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		body := []byte(`{"name":"asd","rating":10}`)
		s := setUpArticle(t)
		defer s.tearDown()

		s.mArticle.EXPECT().Add(gomock.Any(), fixtures.Article().Name("asd").Rating(10).P()).Return(fixtures.Article().AllData().P(), nil)

		data, status, err := s.server.CreateArticle(ctx, body)
		require.Equal(t, http.StatusOK, status)
		require.Nil(t, err)
		assert.Equal(
			t,
			`{"ID":1,"Name":"asd","Rating":10,"CreatedAt":"0001-01-01T00:00:00Z"}`,
			string(data))

	})

	t.Run("fail", func(t *testing.T) {
		t.Run("bad request body", func(t *testing.T) {
			t.Parallel()
			body := []byte(`{"name"    "asd" "rate" 10}`)
			s := setUpArticle(t)
			defer s.tearDown()

			data, status, err := s.server.CreateArticle(ctx, body)
			require.NotNil(t, err)
			require.Equal(t, http.StatusBadRequest, status)
			require.Nil(t, data)

		})

		t.Run("internal error", func(t *testing.T) {
			t.Parallel()
			body := []byte(`{"name":"asd","rating":10}`)
			s := setUpArticle(t)
			defer s.tearDown()
			s.mArticle.EXPECT().Add(gomock.Any(), fixtures.Article().Name("asd").Rating(10).P()).Return(nil, errors.New("internal error"))

			data, status, err := s.server.CreateArticle(ctx, body)
			require.NotNil(t, err)
			require.Equal(t, http.StatusInternalServerError, status)
			require.Nil(t, data)

		})

	})

}

func Test_DeleteArticle(t *testing.T) {
	t.Parallel()
	var (
		ctx       context.Context
		articleID int64 = 1
	)
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		s := setUpArticle(t)
		defer s.tearDown()

		s.mArticle.EXPECT().DeleteByID(gomock.Any(), articleID).Return(nil)

		status, err := s.server.DeleteArticle(ctx, articleID)
		require.Nil(t, err)
		require.Equal(t, http.StatusOK, status)

	})

	t.Run("not found articleID", func(t *testing.T) {
		t.Parallel()
		s := setUpArticle(t)
		defer s.tearDown()

		s.mArticle.EXPECT().DeleteByID(gomock.Any(), articleID).Return(repository.ErrObjectNotFound)

		status, err := s.server.DeleteArticle(ctx, articleID)
		require.NotNil(t, err)
		require.Equal(t, http.StatusNotFound, status)

	})

	t.Run("internal error", func(t *testing.T) {
		t.Parallel()
		s := setUpArticle(t)
		defer s.tearDown()

		s.mArticle.EXPECT().DeleteByID(gomock.Any(), articleID).Return(errors.New("internal error"))

		status, err := s.server.DeleteArticle(ctx, articleID)
		require.NotNil(t, err)
		require.Equal(t, http.StatusInternalServerError, status)

	})
}

func Test_UpdateArticle(t *testing.T) {
	t.Parallel()
	var (
		ctx       context.Context
		articleID int64 = 1
	)
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		body := []byte(`{"id": 1, "name":"asd","rating":10}`)
		s := setUpArticle(t)
		defer s.tearDown()

		s.mArticle.EXPECT().Update(gomock.Any(), fixtures.Article().ID(articleID).Name("asd").Rating(10).P()).Return(nil)

		status, err := s.server.UpdateArticle(ctx, body)
		require.Nil(t, err)
		require.Equal(t, http.StatusOK, status)

	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()
		body := []byte(`{"id": 1, "name":"asd","rating":10}`)
		s := setUpArticle(t)
		defer s.tearDown()

		s.mArticle.EXPECT().Update(gomock.Any(), fixtures.Article().ID(articleID).Name("asd").Rating(10).P()).Return(repository.ErrObjectNotFound)

		status, err := s.server.UpdateArticle(ctx, body)
		require.NotNil(t, err)
		require.Equal(t, http.StatusNotFound, status)

	})

	t.Run("bad request", func(t *testing.T) {
		t.Parallel()
		body := []byte(`{"id" 1  "name":"asd","rating":10}`)
		s := setUpArticle(t)
		defer s.tearDown()

		status, err := s.server.UpdateArticle(ctx, body)
		require.NotNil(t, err)
		require.Equal(t, http.StatusBadRequest, status)

	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()
		body := []byte(`{"id": 1, "name":"asd","rating":10}`)
		s := setUpArticle(t)
		defer s.tearDown()

		s.mArticle.EXPECT().Update(gomock.Any(), fixtures.Article().ID(articleID).Name("asd").Rating(10).P()).Return(errors.New("internal error"))

		status, err := s.server.UpdateArticle(ctx, body)
		require.NotNil(t, err)
		require.Equal(t, http.StatusInternalServerError, status)

	})
}

func Test_CreateNewComment(t *testing.T) {
	t.Parallel()
	var (
		ctx       context.Context
		articleID int64 = 1
	)
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		body := []byte(`{"article_id":1,"text":"comment1"}`)
		s := setUpComment(t)
		defer s.tearDown()

		s.mComments.EXPECT().AddComment(gomock.Any(), fixtures.Comments().ArticleID(articleID).Text("comment1").P()).Return(fixtures.Comments().DataComment(articleID), nil)

		data, status, err := s.server.CreateNewComment(ctx, body)
		require.Nil(t, err)
		require.Equal(t, http.StatusOK, status)
		assert.Equal(t, `{"ID":1,"ArticleID":1,"Text":"comment1","CreatedAt":"0001-01-01T00:00:00Z"}`, string(data))
	})

	t.Run("bad request", func(t *testing.T) {
		t.Parallel()
		body := []byte(`{"article_id" 1 "text" "comment1"}`)
		s := setUpComment(t)
		defer s.tearDown()

		data, status, err := s.server.CreateNewComment(ctx, body)
		require.NotNil(t, err)
		require.Equal(t, http.StatusBadRequest, status)
		require.Nil(t, data)
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()
		body := []byte(`{"article_id":1,"text":"comment1"}`)
		s := setUpComment(t)
		defer s.tearDown()

		s.mComments.EXPECT().AddComment(gomock.Any(), fixtures.Comments().ArticleID(articleID).Text("comment1").P()).Return(nil, repository.ErrObjectNotFound)

		data, status, err := s.server.CreateNewComment(ctx, body)
		require.NotNil(t, err)
		require.Equal(t, http.StatusNotFound, status)
		require.Nil(t, data)
	})

	t.Run("internal error", func(t *testing.T) {
		t.Parallel()
		body := []byte(`{"article_id":1,"text":"comment1"}`)
		s := setUpComment(t)
		defer s.tearDown()
		s.mComments.EXPECT().AddComment(gomock.Any(), fixtures.Comments().ArticleID(articleID).Text("comment1").P()).Return(nil, errors.New("internal error"))

		data, status, err := s.server.CreateNewComment(ctx, body)
		require.NotNil(t, err)
		require.Equal(t, http.StatusInternalServerError, status)
		require.Nil(t, data)
	})
}
