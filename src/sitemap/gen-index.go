package sitemap

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type SmIndex struct {
	f            *os.File
	w            *bufio.Writer
	Lastmod      string
	Domain       string
	FileName     string
	FileFullName string
}

func (sf *SmIndex) Init(domain, outDir, fileName string) (err error) {
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
	fmt.Fprint(sf.w, `<sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`)
	return nil
}

func (sf *SmIndex) Add(name, lastmod string) error {
	if lastmod == "" {
		lastmod = sf.Lastmod
	}
	fmt.Fprintf(sf.w, "<sitemap><loc>https://%s/%s</loc><lastmod>%s</lastmod></sitemap>", sf.Domain, name, lastmod)
	return nil
}

func (sf *SmIndex) Close() error {
	defer sf.f.Close()
	fmt.Fprintln(sf.w, `</sitemapindex>`)
	// log.Print("Finished file", sf.FileFullName)
	return sf.w.Flush()
}
