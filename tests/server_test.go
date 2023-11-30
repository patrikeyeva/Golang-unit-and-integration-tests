//go:build integrationHandler

package tests

import (
	"bytes"
	"encoding/json"
	"homework5/internal/pkg/repository"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"gotest.tools/v3/assert"
)

func TestCreate(t *testing.T) {

	t.Run("success", func(t *testing.T) {

		testDB.SetUp(t)
		defer testDB.TearDown()

		req := httptest.NewRequest("POST", "/article", bytes.NewReader([]byte(`{"name":"cats","rating":5}`)))
		recorder := httptest.NewRecorder()

		testServer.Create(recorder, req)

		resp := recorder.Result()
		require.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}
		var article repository.Article
		if err := json.Unmarshal(body, &article); err != nil {
			t.Fatalf("Failed to unmarshal JSON : %v", err)
		}
		assert.Equal(t, "cats", article.Name)
		assert.Equal(t, int64(5), article.Rating)

	})

	t.Run("fail", func(t *testing.T) {

		testDB.SetUp(t)
		defer testDB.TearDown()

		req := httptest.NewRequest("POST", "/article", bytes.NewReader([]byte(`{"naming":"cats" "rate":5}`)))
		recorder := httptest.NewRecorder()

		testServer.Create(recorder, req)

		resp := recorder.Result()
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	})
}
