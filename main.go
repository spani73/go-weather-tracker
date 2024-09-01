package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"fmt"
)

type apiConfigData struct{
	OpenWeatherMapApiKey string `json:"OpenWeatherMapApiKey"`
}

type weatherData struct{
	Name string `json:"name"`
	Main struct{
		Kelvin float64 `json:"temp"`
	}`json:"main"`
}

func loadApiConfig(filename string) (apiConfigData, error){
	bytes , err := ioutil.ReadFile(filename)

	if err != nil {
		return apiConfigData{} , err
	}

	var c apiConfigData
	err = json.Unmarshal(bytes , &c)
	if err != nil{
		return apiConfigData{},err
	}
	return c,nil
}

func hello(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("hello from go!\n"))
}

func query(city string)(weatherData, error){
	apiConfig , err := loadApiConfig(".apiConfig")
	if err != nil{
		return weatherData{},err
	}

	resp , err := http.Get("http://api.openweathermap.org/data/2.5/weather?"+ "q=" + city + "&appid="+apiConfig.OpenWeatherMapApiKey )
	if err != nil{ 
		return weatherData{},err
	}
	fmt.Println(resp ,apiConfig.OpenWeatherMapApiKey  )
	defer resp.Body.Close()
	var d weatherData
	if err := json.NewDecoder(resp.Body).Decode(&d);err != nil{
		return weatherData{},err
	}
	return d,nil
}

func main()  {
	http.HandleFunc("/hello" , hello)
	http.HandleFunc("/weather/",
	func(w http.ResponseWriter , r *http.Request){
		city:= strings.Split(r.URL.Path, "/" )[2]
		data , err := query(city)
		if err != nil{
			http.Error(w, err.Error() , http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type" , "appication/json; charset=utf-8")
		json.NewEncoder(w).Encode(data)
	})

	http.ListenAndServe(":8080", nil)
}

