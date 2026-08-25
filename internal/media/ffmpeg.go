package media

import (
	"errors"
	"fmt"
	"petcolor/internal/domain"
	"petcolor/internal/grading"
	"strconv"
)

type Command struct {
	Program string   `json:"program"`
	Args    []string `json:"args"`
}

func PreviewCommand(input, output string, frame domain.PreviewFrame, grade domain.GradeSession) (Command, error) {
	if input == "" || output == "" {
		return Command{}, errors.New("input and output are required")
	}
	if frame.TimestampMS < 0 {
		return Command{}, errors.New("timestamp must not be negative")
	}
	graph := grading.BuildFilter(grade)
	scale := fmt.Sprintf("scale=%d:%d:force_original_aspect_ratio=decrease", frame.Width, frame.Height)
	filter := graph.Expression + "," + scale
	args := []string{"-hide_banner", "-loglevel", "error", "-ss", seconds(frame.TimestampMS), "-i", input, "-frames:v", "1", "-vf", filter, "-q:v", "3", "-y", output}
	return Command{Program: "ffmpeg", Args: args}, nil
}

func ProbeCommand(input string) (Command, error) {
	if input == "" {
		return Command{}, errors.New("input is required")
	}
	return Command{Program: "ffprobe", Args: []string{"-v", "error", "-show_entries", "format=duration:stream=width,height,codec_type", "-of", "json", input}}, nil
}

func ExportCommand(input, output string, grade domain.GradeSession) (Command, error) {
	if input == "" || output == "" {
		return Command{}, errors.New("input and output are required")
	}
	graph := grading.BuildFilter(grade)
	return Command{Program: "ffmpeg", Args: []string{"-hide_banner", "-loglevel", "error", "-i", input, "-vf", graph.Expression, "-c:v", "libx264", "-preset", "veryfast", "-crf", "20", "-c:a", "copy", "-movflags", "+faststart", "-y", output}}, nil
}

func seconds(milliseconds int64) string {
	return strconv.FormatFloat(float64(milliseconds)/1000, 'f', 3, 64)
}

func EstimatePreviewCost(frame domain.PreviewFrame, grade domain.GradeSession) int64 {
	pixels := int64(frame.Width * frame.Height)
	nodes := int64(len(grading.BuildFilter(grade).Nodes))
	if nodes < 1 {
		nodes = 1
	}
	return pixels * nodes
}

func BatchPreviewCommands(input, root string, frames []domain.PreviewFrame, grade domain.GradeSession) ([]Command, error) {
	commands := make([]Command, 0, len(frames))
	for _, frame := range frames {
		output := fmt.Sprintf("%s/%s.jpg", root, frame.ID)
		command, err := PreviewCommand(input, output, frame, grade)
		if err != nil {
			return nil, err
		}
		commands = append(commands, command)
	}
	return commands, nil
}
