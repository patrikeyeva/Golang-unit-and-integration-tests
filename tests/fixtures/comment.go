package fixtures

import (
	"homework5/internal/pkg/repository"
	"strconv"
	"time"
)

type CommentBuilder struct {
	instances      []repository.Comment
	currentComment *repository.Comment
}

func Comments() *CommentBuilder {
	return &CommentBuilder{
		instances:      []repository.Comment{},
		currentComment: &repository.Comment{},
	}
}

func (b *CommentBuilder) Add() *CommentBuilder {
	comment := &repository.Comment{}
	b.instances = append(b.instances, *comment)
	b.currentComment = &b.instances[len(b.instances)-1]
	return b
}

func (b *CommentBuilder) Id(v int64) *CommentBuilder {
	b.currentComment.ID = v
	return b
}

func (b *CommentBuilder) ArticleID(v int64) *CommentBuilder {
	b.currentComment.ArticleID = v
	return b
}

func (b *CommentBuilder) Text(v string) *CommentBuilder {
	b.currentComment.Text = v
	return b
}

func (b *CommentBuilder) CreatedAt(v time.Time) *CommentBuilder {
	b.currentComment.CreatedAt = v
	return b
}

func (b *CommentBuilder) SetComments(id int64, num int) *CommentBuilder {
	comment := "comment"
	for i := 1; i <= num; i++ {
		idx := strconv.Itoa(i)
		b.Add().Id(int64(i)).ArticleID(id).Text(comment + idx).CreatedAt(time.Time{})
	}
	return b
}

func (b *CommentBuilder) P() *repository.Comment {
	return b.currentComment
}

func (b *CommentBuilder) AllComments() []repository.Comment {
	return b.instances
}

func (b *CommentBuilder) DataComment(articleID int64) *repository.Comment {
	return b.SetComments(articleID, 1).P()
}
