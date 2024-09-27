package utils

import "strings"

type ProblemDefinition struct {
	Name            string     `json:"name"`
	FunctionName    string     `json:"function_name"`
	InputStructure  string     `json:"input_structure"`
	OutputStructure string     `json:"output_structure"`
	Description     string     `json:"description"`
	TestCases       []TestCase `json:"test_cases"`
}

type TestCase struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

func GetInputNames(inputStructure string) []string {
	inputs := strings.Split(inputStructure, ", ")
	var names []string
	for _, input := range inputs {
		names = append(names, strings.Split(input, " ")[1])
	}
	return names
}