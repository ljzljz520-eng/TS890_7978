package grading

import (
	"fmt"
	"petcolor/internal/domain"
	"strconv"
	"strings"
)

type FilterGraph struct {
	Nodes      []string `json:"nodes"`
	Expression string   `json:"expression"`
}

func normalized(value int) float64 { return float64(value) / 100 }
func signed(value float64) string  { return strconv.FormatFloat(value, 'f', 3, 64) }

func BuildFilter(session domain.GradeSession) FilterGraph {
	nodes := []string{}
	exposure := normalized(session.Exposure)
	if exposure != 0 {
		nodes = append(nodes, "eq=brightness="+signed(exposure))
	}
	saturation := 1 + normalized(session.Saturation)
	if saturation != 1 {
		nodes = append(nodes, "eq=saturation="+signed(saturation))
	}
	contrast := 1 + normalized(session.Contrast)*0.8
	if contrast != 1 {
		nodes = append(nodes, "eq=contrast="+signed(contrast))
	}
	temperature := normalized(session.Temperature)
	if temperature != 0 {
		red := 1 + temperature*0.15
		blue := 1 - temperature*0.15
		nodes = append(nodes, fmt.Sprintf("colorbalance=rs=%s:bs=%s", signed(red-1), signed(blue-1)))
	}
	if session.Highlights != 0 || session.Shadows != 0 {
		nodes = append(nodes, fmt.Sprintf("curves=all='0/0 0.25/%s 0.75/%s 1/1'", signed(0.25+normalized(session.Shadows)*0.08), signed(0.75+normalized(session.Highlights)*0.08)))
	}
	if len(nodes) == 0 {
		nodes = append(nodes, "null")
	}
	return FilterGraph{Nodes: nodes, Expression: strings.Join(nodes, ",")}
}

func Describe(session domain.GradeSession) []string {
	out := []string{string(session.Preset)}
	pairs := []struct {
		name  string
		value int
	}{{"曝光", session.Exposure}, {"饱和度", session.Saturation}, {"色温", session.Temperature}, {"对比度", session.Contrast}, {"高光", session.Highlights}, {"阴影", session.Shadows}}
	for _, pair := range pairs {
		if pair.value != 0 {
			out = append(out, fmt.Sprintf("%s %+d", pair.name, pair.value))
		}
	}
	return out
}

func HistogramTransform(session domain.GradeSession, input [16]int) [16]int {
	var out [16]int
	shift := session.Exposure / 12
	scale := 1 + float64(session.Contrast)/125
	for index, count := range input {
		target := int((float64(index-8)*scale)+8) + shift
		if target < 0 {
			target = 0
		}
		if target > 15 {
			target = 15
		}
		out[target] += count
	}
	return out
}
