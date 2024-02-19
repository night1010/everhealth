package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/night1010/everhealth/config"
	"github.com/night1010/everhealth/dto"
	"github.com/night1010/everhealth/entity"
	"github.com/night1010/everhealth/valueobject"
	"github.com/go-resty/resty/v2"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type ShippingMethodRepository interface {
	BaseRepository[entity.ShippingMethod]
	GetThirdPartyShipping(ctx context.Context, origin, destination, weight string) ([]*entity.CalculatedShippingMethod, error)
	FindDistanceBetween(ctx context.Context, coordinate1 *valueobject.Coordinate, coordinate2 *valueobject.Coordinate) (decimal.Decimal, error)
}

type shippingMethodRepository struct {
	*baseRepository[entity.ShippingMethod]
	db         *gorm.DB
	httpClient *resty.Client
}

func NewShippingMethodRepository(db *gorm.DB, client *resty.Client) ShippingMethodRepository {
	return &shippingMethodRepository{
		db:             db,
		baseRepository: &baseRepository[entity.ShippingMethod]{db: db},
		httpClient:     client,
	}
}

func (r *shippingMethodRepository) GetThirdPartyShipping(ctx context.Context, origin, destination, weight string) ([]*entity.CalculatedShippingMethod, error) {
	var result dto.RajaOngkirResult
	rajaOngkirCfg := config.NewRajaOngkirConfig()
	url := "https://api.rajaongkir.com/starter/cost"
	body, _ := json.Marshal(map[string]string{
		"origin":      origin,
		"destination": destination,
		"weight":      weight,
		"courier":     "jne",
	})

	var shippingMethods []*entity.CalculatedShippingMethod
	_, err := r.httpClient.R().
		SetHeader("key", rajaOngkirCfg.Token).
		SetBody(body).
		SetResult(&result).
		Post(url)
	if err != nil {
		return nil, err
	}

	if len(result.RajaOngkir.Results) < 1 {
		return nil, nil
	}

	for _, provider := range result.RajaOngkir.Results[0].Costs {
		shippingMethods = append(shippingMethods, &entity.CalculatedShippingMethod{
			Name:              fmt.Sprintf("%s | %s", provider.Services, provider.Description),
			EstimatedDuration: provider.Costs[0].EstimatedTime,
			Cost:              strconv.Itoa(provider.Costs[0].Value),
		})
	}

	return shippingMethods, nil
}

func (r *shippingMethodRepository) FindDistanceBetween(ctx context.Context, coordinate1 *valueobject.Coordinate, coordinate2 *valueobject.Coordinate) (decimal.Decimal, error) {
	query := fmt.Sprintf("SELECT st_distance(st_geogfromtext('SRID=4326;POINT(%s %s)'), st_geogfromtext('SRID=4326;POINT(%s %s)'))/1000", coordinate1.Longitude.String(), coordinate1.Latitude.String(), coordinate2.Longitude.String(), coordinate2.Latitude.String())
	var distance float64
	err := r.db.Raw(query).Scan(&distance).Error
	if err != nil {
		return decimal.Decimal{}, err
	}
	return decimal.NewFromFloat(distance), nil
}
