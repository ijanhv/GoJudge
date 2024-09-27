package generators

import (
	"bolierplate/utils"
	"fmt"
	"strings"
)

func GeneratePartialPython(problem utils.ProblemDefinition) string {
	inputs := generatePythonInputs(problem.InputStructure)
	output := generatePythonOutput(problem.OutputStructure)
	return fmt.Sprintf("def %s(%s) -> %s:\n    # Implementation goes here\n    return result", 
		problem.FunctionName, 
		inputs,
		output)
}

func GenerateFullPython(problem utils.ProblemDefinition) string {
	return fmt.Sprintf(`import sys
import os
from typing import List

##USER_CODE_HERE##

def main():
    with open(os.path.join(os.path.dirname(__file__), "..", "tests", "inputs", "##INPUT_FILE_INDEX##.txt")) as f:
        lines = f.readlines()
    
    %s

    result = %s(%s)
    print(result)

if __name__ == "__main__":
    main()`, 
		generatePythonInputParsing(problem.InputStructure),
		problem.FunctionName,
		strings.Join(utils.GetInputNames(problem.InputStructure), ", "))
}

func generatePythonInputs(inputStructure string) string {
	inputs := strings.Split(inputStructure, ", ")
	var result []string
	for _, input := range inputs {
		parts := strings.Split(input, " ")
		if strings.HasPrefix(parts[0], "list<") {
			result = append(result, fmt.Sprintf("%s: List[%s]", parts[1], mapTypeToPython(strings.Trim(parts[0], "list<>"))))
		} else {
			result = append(result, fmt.Sprintf("%s: %s", parts[1], mapTypeToPython(parts[0])))
		}
	}
	return strings.Join(result, ", ")
}

func generatePythonOutput(outputStructure string) string {
	if strings.HasPrefix(outputStructure, "list<") {
		return fmt.Sprintf("List[%s]", mapTypeToPython(strings.Trim(outputStructure, "list<>")))
	}
	return mapTypeToPython(outputStructure)
}

func generatePythonInputParsing(inputStructure string) string {
	inputs := strings.Split(inputStructure, ", ")
	var result []string
	for i, input := range inputs {
		parts := strings.Split(input, " ")
		if strings.HasPrefix(parts[0], "list<") {
			result = append(result, fmt.Sprintf("%s = [%s(x) for x in lines[%d].strip().split()]", parts[1], mapTypeToPython(strings.Trim(parts[0], "list<>")), i))
		} else {
			result = append(result, fmt.Sprintf("%s = %s(lines[%d].strip())", parts[1], mapTypeToPython(parts[0]), i))
		}
	}
	return strings.Join(result, "\n    ")
}

func mapTypeToPython(goType string) string {
	switch goType {
	case "int":
		return "int"
	case "float64":
		return "float"
	case "string":
		return "str"
	case "bool":
		return "bool"
	default:
		return goType // For custom types, we'll just use the same name
	}
}