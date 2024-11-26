package sitemap

import (
	"dle-sitemap/database"
	"log"
	"strings"

	"github.com/ikeikeikeike/go-sitemap-generator/v2/stm"
)

type Service struct {
	dbService *database.Service
}

func (s *Service) Sitemap(dom string) (data []byte, err error) {
	sm := stm.NewSitemap(1)
	sm.Create()
	sm.SetDefaultHost("https://" + dom)
	sm.Add(stm.URL{{"loc", ""}, {"changefreq", "always"}, {"mobile", true}})
	domainPrefix := "https://" + dom

	if rootCats, err := s.dbService.Cats(0); err != nil {
		return nil, err
	} else {
		for _, rc := range rootCats {
			sm.Add(stm.URL{{"loc", "/" + rc.AltName}, {"changefreq", "daily"}})
			log.Printf("%+v", rc.Parentid)
			if cats, err := s.dbService.Cats(rc.ID); err != nil {
				return nil, err
			} else {
				for _, c := range cats {
					sm.Add(stm.URL{{"loc", domainPrefix + "/" + rc.AltName + "/" + c.AltName}, {"changefreq", "daily"}})
				}
			}
		}
	}

	if posts, err := s.dbService.Posts(); err != nil {
		return nil, err
	} else {
		for _, p := range posts {
			// movies/komediya/21-kapkarashka-kubinskaja-istorija.html
			if p.URL != "" {
				if strings.HasPrefix(p.Category, "4,") || strings.HasPrefix(p.Category, "5,") || strings.HasPrefix(p.Category, "6,") || strings.HasPrefix(p.Category, "7,") || strings.HasPrefix(p.Category, "8,") {
					sm.Add(stm.URL{{"loc", domainPrefix + p.URL}, {"changefreq", "daily"}})
				}
			}
		}
	}

	// sm.Finalize().PingSearchEngines()
	// sm.PingSearchEngines()
	data = sm.XMLContent()
	return
}

func NewService(dbService *database.Service) (s *Service, err error) {
	s = &Service{
		dbService: dbService,
	}
	return
}
