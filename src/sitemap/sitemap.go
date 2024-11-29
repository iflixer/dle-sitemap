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

	domainID := s.dbService.FlixDomainIDByHost(dom)
	sm := stm.NewSitemap(1)
	sm.SetCompress(false)
	sm.SetSitemapsPath("")
	sm.SetPublicPath("")

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

	if posts, err := s.dbService.Posts(domainID); err != nil {
		return nil, err
	} else {
		log.Println("loaded posts:", len(posts))
		emptyURLQty := 0
		for _, p := range posts {
			// movies/komediya/21-kapkarashka-kubinskaja-istorija.html
			if p.URL != "" {
				sm.Add(stm.URL{{"loc", domainPrefix + p.URL}, {"changefreq", "daily"}})
			} else {
				emptyURLQty++
			}
		}

		log.Println("Empty URL qty", emptyURLQty)
	}

	sm.Finalize()
	// sm.PingSearchEngines()
	//data = sm.XMLContent()

	return
}

func NewService(dbService *database.Service) (s *Service, err error) {
	s = &Service{
		dbService: dbService,
	}
	return
}
