package migration

import "github.com/night1010/everhealth/entity"

func getProvinces() []*entity.Province {
	return []*entity.Province{
		{
			ProvinceGid: 2,
			Code:        "1",
			Name:        "Bali",
		},
		{
			ProvinceGid: 3,
			Code:        "2",
			Name:        "Bangka Belitung",
		},
		{
			ProvinceGid: 4,
			Code:        "3",
			Name:        "Banten",
		},
		{
			ProvinceGid: 5,
			Code:        "4",
			Name:        "Bengkulu",
		},
		{
			ProvinceGid: 34,
			Code:        "5",
			Name:        "DI Yogyakarta",
		},
		{
			ProvinceGid: 8,
			Code:        "6",
			Name:        "DKI Jakarta",
		},
		{
			ProvinceGid: 6,
			Code:        "7",
			Name:        "Gorontalo",
		},
		{
			ProvinceGid: 9,
			Code:        "8",
			Name:        "Jambi",
		},
		{
			ProvinceGid: 10,
			Code:        "9",
			Name:        "Jawa Barat",
		},
		{
			ProvinceGid: 11,
			Code:        "10",
			Name:        "Jawa Tengah",
		},
		{
			ProvinceGid: 12,
			Code:        "11",
			Name:        "Jawa Timur",
		},
		{
			ProvinceGid: 13,
			Code:        "12",
			Name:        "Kalimantan Barat",
		},
		{
			ProvinceGid: 14,
			Code:        "13",
			Name:        "Kalimantan Selatan",
		},
		{
			ProvinceGid: 15,
			Code:        "14",
			Name:        "Kalimantan Tengah",
		},
		{
			ProvinceGid: 16,
			Code:        "15",
			Name:        "Kalimantan Timur",
		},
		{
			ProvinceGid: 17,
			Code:        "16",
			Name:        "Kalimantan Utara",
		},
		{
			ProvinceGid: 18,
			Code:        "17",
			Name:        "Kepulauan Riau",
		},
		{
			ProvinceGid: 19,
			Code:        "18",
			Name:        "Lampung",
		},
		{
			ProvinceGid: 21,
			Code:        "19",
			Name:        "Maluku",
		},
		{
			ProvinceGid: 20,
			Code:        "20",
			Name:        "Maluku Utara",
		},
		{
			ProvinceGid: 1,
			Code:        "21",
			Name:        "Nanggroe Aceh Darussalam (NAD)",
		},
		{
			ProvinceGid: 22,
			Code:        "22",
			Name:        "Nusa Tenggara Barat (NTB)",
		},
		{
			ProvinceGid: 23,
			Code:        "23",
			Name:        "Nusa Tenggara Timur (NTT)",
		},
		{
			ProvinceGid: 24,
			Code:        "24",
			Name:        "Papua",
		},
		{
			ProvinceGid: 7,
			Code:        "25",
			Name:        "Papua Barat",
		},
		{
			ProvinceGid: 25,
			Code:        "26",
			Name:        "Riau",
		},
		{
			ProvinceGid: 26,
			Code:        "27",
			Name:        "Sulawesi Barat",
		},
		{
			ProvinceGid: 27,
			Code:        "28",
			Name:        "Sulawesi Selatan",
		},
		{
			ProvinceGid: 28,
			Code:        "29",
			Name:        "Sulawesi Tengah",
		},
		{
			ProvinceGid: 29,
			Code:        "30",
			Name:        "Sulawesi Tenggara",
		},
		{
			ProvinceGid: 30,
			Code:        "31",
			Name:        "Sulawesi Utara",
		},
		{
			ProvinceGid: 31,
			Code:        "32",
			Name:        "Sumatera Barat",
		},
		{
			ProvinceGid: 32,
			Code:        "33",
			Name:        "Sumatera Selatan",
		},
		{
			ProvinceGid: 33,
			Code:        "34",
			Name:        "Sumatera Utara",
		},
	}
}

func getCities() []*entity.City {
	return ImportCities()
}
