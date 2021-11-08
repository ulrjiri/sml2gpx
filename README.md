# SML 2 GPX

## Structure of the program
```
.
`-- Sml2GpxGo
    |-- sml
	|-- gpx
	`-- main.go
```

## Usage

1. Start the Moveslink app
2. Connect the Suunto POD
3. Let dowload the workouts. Connection errors will be reported - ignore it.
4. Files with workouts are *.sml files in the C:\Users\USERNAME\AppData\Roaming\Suunto\Moveslink2 folder.
5. Copy or move sml files which you need to convert to the sml folder (input).
6. Run main.go program
```
go run main.go
```

## Documentation

[Parsing XML Files With Golang](https://tutorialedge.net/golang/parsing-xml-with-golang/)

## Future Improvements

- Configuration file (yml/yaml) - input, output folders, input file name filter
- Options to delete source file after processing
- Ask vs. overwrite existing output file
- Interface to Strava or Movescount