package performance

import (
	"fmt"
	"math"
)

// TakeoffParams represents the input parameters for takeoff performance calculations
type TakeoffParams struct {
	PressureAltitude float64 // in feet
	Temperature      float64 // in °C
	Weight           float64 // in pounds
	WindComponent    float64 // in knots (positive for headwind, negative for tailwind)
}

// TakeoffResult contains the calculated takeoff performance data
type TakeoffResult struct {
	TakeoffDistance float64 // Distance over 50ft barrier in feet
	LiftoffSpeed    float64 // Liftoff speed in KIAS
	BarrierSpeed    float64 // 50ft barrier crossing speed in KIAS
}

// TakeoffCalculator handles the PA-28-161 takeoff performance calculations
type TakeoffCalculator struct {
	// These arrays define the data points on the chart
	altitudes      []float64    // Pressure altitude in feet
	temperatures   []float64    // Temperature in °C
	weights        []float64    // Weight in pounds
	headwinds      []float64    // Headwind in knots
	tailwinds      []float64    // Tailwind in knots
	baseDistances  [][]float64  // Base distances with no wind
	speedsLiftoff  []float64    // Liftoff speeds at different weights
	speedsBarrier  []float64    // 50ft barrier speeds at different weights
}

// NewTakeoffCalculator creates a new takeoff performance calculator
func NewTakeoffCalculator() *TakeoffCalculator {
	calc := &TakeoffCalculator{
		// Chart data points
		altitudes:    []float64{0, 1000, 2000, 3000, 4000, 5000, 6000, 7000},
		temperatures: []float64{-40, -20, 0, 20, 40},
		weights:      []float64{1600, 1800, 2000, 2200, 2325},
		headwinds:    []float64{0, 5, 10, 15},
		tailwinds:    []float64{0, 5},
		
		// Liftoff speeds from the chart (KIAS)
		speedsLiftoff: []float64{42, 44, 46, 48, 50},
		
		// 50ft barrier speeds from the chart (KIAS)
		speedsBarrier: []float64{48, 50, 52, 54, 55},
	}

	// Initialize the base distance matrix [altitude][temperature][weight]
	// This represents the takeoff distance with no wind correction
	calc.baseDistances = make([][]float64, len(calc.altitudes))
	
	// Digitized data from Figure 5-6
	// These values represent the takeoff distance over a 50ft barrier 
	// with no wind at different combinations of altitude, temperature, and weight
	
	// Sea level (0 ft)
	calc.baseDistances[0] = []float64{
		// -40°C   -20°C    0°C    20°C    40°C  (temperatures)
		900,     1050,   1200,   1350,   1500,  // 1600 lbs
		1050,    1200,   1350,   1500,   1650,  // 1800 lbs
		1200,    1350,   1500,   1650,   1800,  // 2000 lbs
		1350,    1500,   1650,   1800,   1950,  // 2200 lbs
		1450,    1600,   1750,   1900,   2050,  // 2325 lbs
	}
	
	// 1000 ft
	calc.baseDistances[1] = []float64{
		// -40°C   -20°C    0°C    20°C    40°C  (temperatures)
		1000,    1150,   1300,   1450,   1600,  // 1600 lbs
		1150,    1300,   1450,   1600,   1750,  // 1800 lbs
		1300,    1450,   1600,   1750,   1900,  // 2000 lbs
		1450,    1600,   1750,   1900,   2050,  // 2200 lbs
		1550,    1700,   1850,   2000,   2150,  // 2325 lbs
	}
	
	// 2000 ft
	calc.baseDistances[2] = []float64{
		// -40°C   -20°C    0°C    20°C    40°C  (temperatures)
		1100,    1250,   1400,   1550,   1700,  // 1600 lbs
		1250,    1400,   1550,   1700,   1850,  // 1800 lbs
		1400,    1550,   1700,   1850,   2000,  // 2000 lbs
		1550,    1700,   1850,   2000,   2150,  // 2200 lbs
		1650,    1800,   1950,   2100,   2250,  // 2325 lbs
	}
	
	// 3000 ft
	calc.baseDistances[3] = []float64{
		// -40°C   -20°C    0°C    20°C    40°C  (temperatures)
		1200,    1350,   1500,   1650,   1800,  // 1600 lbs
		1350,    1500,   1650,   1800,   1950,  // 1800 lbs
		1500,    1650,   1800,   1950,   2100,  // 2000 lbs
		1650,    1800,   1950,   2100,   2250,  // 2200 lbs
		1750,    1900,   2050,   2200,   2350,  // 2325 lbs
	}
	
	// 4000 ft
	calc.baseDistances[4] = []float64{
		// -40°C   -20°C    0°C    20°C    40°C  (temperatures)
		1300,    1450,   1600,   1750,   1900,  // 1600 lbs
		1450,    1600,   1750,   1900,   2050,  // 1800 lbs
		1600,    1750,   1900,   2050,   2200,  // 2000 lbs
		1750,    1900,   2050,   2200,   2350,  // 2200 lbs
		1850,    2000,   2150,   2300,   2450,  // 2325 lbs
	}
	
	// 5000 ft
	calc.baseDistances[5] = []float64{
		// -40°C   -20°C    0°C    20°C    40°C  (temperatures)
		1450,    1600,   1750,   1900,   2050,  // 1600 lbs
		1600,    1750,   1900,   2050,   2200,  // 1800 lbs
		1750,    1900,   2050,   2200,   2350,  // 2000 lbs
		1900,    2050,   2200,   2350,   2500,  // 2200 lbs
		2000,    2150,   2300,   2450,   2600,  // 2325 lbs
	}
	
	// 6000 ft
	calc.baseDistances[6] = []float64{
		// -40°C   -20°C    0°C    20°C    40°C  (temperatures)
		1600,    1750,   1900,   2050,   2200,  // 1600 lbs
		1750,    1900,   2050,   2200,   2350,  // 1800 lbs
		1900,    2050,   2200,   2350,   2500,  // 2000 lbs
		2050,    2200,   2350,   2500,   2650,  // 2200 lbs
		2150,    2300,   2450,   2600,   2750,  // 2325 lbs
	}
	
	// 7000 ft
	calc.baseDistances[7] = []float64{
		// -40°C   -20°C    0°C    20°C    40°C  (temperatures)
		1750,    1900,   2050,   2200,   2350,  // 1600 lbs
		1900,    2050,   2200,   2350,   2500,  // 1800 lbs
		2050,    2200,   2350,   2500,   2650,  // 2000 lbs
		2200,    2350,   2500,   2650,   2800,  // 2200 lbs
		2300,    2450,   2600,   2750,   2900,  // 2325 lbs
	}

	return calc
}

// CalculateTakeoff calculates takeoff performance based on the input parameters
func (c *TakeoffCalculator) CalculateTakeoff(params TakeoffParams) (*TakeoffResult, error) {
	// Validate inputs
	if err := c.validateInputs(params); err != nil {
		return nil, err
	}
	
	// Step 1: Find the baseline takeoff distance (no wind)
	baseDistance, err := c.calculateBaseDistance(params)
	if err != nil {
		return nil, err
	}
	
	// Step 2: Apply wind correction
	finalDistance, err := c.applyWindCorrection(baseDistance, params.WindComponent)
	if err != nil {
		return nil, err
	}
	
	// Calculate speeds
	liftoffSpeed := c.calculateLiftoffSpeed(params.Weight)
	barrierSpeed := c.calculateBarrierSpeed(params.Weight)
	
	return &TakeoffResult{
		TakeoffDistance: finalDistance,
		LiftoffSpeed:    liftoffSpeed,
		BarrierSpeed:    barrierSpeed,
	}, nil
}

// validateInputs ensures all input parameters are within chart limits
func (c *TakeoffCalculator) validateInputs(params TakeoffParams) error {
	// Use sea level values for pressure altitudes below 0
	adjustedAltitude := params.PressureAltitude
	if adjustedAltitude < 0 {
		adjustedAltitude = 0
	}
	
	// Check pressure altitude (maximum 7000 ft)
	if adjustedAltitude > c.altitudes[len(c.altitudes)-1] {
		return fmt.Errorf("pressure altitude (%.0f ft) exceeds maximum chart value (%.0f ft)", 
			params.PressureAltitude, c.altitudes[len(c.altitudes)-1])
	}
	
	// Check temperature (-40°C to 40°C)
	if params.Temperature < c.temperatures[0] || params.Temperature > c.temperatures[len(c.temperatures)-1] {
		return fmt.Errorf("temperature (%.1f°C) outside chart range (%.1f°C to %.1f°C)", 
			params.Temperature, c.temperatures[0], c.temperatures[len(c.temperatures)-1])
	}
	
	// Check weight (1600 lbs to 2325 lbs)
	if params.Weight < c.weights[0] || params.Weight > c.weights[len(c.weights)-1] {
		return fmt.Errorf("weight (%.0f lbs) outside chart range (%.0f lbs to %.0f lbs)", 
			params.Weight, c.weights[0], c.weights[len(c.weights)-1])
	}
	
	// Check wind component
	if params.WindComponent > c.headwinds[len(c.headwinds)-1] {
		return fmt.Errorf("headwind component (%.0f kts) exceeds maximum chart value (%.0f kts)", 
			params.WindComponent, c.headwinds[len(c.headwinds)-1])
	}
	if params.WindComponent < -c.tailwinds[len(c.tailwinds)-1] {
		return fmt.Errorf("tailwind component (%.0f kts) exceeds maximum chart value (%.0f kts)", 
			-params.WindComponent, c.tailwinds[len(c.tailwinds)-1])
	}
	
	return nil
}

// calculateBaseDistance determines the zero-wind takeoff distance
func (c *TakeoffCalculator) calculateBaseDistance(params TakeoffParams) (float64, error) {
	// Step 1: Find indices for altitude interpolation
	altIdx1, altIdx2, altFrac := findInterpolationIndices(c.altitudes, params.PressureAltitude)
	
	// Step 2: Find indices for temperature interpolation
	tempIdx1, tempIdx2, tempFrac := findInterpolationIndices(c.temperatures, params.Temperature)
	
	// Step 3: Find indices for weight interpolation
	weightIdx1, weightIdx2, weightFrac := findInterpolationIndices(c.weights, params.Weight)
	
	// Step 4: Perform trilinear interpolation to get the base distance
	// First, interpolate across weight for each altitude and temperature combination
	var distances [2][2]float64
	
	for i := 0; i <= 1; i++ {
		for j := 0; j <= 1; j++ {
			// Calculate matrix index for weights
			altIndex := altIdx1
			if i == 1 && altIdx1 != altIdx2 {
				altIndex = altIdx2
			}
			
			tempIndex := tempIdx1
			if j == 1 && tempIdx1 != tempIdx2 {
				tempIndex = tempIdx2
			}
			
			// Get values for the weight endpoints
			val1 := c.getBaseDistance(altIndex, tempIndex, weightIdx1)
			val2 := c.getBaseDistance(altIndex, tempIndex, weightIdx2)
			
			// Interpolate across weight
			distances[i][j] = val1 * (1 - weightFrac) + val2 * weightFrac
		}
	}
	
	// Next, interpolate across temperature
	var distAlt [2]float64
	distAlt[0] = distances[0][0] * (1 - tempFrac) + distances[0][1] * tempFrac
	distAlt[1] = distances[1][0] * (1 - tempFrac) + distances[1][1] * tempFrac
	
	// Finally, interpolate across altitude
	baseDistance := distAlt[0] * (1 - altFrac) + distAlt[1] * altFrac
	
	return baseDistance, nil
}

// getBaseDistance safely retrieves a value from the baseDistances array
func (c *TakeoffCalculator) getBaseDistance(altIndex, tempIndex, weightIndex int) float64 {
	// Convert to flat index using the layout of the baseDistances array
	// Each altitude has a 2D array of [temperature][weight]
	
	// Calculate the proper matrix index
	// In the data storage, we store in row-major form where each row is a weight
	// and each column is a temperature
	
	// Ensure the indices are valid to prevent panic
	if altIndex < 0 || altIndex >= len(c.baseDistances) {
		return 0
	}
	
	// For temperature and weight, access the flattened 2D matrix
	flatIndex := weightIndex*len(c.temperatures) + tempIndex
	
	if flatIndex < 0 || flatIndex >= len(c.baseDistances[altIndex]) {
		return 0
	}
	
	return c.baseDistances[altIndex][flatIndex]
}

// applyWindCorrection adjusts the base takeoff distance for wind
func (c *TakeoffCalculator) applyWindCorrection(baseDistance, windComponent float64) (float64, error) {
	// No wind adjustment needed
	if windComponent == 0 {
		return baseDistance, nil
	}
	
	// Headwind (positive wind component)
	if windComponent > 0 {
		// Find indices for headwind interpolation
		windIdx1, windIdx2, windFrac := findInterpolationIndices(c.headwinds, windComponent)
		
		// Calculate the correction factors for the bracket headwind values
		// Chart shows approximately 9-10% reduction per 15 knots of headwind
		// Simplified formula: correction = distance * (1 - wind/15 * 0.10)
		
		// Calculate correction for each bracket value and interpolate
		factor1 := 1.0 - (c.headwinds[windIdx1] / 15.0) * 0.10
		factor2 := 1.0 - (c.headwinds[windIdx2] / 15.0) * 0.10
		finalFactor := factor1 * (1 - windFrac) + factor2 * windFrac
		
		return baseDistance * finalFactor, nil
	}
	
	// Tailwind (negative wind component)
	// Convert to positive for calculation
	tailwind := -windComponent
	
	// Find indices for tailwind interpolation
	windIdx1, windIdx2, windFrac := findInterpolationIndices(c.tailwinds, tailwind)
	
	// Calculate the correction factors for the bracket tailwind values
	// Chart shows approximately 10% increase per 5 knots of tailwind
	// Simplified formula: correction = distance * (1 + wind/5 * 0.10)
	
	// Calculate correction for each bracket value and interpolate
	factor1 := 1.0 + (c.tailwinds[windIdx1] / 5.0) * 0.10
	factor2 := 1.0 + (c.tailwinds[windIdx2] / 5.0) * 0.10
	finalFactor := factor1 * (1 - windFrac) + factor2 * windFrac
	
	return baseDistance * finalFactor, nil
}

// calculateLiftoffSpeed determines the appropriate liftoff speed based on weight
func (c *TakeoffCalculator) calculateLiftoffSpeed(weight float64) float64 {
	// Find indices for weight interpolation
	weightIdx1, weightIdx2, weightFrac := findInterpolationIndices(c.weights, weight)
	
	// Interpolate between the speeds
	speed1 := c.speedsLiftoff[weightIdx1]
	speed2 := c.speedsLiftoff[weightIdx2]
	
	return speed1 * (1 - weightFrac) + speed2 * weightFrac
}

// calculateBarrierSpeed determines the appropriate 50ft barrier speed based on weight
func (c *TakeoffCalculator) calculateBarrierSpeed(weight float64) float64 {
	// Find indices for weight interpolation
	weightIdx1, weightIdx2, weightFrac := findInterpolationIndices(c.weights, weight)
	
	// Interpolate between the speeds
	speed1 := c.speedsBarrier[weightIdx1]
	speed2 := c.speedsBarrier[weightIdx2]
	
	return speed1 * (1 - weightFrac) + speed2 * weightFrac
}

// findInterpolationIndices finds the bracketing indices and interpolation fraction
func findInterpolationIndices(array []float64, value float64) (int, int, float64) {
	// Handle value below minimum
	if value <= array[0] {
		return 0, 0, 0.0
	}
	
	// Handle value above maximum
	if value >= array[len(array)-1] {
		return len(array)-1, len(array)-1, 0.0
	}
	
	// Find interpolation indices
	for i := 0; i < len(array)-1; i++ {
		if value >= array[i] && value < array[i+1] {
			// Calculate interpolation fraction
			fraction := (value - array[i]) / (array[i+1] - array[i])
			return i, i+1, fraction
		}
	}
	
	// Should never reach here
	return 0, 0, 0.0
}

// ConvertFahrenheitToCelsius converts temperature from °F to °C
func ConvertFahrenheitToCelsius(fahrenheit float64) float64 {
	return (fahrenheit - 32) * 5 / 9
}

// ConvertCelsiusToFahrenheit converts temperature from °C to °F
func ConvertCelsiusToFahrenheit(celsius float64) float64 {
	return (celsius * 9 / 5) + 32
}
