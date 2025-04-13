package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func (s *Service) handler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	// domain should be in header X-Domain-Host
	dom := r.Header.Get("X-Domain-Host")
	if dom == "" {
		log.Println("domain not set")
		http.Error(w, "domain not set", http.StatusBadRequest)
		return
	}

	log.Println("request: ", r.RequestURI)

	// r.RequestURI = fex /sitemap.xml
	returnFile := os.Getenv("STORAGE_PATH") + "/" + dom + r.RequestURI

	w.Header().Add("X-Proxy-tm", fmt.Sprintf("%d", time.Since(start).Milliseconds()))

	// открываем файл
	f, err := os.Open(returnFile)
	if err != nil {
		log.Println("error opening file:", returnFile)
		log.Println(err)
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}
	defer f.Close()

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	_, err = io.Copy(w, f)
	if err != nil {
		log.Println(err)
	}
}
