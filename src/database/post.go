package database

import (
	"dle-sitemap/helper"
	"fmt"
	"log"
	"strings"
)

type Post struct {
	ID       int
	AltName  string
	Category string
	URL      string
}

func (c *Post) TableName() string {
	return "dle_post"
}

func (s *Service) Posts() (res []*Post, err error) {

	if err = s.DB.Where("approve=?", 1).Find(&res).Error; err != nil {
		log.Println("Cannot load posts", err)
	}

	cats := []*Category{}
	if err = s.DB.Where("active=?", 1).Find(&cats).Error; err != nil {
		log.Println("Cannot load categories", err)
	}

	for i, p := range res {
		firstCat := 0
		postCats := strings.Split(p.Category, ",")
		if len(postCats) == 0 {
			continue
		}
		firstCat = helper.StrToInt(postCats[0])

		res[i].URL = s.makeUrl(cats, firstCat, p.ID, p.AltName)
	}

	return
}

func (s *Service) makeUrl(cats []*Category, catId int, postId int, altName string) (res string) {
	catAlt := ""
	catParentAlt := ""
	parentCatId := 0
	for _, c := range cats {
		if c.ID == catId {
			parentCatId = c.Parentid
			catAlt = c.AltName
		}
	}
	for _, c := range cats {
		if c.ID == parentCatId {
			catParentAlt = c.AltName
		}
	}
	res = fmt.Sprintf("%s/%s/%d-%s.html", catParentAlt, catAlt, postId, altName)
	return
}
