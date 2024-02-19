package migration

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/night1010/everhealth/entity"
	"github.com/night1010/everhealth/logger"
	"github.com/night1010/everhealth/valueobject"
	"github.com/shopspring/decimal"
)

type productCategory struct {
	Id       uint   `json:"id"`
	Name     string `json:"name"`
	IsDrug   bool   `json:"is_drug"`
	Image    string `json:"image"`
	ImageKey string `json:"image_key"`
}

type product struct {
	Id                uint   `json:"id"`
	Name              string `json:"name"`
	Manufacture       string `json:"manufacture"`
	Detail            string `json:"detail"`
	ProductCategoryId uint   `json:"product_category_id"`
	UnitInPack        string `json:"unit_in_pack"`
	Price             string `json:"price"`
	SellingUnit       string `json:"selling_unit"`
	Weight            string `json:"weight"`
	Height            string `json:"height"`
	Length            string `json:"length"`
	Width             string `json:"width"`
	Image             string `json:"image"`
	ImageKey          string `json:"image_key"`
}

type city struct {
	Gid        uint            `json:"gid"`
	CityName   string          `json:"city_name"`
	CityId     string          `json:"city_id"`
	ProvinceId string          `json:"province_id"`
	Province   string          `json:"province"`
	Latitude   decimal.Decimal `json:"latitude"`
	Longitude  decimal.Decimal `json:"longitude"`
}

type pharmacy struct {
	Name       string  `json:"name"`
	AdminId    uint    `json:"admin_id"`
	Address    string  `json:"address"`
	CityId     uint    `json:"city_id"`
	ProvinceId uint    `json:"province_id"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
}

func ImportProduct() []*entity.Product {
	var p []*product
	var products []*entity.Product
	data, _ := os.ReadFile("./migration/data/product-details.json")

	_ = json.Unmarshal(data, &p)

	for _, p2 := range p {
		price, _ := decimal.NewFromString(p2.Price)
		width, _ := decimal.NewFromString(p2.Width)
		length, _ := decimal.NewFromString(p2.Length)
		height, _ := decimal.NewFromString(p2.Height)
		weight, _ := decimal.NewFromString(p2.Weight)
		products = append(products, &entity.Product{
			Name:              p2.Name,
			Manufacture:       p2.Manufacture,
			Detail:            p2.Detail,
			ProductCategoryId: p2.ProductCategoryId,
			UnitInPack:        p2.UnitInPack,
			Price:             price,
			SellingUnit:       p2.SellingUnit,
			Weight:            weight,
			Height:            height,
			Length:            length,
			Width:             width,
			Image:             p2.Image,
			ImageKey:          p2.ImageKey,
		})
	}

	return products
}

func ImportProductCategories() []*entity.ProductCategory {
	var c []*productCategory
	var categories []*entity.ProductCategory
	data, err := os.ReadFile("./migration/data/product-categories.json")
	if err != nil {
		log.Println(err)
	}

	_ = json.Unmarshal(data, &c)

	for _, category := range c {
		categories = append(categories, &entity.ProductCategory{
			Name:     category.Name,
			IsDrug:   category.IsDrug,
			Image:    category.Image,
			ImageKey: category.ImageKey,
		})
	}

	return categories
}

func ImportCities() []*entity.City {
	var jsonCities []*city
	var cities []*entity.City

	data, err := os.ReadFile("./migration/data/cities.json")
	if err != nil {
		logger.Log.Error(err)
	}

	_ = json.Unmarshal(data, &jsonCities)

	for _, jsonCity := range jsonCities {
		provinceId, _ := strconv.Atoi(jsonCity.ProvinceId)
		cities = append(cities, &entity.City{
			CityGid:    jsonCity.Gid,
			Name:       jsonCity.CityName,
			Code:       jsonCity.CityId,
			ProvinceId: uint(provinceId),
			Location: &valueobject.Coordinate{
				Latitude:  jsonCity.Latitude,
				Longitude: jsonCity.Longitude,
			},
		})
	}
	return cities
}

func ImportPharmacies() []*entity.Pharmacy {
	var p []*pharmacy
	var pharmacies []*entity.Pharmacy

	data, err := os.ReadFile("./migration/data/pharmacies.json")
	if err != nil {
		logger.Log.Error(err)
	}

	_ = json.Unmarshal(data, &p)

	for _, p2 := range p {
		pharmacies = append(pharmacies, &entity.Pharmacy{
			Name:       p2.Name,
			AdminId:    p2.AdminId,
			Address:    p2.Address,
			CityId:     p2.CityId,
			ProvinceId: p2.ProvinceId,
			Location: &valueobject.Coordinate{
				Latitude:  decimal.NewFromFloat(p2.Latitude),
				Longitude: decimal.NewFromFloat(p2.Longitude),
			},
			StartTime:               time.Date(1, 1, 1, 10, 0, 0, 0, time.Local),
			EndTime:                 time.Date(1, 1, 1, 23, 59, 0, 0, time.Local),
			OperationalDay:          "Monday,Tuesday,Wednesday,Saturday",
			PharmacistLicenseNumber: "9999-1010-2020-0101",
			PharmacistPhoneNumber:   "089654749370",
		})
	}

	return pharmacies
}
