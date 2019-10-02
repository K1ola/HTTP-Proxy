package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	. "http-proxy/internal"
	"log"
)

//func handleTunneling(w http.ResponseWriter, r *http.Request) {
//	dest_conn, err := net.DialTimeout("tcp", r.Host, 10*time.Second)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusServiceUnavailable)
//		return
//	}
//	w.WriteHeader(http.StatusOK)
//	hijacker, ok := w.(http.Hijacker)
//	if !ok {
//		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
//		return
//	}
//	client_conn, _, err := hijacker.Hijack()
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusServiceUnavailable)
//	}
//	go transfer(dest_conn, client_conn)
//	go transfer(client_conn, dest_conn)
//}
//
//func transfer(destination io.WriteCloser, source io.ReadCloser) {
//	defer destination.Close()
//	defer source.Close()
//	io.Copy(destination, source)
//}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////


func main() {
	proxy := InitProxy()
	proxy.DB, _ = sql.Open("sqlite3", "./data/data.db")
	defer proxy.DB.Close()

	log.Fatal(proxy.Server.ListenAndServe())

	//log.Fatal(server.ListenAndServeTLS(pemPath, keyPath))
}