package dto

import "github.com/night1010/everhealth/entity"

type CityRes struct {
	Id         uint   `json:"id"`
	Name       string `json:"name"`
	Code       string `json:"code"`
	Latitude   string `json:"latitude"`
	Longitude  string `json:"longitude"`
	ProvinceId uint   `json:"province_id,omitempty"`
}

type ProvinceRes struct {
	Id     uint       `json:"id"`
	Name   string     `json:"name"`
	Code   string     `json:"code"`
	Cities []*CityRes `json:"cities,omitempty"`
}

func NewProvinceRes(p *entity.Province) *ProvinceRes {
	citiesRes := []*CityRes{}
	for _, city := range p.Cities {
		cityRes := NewCityProvinceRes(city)
		citiesRes = append(citiesRes, cityRes)
	}
	return &ProvinceRes{Id: p.Id, Name: p.Name, Code: p.Code, Cities: citiesRes}
}

func NewCityProvinceRes(c *entity.City) *CityRes {
	return &CityRes{
		Id:        c.Id,
		Name:      c.Name,
		Code:      c.Code,
		Latitude:  c.Location.Latitude.String(),
		Longitude: c.Location.Longitude.String(),
	}
}

type ProvinceUri struct {
	Id int64 `uri:"id" binding:"required,numeric"`
}
