package valueobject

import (
	"context"
	"fmt"

	"github.com/shopspring/decimal"
	"github.com/twpayne/go-geom/encoding/ewkbhex"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Coordinate struct {
	Latitude  decimal.Decimal
	Longitude decimal.Decimal
}

func (c Coordinate) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	return clause.Expr{
		SQL:  "ST_MakePoint(?,?)",
		Vars: []interface{}{c.Longitude, c.Latitude},
	}
}

func (c *Coordinate) Scan(input any) error {
	point, ok := input.(string)
	if !ok {
		return fmt.Errorf("unexpected type for location: %T", input)
	}

	location := ""
	_, err := fmt.Sscanf(point, "%s", &location)
	if err != nil {
		return err
	}
	a, _ := ewkbhex.Decode(location)
	c.Longitude = decimal.NewFromFloat(a.FlatCoords()[0])
	c.Latitude = decimal.NewFromFloat(a.FlatCoords()[1])
	return err
}
