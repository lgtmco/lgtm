package datastore

import (
	"database/sql"
	"os"
	"time"

	"github.com/lgtmco/lgtm/store"
	"github.com/lgtmco/lgtm/store/migration"

	"github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rubenv/sql-migrate"
	"github.com/russross/meddler"
)

type datastore struct {
	*sql.DB
}

// New creates a database connection for the given driver and datasource
// and returns a new Store.
func New(driver, config string) store.Store {
	db := Open(driver, config)
	return From(db)
}

// From returns a Store using an existing database connection.
func From(db *sql.DB) store.Store {
	return &datastore{db}
}

// Open opens a new database connection with the specified
// driver and connection string and returns a store.
func Open(driver, config string) *sql.DB {
	db, err := sql.Open(driver, config)
	if err != nil {
		logrus.Errorln(err)
		logrus.Fatalln("database connection failed")
	}
	if driver == "mysql" {
		// per issue https://github.com/go-sql-driver/mysql/issues/257
		db.SetMaxIdleConns(0)
	}

	setupMeddler(driver)

	logrus.Debugf("Driver %s", driver)
	logrus.Debugf("Data Source %s", config)

	if err := pingDatabase(db); err != nil {
		logrus.Errorln(err)
		logrus.Fatalln("database ping attempts failed")
	}

	if err := setupDatabase(driver, db); err != nil {
		logrus.Errorln(err)
		logrus.Fatalln("migration failed")
	}
	return db
}

// OpenTest opens a new database connection for testing purposes.
// The database driver and connection string are provided by
// environment variables, with fallback to in-memory sqlite.
func openTest() *sql.DB {
	var (
		driver = "sqlite3"
		config = ":memory:"
	)
	if os.Getenv("DATABASE_DRIVER") != "" {
		driver = os.Getenv("DATABASE_DRIVER")
		config = os.Getenv("DATABASE_DATASOURCE")
	}
	return Open(driver, config)
}

// helper function to ping the database with backoff to ensure
// a connection can be established before we proceed with the
// database setup and migration.
func pingDatabase(db *sql.DB) (err error) {
	for i := 0; i < 30; i++ {
		err = db.Ping()
		if err == nil {
			return
		}
		logrus.Infof("database ping failed. retry in 1s")
		time.Sleep(time.Second)
	}
	return
}

// helper function to setup the databsae by performing
// automated database migration steps.
func setupDatabase(driver string, db *sql.DB) error {
	var migrations = &migrate.AssetMigrationSource{
		Asset:    migration.Asset,
		AssetDir: migration.AssetDir,
		Dir:      driver,
	}
	_, err := migrate.Exec(db, driver, migrations, migrate.Up)
	return err
}

// helper function to setup the meddler default driver
// based on the selected driver name.
func setupMeddler(driver string) {
	switch driver {
	case "sqlite3":
		meddler.Default = meddler.SQLite
	case "mysql":
		meddler.Default = meddler.MySQL
	case "postgres":
		meddler.Default = meddler.PostgreSQL
	}
}
