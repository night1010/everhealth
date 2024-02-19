package migration

import (
	"encoding/json"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/night1010/everhealth/entity"
	"github.com/shopspring/decimal"
)

func generateProfile(id uint, name string) *entity.Profile {
	return &entity.Profile{
		UserId:    id,
		Name:      name,
		Image:     "https://everhealth-asset.irfancen.com/assets/avatar-placeholder.png",
		ImageKey:  "key",
		Birthdate: time.Date(2000, 01, 01, 0, 0, 0, 0, time.Local),
	}
}

func generateDoctorProfile(id uint, specialistId uint, fee int, status entity.StatusDoctor) *entity.DoctorProfile {
	return &entity.DoctorProfile{
		ProfileId:        id,
		Certificate:      "https://everhealth-asset.irfancen.com/doctor-certificate/sertif-dokter.pdf",
		CertificateKey:   "key",
		SpecialistId:     specialistId,
		YearOfExperience: uint(randomNumber(1, 5)),
		Fee:              decimal.NewFromInt(int64(fee)),
		Status:           status,
	}
}

func generateUser(email, password string, roleId entity.RoleId, isVerified bool) *entity.User {
	return &entity.User{Email: email, Password: hashPassword(password), RoleId: roleId, IsVerified: isVerified}
}

func generateFakerDoctor() ([]*entity.User, []*entity.Profile, []*entity.DoctorProfile) {
	var users []*entity.User
	var profiles []*entity.Profile
	var doctors []*entity.DoctorProfile
	counter := 73
	for i := 3; i < 11; i++ {
		for j := 1; j < 26; j++ {
			temp := generateUser("doctor"+strconv.Itoa(counter)+"@gmail.com", "Doctor12345", entity.RoleDoctor, true)
			users = append(users, temp)
			profileTemp := generateProfile(uint(counter), "doctor"+strconv.Itoa(counter))
			profiles = append(profiles, profileTemp)
			doctorTemp := generateDoctorProfile(uint(counter), uint(i), randomNumber500(20000, 10, 5000), entity.Online)
			doctors = append(doctors, doctorTemp)
			counter += 1
		}
	}
	return users, profiles, doctors
}

func generateRealDoctorProfile() []*entity.Profile {
	return []*entity.Profile{
		{UserId: 25, Name: "Ariaxina ramadhani", Image: "https://everhealth-asset.irfancen.com/seed-doctor-profile/ariacina.png"},
		{UserId: 26, Name: "Siti Rahayu", Image: "https://everhealth-asset.irfancen.com/seed-doctor-profile/siti-rahayu.png"},
		{UserId: 27, Name: "Budi Saturday", Image: "https://everhealth-asset.irfancen.com/seed-doctor-profile/budi-saturday.png"},
		{UserId: 28, Name: "Dewi Indah", Image: "https://everhealth-asset.irfancen.com/seed-doctor-profile/dewi-indah.png"},
		{UserId: 29, Name: "Adi Nugroho", Image: "https://everhealth-asset.irfancen.com/seed-doctor-profile/adi-nugroho.png"},
		{UserId: 30, Name: "Lina Setiawan", Image: "https://everhealth-asset.irfancen.com/seed-doctor-profile/lina-setiawan.png"},
		{UserId: 31, Name: "Wahyu Kurniawan", Image: "https://everhealth-asset.irfancen.com/seed-doctor-profile/wahyu-kurniawan.png"},
		{UserId: 32, Name: "Rini Permata", Image: "https://everhealth-asset.irfancen.com/seed-doctor-profile/rini-permata.png"},
		{UserId: 33, Name: "Arif Saputra", Image: "https://everhealth-asset.irfancen.com/seed-doctor-profile/arif-saputra.png"},
		{UserId: 34, Name: "Maya Purnama", Image: "https://everhealth-asset.irfancen.com/seed-doctor-profile/maya-purnama.png"},
		{UserId: 35, Name: "Dito Pratama", Image: "https://everhealth-asset.irfancen.com/seed-doctor-profile/dito-pratama.png"},
		{UserId: 36, Name: "Putri Anggraini", Image: "https://everhealth-asset.irfancen.com/seed-doctor-profile/putri-anggraini.png"},
		{UserId: 37, Name: "Fajar Wibowo", Image: "https://everhealth-asset.irfancen.com/seed-doctor-profile/fajar-wibowo.png"},
		{UserId: 38, Name: "Nia Astuti", Image: "https://everhealth-asset.irfancen.com/seed-doctor-profile/nia-astuti.png"},
		{UserId: 39, Name: "Agus Setiadi", Image: "https://everhealth-asset.irfancen.com/seed-doctor-profile/agus-setiadi.png"},
		{UserId: 40, Name: "Dina Rahmawati", Image: "https://everhealth-asset.irfancen.com/seed-doctor-profile/dina-rahmawati.png"},
		{UserId: 41, Name: "Joko Susanto", Image: "https://everhealth-asset.irfancen.com/seed-doctor-profile/joko-susanto.png"},
		{UserId: 42, Name: "Eka Fitriani", Image: "https://everhealth-asset.irfancen.com/seed-doctor-profile/eka-fitriani.png"},
		{UserId: 43, Name: "Rima Utami", Image: "https://everhealth-asset.irfancen.com/seed-doctor-profile/rima-utami.png"},
		{UserId: 44, Name: "Rudi Hartono", Image: "https://everhealth-asset.irfancen.com/seed-doctor-profile/rudi-hartono.png"},
		{UserId: 45, Name: "Siska Mariani", Image: "https://everhealth-asset.irfancen.com/seed-doctor-profile/siska-mariani.png"},
		{UserId: 46, Name: "Bayu Nugroho", Image: "https://everhealth-asset.irfancen.com/seed-doctor-profile/bayu-nugroho.png"},
		{UserId: 47, Name: "Lita Cahyani", Image: "https://everhealth-asset.irfancen.com/seed-doctor-profile/lita-cahyani.png"},
		{UserId: 48, Name: "Yasmin Wijaya", Image: "https://everhealth-asset.irfancen.com/seed-doctor-profile/yasmin-wijaya.png"},
		{UserId: 49, Name: "rudi Puspitasari", Image: "https://everhealth-asset.irfancen.com/seed-doctor-profile/ariacina.png"},
		{UserId: 50, Name: "siti Maulana", Image: "https://everhealth-asset.irfancen.com/seed-doctor-profile/siti-rahayu.png"},
		{UserId: 51, Name: "Wira Saputra", Image: "https://everhealth-asset.irfancen.com/seed-doctor-profile/budi-saturday.png"},
		{UserId: 52, Name: "Nana Fitri", Image: "https://everhealth-asset.irfancen.com/seed-doctor-profile/dewi-indah.png"},
		{UserId: 53, Name: "Dian Kusuma", Image: "https://everhealth-asset.irfancen.com/seed-doctor-profile/adi-nugroho.png"},
		{UserId: 54, Name: "Irma Utami", Image: "https://everhealth-asset.irfancen.com/seed-doctor-profile/lina-setiawan.png"},
		{UserId: 55, Name: "Hendra Santoso", Image: "https://everhealth-asset.irfancen.com/seed-doctor-profile/wahyu-kurniawan.png"},
		{UserId: 56, Name: "Dian Kusuma", Image: "https://everhealth-asset.irfancen.com/seed-doctor-profile/rini-permata.png"},
		{UserId: 57, Name: "Tioro Anggraeni", Image: "https://everhealth-asset.irfancen.com/seed-doctor-profile/arif-saputra.png"},
		{UserId: 58, Name: "Fita Novianti", Image: "https://everhealth-asset.irfancen.com/seed-doctor-profile/maya-purnama.png"},
		{UserId: 59, Name: "Doni Prasetyo", Image: "https://everhealth-asset.irfancen.com/seed-doctor-profile/dito-pratama.png"},
		{UserId: 60, Name: "Fita Novianti", Image: "https://everhealth-asset.irfancen.com/seed-doctor-profile/putri-anggraini.png"},
		{UserId: 61, Name: "Rizki Ramadhan", Image: "https://everhealth-asset.irfancen.com/seed-doctor-profile/fajar-wibowo.png"},
		{UserId: 62, Name: "Mega Puspita", Image: "https://everhealth-asset.irfancen.com/seed-doctor-profile/nia-astuti.png"},
		{UserId: 63, Name: "Yudi Hermawan", Image: "https://everhealth-asset.irfancen.com/seed-doctor-profile/agus-setiadi.png"},
		{UserId: 64, Name: "Evi Nurhadi", Image: "https://everhealth-asset.irfancen.com/seed-doctor-profile/dina-rahmawati.png"},
		{UserId: 65, Name: "Agung Setiawan", Image: "https://everhealth-asset.irfancen.com/seed-doctor-profile/joko-susanto.png"},
		{UserId: 66, Name: "Hadi Susanto", Image: "https://everhealth-asset.irfancen.com/seed-doctor-profile/eka-fitriani.png"},
		{UserId: 67, Name: "Sinta Rahmawati", Image: "https://everhealth-asset.irfancen.com/seed-doctor-profile/rima-utami.png"},
		{UserId: 68, Name: "Arya Wijaya", Image: "https://everhealth-asset.irfancen.com/seed-doctor-profile/rudi-hartono.png"},
		{UserId: 69, Name: "Ratna Permata", Image: "https://everhealth-asset.irfancen.com/seed-doctor-profile/siska-mariani.png"},
		{UserId: 70, Name: "Dian Setiawan", Image: "https://everhealth-asset.irfancen.com/seed-doctor-profile/bayu-nugroho.png"},
		{UserId: 71, Name: "Nana Fauzi", Image: "https://everhealth-asset.irfancen.com/seed-doctor-profile/lita-cahyani.png"},
		{UserId: 72, Name: "Adel Pratiwi", Image: "https://everhealth-asset.irfancen.com/seed-doctor-profile/yasmin-wijaya.png"},
	}
}

func generateAdminContact(userId uint, name, phone string) *entity.AdminContact {
	return &entity.AdminContact{UserId: userId, Name: name, Phone: phone}
}

func generateAdminPharmacy(start, max int, baseName, password string) ([]*entity.User, []*entity.AdminContact) {
	var users []*entity.User
	var adminContacts []*entity.AdminContact
	for i := start; i+1 < max+2; i++ {
		iString := strconv.Itoa(i - start + 1)
		temp := generateUser(baseName+iString+"@gmail.com", password, entity.RoleAdmin, true)
		users = append(users, temp)
		tempAdmin := generateAdminContact(uint(i+1), baseName+iString, "089654749370")
		adminContacts = append(adminContacts, tempAdmin)
	}
	return users, adminContacts
}

func generateDecimalFromString(decimalString string) decimal.Decimal {
	zero := decimal.NewFromInt(0)
	decimal, err := decimal.NewFromString(decimalString)
	if err != nil {
		return zero
	}
	return decimal
}

func generatePharmacyProduct(max, min int) []*entity.PharmacyProduct {
	var pps []*entity.PharmacyProduct
	var productSlice []*entity.Product
	data, _ := os.ReadFile("./migration/data/product-details.json")

	_ = json.Unmarshal(data, &productSlice)
	for i := 1; i <= 18; i++ {
		for _, product := range productSlice[min:max] {
			pps = append(pps, &entity.PharmacyProduct{
				ProductId:  product.Id,
				PharmacyId: uint(i),
				Price:      decimal.NewFromInt(int64(randomNumber(int(product.Price.IntPart())-2000, int(product.Price.IntPart())+2000))).Abs(),
				Stock:      randomNumber(10, 30),
				IsActive:   true,
			})
		}
	}
	return pps
}

func randomNumber(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

func randomNumber500(min, multiplier, step int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(10)*step + min
}

func generateUniqueRandomArray(length, min, max int) []int {
	rand.Seed(time.Now().UnixNano())
	uniqueNumbers := make(map[int]bool)
	var result []int
	for len(uniqueNumbers) < length {
		randomNumber := rand.Intn(max-min+1) + min
		uniqueNumbers[randomNumber] = true
	}
	for number := range uniqueNumbers {
		result = append(result, number)
	}
	return result
}

func RandomBool() bool {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(2) == 0
}
