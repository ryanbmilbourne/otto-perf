package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	
	"github.com/ryanbmilbourne/otto-perf/performance"
)

func main() {
	// Define CLI flags
	pressureAlt := flag.Float64("altitude", 0, "Pressure altitude in feet")
	
	// Allow temperature to be specified in either Celsius or Fahrenheit
	tempC := flag.Float64("temp-c", 15, "Temperature in °C")
	tempF := flag.Float64("temp-f", 0, "Temperature in °F (overrides temp-c if provided)")
	tempFProvided := false
	
	weight := flag.Float64("weight", 2325, "Aircraft weight in pounds")
	windComponent := flag.Float64("wind", 0, "Wind component in knots (positive for headwind, negative for tailwind)")
	unitSystem := flag.String("units", "imperial", "Unit system for display: 'imperial', 'metric', or 'mixed'")
	showHelp := flag.Bool("help", false, "Show help")
	
	// Custom usage function for better help display
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "PA-28-161 Cherokee Warrior II Takeoff Performance Calculator\n\n")
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExample:\n  %s -altitude 1500 -temp-c 25 -weight 2200 -wind 10\n", os.Args[0])
	}
	
	// Parse command line arguments
	flag.Parse()
	
	// Check if -temp-f was explicitly provided
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "temp-f" {
			tempFProvided = true
		}
	})
	
	// Show help if requested or no arguments provided
	if *showHelp || flag.NFlag() == 0 {
		flag.Usage()
		os.Exit(0)
	}
	
	// Determine temperature in Celsius
	var temperature float64
	if tempFProvided {
		temperature = performance.ConvertFahrenheitToCelsius(*tempF)
	} else {
		temperature = *tempC
	}
	
	// Create params struct with input values
	params := performance.TakeoffParams{
		PressureAltitude: *pressureAlt,
		Temperature:      temperature,
		Weight:           *weight,
		WindComponent:    *windComponent,
	}
	
	// Initialize takeoff calculator
	calculator := performance.NewTakeoffCalculator()
	
	// Calculate takeoff performance
	result, err := calculator.CalculateTakeoff(params)
	if err != nil {
		log.Fatalf("Error calculating takeoff performance: %v", err)
	}
	
	// Display results based on selected unit system
	displayResults(params, result, strings.ToLower(*unitSystem))
}

func displayResults(params performance.TakeoffParams, result *performance.TakeoffResult, unitSystem string) {
	fmt.Printf("\nPA-28-161 Cherokee Warrior II Takeoff Performance\n")
	fmt.Printf("=================================================\n\n")
	
	// Display input parameters
	fmt.Printf("Input Parameters:\n")
	fmt.Printf("----------------\n")
	
	fmt.Printf("Pressure Altitude: %.0f ft\n", params.PressureAltitude)
	
	// Display temperature in appropriate format
	switch unitSystem {
	case "metric":
		fmt.Printf("Temperature: %.1f°C\n", params.Temperature)
	case "imperial":
		fmt.Printf("Temperature: %.1f°F (%.1f°C)\n", 
			performance.ConvertCelsiusToFahrenheit(params.Temperature), params.Temperature)
	case "mixed":
		fmt.Printf("Temperature: %.1f°C (%.1f°F)\n", 
			params.Temperature, performance.ConvertCelsiusToFahrenheit(params.Temperature))
	default:
		fmt.Printf("Temperature: %.1f°C (%.1f°F)\n", 
			params.Temperature, performance.ConvertCelsiusToFahrenheit(params.Temperature))
	}
	
	fmt.Printf("Weight: %.0f lbs\n", params.Weight)
	
	// Display wind in appropriate format
	if params.WindComponent > 0 {
		fmt.Printf("Wind: %.0f knots headwind\n", params.WindComponent)
	} else if params.WindComponent < 0 {
		fmt.Printf("Wind: %.0f knots tailwind\n", -params.WindComponent)
	} else {
		fmt.Printf("Wind: No wind\n")
	}
	
	fmt.Printf("\n")
	
	// Display results
	fmt.Printf("Takeoff Performance:\n")
	fmt.Printf("-------------------\n")
	
	// Display distances in appropriate format
	switch unitSystem {
	case "metric":
		fmt.Printf("Takeoff Distance (over 50 ft obstacle): %.0f m (%.0f ft)\n", 
			feetToMeters(result.TakeoffDistance), result.TakeoffDistance)
	case "imperial":
		fmt.Printf("Takeoff Distance (over 50 ft obstacle): %.0f ft\n", result.TakeoffDistance)
	case "mixed":
		fmt.Printf("Takeoff Distance (over 50 ft obstacle): %.0f ft (%.0f m)\n", 
			result.TakeoffDistance, feetToMeters(result.TakeoffDistance))
	default:
		fmt.Printf("Takeoff Distance (over 50 ft obstacle): %.0f ft\n", result.TakeoffDistance)
	}
	
	// Display speeds
	fmt.Printf("Lift-off Speed: %.0f KIAS\n", result.LiftoffSpeed)
	fmt.Printf("50 ft Barrier Speed: %.0f KIAS\n", result.BarrierSpeed)
	
	// Safety note
	fmt.Printf("\nNOTE: Always verify these calculations against the POH and ensure\n")
	fmt.Printf("      you have adequate runway length with appropriate safety margins.\n")
}

// feetToMeters converts distance from feet to meters
func feetToMeters(feet float64) float64 {
	return feet * 0.3048
}
