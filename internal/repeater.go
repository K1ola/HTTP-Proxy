package internal

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strings"
	"time"
)

type Request struct {
	Method string        `json:"method"`
	Uri    string        `json:"uri"`
	Proto  string        `json:"proto"`
	Body   string        `json:"body"`
	Id     int           `json:"id"`
	//Time   time.Duration `json:"time"`
}

type requestSlice struct {
	Res []Request `json:"res"`
}

type Repeater struct {
	DB *sql.DB
	Server *http.Server
	Router *mux.Router
}

func InitRepeater() *Repeater {
	repeater := Repeater{}
	repeater.Router = mux.NewRouter()
	repeater.Router.HandleFunc("/{id:[0-9]+}", repeater.RepeatRequest)
	repeater.Router.HandleFunc("/showRequests", repeater.ShowRequests)

	repeater.Server = &http.Server{
		Addr:    ":8887",
		Handler: repeater.Router,
	}

	return &repeater
}

func (r *Repeater) ShowRequests(w http.ResponseWriter, req *http.Request) {
	result := make([]Request,0,1024)
	var method string
	var uri string
	var proto string
	var body string
	var id int
	var created time.Duration
	rowsRequest, err := r.DB.Query("select id, method, uri, proto, body, created from requests")
	if err != nil {}
	data := Request{}
	//err := rowsRequest.Scan(&method, &uri, &proto)
	for rowsRequest.Next() {
		err := rowsRequest.Scan(
			&id,
			&method,
			&uri,
			&proto,
			&body,
			&created,
		)
		data.Id = id
		data.Uri = uri
		data.Method = method
		data.Proto = proto
		data.Body = body
		//data.Time = created

		if err != nil {}
		result = append(result, data)
	}

	bytes, err := json.Marshal(result)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(bytes)
}

func (r *Repeater) RepeatRequest(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]
	var method, uri, proto, body string
	var key, value string

	rowsRequest, _ := r.DB.Query("select method, uri, proto, body from requests where id = ?", id)
	//err := rowsRequest.Scan(&method, &uri, &proto)
	for rowsRequest.Next() {
		err := rowsRequest.Scan(
			&method,
			&uri,
			&proto,
			&body,
		)
		if err != nil {}
	}
	req, _ = http.NewRequest(method, uri, strings.NewReader(body))

	rowsHeaders, err := r.DB.Query("select key, value from headers where request_id = ?", id)
	if err != nil {}
	for rowsHeaders.Next() {
		_ = rowsHeaders.Scan(
			&key,
			&value,
		)

		if key != "If-None-Match" && key != "Accept-Encoding" && key != "If-Modified-Since" {
			req.Header.Add(key, value)
		}
	}

	HandleHTTP(w, req)
}

func HandleHTTP(w http.ResponseWriter, req *http.Request) {
	resp, err := http.DefaultClient.Do(req)
	//resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()
	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
