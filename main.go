package main

import (
    "fmt"
    "component/weather"
)

func main() {
    city := "Yan'an"
    weatherInfo, err := weather.GetWeather(city)
    if err != nil {
        fmt.Println("Error fetching weather:", err)
        return
    }

    fmt.Printf("Weather in %s: %.1f°C, %s\n", weatherInfo.Location.Name, weatherInfo.Current.Temperature, weatherInfo.Current.Condition.Text)
}