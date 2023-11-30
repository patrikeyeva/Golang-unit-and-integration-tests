//go:build integrationHandler

package tests

import (
	"homework5/internal/pkg/db"
	"homework5/internal/pkg/repository/postgresql"
	"homework5/internal/pkg/server"
	"homework5/tests/fixtures/postgres"

	"github.com/joho/godotenv"
)

var (
	testDB     *postgres.TestDb
	testServer *server.Server
)

func init() {
	// Загрузка .env файла или переменных окружения
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}

	// Создание фейковой базы данных
	testDB = postgres.NewFromEnv()

	// Создание сервера с фейковой базой данных
	testServer = createTestServer(testDB.DB)
}

func createTestServer(database db.DatabaseOperations) *server.Server {
	// Создание сервера с переданной базой данных
	return &server.Server{
		ArticleRepo: postgresql.NewArticles(database),
		CommentRepo: postgresql.NewComments(database),
	}
}
