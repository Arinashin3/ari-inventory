package api

import (
	"ari-inventory/database"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type jsonForm struct {
	Job       string `json:"job"`
	Uuid      string `json:"uuid"`
	Hostname  string `json:"hostname"`
	HostGroup string `json:"host_group"`
	Ipaddr    string `json:"ipaddr"`
	Port      int    `json:"port"`
	Interval  string `json:"interval"`
}

func HandlerRegister(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
		return
	}
	body, err := io.ReadAll(r.Body)

	var f jsonForm

	err = json.Unmarshal(body, &f)
	if f.HostGroup == "" {
		f.HostGroup = "none"
	}
	if f.Uuid == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("uuid required"))
		return
	}
	if f.Interval == "" {
		f.Interval = "15s"
	}
	if f.Ipaddr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("ipaddr required"))
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	db := database.GetDatabase()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	_ = db.Open()
	defer db.Close()
	rows, err := db.Query("SELECT * FROM job WHERE uuid = :uuid and job = :job", sql.Named("uuid", f.Uuid), sql.Named("job", f.Job))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
	if rows.Next() {
		rows.Close()
		_, err = db.Exec("UPDATE job SET hostname = :hostname, ipaddr = :ipaddr, port = :port, host_group = :hostGroup, interval = :interval, updateAt = CURRENT_TIMESTAMP WHERE uuid = :uuid and job = :job", sql.Named("uuid", f.Uuid), sql.Named("hostname", f.Hostname), sql.Named("ipaddr", f.Ipaddr), sql.Named("port", f.Port), sql.Named("job", f.Job), sql.Named("interval", f.Interval), sql.Named("hostGroup", f.HostGroup))
		log.Println("update job success")
		w.Write([]byte("Successfully Update"))
	} else {
		rows.Close()
		_, err = db.Exec("INSERT INTO job (uuid, hostname, host_group, ipaddr, port, job, interval) VALUES (:uuid, :hostname, :hostGroup,:ipaddr, :port, :job, :interval)", sql.Named("uuid", f.Uuid), sql.Named("hostname", f.Hostname), sql.Named("ipaddr", f.Ipaddr), sql.Named("port", f.Port), sql.Named("job", f.Job), sql.Named("interval", f.Interval), sql.Named("hostGroup", f.HostGroup))
		log.Println("insert job success")
		w.Write([]byte("Successfully Insert"))
	}
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

}
