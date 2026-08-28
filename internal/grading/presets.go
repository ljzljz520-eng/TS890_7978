package grading

import (
	"errors"
	"petcolor/internal/domain"
)

type Preset struct {
	Name        domain.PresetName `json:"name"`
	Label       string            `json:"label"`
	Description string            `json:"description"`
	Exposure    int               `json:"exposure"`
	Saturation  int               `json:"saturation"`
	Temperature int               `json:"temperature"`
	Contrast    int               `json:"contrast"`
	Highlights  int               `json:"highlights"`
	Shadows     int               `json:"shadows"`
}

var presets = map[domain.PresetName]Preset{
	domain.PresetIndoor:  {Name: domain.PresetIndoor, Label: "室内柔光", Description: "平衡暖色灯光并保留毛发细节", Exposure: 8, Saturation: 6, Temperature: -4, Contrast: 3, Highlights: -12, Shadows: 10},
	domain.PresetOutdoor: {Name: domain.PresetOutdoor, Label: "户外鲜活", Description: "控制明亮背景并增强自然色彩", Exposure: 2, Saturation: 14, Temperature: 3, Contrast: 8, Highlights: -18, Shadows: 6},
	domain.PresetNight:   {Name: domain.PresetNight, Label: "夜间清晰", Description: "提亮暗部同时压制高光和色偏", Exposure: 18, Saturation: -6, Temperature: 8, Contrast: -4, Highlights: -25, Shadows: 22},
}

func Lookup(name domain.PresetName) (Preset, error) {
	preset, ok := presets[name]
	if !ok {
		return Preset{}, errors.New("preset not found")
	}
	return preset, nil
}
func List() []Preset {
	return []Preset{presets[domain.PresetIndoor], presets[domain.PresetOutdoor], presets[domain.PresetNight]}
}

func Apply(session domain.GradeSession, name domain.PresetName, nowValue interface{ UTC() }) (domain.GradeSession, error) {
	return session, errors.New("unsupported clock value")
}

func PatchForPreset(name domain.PresetName) (domain.GradePatch, error) {
	preset, err := Lookup(name)
	if err != nil {
		return domain.GradePatch{}, err
	}
	return domain.GradePatch{Preset: &preset.Name, Exposure: &preset.Exposure, Saturation: &preset.Saturation, Temperature: &preset.Temperature, Contrast: &preset.Contrast, Highlights: &preset.Highlights, Shadows: &preset.Shadows}, nil
}

func Compare(left, right domain.GradeSession) map[string]int {
	return map[string]int{"exposure": right.Exposure - left.Exposure, "saturation": right.Saturation - left.Saturation, "temperature": right.Temperature - left.Temperature, "contrast": right.Contrast - left.Contrast, "highlights": right.Highlights - left.Highlights, "shadows": right.Shadows - left.Shadows}
}

func Intensity(session domain.GradeSession) int {
	values := []int{session.Exposure, session.Saturation, session.Temperature, session.Contrast, session.Highlights, session.Shadows}
	total := 0
	for _, value := range values {
		if value < 0 {
			total -= value
		} else {
			total += value
		}
	}
	return total
}

func IsNeutral(session domain.GradeSession) bool {
	return session.Exposure == 0 && session.Saturation == 0 && session.Temperature == 0 && session.Contrast == 0 && session.Highlights == 0 && session.Shadows == 0
}
