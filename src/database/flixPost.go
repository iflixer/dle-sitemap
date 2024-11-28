package database

import (
	"errors"
	"log"
)

type FlixPost struct {
	ID       int
	DomainID int
	PostID   int
	Story    string
	AltName  string
}

func (c *FlixPost) TableName() string {
	return "flix_post"
}

func (s *Service) FlixPostAltNames(domainID int) (posts []*FlixPost) {
	if err := s.DB.Where("domain_id=? AND !isnull(alt_name)", domainID).Find(&posts).Error; err != nil {
		log.Println("Cannot load flix posts", err)
	}
	return
}

func (s *Service) FlixPostFindAltName(flixPosts []*FlixPost, postId int) (altName string, err error) {
	for _, c := range flixPosts {
		if c.PostID == postId {
			return c.AltName, nil
		}
	}
	return "", errors.New("no overrides")
}
