/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/eoin-barr/weatherme/types"

	// "github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	All     string = "all"
	Preview string = "preview"
)

func getSecret() string {
	// env1, err1 := os.LookupEnv("OPEN_WEATHER_API_SECRET")
	// log.Println(env1, err1)
	// if env1 == "" {
	// 	fmt.Println("OPEN_WEATHER_API_SECRET not set")
	// } else {
	// 	fmt.Println("OPEN_WEATHER_API_SECRET set: ", env1)
	// }

	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("KEY:", os.Getenv("OPEN_WEATHER_API_SECRET"))
	return ""
}

func formatPreview(result types.WeatherRes, city string) string {
	temp := result.Main.Temp - 273.15

	return "\n🌆 City:\t" + city + "\n" +
		"🌤  Description:\t" + cases.Title(language.English, cases.Compact).String(result.Weather[0].Description) + "\n" +
		"🌡  Temperature:\t" + strconv.FormatFloat(temp, 'f', 2, 32) + " °C" + "\n" +
		"🌊 Pressure:\t" + strconv.FormatInt(int64(result.Main.Pressure), 10) + " hPa" + "\n" +
		"😰 Humitdity:\t" + strconv.FormatInt(int64(result.Main.Humidity), 10) + " %" + "\n"
}

func formatAll(result types.WeatherRes, city string) string {
	temp := result.Main.Temp - 273.15
	tempFeelsLike := result.Main.Feels_like - 273.15
	tempMin := result.Main.Temp_min - 273.15
	tempMax := result.Main.Temp_max - 273.15

	return "\n🌆  City:\t\t" + city + "\n" +
		"🌍  Country:\t\t" + result.Sys.Country + "\n" +
		"⌚️  Timezone:\t\t" + strconv.Itoa(result.Timezone) + "\n" +
		"🗺   Latitude:\t\t" + strconv.FormatFloat(result.Coord.Lat, 'f', 2, 32) + "\n" +
		"🗺   Longitude:\t\t" + strconv.FormatFloat(result.Coord.Lon, 'f', 2, 32) + "\n\n" +

		"🌤   Description:\t" + cases.Title(language.English, cases.Compact).String(result.Weather[0].Description) + "\n" +
		"🌡   Temperature:\t" + strconv.FormatFloat(temp, 'f', 2, 32) + " °C" + "\n" +
		"💁‍♀️  Temp Feels Like:\t" + strconv.FormatFloat(tempFeelsLike, 'f', 2, 32) + " °C" + "\n" +
		"🔥  Temperature Max:\t" + strconv.FormatFloat(tempMax, 'f', 2, 32) + " °C" + "\n" +
		"🧊  Temperature Min:\t" + strconv.FormatFloat(tempMin, 'f', 2, 32) + " °C" + "\n" +
		"🌊  Pressure:\t\t" + strconv.FormatInt(int64(result.Main.Pressure), 10) + " hPa" + "\n" +
		"😰  Humitdity:\t\t" + strconv.FormatInt(int64(result.Main.Humidity), 10) + " %" + "\n\n" +

		"☁️   Cloudiness:\t\t" + strconv.FormatInt(int64(result.Clouds.All), 10) + " %" + "\n" +
		"🌬   Wind Speed:\t\t" + strconv.FormatFloat(result.Wind.Speed, 'f', 2, 32) + " m/s" + "\n" +
		"🧭  Wind Direction:\t" + strconv.Itoa(result.Wind.Deg) + " °" + "\n" +
		"🌁  Visibility:\t\t" + strconv.Itoa(result.Visibility) + "\n\n" +

		"🌅  Sunrise:\t\t" + strconv.Itoa(result.Sys.Sunrise) + "\n" +
		"🌇  Sunset:\t\t" + strconv.Itoa(result.Sys.Sunset) + "\n"
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
    index =  selectCity(CityDetails) 
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
		getSecret()
		getWeather(args, Preview)
	},
}
// Give the user a chance to select the correct city.
func selectCity(cities types.CityDetails ) int {
    fmt.Println("There are multiple cities with the same name.")
    fmt.Println("Please select the city. ")
    for i, c := range cities{
      fmt.Println(fmt.Sprintf("%d: %s, %s", i, c.Name, c.Country))
    }
    var number int
    var err error

    r := bufio.NewReader(os.Stdin)
    for {
      fmt.Fprint(os.Stderr, fmt.Sprintf("Choose a city between 0 - %d: ",len(cities)-1))
        input, _ := r.ReadString('\n')
        number, err = strconv.Atoi(input[:len(input)-1])
        if err != nil{
          fmt.Println(fmt.Sprintf("%s is not a valid number", input))
        }
        
        if number >= 0 && number < len(cities){
          break;
        }
    }
    return number
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
