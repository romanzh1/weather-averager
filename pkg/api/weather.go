package api

type Weather []struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type WeatherNow struct {
	Main struct {
		Dt        int64   `json:"dt"`
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Pressure  int     `json:"pressure"`
		Humidity  int     `json:"humidity"`
	} `json:"main"`
	Weather

	Visibility int `json:"visibility"`
	Wind       struct {
		Speed float64 `json:"speed"`
	} `json:"wind"`
	Cod int `json:"cod"`
}

type WeatherTwoDays struct {
	Hourly []struct {
		Dt         int64   `json:"dt"`
		Temp       float64 `json:"temp"`
		FeelsLike  float64 `json:"feels_like"`
		Pressure   int     `json:"pressure"`
		Humidity   int     `json:"humidity"`
		Visibility int     `json:"visibility"`
		WindSpeed  float64 `json:"wind_speed"`
		Pop        float32 `json:"pop"`
		Weather
	} `json:"hourly"`

	Daily []struct {
		Dt        int64   `json:"dt"`
		Sunrise   int64   `json:"sunrise"`
		Sunset    int64   `json:"sunset"`
		Moonrise  int64   `json:"moonrise"`
		Moonset   int64   `json:"moonset"`
		MoonPhase float64 `json:"moon_phase"`
		Temp      struct {
			Day   float64 `json:"day"`
			Min   float64 `json:"min"`
			Max   float64 `json:"max"`
			Night float64 `json:"night"`
			Eve   float64 `json:"eve"`
			Morn  float64 `json:"morn"`
		} `json:"temp"`
		FeelsLike struct {
			Day   float64 `json:"day"`
			Night float64 `json:"night"`
			Eve   float64 `json:"eve"`
			Morn  float64 `json:"morn"`
		} `json:"feels_like"`
		Pressure  int     `json:"pressure"`
		Humidity  int     `json:"humidity"`
		WindSpeed float64 `json:"wind_speed"`
		Weather
		Clouds int     `json:"clouds"`
		Pop    float32 `json:"pop"`
		Uvi    float64 `json:"uvi"`
		Rain   float64 `json:"rain,omitempty"`
	} `json:"daily"`
}
