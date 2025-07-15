package api

import (
	"ari-inventory/database"
	"database/sql"
	"net/http"
	"time"
)

func HandlerClear(w http.ResponseWriter, r *http.Request) {

	if r.Method != "DELETE" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
		return
	}

	db := database.GetDatabase()

	db.Open()
	defer db.Close()
	limitDate := time.Now().UTC().Add(-24 * time.Hour)
	db.Exec("DELETE FROM job WHERE updateAt < :limitDate", sql.Named("limitDate", limitDate))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

