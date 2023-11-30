package postgres

import (
	"context"
	"fmt"
	"homework5/internal/pkg/db"
	"os"
	"strings"
	"sync"
	"testing"
)

type TestDb struct {
	DB db.DatabaseOperations
	sync.Mutex
}

func NewFromEnv() *TestDb {
	database, err := db.NewDBWithDSN(context.Background(),
		fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME")))

	if err != nil {
		panic(err)
	}
	return &TestDb{DB: database}
}

func (db *TestDb) SetUp(t *testing.T) {
	t.Helper()
	db.Lock()
	db.Truncate(context.Background())
}

func (db *TestDb) TearDown() {
	defer db.Unlock()
	db.Truncate(context.Background())
}

func (db *TestDb) Truncate(ctx context.Context) {

	var tables []string
	err := db.DB.Select(ctx, &tables, "SELECT table_name FROM information_schema.tables WHERE table_schema='public' AND table_type='BASE TABLE' AND table_name != 'goose_db_version'")
	if err != nil {
		panic(err)
	}
	if len(tables) == 0 {
		panic("run migration plz")
	}
	q := fmt.Sprintf("Truncate table %s", strings.Join(tables, ","))
	if _, err := db.DB.Exec(ctx, q); err != nil {
		panic(err)
	}
}
