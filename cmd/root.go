/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/eoin-barr/weatherme/types"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func getSecret() string {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	return os.Getenv("OPEN_WEATHER_SECRET")
}

func formatString(result types.WeatherRes, city string) string {
	temp := result.Main.Temp - 273.15

	return "\n🌆 City:\t" + city + "\n" +
		"🌤  Description:\t" + cases.Title(language.English, cases.Compact).String(result.Weather[0].Description) + "\n" +
		"🌡  Temperature:\t" + strconv.FormatFloat(temp, 'f', 2, 32) + " °C" + "\n" +
		"🌊 Pressure:\t" + strconv.FormatInt(int64(result.Main.Pressure), 10) + " hPa" + "\n" +
		"😰 Humitdity:\t" + strconv.FormatInt(int64(result.Main.Humidity), 10) + " %" + "\n"
}

func getWeather(args []string) {
	city := strings.Join(args[0:], " ")
	secret := getSecret()
	var u url.URL
	u.Scheme = "https"
	u.Host = "api.openweathermap.org"
	u.Path = "/geo/1.0/direct"

	q := u.Query()
	q.Set("q", strings.ToLower(city))
	q.Set("limit", "1")
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

	lat := CityDetails[0].Lat
	lon := CityDetails[0].Lon

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

	fmt.Println(formatString(WeatherRes, cases.Title(language.English, cases.Compact).String(city)))
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "weatherme",
	Short: "A basic cli weather app",
	Long:  `Type in weatherme and the name of a city to find out the weather in that city.`,

	Run: func(cmd *cobra.Command, args []string) {
		flagVar, err := cmd.Flags().GetBool("differentmessage")
		if err != nil {
			fmt.Println(err)
		}
		if flagVar {
			fmt.Println("This is a different message")
			return
		}

		if len(args) < 1 {
			fmt.Println("Please 'weatherme' keyword followed by a city")
			return
		}
		getWeather(args)
	},
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
	rootCmd.Flags().BoolP("differentmessage", "d", false, "Toggle a different message")
}
