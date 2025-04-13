package main

import (
	"dle-sitemap/database"
	"dle-sitemap/helper"
	"dle-sitemap/server"
	"dle-sitemap/sitemap"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("START")

	log.Println("runtime.GOMAXPROCS:", runtime.GOMAXPROCS(0))

	if err := godotenv.Load("../.env"); err != nil {
		log.Println("Cant load .env: ", err)
	}

	mysqlURL := os.Getenv("MYSQL_URL")
	port := os.Getenv("HTTP_PORT")

	updatePeriod := 600
	if os.Getenv("UPDATE_PERIOD") != "" {
		updatePeriod = helper.StrToInt(os.Getenv("UPDATE_PERIOD"))
	}

	if os.Getenv("MYSQL_URL_FILE") != "" {
		mysqlURL_, err := os.ReadFile(os.Getenv("MYSQL_URL_FILE"))
		if err != nil {
			log.Fatal(err)
		}
		mysqlURL = strings.TrimSpace(string(mysqlURL_))
	}

	dbService, err := database.NewService(mysqlURL)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("dbService OK")
	}

	sitemapService, err := sitemap.NewService(dbService, updatePeriod)
	if err != nil {
		log.Fatal(err)
	}

	serverService, err := server.NewService(port, sitemapService)
	if err != nil {
		log.Fatal(err)
	}
	serverService.Run()
}
