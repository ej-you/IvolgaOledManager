package file

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"regexp"
	"sschmc/internal/app/entity"
	"strings"
)

const (
	_sensorSectionPrefix = "sensors:" // text prefix to start parse sensors
	_openBracket         = "{"        // bracket prefix to start parse sensors
	_closeBracket        = "}"        // bracket prefix to stop parse sensors

	_sensorsCap = 10 // init cap for sensors slice (length can be less or more)
)

var _ StationRepoFile = (*repoFile)(nil)

// StationRepoFile implementation.
type repoFile struct {
	configPath            string
	configPermissions     fs.FileMode
	sensorNameRegexp      *regexp.Regexp
	sensorsSettingsRegexp *regexp.Regexp
}

func NewStationRepoFile(configPath string) (StationRepoFile, error) {
	// get config file info
	fileInfo, err := os.Stat(configPath)
	if err != nil {
		return nil, fmt.Errorf("get config file info: %w", err)
	}
	// to parse sensor name from sensor line
	sensorNameRegexp := regexp.MustCompile(`.+\/(.+?)\.conf`)
	// to parse full sensors settings block [use (?s) for single-line mode]
	sensorsSettingsRegexp := regexp.MustCompile(`(?s)(^.*sensors:\s{)(.*?)(\s}.*$)`)
	return &repoFile{
		configPath:            configPath,
		configPermissions:     fileInfo.Mode().Perm(),
		sensorNameRegexp:      sensorNameRegexp,
		sensorsSettingsRegexp: sensorsSettingsRegexp,
	}, nil
}

// ParseSensors returns slice of station sensors.
// It parse station config file with the next layout:
// `_sensorSectionPrefix \n _openBracket ...[config-lines]... \n _closeBracket`.
func (r *repoFile) ParseSensors() (entity.StationSensors, error) {
	sensors := make(entity.StationSensors, 0, _sensorsCap)

	// open config file
	file, err := os.Open(r.configPath)
	if err != nil {
		return nil, fmt.Errorf("open station config: %w", err)
	}
	defer file.Close()

	var line string
	var start bool
	var idx int
	// read file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line = scanner.Text()
		// stop
		if start && strings.HasPrefix(line, _closeBracket) {
			break
		}
		// already start but not parse open bracket
		if start && strings.HasPrefix(line, _openBracket) {
			continue
		}
		// parse line
		if start {
			sensors = append(sensors, r.createSensor(idx, line))
			idx++
			continue
		}
		// start
		if strings.HasPrefix(line, _sensorSectionPrefix) {
			start = true
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read station config: %w", err)
	}
	return sensors, nil
}

// UpdateSensors updates sensor section of config file according to given sensors data.
func (r *repoFile) UpdateSensors(sensors entity.StationSensors) error {
	// read file content
	content, err := os.ReadFile(r.configPath)
	if err != nil {
		return fmt.Errorf("read full station config: %w", err)
	}
	// update config file content
	updatedContent := r.sensorsSettingsRegexp.ReplaceAll(content, sensors.CollectAll("$1", "$3"))

	// write updated content to file
	err = os.WriteFile(r.configPath, updatedContent, r.configPermissions)
	if err != nil {
		return fmt.Errorf("update config file: %w", err)
	}
	return nil
}

// createSensor creates StationSensor instance from raw config line.
func (r *repoFile) createSensor(idx int, rawConfigLine string) *entity.StationSensor {
	sensor := &entity.StationSensor{
		Idx:    idx,
		Line:   rawConfigLine,
		Active: !strings.HasPrefix(rawConfigLine, "#"),
	}
	// parse sensor name
	nameMatches := r.sensorNameRegexp.FindStringSubmatch(rawConfigLine)
	if len(nameMatches) == 0 {
		sensor.Name = rawConfigLine
	} else {
		sensor.Name = nameMatches[len(nameMatches)-1]
	}
	return sensor
}
