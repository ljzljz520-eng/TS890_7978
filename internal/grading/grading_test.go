package grading

import (
	"petcolor/internal/domain"
	"testing"
)

func TestPresetAndFilterGraph(t *testing.T) {
	patch, err := PatchForPreset(domain.PresetNight)
	if err != nil {
		t.Fatal(err)
	}
	grade := domain.GradeSession{ClipID: "c", Preset: domain.PresetIndoor, Revision: 1}
	grade = grade.Apply(patch, grade.UpdatedAt)
	graph := BuildFilter(grade)
	if grade.Exposure != 18 || len(graph.Nodes) < 3 {
		t.Fatalf("grade=%+v graph=%+v", grade, graph)
	}
	if Intensity(grade) == 0 {
		t.Fatal("expected intensity")
	}
}

func TestHistogramTransformPreservesCount(t *testing.T) {
	input := [16]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 8, 7, 6, 5, 4, 3, 2}
	output := HistogramTransform(domain.GradeSession{Exposure: 30, Contrast: 20}, input)
	left, right := 0, 0
	for _, value := range input {
		left += value
	}
	for _, value := range output {
		right += value
	}
	if left != right {
		t.Fatalf("input=%d output=%d", left, right)
	}
}
