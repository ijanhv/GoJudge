package generators

import (
	"bolierplate/utils"
	"fmt"
	"strings"
)

func GeneratePartialRust(problem utils.ProblemDefinition) string {
	inputs := generateRustInputs(problem.InputStructure)
	output := generateRustOutput(problem.OutputStructure)
	return fmt.Sprintf("fn %s(%s) -> %s {\n    // Implementation goes here\n    result\n}", 
		problem.FunctionName, 
		inputs,
		output)
}

func GenerateFullRust(problem utils.ProblemDefinition) string {
	return fmt.Sprintf(`use std::fs::File;
use std::io::{self, BufRead};
use std::path::Path;

##USER_CODE_HERE##

fn main() -> io::Result<()> {
    let input_path = Path::new("../tests/inputs/##INPUT_FILE_INDEX##.txt");
    let file = File::open(&input_path)?;
    let mut lines = io::BufReader::new(file).lines();

    %s

    let result = %s(%s);
    println!("{}", result);

    Ok(())
}`, 
		generateRustInputParsing(problem.InputStructure),
		problem.FunctionName,
		strings.Join(utils.GetInputNames(problem.InputStructure), ", "))
}

func generateRustInputs(inputStructure string) string {
	inputs := strings.Split(inputStructure, ", ")
	var result []string
	for _, input := range inputs {
		parts := strings.Split(input, " ")
		if strings.HasPrefix(parts[0], "list<") {
			result = append(result, fmt.Sprintf("%s: Vec<%s>", parts[1], strings.Trim(parts[0], "list<>")))
		} else {
			result = append(result, fmt.Sprintf("%s: %s", parts[1], mapTypeToRust(parts[0])))
		}
	}
	return strings.Join(result, ", ")
}

func generateRustOutput(outputStructure string) string {
	parts := strings.Split(outputStructure, " ")
	if strings.HasPrefix(parts[0], "list<") {
		return fmt.Sprintf("Vec<%s>", strings.Trim(parts[0], "list<>"))
	}
	return mapTypeToRust(parts[0])
}

func mapTypeToRust(t string) string {
	switch t {
	case "int":
		return "i32"
	case "float":
		return "f64"
	case "string":
		return "String"
	case "bool":
		return "bool"
	default:
		return t
	}
}

func generateRustInputParsing(inputStructure string) string {
	inputs := strings.Split(inputStructure, ", ")
	var result strings.Builder

	for _, input := range inputs {
		parts := strings.Split(input, " ")
		if strings.HasPrefix(parts[0], "list<") {
			result.WriteString(fmt.Sprintf(`let %s_size: usize = lines.next().unwrap()?.parse()?;
    let %s: Vec<%s> = lines.next().unwrap()?
        .split_whitespace()
        .map(|s| s.parse().unwrap())
        .collect();
`, parts[1], parts[1], mapTypeToRust(strings.Trim(parts[0], "list<>"))))
		} else {
			result.WriteString(fmt.Sprintf("let %s: %s = lines.next().unwrap()?.parse()?;\n", parts[1], mapTypeToRust(parts[0])))
		}
	}

	return result.String()
}