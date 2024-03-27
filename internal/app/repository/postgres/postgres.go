// Package postgresql postgresql handling PostgreSQL database.
package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"

	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"
)

// databaseErrorsToRetry errors to retry request to database.
var databaseErrorsToRetry = []error{
	errors.New(pgerrcode.UniqueViolation),
	errors.New(pgerrcode.ConnectionException),
	errors.New(pgerrcode.ConnectionDoesNotExist),
	errors.New(pgerrcode.ConnectionFailure),
	errors.New(pgerrcode.SQLClientUnableToEstablishSQLConnection),
	errors.New(pgerrcode.SQLServerRejectedEstablishmentOfSQLConnection),
	errors.New(pgerrcode.TransactionResolutionUnknown),
	errors.New(pgerrcode.ProtocolViolation),
}

// postgresDB storageManager implementation for PostgreSQL. Consist of database and QueriesHandler.
// Request to database are synchronized by sync.RWMutex. All requests is done on united transaction.
// Multi insert/update/delete is not supported at the moment.
type postgresDB struct {
	database *sql.DB
	reformDB *reform.DB

	log logger.BaseLogger
}

// CreatePostgreDB creates manager implementation. Supports migrations and check connection to database.
func CreatePostgreDB(ctx context.Context, cfg config.Config, log logger.BaseLogger) (manager.BaseStorageManager, error) {
	log.Info("[CreatePostgreDB] open database with settings: '%s'", cfg.DatabaseDSN)
	createDatabaseError := "create db: %w"
	database, err := sql.Open("pgx", cfg.DatabaseDSN)
	if err != nil {
		return nil, fmt.Errorf(createDatabaseError, err)
	}

	driver, err := postgres.WithInstance(database, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf(createDatabaseError, err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://db/migrations", "postgres", driver)
	if err != nil {
		return nil, fmt.Errorf(createDatabaseError, err)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf(createDatabaseError, err)
	}

	dbManager := &postgresDB{
		database: database,
		reformDB: nil,
		log:      log,
	}

	if _, err = dbManager.CheckConnection(ctx); err != nil {
		return nil, fmt.Errorf(createDatabaseError, err)
	}

	database.SetMaxIdleConns(cfg.DatabaseIdleConn)
	database.SetMaxOpenConns(cfg.DatabaseOpenConn)
	dbManager.reformDB = reform.NewDB(database, postgresql.Dialect, reform.NewPrintfLogger(log.Printf))

	log.Info("[CreatePostgreDB] successful")
	return dbManager, nil
}

// CheckConnection checks connection to database.
func (p *postgresDB) CheckConnection(ctx context.Context) (bool, error) {
	exec := func(context context.Context) (int64, []byte, error) {
		return 0, []byte{}, p.database.PingContext(context)
	}
	_, _, err := retryer.RetryCallWithTimeout(ctx, p.log, nil, databaseErrorsToRetry, exec)
	if err != nil {
		return false, fmt.Errorf("check connection: %w", err)
	}
	return true, nil
}

// Close closes database.
func (p *postgresDB) Close() error {
	return p.database.Close()
}

// EditNote modifies existing note in database. Also can affect on news_categories table.
// If note.ID = 0 - inserts new note in news table.
// returns id of created/modified note.
func (p *postgresDB) EditNote(ctx context.Context, note *models.News) error {
	tx, err := p.reformDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("create transaction: %w", err)
	}

	if note.ID == 0 {
		// add new note.
		if err = insertNewNote(tx, note); err != nil {
			_ = tx.Rollback()
			return err
		}
	} else {
		// update if exists.
		if err = updateExistingNote(tx, note); err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	if err = updateCategoriesIfNeed(tx, note); err != nil {
		_ = tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func insertNewNote(tx *reform.TX, note *models.News) error {
	if note.Title == constants.MissingStringFlag || note.Content == constants.MissingStringFlag {
		return errStrg.ErrIncorrectNewNote
	}

	if err := tx.Save(note); err != nil {
		return fmt.Errorf("save new note: %w", err)
	}

	return nil
}

func updateExistingNote(tx *reform.TX, note *models.News) error {
	noteRaw, err := tx.SelectOneFrom(models.NewsTable, "WHERE id = $1", note.ID)
	if err != nil {
		return fmt.Errorf("Not found id: %d: %w", note.ID, err)
	}

	noteExisting := noteRaw.(*models.News)
	if note.Title != constants.MissingStringFlag {
		noteExisting.Title = note.Title
	}
	if note.Content != constants.MissingStringFlag {
		noteExisting.Content = note.Content
	}
	if !reflect.DeepEqual(note.Categories, constants.MissingInt64ArrayFlag) {
		noteExisting.Categories = note.Categories
	}

	if err = tx.Save(noteExisting); err != nil {
		return fmt.Errorf("save/update news: %w", err)
	}

	return nil
}

func updateCategoriesIfNeed(tx *reform.TX, note *models.News) error {
	if !reflect.DeepEqual(note.Categories, constants.MissingInt64ArrayFlag) {
		_, err := tx.DeleteFrom(models.NewsCategoriesTable, "WHERE news_id = $1", note.ID)
		if err != nil {
			return fmt.Errorf("remove existing categories for news_id '%d': %w", note.ID, err)
		}

		for _, category := range note.Categories {
			categoryModel := &models.NewsCategories{
				NewsID:     note.ID,
				CategoryID: category,
			}

			if err = tx.Save(categoryModel); err != nil {
				return fmt.Errorf("save/update news: %w", err)
			}
		}
	}

	return nil
}

func (p *postgresDB) FindAll(ctx context.Context, page int64, perPage int64) ([]models.News, error) {
	tx, err := p.reformDB.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("create transaction: %w", err)
	}

	if notes, err := tx.SelectAllFrom(models.NewsTable, "LIMIT $1 OFFSET $2", perPage, (page-1)*perPage); err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("get news from db: %w", err)
	}

	res := make([]models.News, 0, len(notes))
	for _, noteStruct := range notes {
		note := *noteStruct.(*models.News)

		categories, err := tx.SelectAllFrom(models.NewsCategoriesTable, "WHERE news_id = $1", note.ID)
		if err != nil {
			_ = tx.Rollback()
			return nil, fmt.Errorf("create transaction: %w", err)
		}

		note.Categories = []int64{}
		for _, category := range categories {
			note.Categories = append(note.Categories, category.(*models.NewsCategories).CategoryID)
		}

		res = append(res, note)
	}

	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return res, nil
}
