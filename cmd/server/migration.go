package main

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	"k8s.io/klog"
)

func migrateSQL(dsn string, database string) {
	migrationDSN := fmt.Sprintf("%s&multiStatements=true", dsn)
	db1, err := sql.Open("mysql", migrationDSN)
	if err != nil {
		klog.Fatalf("failed to open migration connection: %v", err)
	}
	driver, err := mysql.WithInstance(db1, &mysql.Config{})
	if err != nil {
		klog.Fatalf("failed to create instance : %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		database,
		driver,
	)
	if err != nil {
		klog.Fatalf("failed to deploy sql: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		klog.Errorf("up migration failed: %v", err)

		err2 := m.Down()
		if err2 != nil {
			klog.Errorf("migration rollback failed!: %v", err)
		} else {
			klog.Errorf("rolled back to last version")
		}

		panic(err)
	}
}
