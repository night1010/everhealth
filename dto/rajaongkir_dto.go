package dto

type RajaOngkirResult struct {
	RajaOngkir struct {
		Results []RajaOngkirProvider `json:"results"`
	} `json:"rajaongkir"`
}

type RajaOngkirProvider struct {
	Code  string               `json:"code"`
	Name  string               `json:"name"`
	Costs []RajaOngkirServices `json:"costs"`
}

type RajaOngkirServices struct {
	Services    string           `json:"service"`
	Description string           `json:"description"`
	Costs       []RajaOngkirCost `json:"cost"`
}

type RajaOngkirCost struct {
	Value         int    `json:"value"`
	EstimatedTime string `json:"etd"`
}
