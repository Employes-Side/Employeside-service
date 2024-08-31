package handlers

import (
	"database/sql"
	"net/http"
)

func UserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
