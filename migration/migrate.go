package migration

import (
	"github.com/night1010/everhealth/entity"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	u := &entity.User{}
	drugForm := &entity.DrugForm{}
	pc := &entity.ProductCategory{}
	p := &entity.Product{}
	d := &entity.Drug{}
	dc := &entity.DrugClassification{}
	r := &entity.Role{}
	dp := &entity.DoctorProfile{}
	ds := &entity.DoctorSpecialist{}
	profile := &entity.Profile{}
	ftp := &entity.ForgotPasswordToken{}
	pharmacy := &entity.Pharmacy{}
	pharmacyProduct := &entity.PharmacyProduct{}
	pr := &entity.Province{}
	ct := &entity.City{}
	ors := &entity.OrderStatus{}
	spm := &entity.ShippingMethod{}
	a := &entity.Address{}
	c := &entity.Cart{}
	ci := &entity.CartItem{}
	ts := &entity.StockRecord{}
	sm := &entity.StockMutation{}
	po := &entity.ProductOrder{}
	oi := &entity.OrderItem{}
	ac := &entity.AdminContact{}
	telemedicine := &entity.Telemedicine{}
	chat := &entity.Chat{}
	// _ = db.SetupJoinTable(pharmacy, "Products", pharmacyProduct)

	_ = db.Migrator().DropTable(u, drugForm, pc, p, d, dc, r, dp, ds, profile, ftp, pharmacy, pharmacyProduct, pr, ct, ors, spm, a, c, ci, ts, sm, ac, po, oi, telemedicine, chat)

	_ = db.AutoMigrate(u, drugForm, pc, p, d, dc, r, dp, ds, profile, ftp, pharmacy, pharmacyProduct, pr, ct, ors, spm, a, c, ci, ts, sm, ac, po, oi, telemedicine, chat)
}
