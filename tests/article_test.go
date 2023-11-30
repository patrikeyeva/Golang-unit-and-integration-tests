//go:build integrationDB

package tests

import (
	"context"
	"homework5/internal/pkg/repository/postgresql"
	"homework5/tests/fixtures"
	"testing"

	"github.com/stretchr/testify/require"
	"gotest.tools/v3/assert"
)

func TestCreateArticle(t *testing.T) {
	t.Parallel()
	var (
		ctx = context.Background()
	)
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		database.SetUp(t)
		defer database.TearDown()

		article := postgresql.NewArticles(database.DB)

		resp, err := article.Add(ctx, fixtures.Article().Name("asd").Rating(10).P())

		require.NoError(t, err)
		require.NotNil(t, resp.ID, resp.CreatedAt)
		assert.Equal(t, "asd", resp.Name)
		assert.Equal(t, int64(10), resp.Rating)

	})
}

func TestGetArticle(t *testing.T) {
	t.Parallel()
	var (
		ctx = context.Background()
	)
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		database.SetUp(t)
		defer database.TearDown()

		article := postgresql.NewArticles(database.DB)
		resp, err := article.Add(ctx, fixtures.Article().Name("asd").Rating(10).P())
		require.NoError(t, err)
		require.NotZero(t, resp)

		respGet, err := article.GetByID(ctx, resp.ID)

		require.NoError(t, err)
		require.NotNil(t, respGet.ID, resp.ID)
		assert.Equal(t, "asd", respGet.Name)
		assert.Equal(t, int64(10), respGet.Rating)

	})

	t.Run("fail not found", func(t *testing.T) {
		t.Parallel()
		var id int64 = 1
		database.SetUp(t)
		defer database.TearDown()

		article := postgresql.NewArticles(database.DB)

		_, err := article.GetByID(ctx, id)

		require.Error(t, err)
	})
}

func TestDeleteArticle(t *testing.T) {
	t.Parallel()
	var (
		ctx = context.Background()
	)
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		database.SetUp(t)
		defer database.TearDown()

		article := postgresql.NewArticles(database.DB)
		resp, err := article.Add(ctx, fixtures.Article().Name("asd").Rating(10).P())
		require.NoError(t, err)
		require.NotZero(t, resp)

		err = article.DeleteByID(ctx, resp.ID)
		require.NoError(t, err)

		_, err = article.GetByID(ctx, resp.ID)
		require.Error(t, err)

	})

	t.Run("fail", func(t *testing.T) {
		t.Parallel()
		var id int64 = 1
		database.SetUp(t)
		defer database.TearDown()

		article := postgresql.NewArticles(database.DB)

		err := article.DeleteByID(ctx, id)
		require.Error(t, err)
	})
}

func TestUpdataArticle(t *testing.T) {
	t.Parallel()
	var (
		ctx = context.Background()
	)
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		database.SetUp(t)
		defer database.TearDown()

		article := postgresql.NewArticles(database.DB)
		resp, err := article.Add(ctx, fixtures.Article().Name("asd").Rating(10).P())
		require.NoError(t, err)
		require.NotZero(t, resp)

		err = article.Update(ctx, fixtures.Article().ID(resp.ID).Name("new").Rating(11).P())
		require.NoError(t, err)

		respGet, err := article.GetByID(ctx, resp.ID)
		require.NoError(t, err)
		assert.Equal(t, resp.ID, respGet.ID)
		assert.Equal(t, "new", respGet.Name)
		assert.Equal(t, int64(11), respGet.Rating)

	})
	t.Run("fail", func(t *testing.T) {
		t.Parallel()
		database.SetUp(t)
		defer database.TearDown()

		article := postgresql.NewArticles(database.DB)

		err := article.Update(ctx, fixtures.Article().P())
		require.Error(t, err)
	})
}
