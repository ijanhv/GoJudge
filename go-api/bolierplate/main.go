package main

import (
	"bolierplate/generators"
	"bolierplate/utils"
	"encoding/json"
	"fmt"
)

func main() {
	// Example problem definition
	problemJSON := `{
    "name": "Calculate GCD",
    "function_name": "gcd",
    "input_structure": "int a, int b",
    "output_structure": "int result",
    "description": "Given two integers, find their greatest common divisor.",
    "test_cases": [
        {
            "input": "48\n18",
            "output": "6"
        },
        {
            "input": "101\n103",
            "output": "1"
        },
        {
            "input": "56\n98",
            "output": "14"
        },
        {
            "input": "100\n25",
            "output": "25"
        }
    ]
}

`

	var problem utils.ProblemDefinition
	err := json.Unmarshal([]byte(problemJSON), &problem)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	// Generate partial boilerplate
	fmt.Println("Partial Boilerplate:")
	fmt.Println("C++:")
	fmt.Println(generators.GeneratePartialCpp(problem))
	fmt.Println("\nJavaScript:")
	fmt.Println(generators.GeneratePartialJs(problem))
	fmt.Println("\nRust:")
	fmt.Println(generators.GeneratePartialRust(problem))

	fmt.Println("\nPython:")
	fmt.Println(generators.GeneratePartialPython(problem))

	// Generate full boilerplate
	fmt.Println("\nFull Boilerplate:")
	fmt.Println("C++:")
	fmt.Println(generators.GenerateFullCpp(problem))
	fmt.Println("\nJavaScript:")
	fmt.Println(generators.GenerateFullJs(problem))
	fmt.Println("\nRust:")
	fmt.Println(generators.GenerateFullRust(problem))

	fmt.Println("\nPython:")
	fmt.Println(generators.GenerateFullPython(problem))
}