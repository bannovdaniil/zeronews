package database

import (
	"database/sql"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/sirupsen/logrus"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"
	"os"
)

func ConnectToDB(dbUrl string, log *logrus.Logger) (*reform.DB, error) {
	logger := logrus.New()
	logger.SetOutput(os.Stderr)

	sqlDB, err := sql.Open("pgx", dbUrl)
	if err != nil {
		return nil, err
	}

	reformDB := reform.NewDB(sqlDB, postgresql.Dialect, reform.NewPrintfLogger(log.Printf))

	return reformDB, nil
}
