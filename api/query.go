package api

import (
	"ari-inventory/database"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

type Job struct {
	Uuid      string
	Hostname  string
	HostGroup string
	Ipaddr    string
	Port      string
	Interval  string
}

type ResponseJob struct {
	Targets []string `json:"targets"`
	Labels  struct {
		MetaInvHostname  string `json:"__meta_inv_hostname"`
		MetaInvHostGroup string `json:"__meta_inv_host_group"`
		MetaInvIpaddr    string `json:"__meta_inv_ipaddr"`
		MetaInvUuid      string `json:"__meta_inv_uuid"`
		ScrapeInterval   string `json:"__scrape_interval__"`
	} `json:"labels"`
}

func HandlerQuery(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		w.Write([]byte("Method not allowed"))
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	req := r.URL.Query().Get("job")
	if req == "" {
		w.Write([]byte("Job not provided"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db := database.GetDatabase()
	db.Open()
	defer db.Close()
	rows, err := db.Query("SELECT uuid, hostname, host_group, ipaddr, port, interval FROM job WHERE job = :job", sql.Named("job", req))
	if err != nil {
		log.Println(err)
		w.Write([]byte("Query Error"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var jobs []Job
	for i := 0; rows.Next(); i++ {
		job := Job{}
		rows.Scan(&job.Uuid, &job.Hostname, &job.HostGroup, &job.Ipaddr, &job.Port, &job.Interval)
		jobs = append(jobs, job)

	}
	var rjs []ResponseJob
	for _, job := range jobs {
		rj := ResponseJob{}
		targets := []string{job.Ipaddr + ":" + job.Port}
		rj.Targets = targets
		rj.Labels.MetaInvHostname = job.Hostname
		rj.Labels.MetaInvHostGroup = job.HostGroup
		rj.Labels.MetaInvIpaddr = job.Ipaddr
		rj.Labels.MetaInvUuid = job.Uuid
		rj.Labels.ScrapeInterval = job.Interval
		rjs = append(rjs, rj)
	}
	b, err := json.Marshal(rjs)
	if err != nil {
		log.Println(err)
		w.Write([]byte("Query Error"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)

}
