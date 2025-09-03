// Package file contains implementations of File repo interfaces.
package fs

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"regexp"
	"strings"

	"IvolgaOledManager/internal/app/entity"
	"IvolgaOledManager/internal/app/repo"
)

const (
	_sensSectionPrefix = "sensors:" // text prefix to start parse sensors
	_openBracket       = "{"        // bracket prefix to start parse sensors
	_closeBracket      = "}"        // bracket prefix to stop parse sensors

	_sensCap = 10 // init cap for sensors slice (length can be less or more)
)

// Ensure sensors' config repo implementats interface.
var _ repo.SensconfRepoFS = (*SensconfRepo)(nil)

// SensconfRepo represents a file repo for sensors' config.
type SensconfRepo struct {
	configPath            string
	configPermissions     fs.FileMode
	sensorNameRegexp      *regexp.Regexp
	sensorsSettingsRegexp *regexp.Regexp
}

// NewSensconfRepoFS returns a new instance of SensconfRepo.
func NewSensconfRepoFS(configPath string) (*SensconfRepo, error) {
	// get config file info
	fileInfo, err := os.Stat(configPath)
	if err != nil {
		return nil, fmt.Errorf("get config file info: %w", err)
	}
	// to parse sensor name from sensor line
	sensorNameRegexp := regexp.MustCompile(`.+/(.+?)\.conf`)
	// to parse full sensors settings block [use (?s) for single-line mode]
	sensorsSettingsRegexp := regexp.MustCompile(`(?s)(^.*sensors:\s{)(.*?)(\s}.*$)`)
	return &SensconfRepo{
		configPath:            configPath,
		configPermissions:     fileInfo.Mode().Perm(),
		sensorNameRegexp:      sensorNameRegexp,
		sensorsSettingsRegexp: sensorsSettingsRegexp,
	}, nil
}

// ParseSensconf returns slice of station sensors.
// It parse station config file with the next layout:
// `_sensorSectionPrefix \n _openBracket ...[config-lines]... \n _closeBracket`.
func (r *SensconfRepo) ParseSensconf() (entity.Sensconf, error) {
	sensors := make(entity.Sensconf, 0, _sensCap)

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
			sensors = append(sensors, r.createSensconfItem(idx, line))
			idx++
			continue
		}
		// start
		if strings.HasPrefix(line, _sensSectionPrefix) {
			start = true
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read station config: %w", err)
	}
	return sensors, nil
}

// UpdateSensconf updates sensor section of config file according to given sensors data.
func (r *SensconfRepo) UpdateSensconf(sensors entity.Sensconf) error {
	// read file content
	content, err := os.ReadFile(r.configPath)
	if err != nil {
		return fmt.Errorf("read full station config: %w", err)
	}
	// update config file content
	updatedContent := r.sensorsSettingsRegexp.ReplaceAll(
		content, collectSensorconfToBytes(sensors, "$1", "$3"))

	// write updated content to file
	err = os.WriteFile(r.configPath, updatedContent, r.configPermissions)
	if err != nil {
		return fmt.Errorf("update config file: %w", err)
	}
	return nil
}

// createSensconfItem creates SensconfItem instance from raw config line.
func (r *SensconfRepo) createSensconfItem(idx int, rawConfigLine string) *entity.SensconfItem {
	sensor := &entity.SensconfItem{
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

// CollectAll collects all sensor lines to a slice of bytes and returns it.
func collectSensorconfToBytes(data entity.Sensconf, prefix, suffix string) []byte {
	var builder strings.Builder
	builder.WriteString(prefix)
	// collect sensor lines
	for _, sensor := range data {
		builder.WriteString("\n")
		builder.WriteString(sensor.Line)
	}
	builder.WriteString(suffix)
	return []byte(builder.String())
}
