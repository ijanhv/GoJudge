package generators

import (
	"bolierplate/utils"
	"fmt"
	"strings"
)

func GeneratePartialJs(problem utils.ProblemDefinition) string {
	inputs := utils.GetInputNames(problem.InputStructure)
	return fmt.Sprintf("function %s(%s) {\n    // Implementation goes here\n    return result;\n}", 
		problem.FunctionName, 
		strings.Join(inputs, ", "))
}

func GenerateFullJs(problem utils.ProblemDefinition) string {
	inputs := utils.GetInputNames(problem.InputStructure)
	return fmt.Sprintf(`const fs = require('fs');
const path = require('path');

##USER_CODE_HERE##

const inputFile = path.join(__dirname, '..', 'tests', 'inputs', '##INPUT_FILE_INDEX##.txt');
const input = fs.readFileSync(inputFile, 'utf8').trim().split('\n').join(' ').split(' ');

// Parse input
%s

// Call function
const result = %s(%s);

// Print output
console.log(result);`, 
		generateJsInputParsing(problem.InputStructure),
		problem.FunctionName,
		strings.Join(inputs, ", "))
}

func generateJsInputParsing(inputStructure string) string {
	inputs := strings.Split(inputStructure, ", ")
	var result strings.Builder

	for _, input := range inputs {
		parts := strings.Split(input, " ")
		if strings.HasPrefix(parts[0], "list<") {
			result.WriteString(fmt.Sprintf(`const %s_size = parseInt(input.shift());
const %s = input.splice(0, %s_size).map(Number);
`, parts[1], parts[1], parts[1]))
		} else {
			result.WriteString(fmt.Sprintf("const %s = parseInt(input.shift());\n", parts[1]))
		}
	}

	return result.String()
}