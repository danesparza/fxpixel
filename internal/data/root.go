package data

import (
	"context"
	"github.com/danesparza/fxpixel/scripts"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file" // Add the file migrations source to golang-migrate
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/golang-migrate/migrate/v4/source/iofs" // Add the iofs migrations source (for embed)
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	_ "modernc.org/sqlite" // Register relevant drivers.
	"os"
	"path/filepath"
)

type AppDataService interface {
	InitConfig() error
	AddTimeline(ctx context.Context, source Timeline) (Timeline, error)
	GetTimeline(ctx context.Context, id string) (Timeline, error)
	GetAllTimelines(ctx context.Context) ([]Timeline, error)
	GetAllTimelinesWithTag(ctx context.Context, tag string) ([]Timeline, error)
	DeleteTimeline(ctx context.Context, id string) error
	UpdateTags(ctx context.Context, id string, tags []string) error
}

type appDataService struct {
	*sqlx.DB
}

// InitConfig performs runtime config
func (a appDataService) InitConfig() error {
	//	Set the pragma for delete cascade
	_, err := a.DB.Exec(`PRAGMA foreign_keys = ON;`)
	if err != nil {
		return err
	}

	return nil
}

func NewAppDataService(db *sqlx.DB) AppDataService {
	svc := &appDataService{db}

	err := svc.InitConfig()
	if err != nil {
		log.Err(err).Msg("There was a problem initializing the SQLite config")
	}

	return svc
}

// InitSqlite initializes SQLite and returns a pointer to the db
// object, or an error
func InitSqlite(datasource string) (*sqlx.DB, error) {
	log.Info().Msg("InitSqlite")

	//	Make sure the path is created:
	err := os.MkdirAll(filepath.Dir(datasource), 0777)
	if err != nil {
		log.Fatal().Err(err).Str("datasource", datasource).Msg("Problem creating datasource directory")
	}

	//	Connect to the datasource
	db, err := sqlx.Open("sqlite", datasource)
	if err != nil {
		log.Fatal().Err(err).Str("datasource", datasource).Msg("Problem connecting to SQLite database")
	}

	//	Create a 'databaseDriver' object from the existing connection
	databaseDriver, err := sqlite.WithInstance(db.DB, &sqlite.Config{})
	if err != nil {
		log.Fatal().Str("datasource", datasource).Err(err).Msg("problem setting up databaseDriver for migrations")
	}

	//	Create the new migration source
	sourceDriver, err := iofs.New(scripts.FS, "sqlite/migrations")
	if err != nil {
		log.Fatal().Err(err).
			Str("datasource", datasource).
			Str("migrationsource", "(iofs)sqlite/migrations").
			Msg("problem setting up migration source")
	}

	//	Create a new migrator with the databaseDriver (and existing connection)
	migrator, err := migrate.NewWithInstance("iofs", sourceDriver, datasource, databaseDriver)
	if err != nil {
		log.Fatal().
			Str("datasource", datasource).
			Str("migrationsource", "(iofs)sqlite/migrations").
			Err(err).Msg("problem creating migrator config")
	}

	//	Run the migrations
	err = migrator.Up()
	switch err {
	case migrate.ErrNoChange:
		log.Info().Msg("sqlite schema is up-to-date")
	case nil:
		log.Info().Msg("sqlite schema was updated successfully")
	default:
		log.Err(err).Msg("problem running migrations")
		return db, err
	}

	return db, nil
}
