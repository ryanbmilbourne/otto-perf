package performance

import (
	"math"
	"testing"
)

func TestTakeoffPerformance(t *testing.T) {
	calculator := NewTakeoffCalculator()

	// Test cases based on our chart analysis
	testCases := []struct {
		name           string
		params         TakeoffParams
		expectedDist   float64
		expectedLiftoff float64
		expectedBarrier float64
		tolerance      float64
	}{
		{
			name: "POH Example Case",
			params: TakeoffParams{
				PressureAltitude: 1500,
				Temperature:      ConvertFahrenheitToCelsius(80),
				Weight:           2325,
				WindComponent:    15,
			},
			expectedDist:    2100,
			expectedLiftoff: 50,
			expectedBarrier: 55,
			tolerance:       50, // Allow for some interpolation differences
		},
		{
			name: "Lower Weight Example",
			params: TakeoffParams{
				PressureAltitude: 1500,
				Temperature:      ConvertFahrenheitToCelsius(80),
				Weight:           2200,
				WindComponent:    15,
			},
			expectedDist:    1875,
			expectedLiftoff: 48,
			expectedBarrier: 54,
			tolerance:       50,
		},
		{
			name: "No Wind Example",
			params: TakeoffParams{
				PressureAltitude: 1500,
				Temperature:      ConvertFahrenheitToCelsius(80),
				Weight:           2200,
				WindComponent:    0,
			},
			expectedDist:    2250,
			expectedLiftoff: 48,
			expectedBarrier: 54,
			tolerance:       50,
		},
		{
			name: "Tailwind Example",
			params: TakeoffParams{
				PressureAltitude: 1500,
				Temperature:      ConvertFahrenheitToCelsius(80),
				Weight:           2200,
				WindComponent:    -5, // 5kt tailwind
			},
			expectedDist:    2500,
			expectedLiftoff: 48,
			expectedBarrier: 54,
			tolerance:       50,
		},
		{
			name: "Sea Level Standard Day",
			params: TakeoffParams{
				PressureAltitude: 0,
				Temperature:      15, // 15°C is standard day at sea level
				Weight:           2000,
				WindComponent:    0,
			},
			expectedDist:    1425, // Estimated from chart
			expectedLiftoff: 46,
			expectedBarrier: 52,
			tolerance:       50,
		},
		{
			name: "High Altitude Cold",
			params: TakeoffParams{
				PressureAltitude: 6000,
				Temperature:      -20,
				Weight:           1800,
				WindComponent:    0,
			},
			expectedDist:    1900, // Estimated from chart
			expectedLiftoff: 44,
			expectedBarrier: 50,
			tolerance:       50,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := calculator.CalculateTakeoff(tc.params)
			if err != nil {
				t.Fatalf("Error calculating takeoff: %v", err)
			}
			
			// Check takeoff distance
			if math.Abs(result.TakeoffDistance-tc.expectedDist) > tc.tolerance {
				t.Errorf("Takeoff distance incorrect: got %.0f, expected %.0f (±%.0f)",
					result.TakeoffDistance, tc.expectedDist, tc.tolerance)
			}
			
			// Check liftoff speed
			if math.Abs(result.LiftoffSpeed-tc.expectedLiftoff) > 1 {
				t.Errorf("Liftoff speed incorrect: got %.1f, expected %.1f",
					result.LiftoffSpeed, tc.expectedLiftoff)
			}
			
			// Check barrier speed
			if math.Abs(result.BarrierSpeed-tc.expectedBarrier) > 1 {
				t.Errorf("Barrier speed incorrect: got %.1f, expected %.1f",
					result.BarrierSpeed, tc.expectedBarrier)
			}
		})
	}
}

func TestInputValidation(t *testing.T) {
	calculator := NewTakeoffCalculator()
	
	testCases := []struct {
		name           string
		params         TakeoffParams
		shouldError    bool
	}{
		{
			name: "Valid Inputs",
			params: TakeoffParams{
				PressureAltitude: 3000,
				Temperature:      20,
				Weight:           2000,
				WindComponent:    10,
			},
			shouldError: false,
		},
		{
			name: "Altitude Too High",
			params: TakeoffParams{
				PressureAltitude: 8000, // Above 7000 ft limit
				Temperature:      20,
				Weight:           2000,
				WindComponent:    10,
			},
			shouldError: true,
		},
		{
			name: "Temperature Too Low",
			params: TakeoffParams{
				PressureAltitude: 3000,
				Temperature:      -50, // Below -40°C limit
				Weight:           2000,
				WindComponent:    10,
			},
			shouldError: true,
		},
		{
			name: "Temperature Too High",
			params: TakeoffParams{
				PressureAltitude: 3000,
				Temperature:      50, // Above 40°C limit
				Weight:           2000,
				WindComponent:    10,
			},
			shouldError: true,
		},
		{
			name: "Weight Too Low",
			params: TakeoffParams{
				PressureAltitude: 3000,
				Temperature:      20,
				Weight:           1500, // Below 1600 lb limit
				WindComponent:    10,
			},
			shouldError: true,
		},
		{
			name: "Weight Too High",
			params: TakeoffParams{
				PressureAltitude: 3000,
				Temperature:      20,
				Weight:           2400, // Above 2325 lb limit
				WindComponent:    10,
			},
			shouldError: true,
		},
		{
			name: "Headwind Too High",
			params: TakeoffParams{
				PressureAltitude: 3000,
				Temperature:      20,
				Weight:           2000,
				WindComponent:    20, // Above 15 kt limit
			},
			shouldError: true,
		},
		{
			name: "Tailwind Too High",
			params: TakeoffParams{
				PressureAltitude: 3000,
				Temperature:      20,
				Weight:           2000,
				WindComponent:    -10, // Beyond -5 kt limit
			},
			shouldError: true,
		},
		{
			name: "Below Sea Level",
			params: TakeoffParams{
				PressureAltitude: -500, // Should use sea level values
				Temperature:      20,
				Weight:           2000,
				WindComponent:    10,
			},
			shouldError: false, // Should not error, use sea level values
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := calculator.CalculateTakeoff(tc.params)
			
			if tc.shouldError && err == nil {
				t.Errorf("Expected error for invalid inputs, but got none")
			}
			
			if !tc.shouldError && err != nil {
				t.Errorf("Expected no error for valid inputs, but got: %v", err)
			}
		})
	}
}

func TestInterpolationFunctions(t *testing.T) {
	testCases := []struct {
		array    []float64
		value    float64
		idx1     int
		idx2     int
		fraction float64
	}{
		{[]float64{0, 1000, 2000, 3000}, 1500, 1, 2, 0.5},
		{[]float64{0, 1000, 2000, 3000}, 0, 0, 0, 0.0},
		{[]float64{0, 1000, 2000, 3000}, 3000, 3, 3, 0.0},
		{[]float64{0, 1000, 2000, 3000}, -100, 0, 0, 0.0}, // Below min
		{[]float64{0, 1000, 2000, 3000}, 4000, 3, 3, 0.0}, // Above max
	}
	
	for i, tc := range testCases {
		idx1, idx2, frac := findInterpolationIndices(tc.array, tc.value)
		
		if idx1 != tc.idx1 || idx2 != tc.idx2 || math.Abs(frac-tc.fraction) > 0.001 {
			t.Errorf("Case %d: Got (%d, %d, %.3f), expected (%d, %d, %.3f)",
				i, idx1, idx2, frac, tc.idx1, tc.idx2, tc.fraction)
		}
	}
}

func TestTemperatureConversion(t *testing.T) {
	testCases := []struct {
		fahrenheit float64
		celsius    float64
	}{
		{32, 0},
		{-40, -40}, // Interesting case where F = C
		{68, 20},
		{-4, -20},
		{104, 40},
	}
	
	for _, tc := range testCases {
		// Test F to C
		gotC := ConvertFahrenheitToCelsius(tc.fahrenheit)
		if math.Abs(gotC-tc.celsius) > 0.1 {
			t.Errorf("F to C conversion: got %.1f°C, expected %.1f°C for %.1f°F", 
				gotC, tc.celsius, tc.fahrenheit)
		}
		
		// Test C to F
		gotF := ConvertCelsiusToFahrenheit(tc.celsius)
		if math.Abs(gotF-tc.fahrenheit) > 0.1 {
			t.Errorf("C to F conversion: got %.1f°F, expected %.1f°F for %.1f°C", 
				gotF, tc.fahrenheit, tc.celsius)
		}
	}