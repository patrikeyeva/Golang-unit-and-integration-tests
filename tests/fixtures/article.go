package fixtures

import (
	"time"

	"homework5/internal/pkg/repository"
)

type ArticleBuilder struct {
	instance *repository.Article
}

func Article() *ArticleBuilder {
	return &ArticleBuilder{instance: &repository.Article{}}
}

func (b *ArticleBuilder) ID(v int64) *ArticleBuilder {
	b.instance.ID = v
	return b
}
func (b *ArticleBuilder) Name(v string) *ArticleBuilder {
	b.instance.Name = v
	return b
}

func (b *ArticleBuilder) Rating(v int64) *ArticleBuilder {
	b.instance.Rating = v
	return b
}

func (b *ArticleBuilder) CreatedAt(v time.Time) *ArticleBuilder {
	b.instance.CreatedAt = v
	return b
}

func (b *ArticleBuilder) P() *repository.Article {
	return b.instance
}

func (b *ArticleBuilder) V() repository.Article {
	return *b.instance
}

func (b *ArticleBuilder) AllData() *ArticleBuilder {
	return Article().ID(1).Name("asd").Rating(10).CreatedAt(time.Time{})
}
