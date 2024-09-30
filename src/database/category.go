package database

import "log"

type Category struct {
	ID       int
	Parentid int
	Name     string
	AltName  string
	Active   int
}

func (c *Category) TableName() string {
	return "dle_category"
}

func (s *Service) Cats(parentId int) (res []*Category, err error) {
	if err = s.DB.Where("active=? AND parentid=?", 1, parentId).Find(&res).Error; err != nil {
		log.Println("Cannot load categories", err)
	}
	return
}
