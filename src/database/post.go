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

func (s *Service) PostsAll() (posts []*Post, err error) {

	if err = s.DB.Where("approve=?", 1).Find(&posts).Error; err != nil {
		log.Println("Cannot load posts", err)
	}
	return
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
			posts[i].URL = s.MakeUrl(cats, p.Category, p.ID, altName)
		} else {
			posts[i].URL = s.MakeUrl(cats, p.Category, p.ID, p.AltName)
		}
	}

	return
}

func (s *Service) MakeUrl(catsAll []*Category, catsPostStr string, postId int, altName string) (res string) {

	catsPost := strings.Split(catsPostStr, ",")
	if len(catsPost) == 0 {
		return
	}

	genreID := helper.StrToInt(catsPost[0])

	genre := &Category{}
	for _, c := range catsAll {
		if c.ID == genreID {
			genre = c
		}
	}

	mainCat := &Category{}
	for _, c := range catsAll {
		if c.ID == genre.Parentid {
			mainCat = c
		}
	}

	res = fmt.Sprintf("/%s/%s/%d-%s.html", mainCat.AltName, genre.AltName, postId, altName)
	return
}
