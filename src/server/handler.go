package server

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func (s *Service) handler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	dom := r.URL.Query().Get("dom")
	if dom == "" {
		log.Println("domain not set")
		http.Error(w, "domain not set", http.StatusInternalServerError)
		return
	}
	sm, err := s.sitemapService.Sitemap(dom)
	//log.Printf("%+v", sm)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("X-Proxy-tm", fmt.Sprintf("%d", time.Since(start).Milliseconds()))
	w.WriteHeader(http.StatusOK)
	w.Write(sm)
}
