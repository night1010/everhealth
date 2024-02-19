package migration

import (
	"strconv"
	"time"

	"github.com/night1010/everhealth/entity"
	"github.com/night1010/everhealth/hasher"
	"github.com/night1010/everhealth/valueobject"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	roles := []*entity.Role{
		{Id: entity.RoleUser, Name: "user"},
		{Id: entity.RoleDoctor, Name: "doctor"},
		{Id: entity.RoleAdmin, Name: "pharmacist-admin"},
		{Id: entity.RoleSuperAdmin, Name: "super-admin"},
	}

	users := []*entity.User{
		{Email: "alice@example.com", Password: hashPassword("Alice12345"), RoleId: entity.RoleUser, IsVerified: true},
		{Email: "bob@example.com", Password: hashPassword("Bob12345"), RoleId: entity.RoleDoctor, IsVerified: true},
		{Email: "charlie@example.com", Password: hashPassword("Charlie12345"), RoleId: entity.RoleAdmin, IsVerified: true},
		{Email: "doni@example.com", Password: hashPassword("Doni12345"), RoleId: entity.RoleAdmin, IsVerified: true},
		{Email: "david@example.com", Password: hashPassword("David12345"), RoleId: entity.RoleSuperAdmin, IsVerified: true},
		{Email: "daniel@example.com", Password: hashPassword("Daniel12345"), RoleId: entity.RoleUser, IsVerified: true},
	}
	adminContact := []*entity.AdminContact{
		{UserId: 3, Name: "charlie", Phone: "089654749370"},
		{UserId: 4, Name: "doni", Phone: "089654749370"},
	}

	profiles := []*entity.Profile{
		generateProfile(1, "Alice"),
		generateProfile(2, "Bob"),
		generateProfile(6, "Daniel"),
	}
	doctorProfiles := []*entity.DoctorProfile{
		generateDoctorProfile(2, 2, 50000, entity.Online)}
	users2, adminContact2 := generateAdminPharmacy(len(users), len(users)+17, "admin", "Admin12345")
	var users3 []*entity.User
	for i := 1; i < 49; i++ {
		iString := strconv.Itoa(i)
		temp := generateUser("doctor"+iString+"@gmail.com", "Doctor12345", entity.RoleDoctor, true)
		users3 = append(users3, temp)
	}
	var doctorProfile2 []*entity.DoctorProfile
	for i := 25; i < 49; i++ {
		temp := generateDoctorProfile(uint(i), 1, randomNumber500(20000, 10, 5000), entity.Online)
		doctorProfile2 = append(doctorProfile2, temp)
	}
	var doctorProfile3 []*entity.DoctorProfile
	for i := 49; i < 73; i++ {
		temp := generateDoctorProfile(uint(i), 2, randomNumber500(20000, 10, 5000), entity.Online)
		doctorProfile3 = append(doctorProfile2, temp)
	}
	profile3 := generateRealDoctorProfile()
	users4, profile4, doctorProfile4 := generateFakerDoctor()
	doctorSpecialist := []*entity.DoctorSpecialist{
		{Name: "Spesialis Anak", Image: "https://everhealth-asset.irfancen.com/doctor-specialist/anak.png"},
		{Name: "Spesialis Saraf", Image: "https://everhealth-asset.irfancen.com/doctor-specialist/saraf.png"},
		{Name: "Dokter Umum", Image: "https://everhealth-asset.irfancen.com/doctor-specialist/umum.png"},
		{Name: "Spesialis Kulit", Image: "https://everhealth-asset.irfancen.com/doctor-specialist/kulit.png"},
		{Name: "Spesialis THT", Image: "https://everhealth-asset.irfancen.com/doctor-specialist/tht.png"},
		{Name: "Dokter Gigi", Image: "https://everhealth-asset.irfancen.com/doctor-specialist/gigi.png"},
		{Name: "Psikiater", Image: "https://everhealth-asset.irfancen.com/doctor-specialist/psikiater.png"},
		{Name: "Spesialis Mata", Image: "https://everhealth-asset.irfancen.com/doctor-specialist/mata.png"},
		{Name: "Spesialis Bedah", Image: "https://everhealth-asset.irfancen.com/doctor-specialist/bedah.png"},
		{Name: "Spesialis Penyakit Dalam", Image: "https://everhealth-asset.irfancen.com/doctor-specialist/penyakit-dalam.png"},
	}

	drugForms := []*entity.DrugForm{
		{Name: "tablet"},
		{Name: "kapsul"},
		{Name: "pil"},
		{Name: "serbuk"},
		{Name: "salep"},
		{Name: "krim"},
		{Name: "gel"},
		{Name: "sirup"},
		{Name: "injeksi"},
		{Name: "tetes"},
		{Name: "inhalasi"},
		{Name: "aerosol"},
	}

	productCategories := ImportProductCategories()

	drugClassifications := []*entity.DrugClassification{
		{Name: "obat bebas"},
		{Name: "obat keras"},
		{Name: "obat bebas terbatas"},
	}

	products := ImportProduct()

	drugs := []*entity.Drug{
		{
			ProductId:            1,
			GenericName:          "Tremenza",
			DrugClassificationId: 2,
			DrugFormId:           1,
			Content:              "Pseudoephedrine HCl 60 mg, Triprolidine HCl 2.5 mg",
		},
		{
			ProductId:            2,
			GenericName:          "Rhinos SR",
			DrugClassificationId: 2,
			DrugFormId:           2,
			Content:              "Loratadine 5 mg Pseudoephedrine HCI 120 mg",
		},
	}

	provinces := getProvinces()

	cities := getCities()

	orderStatuses := []*entity.OrderStatus{
		{Name: "waiting for payment"},
		{Name: "waiting for payment confirmation"},
		{Name: "processed"},
		{Name: "sent"},
		{Name: "order confirmed"},
		{Name: "canceled"},
	}

	shippingMethods := []*entity.ShippingMethod{
		{
			Name:       "Official Instant",
			Duration:   "1-2 Hours",
			PricePerKM: decimal.NewFromInt(2500),
		},
		{
			Name:       "Official SameDay",
			Duration:   "8-12 Hours",
			PricePerKM: decimal.NewFromInt(1000),
		},
	}

	addresses := []*entity.Address{
		{
			Name:       "Alice",
			StreetName: "Jalan Mega Kuningan Barat",
			PostalCode: "12950",
			Phone:      "08772348585",
			Detail:     "",
			Location: &valueobject.Coordinate{
				Latitude:  generateDecimalFromString("-6.230835326032342"),
				Longitude: generateDecimalFromString("106.82413596846786"),
			},
			IsDefault:  false,
			ProfileId:  1,
			ProvinceId: 6,
			CityId:     42,
		},
		{
			Name:       "Alice",
			StreetName: "Jalan Cihampelas No 160",
			PostalCode: "12950",
			Phone:      "08772348585",
			Detail:     "",
			Location: &valueobject.Coordinate{
				Latitude:  generateDecimalFromString("-6.894393363416537"),
				Longitude: generateDecimalFromString("107.60773740044303"),
			},
			IsDefault:  false,
			ProfileId:  1,
			ProvinceId: 9,
			CityId:     64,
		},
		{
			Name:       "Alice",
			StreetName: "Jalan Dharmahusada 144",
			PostalCode: "60285",
			Phone:      "08772348585",
			Detail:     "",
			Location: &valueobject.Coordinate{
				Latitude:  generateDecimalFromString("-7.266614744067581"),
				Longitude: generateDecimalFromString("112.77033130376155"),
			},
			IsDefault:  true,
			ProfileId:  1,
			ProvinceId: 11,
			CityId:     161,
		},
	}
	pharmacies := ImportPharmacies()

	carts := []entity.Cart{
		{
			UserId: 1,
		},
	}

	orders := []entity.ProductOrder{
		{
			OrderedAt:     time.Now(),
			OrderStatusId: 1,
			ProfileId:     1,
			ExpiredAt:     time.Now().Add(24 * time.Hour),
			ShippingName:  "Fast",
			ShippingPrice: decimal.NewFromInt(209000),
			ShippingEta:   "1-2 hour",
			AddressId:     1,
			ItemOrderQty:  2,
			TotalPayment:  decimal.NewFromInt(309000),
			PaymentMethod: "NIN",
			CreatedAt:     time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt:     time.Now(),
		},
		{
			OrderedAt:     time.Now(),
			OrderStatusId: 2,
			ProfileId:     1,
			ExpiredAt:     time.Now().Add(24 * time.Hour),
			ShippingName:  "Fast",
			ShippingPrice: decimal.NewFromInt(209000),
			ShippingEta:   "1-2 hour",
			AddressId:     1,
			ItemOrderQty:  2,
			TotalPayment:  decimal.NewFromInt(309000),
			PaymentMethod: "NIN",
			CreatedAt:     time.Date(2024, time.March, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt:     time.Now(),
		},
		{
			OrderedAt:     time.Now(),
			OrderStatusId: 3,
			ProfileId:     1,
			ExpiredAt:     time.Now().Add(24 * time.Hour),
			ShippingName:  "Fast",
			ShippingPrice: decimal.NewFromInt(309000),
			ShippingEta:   "1-2 hour",
			AddressId:     1,
			ItemOrderQty:  2,
			TotalPayment:  decimal.NewFromInt(209000),
			PaymentMethod: "NIN",
			CreatedAt:     time.Date(2024, time.August, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt:     time.Now(),
		},
		{
			OrderedAt:     time.Now(),
			OrderStatusId: 4,
			ProfileId:     1,
			ExpiredAt:     time.Now().Add(24 * time.Hour),
			ShippingName:  "Fast",
			ShippingPrice: decimal.NewFromInt(209000),
			ShippingEta:   "1-2 hour",
			AddressId:     1,
			ItemOrderQty:  1,
			TotalPayment:  decimal.NewFromInt(309000),
			PaymentMethod: "NIN",
			CreatedAt:     time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt:     time.Now(),
		},
		{
			OrderedAt:     time.Now(),
			OrderStatusId: 5,
			ProfileId:     1,
			ExpiredAt:     time.Now().Add(24 * time.Hour),
			ShippingName:  "Fast",
			ShippingPrice: decimal.NewFromInt(209000),
			ShippingEta:   "1-2 hour",
			AddressId:     1,
			ItemOrderQty:  2,
			TotalPayment:  decimal.NewFromInt(309000),
			PaymentMethod: "NIN",
			CreatedAt:     time.Date(2024, time.June, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt:     time.Now(),
		},
		{
			OrderedAt:     time.Now(),
			OrderStatusId: 6,
			ProfileId:     1,
			ExpiredAt:     time.Now().Add(24 * time.Hour),
			ShippingName:  "Fast",
			ShippingPrice: decimal.NewFromInt(209000),
			ShippingEta:   "1-2 hour",
			AddressId:     1,
			ItemOrderQty:  1,
			TotalPayment:  decimal.NewFromInt(309000),
			PaymentMethod: "NIN",
			CreatedAt:     time.Date(2024, time.December, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt:     time.Now(),
		},
		{
			OrderedAt:     time.Now(),
			OrderStatusId: 1,
			ProfileId:     2,
			ExpiredAt:     time.Now().Add(24 * time.Hour),
			ShippingName:  "Fast",
			ShippingPrice: decimal.NewFromInt(209000),
			ShippingEta:   "1-2 hour",
			AddressId:     1,
			ItemOrderQty:  1,
			TotalPayment:  decimal.NewFromInt(309000),
			PaymentMethod: "Mandiri",
			CreatedAt:     time.Date(2024, time.December, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt:     time.Now(),
		},
	}

	orderItems := []entity.OrderItem{
		{
			OrderId:           1,
			PharmacyProductId: 1,
			Quantity:          12,
			SubTotal:          decimal.NewFromInt(1196400),
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		},
		{
			OrderId:           1,
			PharmacyProductId: 2,
			Quantity:          4,
			SubTotal:          decimal.NewFromInt(398800),
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		},
		{
			OrderId:           2,
			PharmacyProductId: 4,
			Quantity:          4,
			SubTotal:          decimal.NewFromInt(132400),
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		},
		{
			OrderId:           2,
			PharmacyProductId: 5,
			Quantity:          4,
			SubTotal:          decimal.NewFromInt(372000),
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		},
		{
			OrderId:           3,
			PharmacyProductId: 3,
			Quantity:          8,
			SubTotal:          decimal.NewFromInt(797600),
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		},
		{
			OrderId:           3,
			PharmacyProductId: 3,
			Quantity:          4,
			SubTotal:          decimal.NewFromInt(102400),
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		},
		{
			OrderId:           4,
			PharmacyProductId: 290,
			Quantity:          4,
			SubTotal:          decimal.NewFromInt(102400),
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		},
		{
			OrderId:           5,
			PharmacyProductId: 6,
			Quantity:          8,
			SubTotal:          decimal.NewFromInt(204800),
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		},
		{
			OrderId:           5,
			PharmacyProductId: 4,
			Quantity:          8,
			SubTotal:          decimal.NewFromInt(548000),
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		},
		{
			OrderId:           6,
			PharmacyProductId: 6,
			Quantity:          4,
			SubTotal:          decimal.NewFromInt(398800),
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		},
		{
			OrderId:           7,
			PharmacyProductId: 10,
			Quantity:          4,
			SubTotal:          decimal.NewFromInt(398800),
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		},
	}

	pharmacyProduct := generatePharmacyProduct(301, 0)
	pharmacyProduct2 := generatePharmacyProduct(601, 300)
	pharmacyProduct3 := generatePharmacyProduct(908, 600)
	telemedicines := []*entity.Telemedicine{
		{ProfileId: 1, DoctorId: 25, OrderedAt: time.Now(), ExpiredAt: time.Now().Add(1 * time.Second), Status: entity.Waiting, TotalPayment: doctorProfile2[0].Fee},
		{ProfileId: 1, DoctorId: 25, OrderedAt: time.Now(), ExpiredAt: time.Now().Add(10 * time.Minute), Status: entity.Ongoing, TotalPayment: doctorProfile2[0].Fee, PrescriptionPdf: "https://everhealth-asset.irfancen.com/prescription/prescription-UJagIV", SickLeavePdf: "https://everhealth-asset.irfancen.com/sick-leave/sick-leave-AeCRNM", Proof: "https://everhealth-asset.irfancen.com/telemedicine-proof/telemedicine-proof-DBgvXY"},
		{ProfileId: 1, DoctorId: 26, OrderedAt: time.Now(), ExpiredAt: time.Now().Add(10 * time.Minute), Status: entity.End, TotalPayment: doctorProfile2[1].Fee, PrescriptionPdf: "https://everhealth-asset.irfancen.com/prescription/prescription-UJagIV", SickLeavePdf: "https://everhealth-asset.irfancen.com/sick-leave/sick-leave-AeCRNM", Proof: "https://everhealth-asset.irfancen.com/telemedicine-proof/telemedicine-proof-DBgvXY"},
		{ProfileId: 6, DoctorId: 26, OrderedAt: time.Now(), ExpiredAt: time.Now().Add(10 * time.Minute), Status: entity.Ongoing, TotalPayment: doctorProfile2[1].Fee, PrescriptionPdf: "https://everhealth-asset.irfancen.com/prescription/prescription-UJagIV", SickLeavePdf: "https://everhealth-asset.irfancen.com/sick-leave/sick-leave-AeCRNM", Proof: "https://everhealth-asset.irfancen.com/telemedicine-proof/telemedicine-proof-DBgvXY"},
	}
	chats := []*entity.Chat{
		{TelemedicineId: 2, UserId: 1, ChatTime: time.Date(2023, 12, 1, 10, 0, 0, 0, time.UTC), Message: "Hi", MessageType: entity.MessageTypeText},
		{TelemedicineId: 2, UserId: 2, ChatTime: time.Date(2023, 12, 1, 11, 0, 0, 0, time.UTC), Message: "Hello", MessageType: entity.MessageTypeText},
		{TelemedicineId: 2, UserId: 1, ChatTime: time.Date(2023, 12, 1, 12, 0, 0, 0, time.UTC), Message: "How are you", MessageType: entity.MessageTypeText},
		{TelemedicineId: 2, UserId: 2, ChatTime: time.Date(2023, 12, 1, 13, 0, 0, 0, time.UTC), Message: "I'm fine", MessageType: entity.MessageTypeText},
		{TelemedicineId: 3, UserId: 1, ChatTime: time.Date(2023, 12, 1, 12, 0, 0, 0, time.UTC), Message: "I'm sick yo", MessageType: entity.MessageTypeText},
		{TelemedicineId: 3, UserId: 8, ChatTime: time.Date(2023, 12, 1, 13, 0, 0, 0, time.UTC), Message: "What, I don't understand", MessageType: entity.MessageTypeText},
		{TelemedicineId: 3, UserId: 1, ChatTime: time.Date(2023, 12, 1, 14, 0, 0, 0, time.UTC), Message: "I've a headache", MessageType: entity.MessageTypeText},
		{TelemedicineId: 3, UserId: 8, ChatTime: time.Date(2023, 12, 1, 15, 0, 0, 0, time.UTC), Message: "So do I", MessageType: entity.MessageTypeText},
	}
	StockRecords := []*entity.StockRecord{
		{PharmacyProductId: 1, Quantity: 11, IsReduction: false, ChangeAt: time.Date(2024, time.January, 10, 0, 0, 0, 0, time.UTC)},
		{PharmacyProductId: 1, Quantity: 5, IsReduction: true, ChangeAt: time.Date(2024, time.January, 10, 0, 0, 0, 0, time.UTC)},
		{PharmacyProductId: 41, Quantity: 15, IsReduction: false, ChangeAt: time.Date(2024, time.January, 10, 0, 0, 0, 0, time.UTC)},
		{PharmacyProductId: 41, Quantity: 11, IsReduction: true, ChangeAt: time.Date(2024, time.January, 10, 0, 0, 0, 0, time.UTC)},
		{PharmacyProductId: 31, Quantity: 10, IsReduction: false, ChangeAt: time.Date(2024, time.January, 10, 0, 0, 0, 0, time.UTC)},
		{PharmacyProductId: 31, Quantity: 2, IsReduction: true, ChangeAt: time.Date(2024, time.January, 10, 0, 0, 0, 0, time.UTC)},
		{PharmacyProductId: 21, Quantity: 15, IsReduction: false, ChangeAt: time.Date(2024, time.January, 10, 0, 0, 0, 0, time.UTC)},
		{PharmacyProductId: 21, Quantity: 10, IsReduction: true, ChangeAt: time.Date(2024, time.January, 10, 0, 0, 0, 0, time.UTC)},
	}
	db.Create(roles)
	db.Create(doctorSpecialist)
	db.Create(users)
	db.Create(adminContact)
	db.Create(profiles)
	db.Create(doctorProfiles)
	db.Create(users2)
	db.Create(adminContact2)
	db.Create(users3)
	db.Create(profile3)
	db.Create(doctorProfile2)
	db.Create(doctorProfile3)
	db.Create(users4)
	db.Create(profile4)
	db.Create(doctorProfile4)
	db.Create(drugForms)
	db.Create(productCategories)
	db.Create(drugClassifications)
	db.Create(products)
	db.Create(drugs)
	db.Create(orderStatuses)
	db.Create(shippingMethods)
	db.Create(provinces)
	db.Create(cities)
	db.Create(addresses)
	db.Create(pharmacies)
	db.Create(carts)
	db.Create(orders)
	db.Create(pharmacyProduct)
	db.Create(pharmacyProduct2)
	db.Create(pharmacyProduct3)
	db.Create(orderItems)
	db.Create(telemedicines)
	db.Create(chats)
	db.Create(StockRecords)
}

func hashPassword(text string) string {
	h := hasher.NewHasher()
	hashedText, err := h.Hash(text)
	if err != nil {
		return ""
	}
	return string(hashedText)
}
