package main

import (
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"
)

type Request struct {
	gorm.Model
	Host      string
	Path      string
	Ip        string
	TimeTaken int32
}

func main() {
	initDB()
	proxyTarget := os.Getenv("PROXY_TARGET")
	target, err := url.Parse(proxyTarget)
	if err != nil {
		log.Fatal(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Director = func(req *http.Request) {
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Add("X-Origin-Host", target.Host)
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
	}

	http.Handle("/", loggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		print("Request received")
		proxy.ServeHTTP(w, r)
	})))

	log.Fatal(http.ListenAndServe("0.0.0.0:3001", nil))
}

var db *gorm.DB

func initDB() {
	dbHost := os.Getenv("DB_HOST")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbUser := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	dsn := "host=" + dbHost + " user=" + dbUser + " password=" + dbPassword + " dbname=" + dbName + " port=" + dbPort + " sslmode=disable TimeZone=Europe/Paris"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		return
	}
	db = database
	err = db.AutoMigrate(&Request{})
	if err != nil {
		log.Fatal(err)
		return
	}
}

func logRequest(r *http.Request, duration int32) {
	ip := strings.Split(r.RemoteAddr, ":")[0]
	request := Request{Host: r.Host, Path: r.URL.Path, TimeTaken: duration, Ip: ip}
	db.Create(&request)

	log.Printf("Request: %s %s %s Duration: %s", r.Host, r.URL.Path, r.URL.RawQuery, duration)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		var duration = int32(time.Since(start).Microseconds())
		go logRequest(r, duration)
	})
}
