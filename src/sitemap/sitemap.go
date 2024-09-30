package sitemap

import (
	"dle-sitemap/database"
	"log"

	"github.com/ikeikeikeike/go-sitemap-generator/v2/stm"
)

type Service struct {
	dbService *database.Service
}

func (s *Service) Sitemap(dom string) (data []byte, err error) {
	sm := stm.NewSitemap(1)
	sm.Create()
	sm.SetDefaultHost("https://" + dom)
	// sm.Add(stm.URL{{"loc", ""}, {"changefreq", "always"}, {"mobile", true}})

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
					sm.Add(stm.URL{{"loc", "/" + rc.AltName + "/" + c.AltName}, {"changefreq", "daily"}})
				}
			}
		}
	}

	if posts, err := s.dbService.Posts(); err != nil {
		return nil, err
	} else {
		for _, p := range posts {
			// movies/komediya/21-kapkarashka-kubinskaja-istorija.html
			sm.Add(stm.URL{{"loc", "/" + p.URL}, {"changefreq", "daily"}})
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
