package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"homework5/internal/pkg/repository"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Server struct {
	ArticleRepo repository.ArticlesRepo
	CommentRepo repository.CommentsRepo
}

type articleRequest struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Rating int64  `json:"rating"`
}

type commentRequest struct {
	ArticleID int64  `json:"article_id"`
	Text      string `json:"text"`
}

func mapArticleRequest(article articleRequest) *repository.Article {
	return &repository.Article{
		ID:     article.ID,
		Name:   article.Name,
		Rating: article.Rating,
	}
}

func mapCommentRequest(comment commentRequest) *repository.Comment {
	return &repository.Comment{
		ArticleID: comment.ArticleID,
		Text:      comment.Text,
	}
}

func CreateRouter(server Server) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/article", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			server.Create(w, req)
		case http.MethodPut:
			server.Update(w, req)
		case http.MethodGet:
			server.Get(w, req)
		case http.MethodDelete:
			server.Delete(w, req)
		default:
			fmt.Println("error")
		}
	})

	router.HandleFunc("/comment", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			server.CreateComment(w, req)
		default:
			fmt.Println("error")
		}
	})
	return router
}

func (server *Server) Create(w http.ResponseWriter, req *http.Request) {
	body, status, err := parseGetBody(req.Body)
	if err != nil {
		http.Error(w, err.Error(), status)
		return
	}

	data, status, err := server.CreateArticle(req.Context(), body)
	if err != nil {
		http.Error(w, err.Error(), status)
		return
	}
	w.Write(data)
}

func parseGetBody(reqBody io.ReadCloser) ([]byte, int, error) {
	body, err := io.ReadAll(reqBody)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return body, 0, nil

}

func (server *Server) CreateArticle(ctx context.Context, body []byte) ([]byte, int, error) {
	var articleReq articleRequest
	if errJson := json.Unmarshal(body, &articleReq); errJson != nil {
		return nil, http.StatusBadRequest, errJson
	}
	article := mapArticleRequest(articleReq)
	article, err := server.ArticleRepo.Add(ctx, article)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	articleJson, err := json.Marshal(&article)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return articleJson, http.StatusOK, nil
}

func parseGetID(req *http.Request) (int64, int, error) {
	articleID, err := strconv.ParseInt(req.URL.Query().Get("id"), 10, 64)
	if err != nil {
		return 0, http.StatusBadRequest, err
	}

	return articleID, http.StatusOK, nil
}

func (server *Server) GetArticle(ctx context.Context, articleID int64) ([]byte, int, error) {
	article, err := server.ArticleRepo.GetByID(ctx, articleID)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return nil, http.StatusNotFound, err
		}
		return nil, http.StatusInternalServerError, err
	}
	articleJson, err := json.Marshal(&article)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return articleJson, http.StatusOK, nil

}

func (server *Server) GetComments(ctx context.Context, articleID int64) ([]byte, int, error) {
	comments, err := server.CommentRepo.GetCommentsForArticle(ctx, articleID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	var data []byte

	for idx, comment := range comments {
		commentJson, err := json.Marshal(&comment)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		if idx > 0 {
			data = append(data, []byte("\n")...)
		}

		data = append(data, commentJson...)
	}

	return data, http.StatusOK, nil

}

func (server *Server) Get(w http.ResponseWriter, req *http.Request) {
	articleID, statusArticle, errArticle := parseGetID(req)
	if statusArticle != http.StatusOK {
		http.Error(w, errArticle.Error(), statusArticle)
		return
	}

	dataArticle, statusArticle, errArticle := server.GetArticle(req.Context(), articleID)
	if statusArticle != http.StatusOK {
		http.Error(w, errArticle.Error(), statusArticle)
		return
	}
	w.Write([]byte("Article json: "))
	w.Write(dataArticle)
	w.Write([]byte("\n"))

	dataComment, statusComment, errComment := server.GetComments(req.Context(), articleID)
	if statusArticle != http.StatusOK {
		http.Error(w, errComment.Error(), statusComment)
		return
	}
	w.Write([]byte("Comments: \n"))
	if len(dataComment) == 0 {
		w.Write([]byte("no comments"))
	}
	w.Write(dataComment)
	w.Write([]byte("\n"))

}

func (server *Server) Delete(w http.ResponseWriter, req *http.Request) {
	articleID, status, err := parseGetID(req)
	if err != nil {
		http.Error(w, err.Error(), status)
		return
	}
	status, err = server.DeleteArticle(req.Context(), articleID)
	if err != nil {
		http.Error(w, err.Error(), status)
	}

}

func (server *Server) DeleteArticle(ctx context.Context, articleID int64) (int, error) {
	err := server.ArticleRepo.DeleteByID(ctx, articleID)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return http.StatusNotFound, err
		} else {
			return http.StatusInternalServerError, err
		}
	}
	return http.StatusOK, nil
}

func (server *Server) Update(w http.ResponseWriter, req *http.Request) {
	body, status, err := parseGetBody(req.Body)
	if err != nil {
		http.Error(w, err.Error(), status)
		return
	}

	status, err = server.UpdateArticle(req.Context(), body)
	if err != nil {
		http.Error(w, err.Error(), status)
	}

}

func (server *Server) UpdateArticle(ctx context.Context, body []byte) (int, error) {
	var updateData articleRequest
	if err := json.Unmarshal(body, &updateData); err != nil {
		return http.StatusBadRequest, err
	}

	articleRepo := mapArticleRequest(updateData)

	if err := server.ArticleRepo.Update(ctx, articleRepo); err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return http.StatusNotFound, err
		} else {
			return http.StatusInternalServerError, err
		}
	}
	return http.StatusOK, nil

}

func (server *Server) CreateComment(w http.ResponseWriter, req *http.Request) {
	body, status, err := parseGetBody(req.Body)
	if err != nil {
		http.Error(w, err.Error(), status)
		return
	}

	data, status, err := server.CreateNewComment(req.Context(), body)
	if err != nil {
		http.Error(w, err.Error(), status)
		return
	}
	w.Write(data)
}

func (server *Server) CreateNewComment(ctx context.Context, body []byte) ([]byte, int, error) {
	var comment commentRequest
	if err := json.Unmarshal(body, &comment); err != nil {
		return nil, http.StatusBadRequest, err
	}
	commentRepo := mapCommentRequest(comment)
	commentRepo, err := server.CommentRepo.AddComment(ctx, commentRepo)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return nil, http.StatusNotFound, err
		}
		return nil, http.StatusInternalServerError, err

	}
	commentJson, err := json.Marshal(commentRepo)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return commentJson, http.StatusOK, nil
}
