# PA-28-161 Warrior II Performance Calculator

A performance calculator for the Piper PA-28-161 Cherokee Warrior II aircraft, designed to digitize and automate the performance charts from the Pilot's Operating Handbook (POH).

## Overview

This tool provides an easy way to calculate takeoff, climb, cruise, and landing performance for your PA-28-161 flights. The calculations are based on the performance charts found in the aircraft's POH, carefully digitized and implemented with accurate interpolation.

## Features

Currently implemented:
- Takeoff performance calculator (Figure 5-6: Normal Short Field Takeoff Distance)
  - Ground roll distance
  - Distance over 50ft obstacle
  - Lift-off and 50ft speeds
  - Wind corrections for both headwind and tailwind

Coming soon:
- Climb performance calculations
- Cruise performance calculations
- Landing performance calculations
- Web-based user interface

## Installation

### Prerequisites
- Go 1.16 or later

### Building from Source

```bash
# Clone the repository
git clone https://github.com/ryanbmilbourne/otto-perf.git
cd warriorperformance

# Build the takeoff CLI tool
go build -o takeoff ./cmd/takeoff
```

## Usage

### Takeoff Performance Calculator

```bash
# Basic usage with default values
./takeoff

# Calculate with specific values (metric input)
./takeoff -altitude 1500 -temp-c 25 -weight 2200 -wind 10

# Calculate with temperature in Fahrenheit
./takeoff -altitude 1500 -temp-f 77 -weight 2200 -wind 10

# Show results in metric units
./takeoff -altitude 1500 -temp-c 25 -weight 2200 -wind 10 -units metric

# Display help
./takeoff -help
```

### Command-line Options

- `-altitude`: Pressure altitude in feet (Default: 0)
- `-temp-c`: Temperature in degrees Celsius (Default: 15°C)
- `-temp-f`: Temperature in degrees Fahrenheit (overrides -temp-c if provided)
- `-weight`: Aircraft weight in pounds (Default: 2325 lbs)
- `-wind`: Wind component in knots, positive for headwind, negative for tailwind (Default: 0)
- `-units`: Unit system for display: 'imperial', 'metric', or 'mixed' (Default: imperial)
- `-help`: Display help information

### Valid Input Ranges

The calculator enforces the following limits from the POH charts:
- Pressure altitude: 0-7000 ft (sea level values used for altitudes below 0)
- Temperature: -40°C to 40°C (-40°F to 104°F)
- Weight: 1600-2325 lbs
- Headwind: 0-15 KTS
- Tailwind: 0-5 KTS

## How It Works

The calculator implements a mathematical model of the PA-28-161 POH Figure 5-6 chart using trilinear interpolation across pressure altitude, temperature, and weight dimensions, followed by appropriate wind corrections.

The implementation accurately follows the charted values, including:
1. Base takeoff distance calculation from altitude, temperature, and weight
2. Wind correction adjustments
3. Calculation of appropriate airspeeds based on weight

## For Developers

The project is structured as follows:

- `performance/`: Core performance calculation library
  - `takeoff.go`: Implementation of the takeoff performance calculations
  - `takeoff_test.go`: Unit tests for the takeoff calculations
- `cmd/`: Command-line interface tools
  - `takeoff/`: Takeoff performance CLI

To run tests:

```bash
go test ./performance
```

## Safety Notice

While this calculator aims to accurately reproduce the values from the POH charts, it is provided as a convenience tool only. Always verify all performance calculations against the official POH and ensure adequate safety margins in your flight planning.

## License

MIT License

## Acknowledgments

- Based on the PA-28-161 Warrior II Pilot's Operating Handbook charts
- Special thanks to all contributors
