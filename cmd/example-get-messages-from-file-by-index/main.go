package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"runtime/debug"
	"time"

	codes "github.com/scorix/go-eccodes"
)

func main() {
	filename := flag.String("file", "", "io path, e.g. /tmp/dirpw_surface_1.grib2")

	flag.Parse()

	filter := map[string]interface{}{
		"shortName": "dirpw",
	}

	file, err := codes.OpenFileByPathWithFilter(*filename, filter)
	if err != nil {
		log.Fatalf("failed to open file: %s", err.Error())
	}
	defer file.Close()

	n := 0
	for {
		err = process(file, n)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("failed to get message (#%d) from index: %s", n, err.Error())
		}

		n++
	}
}

func process(file codes.File, n int) error {
	start := time.Now()

	msg, err := file.Next()
	if err != nil {
		return err
	}
	defer msg.Close()

	log.Printf("============= BEGIN MESSAGE N%d ==========\n", n)

	startStep, err := msg.GetString("startStep")
	if err != nil {
		return fmt.Errorf("failed to get 'startStep' value: %w", err)
	}
	log.Printf("startStep = %s\n", startStep)

	endStep, err := msg.GetString("endStep")
	if err != nil {
		return fmt.Errorf("failed to get 'endStep' value: %w", err)
	}
	log.Printf("endStep = %s\n", endStep)

	stepRange, err := msg.GetString("stepRange")
	if err != nil {
		return fmt.Errorf("failed to get 'stepRange' value: %w", err)
	}
	log.Printf("stepRange = %s\n", stepRange)

	forecastTime, err := msg.GetString("forecastTime")
	if err != nil {
		return fmt.Errorf("failed to get 'forecastTime' value: %w", err)
	}
	log.Printf("forecastTime = %s\n", forecastTime)

	shortName, err := msg.GetString("shortName")
	if err != nil {
		return fmt.Errorf("failed to get 'shortName' value: %w", err)
	}
	name, err := msg.GetString("name")
	if err != nil {
		return fmt.Errorf("failed to get 'name' value: %w", err)
	}

	log.Printf("Variable = [%s](%s)\n", shortName, name)

	size, err := msg.GetLong("numberOfDataPoints")
	if err != nil {
		return fmt.Errorf("failed to get 'numberOfDataPoints' value: %w", err)
	}

	// just to measure timing
	lats, lons, vals, err := msg.Data()
	if err != nil {
		return fmt.Errorf("failed to get data (latitudes, longitudes, values): %w", err)
	}

	for i := int64(0); i < size; i++ {
		if math.IsNaN(vals[i]) {
			continue
		}

		log.Printf("(%.02f,%.02f): %.02f", lats[i], lons[i]-180, vals[i])
	}

	log.Printf("elapsed=%.0f ms", time.Since(start).Seconds()*1000)
	log.Printf("============= END MESSAGE N%d ============\n\n", n)

	debug.FreeOSMemory()

	return nil
}
