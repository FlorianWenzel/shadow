package main

import (
	"encoding/csv"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io"
	"log"
	"math/big"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type Request struct {
	gorm.Model
	Host      string
	Path      string
	Ip        string
	TimeTaken int32
	Continent string
	Country   string
	Region    string
	City      string
	Latitude  float64
	Longitude float64
}

type IpToLocation struct {
	Type      string
	IpFrom    big.Int
	IpTo      big.Int
	Continent string
	Country   string
	Region    string
	City      string
	Latitude  float64
	Longitude float64
}

func main() {
	initDB()
	readIpToLocationCsv()
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

var ipToLocations = []IpToLocation{}

func readIpToLocationCsv() {
	dat, err := os.ReadFile("dbip-city-lite-2024-06.csv")
	if err != nil {
		panic(err)
	}
	reader := csv.NewReader(strings.NewReader(string(dat)))
	totalLines := len(strings.Split(string(dat), "\n"))
	processedLines := 0
	processedPercent := 0

	ipToLocations = make([]IpToLocation, 0, totalLines)

	reader.LazyQuotes = true

	for {
		fields, err := reader.Read()
		processedLines++
		percent := int((float64(processedLines) / float64(totalLines)) * 100)
		if percent > processedPercent && percent%10 == 0 {
			processedPercent = percent
			fmt.Printf("Processed %d%%%%", processedPercent)
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}

		if len(fields) < 8 {
			continue
		}
		format := "ipv4"
		if strings.Contains(fields[0], ":") {
			format = "ipv6"
		}

		from := big.NewInt(0)
		to := big.NewInt(0)

		if format == "ipv6" {
			from, err = ipv6ToBigInt(fields[0])
			to, err = ipv6ToBigInt(fields[1])
		} else {
			from, err = ipv4ToBigInt(fields[0])
			to, err = ipv4ToBigInt(fields[1])
		}
		continent := strings.Trim(fields[2], "\"")
		country := strings.Trim(fields[3], "\"")
		region := strings.Trim(fields[4], "\"")
		city := strings.Trim(fields[5], "\"")
		latitude, err := strconv.ParseFloat(fields[6], 64)
		if err != nil {
			log.Fatal(err)
		}
		longitude, err := strconv.ParseFloat(fields[7], 64)
		if err != nil {
			log.Fatal(err)
		}

		ipToLocation := IpToLocation{
			Type:      format,
			Continent: continent,
			Country:   country,
			Region:    region,
			City:      city,
			Latitude:  latitude,
			Longitude: longitude,
			IpFrom:    *from,
			IpTo:      *to,
		}
		ipToLocations = append(ipToLocations, ipToLocation)
	}
	fmt.Printf("IpToLocation data loaded")
}
func ipv6ToBigInt(ipv6 string) (*big.Int, error) {
	ip := net.ParseIP(ipv6)
	if ip == nil {
		return nil, fmt.Errorf("invalid IPv6 address")
	}

	ip = ip.To16()
	if ip == nil {
		return nil, fmt.Errorf("not a valid IPv6 address")
	}

	ipInt := new(big.Int)
	ipInt.SetBytes(ip)

	return ipInt, nil
}

func ipv4ToBigInt(ipv4 string) (*big.Int, error) {
	ip := net.ParseIP(ipv4)
	if ip == nil {
		return nil, fmt.Errorf("invalid IPv4 address")
	}

	ip = ip.To4()
	if ip == nil {
		return nil, fmt.Errorf("not a valid IPv4 address")
	}

	ipInt := new(big.Int)
	ipInt.SetBytes(ip)

	return ipInt, nil
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
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = strings.Split(r.RemoteAddr, ":")[0]
	}

	location := IpToLocation{}
	format := "ipv4"
	if strings.Contains(ip, ":") {
		format = "ipv6"
	}
	ipInt := big.NewInt(0)
	if format == "ipv6" {
		ipInt, _ = ipv6ToBigInt(ip)
	} else {
		ipInt, _ = ipv4ToBigInt(ip)
	}
	location = IpToLocation{}
	for _, ipToLocation := range ipToLocations {
		if ipToLocation.IpFrom.Cmp(ipInt) < 1 && ipToLocation.IpTo.Cmp(ipInt) > -1 {
			location = ipToLocation
			break
		}
	}

	if location.City == "" {
		log.Printf("No location found for IP: %s", ip)
	} else {
		log.Printf("Location found for IP: %s, City: %s, Country: %s, Continent: %s", ip, location.City, location.Country, location.Continent)
	}

	request := Request{Host: r.Host, Path: r.URL.Path, TimeTaken: duration, Ip: ip, Continent: location.Continent, Country: location.Country, Region: location.Region, City: location.City, Latitude: location.Latitude, Longitude: location.Longitude}
	db.Create(&request)

	log.Printf("Request: %s %s %s Duration: %vns", r.Host, r.URL.Path, r.URL.RawQuery, duration)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		var duration = int32(time.Since(start).Microseconds())
		go logRequest(r, duration)
	})
}
