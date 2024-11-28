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

func (s *Service) Posts(domainID int) (posts []*Post, err error) {

	if err = s.DB.Where("approve=?", 1).Find(&posts).Error; err != nil {
		log.Println("Cannot load posts", err)
	}

	cats := []*Category{}
	if err = s.DB.Where("active=?", 1).Find(&cats).Error; err != nil {
		log.Println("Cannot load categories", err)
	}

	flixPostAltNames := s.FlixPostAltNames(domainID)

	for i, p := range posts {
		if altName, err := s.FlixPostFindAltName(flixPostAltNames, p.ID); err == nil {
			posts[i].URL = s.makeUrl(cats, p.Category, p.ID, altName)
		} else {
			posts[i].URL = s.makeUrl(cats, p.Category, p.ID, p.AltName)
		}
	}

	return
}

func (s *Service) makeUrl(catsAll []*Category, catsPostStr string, postId int, altName string) (res string) {

	catsPost := strings.Split(catsPostStr, ",")
	if len(catsPost) == 0 {
		return
	}
	if len(catsPost) == 2 { // category, genre
		category := helper.StrToInt(catsPost[0])
		genre := helper.StrToInt(catsPost[1])
		categorySlug := s.catSlug(catsAll, category)
		genreSlug := s.catSlug(catsAll, genre)
		res = fmt.Sprintf("/%s/%s/%d-%s.html", categorySlug, genreSlug, postId, altName)
		return
	}

	if len(catsPost) == 1 { // category without genre
		category := helper.StrToInt(catsPost[0])
		categorySlug := s.catSlug(catsAll, category)
		res = fmt.Sprintf("/%s/%d-%s.html", categorySlug, postId, altName)
		return
	}

	// cats := []*Category{}
	// firstCat := 0
	// postCats := strings.Split(p.Category, ",")
	// if len(postCats) == 0 {
	// 	return
	// }
	// firstCat = helper.StrToInt(postCats[0])

	// catAlt := ""
	// catParentAlt := ""
	// parentCatId := 0
	// for _, c := range cats {
	// 	if c.ID == catId {
	// 		parentCatId = c.Parentid
	// 		catAlt = c.AltName
	// 	}
	// }
	// for _, c := range cats {
	// 	if c.ID == parentCatId {
	// 		catParentAlt = c.AltName
	// 	}
	// }
	// res = fmt.Sprintf("%s/%s/%d-%s.html", catParentAlt, catAlt, postId, altName)
	// return
	return
}

func (s *Service) catSlug(catsAll []*Category, catId int) (slug string) {
	for _, c := range catsAll {
		if c.ID == catId {
			return c.AltName
		}
	}
	return ""
}
