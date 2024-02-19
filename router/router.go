package router

import (
	"errors"
	"net/http"

	"github.com/night1010/everhealth/apperror"
	"github.com/night1010/everhealth/entity"
	"github.com/night1010/everhealth/handler"
	"github.com/night1010/everhealth/middleware"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Auth               *handler.AuthHandler
	ProductCategory    *handler.ProductCategoryHandler
	User               *handler.UserHandler
	Product            *handler.ProductHandler
	Province           *handler.ProvinceHandler
	DrugClassification *handler.DrugClassificationHandler
	DrugForm           *handler.DrugFormHandler
	OrderStatus        *handler.OrderStatusHandler
	DoctorSpecialist   *handler.DoctorSpecialistHandler
	Address            *handler.AddressHandler
	ShippingMethod     *handler.ShippingMethodHandler
	Cart               *handler.CartHandler
	Pharmacy           *handler.PharmacyHandler
	PharmacyProduct    *handler.PharmacyProductHandler
	AdminPharmacy      *handler.AdminPharmacyHandler
	StockRecord        *handler.StockRecordHandler
	Order              *handler.OrderHandler
	StockMutation      *handler.StockMutationHandler
	Chat               *handler.ChatHandler
	Telemedicine       *handler.TelemedicineHandler
}

func New(handlers Handlers) http.Handler {
	router := gin.New()

	router.NoRoute(routeNotFoundHandler)

	router.Use(gin.Recovery())
	router.Use(middleware.Cors())
	router.Use(middleware.Timeout())
	router.Use(middleware.Logger())
	router.Use(middleware.Error())

	auth := router.Group("/auth")
	auth.POST("/register", handlers.Auth.Register)
	auth.POST("/verify/:id", middleware.PDFUpload(), handlers.Auth.Verify)
	auth.POST("/login", handlers.Auth.Login)
	auth.POST("/forgot-password", handlers.Auth.RequestForgotPassword)
	auth.PUT("/forgot-password/apply", handlers.Auth.ApplyPassword)

	user := router.Group("/users")
	user.GET("", middleware.Auth(entity.RoleSuperAdmin), handlers.User.GetAllUser)
	user.GET("/profile", middleware.Auth(entity.RoleUser, entity.RoleDoctor), handlers.User.GetProfile)
	user.POST("/reset-password", middleware.Auth(entity.RoleUser, entity.RoleDoctor), handlers.User.ResetPassword)
	user.PUT("/profile", middleware.Auth(entity.RoleUser, entity.RoleDoctor), middleware.ImageUploadMiddleware(), middleware.PDFUpload(), handlers.User.UpdateProfile)
	user.PUT("/status", middleware.Auth(entity.RoleDoctor), handlers.User.UpdateStatus)

	cart := router.Group("/cart")
	cart.GET("", middleware.Auth(entity.RoleUser), handlers.Cart.GetCart)
	cart.POST("", middleware.Auth(entity.RoleUser), handlers.Cart.AddItem)
	cart.PATCH("/:id", middleware.Auth(entity.RoleUser), handlers.Cart.ChangeQty)
	cart.DELETE("/:id", middleware.Auth(entity.RoleUser), handlers.Cart.DeleteItem)
	cart.PATCH("/check/:id", middleware.Auth(entity.RoleUser), handlers.Cart.CheckItem)
	cart.PUT("/check-all", middleware.Auth(entity.RoleUser), handlers.Cart.CheckAllItem)

	router.GET("/doctors/:id", handlers.User.DoctorDetail)
	router.GET("/doctors", handlers.User.ListDoctor)

	productCategory := router.Group("/product-categories")
	productCategory.GET("", handlers.ProductCategory.GetProductCategories)
	productCategory.GET("/:id", handlers.ProductCategory.GetProductCategoriesDetail)
	productCategory.POST("", middleware.Auth(entity.RoleSuperAdmin), middleware.ImageUploadMiddleware(), handlers.ProductCategory.PostProductCategory)
	productCategory.PUT("/:id", middleware.Auth(entity.RoleSuperAdmin), middleware.ImageUploadMiddleware(), handlers.ProductCategory.PutProductCategory)
	productCategory.DELETE("/:id", middleware.Auth(entity.RoleSuperAdmin), handlers.ProductCategory.DeleteProductCategory)

	products := router.Group("/products")
	products.GET("", handlers.Product.ListProduct)
	products.GET("/admin", middleware.Auth(entity.RoleAdmin, entity.RoleSuperAdmin), handlers.Product.ListProductAdmin)
	products.GET("/:id/admin", middleware.Auth(entity.RoleAdmin, entity.RoleSuperAdmin), handlers.Product.GetProductDetailAdmin)
	products.GET("/nearby", middleware.Auth(entity.RoleUser), handlers.Product.ListNearbyProduct)
	products.POST("", middleware.Auth(entity.RoleSuperAdmin), middleware.ImageUploadMiddleware(), handlers.Product.AddProduct)
	products.PUT("/:id", middleware.Auth(entity.RoleSuperAdmin), middleware.ImageUploadMiddleware(), handlers.Product.UpdateProduct)
	products.GET("/:id", handlers.Product.GetProductDetail)

	province := router.Group("/provinces")
	province.GET("", handlers.Province.GetAllProvince)
	province.GET("/:id", handlers.Province.GetDetailProvince)

	drugClassification := router.Group("/drug-classifications")
	drugClassification.GET("", handlers.DrugClassification.GetAllDrugClassification)

	drugForm := router.Group("/drug-forms")
	drugForm.GET("", handlers.DrugForm.GetAllDrugForm)

	orderStatus := router.Group("/order-statuses")
	orderStatus.GET("", handlers.OrderStatus.GetAllOrderStatus)

	doctorSpecialist := router.Group("/doctor-specialists")
	doctorSpecialist.GET("", handlers.DoctorSpecialist.GetAllDoctorSpecialist)

	address := router.Group("/addresses")
	address.GET("", middleware.Auth(entity.RoleUser), handlers.Address.GetAddress)
	address.POST("", middleware.Auth(entity.RoleUser), handlers.Address.CreateAddress)
	address.PUT("/:id", middleware.Auth(entity.RoleUser), handlers.Address.UpdateAddress)
	address.DELETE("/:id", middleware.Auth(entity.RoleUser), handlers.Address.DeleteAddress)
	address.PATCH("/:id/default", middleware.Auth(entity.RoleUser), handlers.Address.ChangeDefaultAddress)
	address.POST("/validate", handlers.Address.ValidateAddress)

	pharmacy := router.Group("/pharmacies")
	pharmacy.GET("", middleware.Auth(entity.RoleAdmin), handlers.Pharmacy.GetAllPharmacy)
	pharmacy.GET("/super-admin", middleware.Auth(entity.RoleSuperAdmin, entity.RoleAdmin), handlers.Pharmacy.GetAllPharmacySuperAdmin)
	pharmacy.GET("/:pharmacy_id", middleware.Auth(entity.RoleAdmin, entity.RoleSuperAdmin), handlers.Pharmacy.GetPharmacyDetail)
	pharmacy.POST("", middleware.Auth(entity.RoleAdmin), handlers.Pharmacy.PostPharmacy)
	pharmacy.PUT("/:pharmacy_id", middleware.Auth(entity.RoleAdmin), handlers.Pharmacy.PutPharmacy)

	pharmacyProduct := pharmacy.Group("/:pharmacy_id/products")
	pharmacyProduct.GET("", middleware.Auth(entity.RoleAdmin), handlers.PharmacyProduct.GetAllPharmacy)
	pharmacyProduct.POST("", middleware.Auth(entity.RoleAdmin), handlers.PharmacyProduct.PostPharmacyProduct)
	pharmacyProduct.GET("/:product_id", middleware.Auth(entity.RoleAdmin), handlers.PharmacyProduct.GetPharmacyProductDetail)
	pharmacyProduct.PUT("/:product_id", middleware.Auth(entity.RoleAdmin), handlers.PharmacyProduct.PutPharmacyProduct)

	adminPharmacy := router.Group("/admins-pharmacy")
	adminPharmacy.GET("", middleware.Auth(entity.RoleSuperAdmin), handlers.AdminPharmacy.GetAllAdminPharmacy)
	adminPharmacy.GET("/:id", middleware.Auth(entity.RoleSuperAdmin), handlers.AdminPharmacy.GetDetailAdminPharmacy)
	adminPharmacy.POST("", middleware.Auth(entity.RoleSuperAdmin), handlers.AdminPharmacy.PostAdminPharmacy)
	adminPharmacy.DELETE("/:id", middleware.Auth(entity.RoleSuperAdmin), handlers.AdminPharmacy.DeleteAdminPharmacy)
	adminPharmacy.PUT("/:id", middleware.Auth(entity.RoleSuperAdmin), handlers.AdminPharmacy.UpdateAdminPharmacy)

	stockRecord := router.Group("/stock-records")
	stockRecord.GET("", middleware.Auth(entity.RoleAdmin), handlers.StockRecord.GetAllStockRecord)
	stockRecord.GET("/monthly", middleware.Auth(entity.RoleAdmin), handlers.StockRecord.GetStockMonthlyReport)
	stockRecord.POST("", middleware.Auth(entity.RoleAdmin), handlers.StockRecord.PostStockRecord)

	stockMutation := router.Group("/stock-mutations")
	stockMutation.GET("", middleware.Auth(entity.RoleAdmin), handlers.StockMutation.GetAllStockMutation)
	stockMutation.GET("/:id", middleware.Auth(entity.RoleAdmin), handlers.StockMutation.GetStockMutationDetail)
	stockMutation.POST("", middleware.Auth(entity.RoleAdmin), handlers.StockMutation.PostStockMutation)
	stockMutation.POST("/:id/change-status", middleware.Auth(entity.RoleAdmin), handlers.StockMutation.ChangeStatusStockMutation)
	stockMutation.GET("/pharmacy", middleware.Auth(entity.RoleAdmin), handlers.StockMutation.GetAllAvailablePharmacyStockMutation)
	router.GET("/shipping-method/:id", middleware.Auth(entity.RoleUser), handlers.ShippingMethod.GetShippingMethod)

	order := router.Group("/order")
	order.POST("", middleware.Auth(entity.RoleUser), handlers.Order.CreateOrder)
	order.GET("", middleware.Auth(entity.RoleUser, entity.RoleAdmin, entity.RoleSuperAdmin), handlers.Order.OrderHistory)
	order.GET("/:id", middleware.Auth(entity.RoleUser, entity.RoleAdmin, entity.RoleSuperAdmin), handlers.Order.OrderDetail)
	order.GET("/items/:id", middleware.Auth(entity.RoleUser), handlers.Order.GetAvailableProduct)
	order.POST("/:id/upload-proof", middleware.Auth(entity.RoleUser), middleware.ImageUploadMiddleware(), handlers.Order.UploadPaymentProof)
	order.GET("/monthly-sales", middleware.Auth(entity.RoleSuperAdmin, entity.RoleAdmin), handlers.Order.GetMonthlySaleReport)
	order.PATCH("/:id/status-admin", middleware.Auth(entity.RoleAdmin), handlers.Order.AdminUpdateOrderStatus)
	order.PATCH("/:id/status-user", middleware.Auth(entity.RoleUser), handlers.Order.UserUpdateOrderStatus)

	chat := router.Group("/chat")
	chat.GET("/:id", handlers.Chat.Handle)

	telemedicine := router.Group("/telemedicines")
	telemedicine.POST("", middleware.Auth(entity.RoleUser), handlers.Telemedicine.PostTelemedicine)
	telemedicine.POST("/:id/payment", middleware.Auth(entity.RoleUser), middleware.ImageUploadMiddleware(), handlers.Telemedicine.PostPaymentTelemedicine)
	telemedicine.POST("/:id/sick-leave", middleware.Auth(entity.RoleDoctor), handlers.Telemedicine.PostSickLeave)
	telemedicine.POST("/:id/prescription", middleware.Auth(entity.RoleDoctor), handlers.Telemedicine.PostPrescription)
	telemedicine.GET("", middleware.Auth(entity.RoleUser, entity.RoleDoctor), handlers.Telemedicine.GetAllTelemedicine)
	telemedicine.GET("/:id", middleware.Auth(entity.RoleUser, entity.RoleDoctor), handlers.Telemedicine.GetTelemedicineDetail)
	telemedicine.PUT("/:id/end", middleware.Auth(entity.RoleUser, entity.RoleDoctor), handlers.Telemedicine.EndChat)
	return router
}

func routeNotFoundHandler(c *gin.Context) {
	var errRouteNotFound = errors.New("route not found")
	_ = c.Error(apperror.NewClientError(errRouteNotFound).NotFound())
}
