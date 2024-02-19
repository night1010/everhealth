package dto

import (
	"github.com/night1010/everhealth/entity"
	"github.com/night1010/everhealth/valueobject"
	"github.com/shopspring/decimal"
)

type GetAddressResponse struct {
	Id         uint   `json:"id"`
	Name       string `json:"name"`
	StreetName string `json:"street"`
	PostalCode string `json:"postal_code"`
	Phone      string `json:"phone"`
	Detail     string `json:"detail"`
	Province   string `json:"province"`
	City       string `json:"city"`
	IsDefault  bool   `json:"is_default"`
	CityId     uint   `json:"city_id"`
	ProvinceId uint   `json:"province_id"`
	Latitude   string `json:"latitude"`
	Longitude  string `json:"longitude"`
}

type CreateAddressRequest struct {
	Name       string `binding:"required" json:"name"`
	StreetName string `binding:"required" json:"street"`
	PostalCode string `binding:"required" json:"postal_code"`
	Phone      string `binding:"required,phonenumberprefix,phonenumberlength" json:"phone"`
	Detail     string `json:"detail"`
	ProvinceId uint   `binding:"required" json:"province_id"`
	CityId     uint   `binding:"required" json:"city_id"`
	Latitude   string `binding:"required,latitude" json:"latitude"`
	Longitude  string `binding:"required,longitude" json:"longitude"`
}

type UpdateAddressRequest struct {
	Name       string `binding:"required" json:"name"`
	StreetName string `binding:"required" json:"street"`
	PostalCode string `binding:"required" json:"postal_code"`
	Phone      string `binding:"required,phonenumberprefix,phonenumberlength" json:"phone"`
	Detail     string `json:"detail"`
	ProvinceId uint   `binding:"required" json:"province_id"`
	CityId     uint   `binding:"required" json:"city_id"`
	Latitude   string `binding:"required,latitude" json:"latitude"`
	Longitude  string `binding:"required,longitude" json:"longitude"`
}

type ShippingMethodResponse struct {
	Name     string `json:"name"`
	Duration string `json:"duration"`
	Cost     string `json:"cost"`
}

type CreateAddressResponse struct {
	AddressId uint `json:"address_id"`
}

type AddressUri struct {
	Id uint `uri:"id" binding:"required,numeric"`
}

func (r *CreateAddressRequest) ToProfile() *entity.Profile {
	return &entity.Profile{
		Name: r.Name,
	}
}

func (r *CreateAddressRequest) ToAddress() *entity.Address {
	address := &entity.Address{
		Name:       r.Name,
		StreetName: r.StreetName,
		PostalCode: r.PostalCode,
		Phone:      r.Phone,
		Detail:     r.Detail,
		ProvinceId: r.ProvinceId,
		CityId:     r.CityId,
	}

	latitude, _ := decimal.NewFromString(r.Latitude)
	longitude, _ := decimal.NewFromString(r.Longitude)
	address.Location = &valueobject.Coordinate{
		Latitude:  latitude,
		Longitude: longitude,
	}
	return address
}

func (r *UpdateAddressRequest) ToAddress() *entity.Address {
	address := &entity.Address{
		Name:       r.Name,
		StreetName: r.StreetName,
		PostalCode: r.PostalCode,
		Phone:      r.Phone,
		Detail:     r.Detail,
		ProvinceId: r.ProvinceId,
		CityId:     r.CityId,
	}

	latitude, _ := decimal.NewFromString(r.Latitude)
	longitude, _ := decimal.NewFromString(r.Longitude)

	address.Location = &valueobject.Coordinate{
		Latitude:  latitude,
		Longitude: longitude,
	}
	return address
}

type ValidateAddressRequest struct {
	CityId    uint   `json:"city_id" binding:"required,numeric,min=1"`
	Latitude  string `json:"latitude" binding:"required,latitude"`
	Longitude string `json:"longitude" binding:"required,longitude"`
}

func (r *ValidateAddressRequest) ToCoordinate() (*valueobject.Coordinate, error) {
	lat, err := decimal.NewFromString(r.Latitude)
	if err != nil {
		return nil, err
	}

	lng, err := decimal.NewFromString(r.Longitude)
	if err != nil {
		return nil, err
	}

	return &valueobject.Coordinate{
		Latitude:  lat,
		Longitude: lng,
	}, nil
}

type ValidateAddressResponse struct {
	IsValid bool `json:"is_valid"`
}
