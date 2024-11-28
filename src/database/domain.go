package database

import "log"

type FlixDomain struct {
	ID         int
	HostPublic string
}

func (c *FlixDomain) TableName() string {
	return "flix_domain"
}

func (s *Service) FlixDomainIDByHost(host string) (id int) {
	res := &FlixDomain{}
	if err := s.DB.Where("host_public=?", host).Find(&res).Error; err != nil {
		log.Println("Cannot load domain", err)
	}
	return res.ID
}
