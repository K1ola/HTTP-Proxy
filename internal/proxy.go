package internal

import (
	"bytes"
	"crypto/tls"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type Proxy struct {
	DB *sql.DB
	Server *http.Server
}

func InitProxy() *Proxy {
	proxy := Proxy{}

	proxy.Server = &http.Server{
		Addr: ":8888",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodConnect {
				//handleTunneling(w, r)
			} else {
				proxy.HandleHTTP(w, r)
			}
		}),
		// Disable HTTP/2.
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

	_ = os.MkdirAll("./data", 0755)
	_, _ = os.Create("./data/data.db")

	proxy.DB, _ = sql.Open("sqlite3", "./data/data.db")

	_, err := proxy.DB.Exec("CREATE TABLE IF NOT EXISTS `requests` (	`id` INTEGER PRIMARY KEY AUTOINCREMENT,	`method` VARCHAR(64) NOT NULL," +
		"`uri` VARCHAR(64) NOT NULL," +
		"`proto` VARCHAR(64) NOT NULL," +
		"`body` VARCHAR(64) NOT NULL," +
		"`created` DATE NULL" +
		");" +
		"CREATE TABLE IF NOT EXISTS  `headers` (" +
		"`id` INTEGER PRIMARY KEY AUTOINCREMENT," +
		"`request_id` INTEGER NOT NULL REFERENCES requests(id)," +
		"`key` VARCHAR(64) NOT NULL," +
		"`value` VARCHAR(64) NOT NULL" +
		");")

	if err != nil {
		log.Fatal(err)
	}

	_ = proxy.DB.Close()


	return &proxy
}

func (p *Proxy) HandleHTTP(w http.ResponseWriter, req *http.Request) {

	body, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()

	request, err := http.NewRequest(req.Method, req.RequestURI, bytes.NewReader(body))
	if err != nil {}
	request.Header = req.Header

	resp, err := http.DefaultTransport.RoundTrip(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()
	p.copyHeader(w.Header(), resp.Header)
	p.insertRequestInDb(req, body)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func (p *Proxy) copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func (p *Proxy) insertRequestInDb(r *http.Request, body []byte) {
	res, err := p.DB.Exec(
		"INSERT INTO requests(method, uri, proto, body, created) VALUES(?, ?, ?, ?, ?) /*RETURNING id*/",
		r.Method,
		r.RequestURI,
		r.Proto,
		body,
		time.Now(),
	)

	if err != nil {}
	id, _ := res.LastInsertId()

	for k, v := range r.Header {
		_, err := p.DB.Exec(
			"INSERT INTO headers(request_id, key, value) VALUES(?, ?, ?)",
			id,
			k,
			v[0],
		)

		if err != nil {
			log.Println(err)
		}
	}
}

