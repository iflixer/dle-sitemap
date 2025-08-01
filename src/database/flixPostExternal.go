package database

import (
	"encoding/json"
	"log"
)

type FlixPostExternal struct {
	ID     int
	PostID int
	Json   string
}

type FlixPostExternalJson struct {
	Name         string `json:"name"`
	OriginalName string `json:"original_name"`
	Seasons      []struct {
		SeasonNumber int     `json:"season_number"`
		Title        string  `json:"title"`
		CoverURL     string  `json:"cover_url"`
		Description  string  `json:"description"`
		VoteAverage  float64 `json:"vote_average"`
		Episodes     []struct {
			EpisodeNumber int     `json:"episode_number"`
			Title         string  `json:"title"`
			CoverURL      string  `json:"cover_url"`
			Description   string  `json:"description"`
			Runtime       int     `json:"runtime"`
			VoteAverage   float64 `json:"vote_average"`
		} `json:"episodes"`
	} `json:"seasons"`
}

func (c *FlixPostExternal) TableName() string {
	return "flix_post_external"
}

func (s *Service) FlixPostExternalGetOne(postID int) (postExternalJson FlixPostExternalJson, err error) {
	var post FlixPostExternal
	if err := s.DB.Where("id=?", postID).First(&post).Error; err != nil {
		log.Println("Cannot load flix post external", err)
	}
	if err := s.DB.Where("post_id=?", post.PostID).First(&post).Error; err != nil {
		log.Println("Cannot load flix post external by post ID", err)
		return postExternalJson, err
	}
	if err := json.Unmarshal([]byte(post.Json), &postExternalJson); err != nil {
		log.Println("Cannot unmarshal flix post external JSON", err)
		return postExternalJson, err
	}
	return
}
