package generators

import (
	"bolierplate/utils"
	"fmt"
	"strings"
)

func GeneratePartialTs(problem utils.ProblemDefinition) string {
	inputs := generateTsInputs(problem.InputStructure)
	output := generateTsOutput(problem.OutputStructure)
	return fmt.Sprintf("function %s(%s): %s {\n    // Implementation goes here\n    return result;\n}", 
		problem.FunctionName, 
		inputs,
		output)
}

func GenerateFullTs(problem utils.ProblemDefinition) string {
	return fmt.Sprintf(`import * as fs from 'fs';
import * as path from 'path';

##USER_CODE_HERE##

const inputFile = path.join(__dirname, '..', 'tests', 'inputs', '##INPUT_FILE_INDEX##.txt');
const input = fs.readFileSync(inputFile, 'utf8').trim().split('\n').join(' ').split(' ');

%s

const result = %s(%s);
console.log(result);`, 
		generateTsInputParsing(problem.InputStructure),
		problem.FunctionName,
		strings.Join(utils.GetInputNames(problem.InputStructure), ", "))
}

func generateTsInputs(inputStructure string) string {
	inputs := strings.Split(inputStructure, ", ")
	var result []string
	for _, input := range inputs {
		parts := strings.Split(input, " ")
		if strings.HasPrefix(parts[0], "list<") {
			result = append(result, fmt.Sprintf("%s: %s[]", parts[1], mapTypeToTs(strings.Trim(parts[0], "list<>"))))
		} else {
			result = append(result, fmt.Sprintf("%s: %s", parts[1], mapTypeToTs(parts[0])))
		}
	}
	return strings.Join(result, ", ")
}

func generateTsOutput(outputStructure string) string {
	parts := strings.Split(outputStructure, " ")
	if strings.HasPrefix(parts[0], "list<") {
		return fmt.Sprintf("%s[]", mapTypeToTs(strings.Trim(parts[0], "list<>")))
	}
	return mapTypeToTs(parts[0])
}

func mapTypeToTs(t string) string {
	switch t {
	case "int", "float":
		return "number"
	case "string":
		return "string"
	case "bool":
		return "boolean"
	default:
		return t
	}
}

func generateTsInputParsing(inputStructure string) string {
	inputs := strings.Split(inputStructure, ", ")
	var result strings.Builder

	for _, input := range inputs {
		parts := strings.Split(input, " ")
		if strings.HasPrefix(parts[0], "list<") {
			result.WriteString(fmt.Sprintf(`const %sSize = parseInt(input.shift()!);
const %s: %s[] = input.splice(0, %sSize).map(Number);
`, parts[1], parts[1], mapTypeToTs(strings.Trim(parts[0], "list<>")), parts[1]))
		} else {
			result.WriteString(fmt.Sprintf("const %s: %s = %s(input.shift()!);\n", parts[1], mapTypeToTs(parts[0]), mapTsParser(parts[0])))
		}
	}

	return result.String()
}

func mapTsParser(t string) string {
	switch t {
	case "int", "float":
		return "Number"
	case "string":
		return ""
	case "bool":
		return "Boolean"
	default:
		return ""
	}
}