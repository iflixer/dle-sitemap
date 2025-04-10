package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
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

	requestFile := r.RequestURI // fex /sitemap.xml
	returnFile := strings.TrimLeft(requestFile, "/")

	_, err := s.sitemapService.Sitemap(dom)
	//log.Printf("%+v", sm)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("X-Proxy-tm", fmt.Sprintf("%d", time.Since(start).Milliseconds()))
	w.WriteHeader(http.StatusOK)

	file, err := os.ReadFile(returnFile)
	if err != nil {
		log.Println(err)
	}

	w.Write(file)
}
