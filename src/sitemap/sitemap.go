package sitemap

import (
	"dle-sitemap/database"
	"dle-sitemap/helper"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/ikeikeikeike/go-sitemap-generator/v2/stm"
)

type Service struct {
	dbService    *database.Service
	updatePeriod time.Duration
}

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
		domainPrefix := "https://" + dom
		tmpFolder := os.TempDir() + "sitemap-generator/" + dom
		os.MkdirAll(tmpFolder, 0755)

		smStatic := stm.NewSitemap(0)
		smStatic.SetCompress(false)
		smStatic.SetDefaultHost("https://" + dom)
		smStatic.SetPublicPath(tmpFolder)
		smStatic.SetSitemapsPath("./sitemap")
		smStatic.SetFilename("static_pages")
		smStatic.Create()
		smStatic.Add(stm.URL{{"loc", "/"}, {"changefreq", "daily"}, {"priority", "1.0"}})

		smCats := stm.NewSitemap(0)
		smCats.SetCompress(false)
		smCats.SetDefaultHost("https://" + dom)
		smCats.SetPublicPath(tmpFolder)
		smCats.SetSitemapsPath("./sitemap")
		smCats.SetFilename("category_pages")
		smCats.Create()

		smPosts := stm.NewSitemap(0)
		smPosts.SetCompress(false)
		smPosts.SetDefaultHost("https://" + dom)
		smPosts.SetPublicPath(tmpFolder)
		smPosts.SetSitemapsPath("./sitemap")
		smPosts.SetFilename("pages")
		smPosts.Create()

		if d.PostID == 0 { // generate sitemap for all posts

			if rootCats, err := s.dbService.Cats(0); err != nil {
				return err
			} else {
				for _, rc := range rootCats {
					smCats.Add(stm.URL{{"loc", "/" + rc.AltName}, {"changefreq", "weekly"}, {"priority", "0.7"}})
					//log.Printf("parent_id: %+v", rc.Parentid)
					if cats, err := s.dbService.Cats(rc.ID); err != nil {
						return err
					} else {
						for _, c := range cats {
							smCats.Add(stm.URL{{"loc", domainPrefix + "/" + rc.AltName + "/" + c.AltName}, {"changefreq", "weekly"}, {"priority", "0.7"}})
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
					smPosts.Add(stm.URL{{"loc", domainPrefix + u}, {"changefreq", "weekly"}, {"priority", "0.9"}})
				}
			}

		} else { // generate sitemap for specific post
			log.Println("Generating sitemap for post ID:", d.PostID)
			for _, p := range posts {
				if p.ID == d.PostID {
					log.Println("Found post ID:", p.ID)
					u := ""
					if altName, err := s.dbService.FlixPostFindAltName(flixPostAltNames, p.ID); err == nil {
						u = s.dbService.MakeUrl(cats, p.Category, p.ID, altName)
					} else {
						u = s.dbService.MakeUrl(cats, p.Category, p.ID, p.AltName)
					}
					if u != "" {
						smPosts.Add(stm.URL{{"loc", domainPrefix + u}, {"changefreq", "weekly"}, {"priority", "0.9"}})
					}

					// add post external data
					if postExternalJson, err := s.dbService.FlixPostExternalGetOne(p.ID); err != nil {
						log.Println("Cannot load flix post external for post ID", p.ID, err)
					} else {
						log.Println("Loaded flix post external for post ID", p.ID)
						for _, season := range postExternalJson.Seasons {
							log.Println("Adding season to sitemap:", season.SeasonNumber)
							smCats.Add(stm.URL{{"loc", domainPrefix + "/season" + helper.IntToString(season.SeasonNumber)}, {"changefreq", "weekly"}, {"priority", "0.7"}})
							for _, episode := range season.Episodes {
								smCats.Add(stm.URL{{"loc", domainPrefix + "/season" + helper.IntToString(season.SeasonNumber) + "/episode" + helper.IntToString(episode.EpisodeNumber)}, {"changefreq", "weekly"}, {"priority", "0.7"}})
							}
						}
					}
				}
			}
		}

		smStatic.Finalize()
		smCats.Finalize()
		smPosts.Finalize()

		// сделать строго 1 файл - вместо sm.Finalize():
		// data := smStatic.XMLContent()
		// if err := os.WriteFile(filepath.Join(tmpFolder, "static_pages.xml"), data, 0644); err != nil {
		// 	log.Fatal(err)
		// }

		// create index file manually
		indexPath := filepath.Join(tmpFolder, "sitemap.xml")
		indexFile, _ := os.Create(indexPath)
		defer indexFile.Close()

		indexFile.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
		indexFile.WriteString(`<sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">` + "\n")

		files := []string{"static_pages.xml", "category_pages.xml", "pages.xml"}
		for _, f := range files {
			fmt.Fprintf(indexFile, "  <sitemap><loc>%s/sitemap/%s</loc></sitemap>\n", domainPrefix, f)
		}

		indexFile.WriteString(`</sitemapindex>`)
		fmt.Println("✅ Sitemaps generated for ", dom)

		targetFolder := os.Getenv("STORAGE_PATH") + "/" + dom

		// cant use rename because of different file systems
		if err := helper.CopyDir(tmpFolder, targetFolder); err != nil {
			log.Printf("cant copy %s to %s: %s\n", tmpFolder, tmpFolder, err)
		}

		if err := os.RemoveAll(tmpFolder); err != nil {
			log.Printf("cant delete %s: %v", tmpFolder, err)
		}

		smPosts.PingSearchEngines()

	}

	return
}
