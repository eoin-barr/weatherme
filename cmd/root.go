/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/eoin-barr/weatherme/types"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

const (
	All     string = "all"
	Preview string = "preview"
)

func calculateDewPoint(temp float64, humidity float64) float64 {
	const a float64 = 17.62
	const b float64 = 243.12
	alpha := math.Log(humidity/100 + a*temp/(b+temp))
	return math.Round(((b*alpha)/(a-alpha))*100) / 100
}

func formatPreview(result types.WeatherRes, city string) string {
	temp := result.Main.Temp - 273.15

	return "\n🌆 City:\t" + city + "\n" +
		"🌤  Description:\t" + cases.Title(language.English, cases.Compact).String(result.Weather[0].Description) + "\n" +
		"🌡  Temperature:\t" + strconv.FormatFloat(temp, 'f', 2, 32) + " °C" + "\n" +
		"🌊 Pressure:\t" + strconv.FormatInt(int64(result.Main.Pressure), 10) + " hPa" + "\n" +
		"😰 Humidity:\t" + strconv.FormatInt(int64(result.Main.Humidity), 10) + " %" + "\n"
}

func formatAll(result types.WeatherRes, city string) string {
	temp := result.Main.Temp - 273.15
	tempFeelsLike := result.Main.Feels_like - 273.15
	tempMin := result.Main.Temp_min - 273.15
	tempMax := result.Main.Temp_max - 273.15

	dewPoint := calculateDewPoint(temp, float64(result.Main.Humidity))

	timeOffset := time.Duration(result.Timezone * int(time.Second))
	sunriseTime := time.UnixMilli(int64(result.Sys.Sunrise) * 1000).Add(timeOffset).UTC()
	sunsetTime := time.UnixMilli(int64(result.Sys.Sunset) * 1000).Add(timeOffset).UTC()

	return "\n🌆  City:\t\t" + city + "\n" +
		"🌍  Country:\t\t" + result.Sys.Country + "\n" +
		"⌚️  Timezone:\t\t" + strconv.Itoa(result.Timezone) + "\n" +
		"🗺   Latitude:\t\t" + strconv.FormatFloat(result.Coord.Lat, 'f', 2, 32) + "\n" +
		"🗺   Longitude:\t\t" + strconv.FormatFloat(result.Coord.Lon, 'f', 2, 32) + "\n\n" +

		"🌤   Description:\t" + cases.Title(language.English, cases.Compact).String(result.Weather[0].Description) + "\n" +
		"🌡   Temperature:\t" + strconv.FormatFloat(temp, 'f', 2, 32) + " °C" + "\n" +
		"💧  Dew point:\t\t" + strconv.FormatFloat(dewPoint, 'f', 2, 32) + " °C" + "\n" +
		"💁‍♀️  Temp Feels Like:\t" + strconv.FormatFloat(tempFeelsLike, 'f', 2, 32) + " °C" + "\n" +
		"🔥  Temperature Max:\t" + strconv.FormatFloat(tempMax, 'f', 2, 32) + " °C" + "\n" +
		"🧊  Temperature Min:\t" + strconv.FormatFloat(tempMin, 'f', 2, 32) + " °C" + "\n" +
		"🌊  Pressure:\t\t" + strconv.FormatInt(int64(result.Main.Pressure), 10) + " hPa" + "\n" +
		"😰  Humidity:\t\t" + strconv.FormatInt(int64(result.Main.Humidity), 10) + " %" + "\n\n" +

		"☁️   Cloudiness:\t\t" + strconv.FormatInt(int64(result.Clouds.All), 10) + " %" + "\n" +
		"🌬   Wind Speed:\t\t" + strconv.FormatFloat(result.Wind.Speed, 'f', 2, 32) + " m/s" + "\n" +
		"🧭  Wind Direction:\t" + strconv.Itoa(result.Wind.Deg) + " °" + "\n" +
		"🌁  Visibility:\t\t" + strconv.Itoa(result.Visibility) + "\n\n" +

		"🌅  Sunrise:\t\t" + sunriseTime.Format("15:04") + " (" + sunriseTime.Format(time.Kitchen) + ") UTC\n" +
		"🌇  Sunset:\t\t" + sunsetTime.Format("15:04") + " (" + sunsetTime.Format(time.Kitchen) + ") UTC\n"
}

func getWeather(args []string, view string) {

	city := strings.Join(args[0:], " ")
	secret := "db2aec76d1d51a968faea300e25e70cc"
	var u url.URL
	u.Scheme = "https"
	u.Host = "api.openweathermap.org"
	u.Path = "/geo/1.0/direct"

	q := u.Query()
	q.Set("q", strings.ToLower(city))
	q.Set("limit", "5")
	q.Set("appid", secret)
	u.RawQuery = q.Encode()

	cityDetailsResp, err := http.Get(u.String())
	if err != nil {
		fmt.Println("Hmmm, something went wrong 😢")
		return
	}

	cityDetailsBody, err := ioutil.ReadAll(cityDetailsResp.Body)
	if err != nil {
		fmt.Println("Hmmm, something went wrong 😢")
		return
	}

	var CityDetails types.CityDetails
	err = json.Unmarshal(cityDetailsBody, &CityDetails)
	if err != nil {
		fmt.Println("Hmmm, something went wrong 😢")
		return
	}

	if len(CityDetails) == 0 {
		return
	}

	var index = 0
	if len(CityDetails) > 1 {
		filtered := uniqueCities(CityDetails)
		if len(filtered) == 1 {
			CityDetails = filtered
		} else {
			index = renderList(filtered)
		}
	}

	lat := CityDetails[index].Lat
	lon := CityDetails[index].Lon

	var u2 url.URL
	u2.Scheme = "https"
	u2.Host = "api.openweathermap.org"
	u2.Path = "/data/2.5/weather"

	q2 := u2.Query()
	q2.Set("lat", fmt.Sprintf("%f", lat))
	q2.Set("lon", fmt.Sprintf("%f", lon))
	q2.Set("appid", secret)
	u2.RawQuery = q2.Encode()

	weatherResp, err := http.Get(u2.String())
	if err != nil {
		fmt.Println("Hmmm, something went wrong 😢")
		return
	}

	weatherBody, err := ioutil.ReadAll(weatherResp.Body)
	if err != nil {
		fmt.Println("Hmmm, something went wrong 😢")
		return
	}

	var WeatherRes types.WeatherRes
	err = json.Unmarshal(weatherBody, &WeatherRes)
	if err != nil {
		fmt.Println("Hmmm, something went wrong 😢")
		return
	}
	if view == All {
		fmt.Println(formatAll(WeatherRes, cases.Title(language.English, cases.Compact).String(city)))
	} else {
		fmt.Println(formatPreview(WeatherRes, cases.Title(language.English, cases.Compact).String(city)))
	}

}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "weatherme",
	Short: "A basic cli weather app",
	Long:  `Type in weatherme and the name of a city to find out the weather in that city.`,

	Run: func(cmd *cobra.Command, args []string) {
		flagVar, err := cmd.Flags().GetBool("all")
		if err != nil {
			fmt.Println(err)
		}
		if len(args) < 1 {
			fmt.Println("Please 'weatherme' keyword followed by a city")
			return
		}

		if flagVar {
			getWeather(args, All)
			return
		}
		getWeather(args, Preview)
	},
}

// Remove duplicates from city list
func uniqueCities(cities types.CityDetails) types.CityDetails {
	var filtered types.CityDetails
	for _, city := range cities {
		duplicate := false
		for _, c := range filtered {
			if city.Name == c.Name && city.Country == c.Country {
				duplicate = true
			}
		}
		if duplicate == false {
			filtered = append(filtered, city)
		}
	}
	return filtered
}

// Give the user a chance to select the correct city.
func renderList(cities types.CityDetails) int {
	if err := ui.Init(); err != nil {
    fmt.Println(err)
	}
	defer ui.Close()

	l := widgets.NewList()
	l.Title = "There are multiple cities with the same name."
  l.Border = false
	for i, c := range cities {
		l.Rows = append(l.Rows, fmt.Sprintf("%d: %s, %s", i, c.Name, c.Country))
	}

	l.WrapText = false
	l.SetRect(0, 0, 100, 100)
  l.SelectedRowStyle.Fg = ui.ColorGreen
	ui.Render(l)
	uiEvents := ui.PollEvents()
	var index int
  selected := false
	for {
		e := <-uiEvents

		switch e.ID {
		case "j", "<Down>":
			l.ScrollDown()
		case "k", "<Up>":
			l.ScrollUp()
		case "<C-d>":
			l.ScrollHalfPageDown()
		case "<C-u>":
			l.ScrollHalfPageUp()
		case "<C-f>":
			l.ScrollPageDown()
		case "<C-b>":
			l.ScrollPageUp()
		case "<Home>":
		case "<C-c>":
			os.Exit(0)
		case "<Enter>":
			index = l.SelectedRow
      selected = true
		}
  
    if(selected){
      return index 
    }
		ui.Render(l)
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("all", "a", false, "Get a granular view of the weather")
}
