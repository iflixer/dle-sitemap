package sitemap

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type SmSitemap struct {
	f            *os.File
	w            *bufio.Writer
	Lastmod      string
	Domain       string
	FileName     string
	FileFullName string
}

type SmSitemapRow struct {
	Lastmod    string
	Loc        string
	ChangeFreq string
	Priority   string
}

func (sf *SmSitemap) Init(domain, outDir, fileName string) (err error) {
	sf.Domain = domain
	sf.FileName = fileName
	tz := time.FixedZone("EET", 2*60*60)
	sf.Lastmod = time.Date(2025, 10, 27, 3, 30, 3, 0, tz).Format(time.RFC3339)
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}
	sf.FileFullName = filepath.Join(outDir, fileName)
	// log.Println("Opening to write file", sf.FileFullName)
	sf.f, err = os.Create(sf.FileFullName)
	if err != nil {
		return err
	}
	sf.w = bufio.NewWriter(sf.f)
	fmt.Fprintln(sf.w, `<?xml version="1.0" encoding="UTF-8"?>`)
	fmt.Fprint(sf.w, `<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`)
	return nil
}

func (sf *SmSitemap) Add(row SmSitemapRow) error {
	if row.Lastmod == "" {
		row.Lastmod = sf.Lastmod
	}
	res := fmt.Sprintf(`<url>
		<loc>%s</loc>
		<changefreq>%s</changefreq>
		<lastmod>%s</lastmod>
		<priority>%s</priority>
	</url>`, row.Loc, row.ChangeFreq, row.Lastmod, row.Priority)
	fmt.Fprint(sf.w, sf.removeBadSymbols(res))
	return nil
}

func (sf *SmSitemap) removeBadSymbols(s string) string {
	return strings.Map(func(r rune) rune {
		if r == '\n' || r == '\r' || r == '\t' || r == ' ' {
			return -1 // выкинуть символ
		}
		return r
	}, s)
}

func (sf *SmSitemap) Close() error {
	defer sf.f.Close()
	fmt.Fprint(sf.w, `</urlset>`)
	// log.Println("Finished file", sf.FileFullName)
	return sf.w.Flush()
}
