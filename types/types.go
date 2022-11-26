package types

type ConfigStruct struct {
	Token             string `json:"token"`
	BotPrefix         string `json:"BotPrefix"`
	OpenWeatherAPIKey string `json:"OpenWeatherAPIKey"`
}

type WeatherRes struct {
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	}
	Weather []struct {
		Id          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	}
	Base string `json:"base"`
	Main struct {
		Temp       float64 `json:"temp"`
		Feels_like float64 `json:"feels_like"`
		Temp_min   float64 `json:"temp_min"`
		Temp_max   float64 `json:"temp_max"`
		Pressure   int     `json:"pressure"`
		Humidity   int     `json:"humidity"`
	}
	Visibility int `json:"visibility"`
	Wind       struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
	}
	Clouds struct {
		All int `json:"all"`
	}
	Dt  int `json:"dt"`
	Sys struct {
		Type    int    `json:"type"`
		Id      int    `json:"id"`
		Country string `json:"country"`
		Sunrise int    `json:"sunrise"`
		Sunset  int    `json:"sunset"`
	}
	Timezone int    `json:"timezone"`
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Cod      int    `json:"cod"`
}

type CityDetails []struct {
	Name       string
	Local_name struct{}
	Lat        float64
	Lon        float64
	Country    string
}

type WeatherView struct {
	Preview string
	All     string
}
