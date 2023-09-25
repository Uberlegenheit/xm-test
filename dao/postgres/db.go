package postgres

import (
	"errors"
	"fmt"
	"net"

	"github.com/golang-migrate/migrate/v4"
	pg "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"xm-task/conf"
	"xm-task/log"
)

const migrationsPath = "./dao/postgres/migrations"

type Postgres struct {
	cfg conf.Postgres
	db  *gorm.DB
}

func NewPostgres(cfg conf.Postgres, migrate bool) (*Postgres, error) {
	conn, err := makeConn(cfg)
	if err != nil {
		return nil, fmt.Errorf("makeConn: %s", err.Error())
	}

	sqlConn, err := conn.DB()
	if err != nil {
		return nil, err
	}

	if err := sqlConn.Ping(); err != nil {
		return nil, err
	}

	db := &Postgres{
		cfg: cfg,
		db:  conn,
	}
	if migrate {
		err = db.makeMigration(conn, migrationsPath)
		if err != nil {
			return nil, fmt.Errorf("makeMigration: %s", err.Error())
		}
	}

	return db, nil
}

func makeConn(cfg conf.Postgres) (*gorm.DB, error) {
	hostPort := net.JoinHostPort(cfg.Host, cfg.Port)

	s := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		cfg.User, cfg.Password, hostPort, cfg.Database, cfg.SSLMode)
	return gorm.Open(postgres.Open(s), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
}

func (db *Postgres) makeMigration(conn *gorm.DB, migrationDir string) error {
	sqlConn, err := conn.DB()
	if err != nil {
		return fmt.Errorf("conn.DB: %s", err.Error())
	}

	driver, err := pg.WithInstance(sqlConn, &pg.Config{
		DatabaseName: db.cfg.Database,
	})
	if err != nil {
		return fmt.Errorf("postgres.WithInstance: %s", err.Error())
	}

	mg, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationDir),
		db.cfg.Database, driver)
	if err != nil {
		return fmt.Errorf("migrate.NewWithDatabaseInstance: %s", err.Error())
	}
	if err = mg.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return err
		}
	}
	return nil
}

func (db *Postgres) CheckDBStatus() bool {
	psg, err := db.db.DB()
	if err != nil {
		log.Error("db status:", zap.Error(err))
		return false
	}

	err = psg.Ping()
	if err != nil {
		log.Error("db ping:", zap.Error(err))
		return false
	}

	return true
}
