package sitemap

import (
	"dle-sitemap/database"
	"dle-sitemap/helper"
	"log"
	"os"
	"time"

	"github.com/ikeikeikeike/go-sitemap-generator/v2/stm"
)

type Service struct {
	dbService    *database.Service
	updatePeriod time.Duration
}

// func (s *Service) Sitemap(dom string) (data []byte, err error) {
// 	return
// 	domainID := s.dbService.FlixDomainIDByHost(dom)
// 	sm := stm.NewSitemap(1)
// 	sm.SetCompress(false)
// 	sm.SetSitemapsPath("")
// 	sm.SetPublicPath("")

// 	sm.Create()
// 	sm.SetDefaultHost("https://" + dom)
// 	sm.Add(stm.URL{{"loc", ""}, {"changefreq", "always"}})
// 	domainPrefix := "https://" + dom

// 	if rootCats, err := s.dbService.Cats(0); err != nil {
// 		return nil, err
// 	} else {
// 		for _, rc := range rootCats {
// 			sm.Add(stm.URL{{"loc", "/" + rc.AltName}, {"changefreq", "daily"}})
// 			log.Printf("parent_id: %+v", rc.Parentid)
// 			if cats, err := s.dbService.Cats(rc.ID); err != nil {
// 				return nil, err
// 			} else {
// 				for _, c := range cats {
// 					sm.Add(stm.URL{{"loc", domainPrefix + "/" + rc.AltName + "/" + c.AltName}, {"changefreq", "daily"}})
// 				}
// 			}
// 		}
// 	}

// 	if posts, err := s.dbService.Posts(domainID); err != nil {
// 		return nil, err
// 	} else {
// 		log.Println("loaded posts:", len(posts))
// 		emptyURLQty := 0
// 		for _, p := range posts {
// 			// movies/komediya/21-kapkarashka-kubinskaja-istorija.html

// 			for i, p := range posts {
// 				if altName, err := s.FlixPostFindAltName(flixPostAltNames, p.ID); err == nil {
// 					posts[i].URL = s.makeUrl(cats, p.Category, p.ID, altName)
// 				} else {
// 					posts[i].URL = s.makeUrl(cats, p.Category, p.ID, p.AltName)
// 				}
// 			}

// 			if p.URL != "" {
// 				sm.Add(stm.URL{{"loc", domainPrefix + p.URL}, {"changefreq", "daily"}})
// 			} else {
// 				emptyURLQty++
// 			}
// 		}

// 		log.Println("Empty URL qty", emptyURLQty)
// 	}

// 	sm.Finalize()
// 	// sm.PingSearchEngines()
// 	//data = sm.XMLContent()

// 	return
// }

func NewService(dbService *database.Service, updatePeriod int) (s *Service, err error) {
	s = &Service{
		dbService:    dbService,
		updatePeriod: time.Duration(updatePeriod),
	}
	err = s.loadData()
	go s.loadWorker()
	return
}

func (s *Service) loadWorker() {
	for {
		time.Sleep(time.Second * s.updatePeriod)
		if err := s.loadData(); err != nil {
			log.Println(err)
		}
	}
}

func (s *Service) loadData() (err error) {
	// load posts
	posts := []*database.Post{}
	if posts, err = s.dbService.PostsAll(); err != nil {
		log.Println("Cannot load posts", err)
		return err
	} else {
		log.Println("loadData: loaded posts ", len(posts))
	}

	// load categories
	cats := []*database.Category{}
	if cats, err = s.dbService.CatsAll(); err != nil {
		log.Println("Cannot load cats", err)
		return err
	} else {
		log.Println("loadData: loaded cats ", len(cats))
	}

	// load domains
	domains := []*database.FlixDomain{}
	if domains, err = s.dbService.DomainsAll(); err != nil {
		log.Println("Cannot load domains", err)
		return err
	} else {
		log.Println("loadData: loaded domains ", len(domains))
	}

	// generate sitemap for each domain

	for _, d := range domains {

		flixPostAltNames := s.dbService.FlixPostAltNames(d.ID)

		dom := d.HostPublic

		tmpFolder := os.TempDir() + "sitemap-generator/" + dom
		os.MkdirAll(tmpFolder, 0777)

		sm := stm.NewSitemap(1)
		sm.SetCompress(false)
		sm.SetSitemapsPath("")
		sm.SetPublicPath(tmpFolder)

		sm.Create()
		sm.SetDefaultHost("https://" + dom)
		sm.Add(stm.URL{{"loc", ""}, {"changefreq", "always"}})
		domainPrefix := "https://" + dom

		if d.PostID == 0 { // generate sitemap for all posts
			// generate categories
			if rootCats, err := s.dbService.Cats(0); err != nil {
				return err
			} else {
				for _, rc := range rootCats {
					sm.Add(stm.URL{{"loc", "/" + rc.AltName}, {"changefreq", "daily"}})
					//log.Printf("parent_id: %+v", rc.Parentid)
					if cats, err := s.dbService.Cats(rc.ID); err != nil {
						return err
					} else {
						for _, c := range cats {
							sm.Add(stm.URL{{"loc", domainPrefix + "/" + rc.AltName + "/" + c.AltName}, {"changefreq", "daily"}})
						}
					}
				}
			}

			// generate posts
			for _, p := range posts {
				// movies/komediya/21-kapkarashka-kubinskaja-istorija.html
				u := ""
				if altName, err := s.dbService.FlixPostFindAltName(flixPostAltNames, p.ID); err == nil {
					u = s.dbService.MakeUrl(cats, p.Category, p.ID, altName)
				} else {
					u = s.dbService.MakeUrl(cats, p.Category, p.ID, p.AltName)
				}
				if u != "" {
					sm.Add(stm.URL{{"loc", domainPrefix + u}, {"changefreq", "daily"}})
				}
			}

		} else { // generate sitemap for specific post
			for _, p := range posts {
				if p.ID == d.PostID {
					u := ""
					if altName, err := s.dbService.FlixPostFindAltName(flixPostAltNames, p.ID); err == nil {
						u = s.dbService.MakeUrl(cats, p.Category, p.ID, altName)
					} else {
						u = s.dbService.MakeUrl(cats, p.Category, p.ID, p.AltName)
					}
					if u != "" {
						sm.Add(stm.URL{{"loc", domainPrefix + u}, {"changefreq", "daily"}})
					}

					// add post external data
					if postExternalJson, err := s.dbService.FlixPostExternalGetOne(p.ID); err != nil {
						log.Println("Cannot load flix post external for post ID", p.ID, err)
					} else {
						for _, season := range postExternalJson.Seasons {
							sm.Add(stm.URL{{"loc", domainPrefix + "/season/" + helper.IntToString(season.SeasonNumber) + ".html"}, {"changefreq", "daily"}})
							for _, episode := range season.Episodes {
								sm.Add(stm.URL{{"loc", domainPrefix + "/season/" + helper.IntToString(season.SeasonNumber) + "/episode/" + helper.IntToString(episode.EpisodeNumber) + ".html"}, {"changefreq", "daily"}})
							}
						}
					}
				}
			}
		}

		sm.Finalize()

		targetFolder := os.Getenv("STORAGE_PATH") + "/" + dom

		// cant use rename because of different file systems
		if err := helper.CopyDir(tmpFolder, targetFolder); err != nil {
			log.Printf("cant copy %s to %s: %s\n", tmpFolder, tmpFolder, err)
		}

		if err := os.RemoveAll(tmpFolder); err != nil {
			log.Printf("cant delete %s: %v", tmpFolder, err)
		}

		//sm.PingSearchEngines()

	}

	return
}
