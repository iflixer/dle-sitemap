package sitemap

import (
	"dle-sitemap/database"
	"dle-sitemap/helper"
	"fmt"
	"log"
	"os"
	"time"
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
		log.Printf("=== Domain %s id %d post_id: %d", d.HostPublic, d.ID, d.PostID)
		flixPostAltNames := s.dbService.FlixPostAltNames(d.ID)
		dom := d.HostPublic
		domainPrefix := "https://" + dom
		tmpFolder := os.TempDir() + "sitemap-generator/" + dom
		// os.MkdirAll(tmpFolder, 0755)

		smIndex := &SmIndex{}
		smIndex.Init(dom, tmpFolder, "sitemap.xml")

		smStatic := &SmSitemap{}
		smStatic.Init(dom, tmpFolder, "sitemap_static.xml")
		smStatic.Add(SmSitemapRow{
			Loc:        domainPrefix + "/",
			ChangeFreq: "daily",
			Priority:   "1.0",
		})

		smNews := &SmSitemap{}
		smNews.Init(dom, tmpFolder, "sitemap_news.xml")

		smPages := &SmSitemap{}
		smPages.Init(dom, tmpFolder, "sitemap_pages.xml")

		smCollections := &SmSitemap{}
		smCollections.Init(dom, tmpFolder, "sitemap_collections.xml")

		smCats := &SmSitemap{}
		smCats.Init(dom, tmpFolder, "sitemap_category.xml")

		smIndex.Add("sitemap_static.xml", "")
		smIndex.Add("sitemap_pages.xml", "")

		if d.PostID == 0 { // generate sitemap for all posts
			smIndex.Add("sitemap_category.xml", "")
			smIndex.Add("sitemap_collections.xml", "")
			if rootCats, err := s.dbService.Cats(0); err != nil {
				return err
			} else {

				for _, rc := range rootCats {
					if rc.ID < 1000 { // categories
						smCats.Add(SmSitemapRow{
							Loc:        domainPrefix + "/" + rc.AltName,
							ChangeFreq: "weekly",
							Priority:   "0.7",
						})
						if cats, err := s.dbService.Cats(rc.ID); err != nil {
							return err
						} else {
							for _, c := range cats {
								smCats.Add(SmSitemapRow{
									Loc:        domainPrefix + "/" + rc.AltName + "/" + c.AltName,
									ChangeFreq: "weekly",
									Priority:   "0.7",
								})
							}
						}
					} else {
						// smCollections.Add(SmSitemapRow{
						// 	Loc:        domainPrefix + "/" + rc.AltName,
						// 	ChangeFreq: "weekly",
						// 	Priority:   "0.7",
						// })
						if cats, err := s.dbService.Cats(rc.ID); err != nil {
							return err
						} else {
							for _, c := range cats {
								smCollections.Add(SmSitemapRow{
									Loc:        domainPrefix + "/" + rc.AltName + "/" + c.AltName,
									ChangeFreq: "daily",
									Priority:   "0.8",
								})
							}
						}
					}
				}

			}

			// posts
			for _, p := range posts {
				u := ""
				if altName, err := s.dbService.FlixPostFindAltName(flixPostAltNames, p.ID); err == nil {
					u = s.dbService.MakeUrl(cats, p.Category, p.ID, altName)
				} else {
					u = s.dbService.MakeUrl(cats, p.Category, p.ID, p.AltName)
				}
				if u != "" {
					smPages.Add(SmSitemapRow{
						Loc:        domainPrefix + u,
						ChangeFreq: "weekly",
						Priority:   "0.9",
					})
				}
			}

		} else { // generate sitemap for specific post
			log.Println("Generating sitemap for post ID:", d.PostID)
			for _, p := range posts {
				if p.ID == d.PostID {
					// add post external data
					if postExternalJson, err := s.dbService.FlixPostExternalGetOne(p.ID); err != nil {
						log.Println("Cannot load flix post external for post ID", p.ID, err)
					} else {
						log.Println("Loaded flix post external for post ID", p.ID)
						for _, season := range postExternalJson.Seasons {
							log.Println("Adding season to sitemap:", season.SeasonNumber)
							smPages.Add(SmSitemapRow{
								Loc:        domainPrefix + "/season" + helper.IntToString(season.SeasonNumber),
								ChangeFreq: "weekly",
								Priority:   "0.9",
							})
							for _, episode := range season.Episodes {
								smPages.Add(SmSitemapRow{
									Loc:        domainPrefix + "/season" + helper.IntToString(season.SeasonNumber) + "/episode" + helper.IntToString(episode.EpisodeNumber),
									ChangeFreq: "weekly",
									Priority:   "0.9",
									Lastmod:    p.UpdatedAt,
								})
							}
						}
					}
				}
			}
		}

		smStatic.Close()
		smIndex.Close()
		smCats.Close()
		smCollections.Close()
		smNews.Close()
		smPages.Close()

		targetFolder := os.Getenv("STORAGE_PATH") + "/" + dom

		// cant use rename because of different file systems
		if err := helper.CopyDir(tmpFolder, targetFolder); err != nil {
			log.Printf("cant copy %s to %s: %s\n", tmpFolder, tmpFolder, err)
		}

		if err := os.RemoveAll(tmpFolder); err != nil {
			log.Printf("cant delete %s: %v", tmpFolder, err)
		}

		fmt.Printf("✅ Sitemaps generated for %s to path %s\n", dom, targetFolder)

		//smPosts.PingSearchEngines()

	}

	return
}

// func saveMainFile(domain, tmpFolder string) {
// 	// Настройки вывода
// 	opts := stm.NewOptions()
// 	opts.SetDefaultHost("https://" + domain) // будет использован при формировании <loc>, если нужно
// 	opts.SetPublicPath(tmpFolder)            // куда писать файлы
// 	opts.SetFilename("sitemap")              // имя файла индекса: sitemap_index.xml
// 	opts.SetPretty(true)

// 	// Конструктор индекс-файла
// 	idx := stm.NewBuilderIndexfile(opts, opts.IndexLocation())

// 	// lastmod как в примере
// 	loc := time.FixedZone("EET", 2*60*60) // +02:00
// 	lastmod := time.Date(2025, 10, 27, 3, 30, 3, 0, loc).Format(time.RFC3339)

// 	add := func(absURL string) {
// 		idx.Add(stm.NewSitemapIndexURL(opts, stm.URL{
// 			{"loc", absURL},
// 			{"lastmod", lastmod},
// 		}))
// 	}

// 	add("https://" + domain + "/static_pages.xml")
// 	add("https://" + domain + "/category_pages.xml")
// 	add("https://" + domain + "/collections.xml")
// 	add("https://" + domain + "/news_pages.xml")
// 	add("https://" + domain + "/pages.xml")

// 	idx.Write()
// }

// func writeIndexManually(domain, outDir string) error {
// 	tz := time.FixedZone("EET", 2*60*60)
// 	lastmod := time.Date(2025, 10, 27, 3, 30, 3, 0, tz).Format(time.RFC3339)

// 	files := []string{
// 		"static_pages.xml",
// 		"category_pages.xml",
// 		"collections.xml",
// 		"news_pages.xml",
// 		"pages.xml",
// 	}

// 	if err := os.MkdirAll(outDir, 0o755); err != nil {
// 		return err
// 	}
// 	f, err := os.Create(filepath.Join(outDir, "sitemap_index.xml"))
// 	if err != nil {
// 		return err
// 	}
// 	defer f.Close()

// 	w := bufio.NewWriter(f)
// 	fmt.Fprintln(w, `<?xml version="1.0" encoding="UTF-8"?>`)
// 	fmt.Fprintln(w, `<sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`)
// 	for _, name := range files {
// 		fmt.Fprintf(w, "  <sitemap>\n    <loc>https://%s/%s</loc>\n    <lastmod>%s</lastmod>\n  </sitemap>\n", domain, name, lastmod)
// 	}
// 	fmt.Fprintln(w, `</sitemapindex>`)
// 	return w.Flush()
// }
