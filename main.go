package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strings"
)

// ------------------------------------------------------------
// CONSTANTS
// ------------------------------------------------------------

const sml_dir = "sml"
const gpx_dir = "gpx"
const separator = "/"

// ------------------------------------------------------------
// SML
// Note: Struct member names have to start with UpperCase
//        letter in order to be processed correctly.
// ------------------------------------------------------------

type Sml struct {
	XMLName xml.Name  `xml:"sml"`
	DevLog  DeviceLog `xml:"DeviceLog"`
}

type DeviceLog struct {
	XMLName xml.Name `xml:"DeviceLog"`
	Smpls   Samples  `xml:"Samples"`
}

type Samples struct {
	XMLName xml.Name `xml:"Samples"`
	Samples []Sample `xml:"Sample"`
}

type Sample struct {
	XMLName       xml.Name `xml:"Sample"`
	VerticalSpeed float64  `xml:"VerticalSpeed"`
	Latitude      float64  `xml:"Latitude"`
	Longitude     float64  `xml:"Longitude"`
	GPSAltitude   float64  `xml:"GPSAltitude"`
	UTC           string   `xml:"UTC"`
}

// ------------------------------------------------------------
// GPX
// Note: Struct member names have to start with UpperCase
//        letter in order to be processed correctly.
// ------------------------------------------------------------

type gpx struct {
	XMLName xml.Name `xml:"gpx"`
	Trk     trk
}

type trk struct {
	XMLName xml.Name `xml:"trk"`
	TrkSeg  trkSeg
}

type trkSeg struct {
	XMLName xml.Name `xml:"trkseg"`
	TrkPt   []trkPt
}

type trkPt struct {
	XMLName xml.Name `xml:"trkpt"`
	Lat     float64  `xml:"lat,attr"`
	Lon     float64  `xml:"lon,attr"`
	Ele     float64  `xml:"ele"`
	Tim     string   `xml:"time"`
}

// ------------------------------------------------------------
// SML 2 GPX
// ------------------------------------------------------------

func sml2gpx(smlFilePath, gpxFilePath string) {
	fmt.Println(smlFilePath, " -> ", gpxFilePath)

	smlFile, _ := os.Open(smlFilePath)
	defer smlFile.Close()

	smlBytes, _ := ioutil.ReadAll(smlFile)

	var sml Sml
	err := xml.Unmarshal(smlBytes, &sml)
	if err != nil {
		fmt.Println("error unmarshaling: ", err)
	}

	var Gpx gpx
	Gpx.Trk.TrkSeg.TrkPt = make([]trkPt, 0)

	for _, sample := range sml.DevLog.Smpls.Samples {
		if sample.Latitude == 0.0 || sample.Longitude == 0.0 || sample.UTC == "" {
			continue
		}

		//fmt.Println("  Latitude=", sample.Latitude, "; Longitude=", sample.Longitude, "; GPSAltitude=", sample.GPSAltitude, "; UTC=", sample.UTC)

		pt := trkPt{
			Lat: 180.0 / math.Pi * sample.Latitude,
			Lon: 180.0 / math.Pi * sample.Longitude,
			Ele: sample.GPSAltitude,
			Tim: sample.UTC,
		}

		Gpx.Trk.TrkSeg.TrkPt = append(Gpx.Trk.TrkSeg.TrkPt, pt)
	}

	//fmt.Println(smlFilePath, ": ", Gpx)

	gpxOut, err := xml.MarshalIndent(&Gpx, " ", "  ")
	if err != nil {
		fmt.Printf("error marshaling: %v\n", err)
	}

	gpxFile, _ := os.Create(gpxFilePath)
	defer gpxFile.Close()

	gpxFile.Write(gpxOut)
}

// ------------------------------------------------------------
// MAIN
// ------------------------------------------------------------

func main() {

	dir, _ := os.Open(sml_dir)
	files, _ := dir.Readdir(-1)
	dir.Close()

	for _, file := range files {
		sml := strings.TrimSpace(file.Name())

		if filepath.Ext(strings.ToLower(sml)) != ".sml" {
			continue
		}

		gpx := sml[:strings.LastIndex(strings.ToLower(sml), ".sml")] + ".gpx"

		sml2gpx(sml_dir+separator+sml, gpx_dir+separator+gpx)
	}
}
