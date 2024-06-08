package main

import (
	"encoding/json"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"time"
)

var db *gorm.DB

type Request struct {
	gorm.Model
	Host      string
	Path      string
	TimeTaken time.Duration
}

type Count struct {
	Label string
	Count float64
}

type ChartData struct {
	Labels   []string  `json:"labels"`
	Datasets []Dataset `json:"datasets"`
}
type Dataset struct {
	Label string    `json:"label"`
	Data  []float64 `json:"data"`
}

func formatCount(title string, results []Count) []byte {
	chartData := ChartData{
		Labels:   []string{},
		Datasets: []Dataset{},
	}

	chartData.Labels = []string{title}

	for _, r := range results {
		dataset := Dataset{
			Label: "",
			Data:  []float64{},
		}
		dataset.Data = append(dataset.Data, r.Count)
		dataset.Label = r.Label
		chartData.Datasets = append(chartData.Datasets, dataset)
	}

	// Marshal the results into JSON
	jsonData, err := json.Marshal(chartData)
	if err != nil {
		log.Println(err)
	}
	return jsonData
}

func main() {
	initDB()

	http.HandleFunc("/api/requests-per-page", func(w http.ResponseWriter, r *http.Request) {
		var results []Count
		result := db.Raw("SELECT path as label, COUNT(*) as count FROM requests GROUP BY path").Scan(&results)
		if result.Error != nil {
			log.Println(result.Error)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		jsonData := formatCount("Requests per Page", results)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	})

	http.HandleFunc("/api/requests-per-ip", func(w http.ResponseWriter, r *http.Request) {
		var results []Count
		result := db.Raw("SELECT ip as label, COUNT(*) as count FROM requests GROUP BY ip").Scan(&results)
		if result.Error != nil {
			log.Println(result.Error)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		jsonData := formatCount("Requests per Page", results)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	})

	http.HandleFunc("/api/average-latency-per-route", func(w http.ResponseWriter, r *http.Request) {
		var results []Count
		result := db.Raw("SELECT path as label, ROUND(AVG(time_taken) / 1000, 2) as count FROM requests GROUP BY path").Scan(&results)
		if result.Error != nil {
			log.Println(result.Error)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		jsonData := formatCount("Average Latency per Route", results)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	})

	log.Fatal(http.ListenAndServe(":3000", nil))

}

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
