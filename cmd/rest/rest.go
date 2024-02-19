package main

import (
	"github.com/night1010/everhealth/appjwt"
	"github.com/night1010/everhealth/appvalidator"
	"github.com/night1010/everhealth/chat"
	"github.com/night1010/everhealth/handler"
	"github.com/night1010/everhealth/hasher"
	"github.com/night1010/everhealth/imagehelper"
	"github.com/night1010/everhealth/logger"
	"github.com/night1010/everhealth/mail"
	"github.com/night1010/everhealth/repository"
	"github.com/night1010/everhealth/router"
	"github.com/night1010/everhealth/server"
	"github.com/night1010/everhealth/transactor"
	"github.com/night1010/everhealth/usecase"
	"github.com/go-resty/resty/v2"
)

func main() {
	logger.SetLogrusLogger()

	db, err := repository.GetConnection()
	if err != nil {
		logger.Log.Error(err)
	}

	client := resty.New()
	client.SetHeader("Content-Type", "application/json")

	manager := transactor.NewManager(db)
	hash := hasher.NewHasher()
	mail := mail.NewSmtpGmail()
	imageHelper, err := imagehelper.NewImageHelper(imagehelper.GoogleMethod)
	if err != nil {
		logger.Log.Error(err)
	}
	jwt := appjwt.NewJwt()
	appvalidator.RegisterCustomValidator()
	ur := repository.NewUserRepository(db)
	pr := repository.NewProfileRepository(db)
	dpr := repository.NewDoctorProfileRepository(db)
	productCategoryRepository := repository.NewProductCategoryRepository(db)
	fr := repository.NewForgotPasswordRepository(db)
	productRepo := repository.NewProductRepository(db)
	cartRepo := repository.NewCartRepository(db)
	pharmacyProductRepository := repository.NewPharmacyProductRepository(db)

	drugRepo := repository.NewDrugRepository(db)
	drugFormRepo := repository.NewDrugFormRepository(db)
	drugClassificationRepo := repository.NewDrugClassificationRepository(db)

	au := usecase.NewAuthUsecase(manager, ur, pr, dpr, fr, cartRepo, mail, hash, jwt, imageHelper)
	uu := usecase.NewUserUsecase(manager, ur, pr, dpr, hash, imageHelper)
	productCategoryUsecase := usecase.NewProductCategoryUsecase(productCategoryRepository, imageHelper, manager)
	productUsecase := usecase.NewProductUsecase(manager, imageHelper, productRepo, productCategoryRepository, drugRepo, drugFormRepo, drugClassificationRepo, pharmacyProductRepository)

	ah := handler.NewAuthHandler(au)
	productCategoryHandler := handler.NewProductCategoryHandler(productCategoryUsecase)
	uh := handler.NewUserHAndler(uu)
	productHandler := handler.NewProductHandler(productUsecase)

	provinceRepository := repository.NewProvinceRepository(db)
	provinceUsecase := usecase.NewProvinceUsecase(provinceRepository)
	provinceHandler := handler.NewProvinceHadnler(provinceUsecase)

	cityRepository := repository.NewCityRepository(db)

	orderStatusRepository := repository.NewOrderStatusRepository(db)
	orderStatusUsecase := usecase.NewOrderStatusUsecase(orderStatusRepository)
	orderStatusHandler := handler.NewOrderStatusHadnler(orderStatusUsecase)

	drugClassificationRepository := repository.NewDrugClassificationRepository(db)
	drugClassificationUsecase := usecase.NewDrugClassificationUsecase(drugClassificationRepository)
	drugClassificationHandler := handler.NewDrugClassificationHandler(drugClassificationUsecase)

	drugFormRepository := repository.NewDrugFormRepository(db)
	drugFormUsecase := usecase.NewDrugFormUsecase(drugFormRepository)
	drugFormHandler := handler.NewDrugFormHandler(drugFormUsecase)

	doctorSpecialistRepository := repository.NewDoctorSpecialistRepository(db)
	doctorSpecialistUsecase := usecase.NewDoctorSpecialistUsecase(doctorSpecialistRepository)
	doctorSpecialistHandler := handler.NewDoctorSpecialistHandler(doctorSpecialistUsecase)

	shippingMethodRepository := repository.NewShippingMethodRepository(db, client)
	addressRepository := repository.NewAddressRepository(db)
	addressUsecase := usecase.NewAddressUsecase(addressRepository, manager, shippingMethodRepository)
	addressHandler := handler.NewAddressHandler(addressUsecase)

	pharmacyRepository := repository.NewPharmacyRepository(db)
	pharmacyUsecase := usecase.NewPharmacyUsecase(pharmacyRepository, provinceRepository, cityRepository)
	pharmacyHandler := handler.NewPharmacyHandler(pharmacyUsecase)

	cartItemRepo := repository.NewCartItemRepository(db)
	cartUsecase := usecase.NewCartUsecase(manager, cartRepo, cartItemRepo, productRepo, pharmacyProductRepository)
	cartHandler := handler.NewCartHandler(cartUsecase)

	pharmacyProductUsecase := usecase.NewPharmacyProductUsecase(pharmacyProductRepository, pharmacyRepository, productRepo)
	pharmacyProductHandler := handler.NewPharmacyProductHandler(pharmacyProductUsecase)

	adminContactRepository := repository.NewAdminContactRepository(db)
	adminPharmacyUsecase := usecase.NewAdminPharmacyUsecase(ur, hash, pharmacyRepository, adminContactRepository, manager)
	adminPharmacyHandler := handler.NewAdminPharmacyHandler(adminPharmacyUsecase)

	stockRecordRepository := repository.NewStockRecordRepository(db)
	stockRecordUsecase := usecase.NewStockRecordUsecase(stockRecordRepository, pharmacyProductRepository, manager)
	stockRecordHandler := handler.NewStockRecordHandler(stockRecordUsecase)

	stockMutationRepository := repository.NewStockMutationRepository(db)
	stockMutationUsecase := usecase.NewStockMutationUsecase(stockMutationRepository, pharmacyProductRepository, stockRecordRepository, manager)
	stockMutationHandler := handler.NewStockMutationHandler(stockMutationUsecase)

	orderItemRepository := repository.NewOrderItemRepository(db)
	productOrderRepository := repository.NewProductOrderRepository(db)
	orderUsecase := usecase.NewOrderUsecase(manager, imageHelper, cartRepo, orderItemRepository, productOrderRepository, cartItemRepo, addressRepository, pharmacyRepository, pharmacyProductRepository, stockMutationRepository, stockRecordRepository)
	orderHandler := handler.NewOrderHandler(orderUsecase)

	shippingMethodRepo := repository.NewShippingMethodRepository(db, client)
	shippingMethodUsecase := usecase.NewShippingMethodUsecase(addressRepository, shippingMethodRepo, pharmacyRepository, orderUsecase)
	shippingMethodHandler := handler.NewShippingMethodHandler(shippingMethodUsecase)

	telemedicineRepo := repository.NewTelemedicineRepository(db)
	telemedicineUsecase := usecase.NewTelemedicineUsecase(telemedicineRepo, dpr, manager, imageHelper)
	telemedicineHandler := handler.NewTelemedicineHadnler(telemedicineUsecase)

	chatRepo := repository.NewChatRepository(db)
	chatUsecase := usecase.NewChatUsecase(chatRepo, telemedicineRepo)
	chatManager := chat.NewManager(chatUsecase)
	chatHandler := handler.NewChatHandler(chatManager, chatUsecase, jwt)

	handlers := router.Handlers{
		Auth:               ah,
		ProductCategory:    productCategoryHandler,
		User:               uh,
		Product:            productHandler,
		Province:           provinceHandler,
		DrugClassification: drugClassificationHandler,
		DrugForm:           drugFormHandler,
		OrderStatus:        orderStatusHandler,
		DoctorSpecialist:   doctorSpecialistHandler,
		Address:            addressHandler,
		Pharmacy:           pharmacyHandler,
		Cart:               cartHandler,
		PharmacyProduct:    pharmacyProductHandler,
		AdminPharmacy:      adminPharmacyHandler,
		StockRecord:        stockRecordHandler,
		Order:              orderHandler,
		StockMutation:      stockMutationHandler,
		ShippingMethod:     shippingMethodHandler,
		Chat:               chatHandler,
		Telemedicine:       telemedicineHandler,
	}

	r := router.New(handlers)

	s := server.New(r)

	server.StartWithGracefulShutdown(s)
}
