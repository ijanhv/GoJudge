package main

import (
	"fmt"
	"strings"
	"time"
)

// Problem represents a coding problem.
// Base model with common fields
type BaseModel struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Problem represents a coding problem.
type Problem struct {
	BaseModel
	Title       string            `gorm:"type:varchar(255);not null" json:"title"`                            // Problem title.
	Description string            `gorm:"type:text;not null" json:"description"`                              // Problem description.
	Difficulty  string            `gorm:"type:varchar(50);not null" json:"difficulty"`                        // Difficulty level (e.g., Easy, Medium, Hard).
	Tags        []string          `gorm:"type:varchar(255);" json:"tags"`                                     // Tags associated with the problem.
	Author      string            `gorm:"type:varchar(255);not null" json:"author"`                           // Author of the problem.
	Function    FunctionSignature `gorm:"foreignKey:ProblemID;constraint:OnDelete:CASCADE;" json:"function"`  // Function signature.
	TestCases   []TestCase        `gorm:"foreignKey:ProblemID;constraint:OnDelete:CASCADE;" json:"testCases"` // Test cases for the problem.
}

// FunctionSignature represents the function signature for a problem.
type FunctionSignature struct {
	BaseModel
	ProblemID    uint        `gorm:"not null" json:"problemId"`                                             // Reference to the problem.
	FunctionName string      `gorm:"type:varchar(100);not null" json:"functionName"`                        // Name of the function.
	Parameters   []Parameter `gorm:"foreignKey:SignatureID;constraint:OnDelete:CASCADE;" json:"parameters"` // List of function parameters.
	ReturnType   string      `gorm:"type:varchar(50);not null" json:"returnType"`                           // Expected return type of the function.
}

// Parameter represents a parameter of the function signature.
type Parameter struct {
	BaseModel
	SignatureID uint   `gorm:"not null" json:"signatureId"`           // Reference to the function signature.
	Name        string `gorm:"type:varchar(50);not null" json:"name"` // Parameter name.
	Type        string `gorm:"type:varchar(50);not null" json:"type"` // Parameter type (e.g., "int[]", "TreeNode").
}

// TestCase represents a test case for a given problem.
type TestCase struct {
	BaseModel
	ProblemID uint                   `gorm:"not null" json:"problemId"`        // Reference to the problem.
	Input     map[string]interface{} `gorm:"type:text;not null" json:"input"`  // Input as a JSON object.
	Output    interface{}            `gorm:"type:text;not null" json:"output"` // Expected output as a JSON object or array.
	IsHidden  bool                   `gorm:"default:false" json:"isHidden"`    // Is the test case hidden from the user.
}

// GenerateCPlusPlusBoilerplate generates C++ boilerplate code for the solution.
func GenerateCPlusPlusBoilerplate(problem Problem) string {
	funcSig := fmt.Sprintf("%s %s(", mapTypeToCpp(problem.Function.ReturnType), problem.Function.FunctionName)
	params := []string{}

	for _, param := range problem.Function.Parameters {
		params = append(params, fmt.Sprintf("%s %s", mapTypeToCpp(param.Type), param.Name))
	}
	funcSig += strings.Join(params, ", ") + ") {\n\t// Your code here\n\treturn result ;\n}\n"

	return fmt.Sprintf("// Problem: %s\n// Difficulty: %s\n// Author: %s\n\n%s\n",
		problem.Title, problem.Difficulty, problem.Author, funcSig)
}

// GenerateJavaScriptBoilerplate generates JavaScript boilerplate code for the solution.
func GenerateJavaScriptBoilerplate(problem Problem) string {
	funcSig := fmt.Sprintf("function %s(", problem.Function.FunctionName)
	params := []string{}

	for _, param := range problem.Function.Parameters {
		params = append(params, param.Name) // JavaScript parameters don't require types
	}
	funcSig += strings.Join(params, ", ") + ") {\n\t// Your code here\n\treturn result;\n}\n"

	return fmt.Sprintf("// Problem: %s\n// Difficulty: %s\n// Author: %s\n\n%s\n",
		problem.Title, problem.Difficulty, problem.Author, funcSig)
}

// GenerateFullCPlusPlusBoilerplate generates full C++ boilerplate code for the solution.
func GenerateFullCPlusPlusBoilerplate(problem Problem) string {
	funcSig := fmt.Sprintf("%s %s(", mapTypeToCpp(problem.Function.ReturnType), problem.Function.FunctionName)
	params := []string{}

	for _, param := range problem.Function.Parameters {
		params = append(params, fmt.Sprintf("%s %s", mapTypeToCpp(param.Type), param.Name))
	}
	funcSig += strings.Join(params, ", ") + ") {\n\t// Your code here\n\treturn result ;\n}\n"

	// Constructing the full boilerplate code
	boilerplate := fmt.Sprintln(`#include <iostream>
#include <vector>
#include <string>
using namespace std;

#USER CODE HERE#


int main() {
	#TEST CASE INPUT#
    // Your code to test the function goes here
    return 0;
}`)

	return boilerplate
}

// GenerateFullJavaScriptBoilerplate generates full JavaScript boilerplate code for the solution.
func GenerateFullJavaScriptBoilerplate(problem Problem) string {
	funcSig := fmt.Sprintf("function %s(", problem.Function.FunctionName)
	params := []string{}

	for _, param := range problem.Function.Parameters {
		params = append(params, param.Name) // JavaScript parameters don't require types
	}
	funcSig += strings.Join(params, ", ") + ") {\n\t // #TEST CASE INPUT# // Your code here\n\treturn result;\n}\n"

	// Constructing the full boilerplate code
	boilerplate := fmt.Sprintln(`
#USER CODE HERE#
// Your code to test the function goes here
`, problem.Title, problem.Difficulty)

	return boilerplate
}

// GenerateJavaBoilerplate generates Java boilerplate code for the solution.
func GenerateJavaBoilerplate(problem Problem) string {
	funcSig := fmt.Sprintf("public %s %s(", mapTypeToJava(problem.Function.ReturnType), problem.Function.FunctionName)
	params := []string{}

	for _, param := range problem.Function.Parameters {
		params = append(params, fmt.Sprintf("%s %s", mapTypeToJava(param.Type), param.Name))
	}
	funcSig += strings.Join(params, ", ") + ") {\n\t// Your code here\n\treturn result;\n}\n"

	// Constructing the Java class boilerplate
	boilerplate := fmt.Sprintf(`// Problem: %s
// Difficulty: %s
// Author: %s

public class Solution {
    %s
}
`, problem.Title, problem.Difficulty, problem.Author, funcSig)

	return boilerplate
}

// GenerateFullJavaBoilerplate generates full Java boilerplate code for the solution.
func GenerateFullJavaBoilerplate(problem Problem) string {
	funcSig := fmt.Sprintf("%s %s(", mapTypeToJava(problem.Function.ReturnType), problem.Function.FunctionName)
	params := []string{}

	for _, param := range problem.Function.Parameters {
		params = append(params, fmt.Sprintf("%s %s", mapTypeToJava(param.Type), param.Name))
	}
	funcSig += strings.Join(params, ", ") + ") {\n\t// Your code here\n\treturn result;\n}\n"

	// Constructing the full boilerplate code
	boilerplate := fmt.Sprintln(`
import java.util.*;

public class Solution {
    #USER CODE HERE#

    public static void main(String[] args) {
		#TEST CASE INPUT#
        // Your code to test the function goes here
    }
}`)

	return boilerplate
}

// GenerateTypescriptBoilerplate generates TypeScript boilerplate code for the solution.
func GenerateTypescriptBoilerplate(problem Problem) string {
	funcSig := fmt.Sprintf("function %s(", problem.Function.FunctionName)
	params := []string{}

	for _, param := range problem.Function.Parameters {
		params = append(params, fmt.Sprintf("%s: %s", param.Name, mapTypeToTypescript(param.Type)))
	}
	funcSig += strings.Join(params, ", ") + "): " + mapTypeToTypescript(problem.Function.ReturnType) + " {\n\t// Your code here\n\treturn result;\n}\n"

	return fmt.Sprintf("// Problem: %s\n// Difficulty: %s\n// Author: %s\n\n%s\n",
		problem.Title, problem.Difficulty, problem.Author, funcSig)
}

// GenerateFullTypescriptBoilerplate generates full TypeScript boilerplate code for the solution.
func GenerateFullTypescriptBoilerplate(problem Problem) string {
	funcSig := fmt.Sprintf("function %s(", problem.Function.FunctionName)
	params := []string{}

	for _, param := range problem.Function.Parameters {
		params = append(params, fmt.Sprintf("%s: %s", param.Name, mapTypeToTypescript(param.Type)))
	}
	funcSig += strings.Join(params, ", ") + "): " + mapTypeToTypescript(problem.Function.ReturnType) + " {\n\t// #TEST CASE INPUT# // Your code here\n\treturn result;\n}\n"

	// Constructing the full boilerplate code
	boilerplate := fmt.Sprintf(`// Problem: %s
// Difficulty: %s
// Author: %s

%s

// Your code to test the function goes here
`, problem.Title, problem.Difficulty, problem.Author, funcSig)

	return boilerplate
}

// generateTestCasesCpp formats test cases for C++.
// func generateTestCasesCpp(testCases []TestCase) string {
// 	var formattedCases []string
// 	for _, tc := range testCases {
// 		formattedCases = append(formattedCases, fmt.Sprintf(`{"%s", "%s"}`, tc.Input, tc.Output))
// 	}
// 	return strings.Join(formattedCases, ", ")
// }

// // generateTestCasesJS formats test cases for JavaScript.
// func generateTestCasesJS(testCases []TestCase) string {
// 	var formattedCases []string
// 	for _, tc := range testCases {
// 		formattedCases = append(formattedCases, fmt.Sprintf(`{input: "%s", output: "%s"}`, tc.Input, tc.Output))
// 	}
// 	return strings.Join(formattedCases, ", ")
// }

// main function to demonstrate boilerplate generation
func main() {
	// Example problem data

	problem := Problem{
		Title:       "Find the Duplicate Number",
		Description: "Prove that at least one duplicate number must exist.",
		Difficulty:  "Medium",
		Tags:        []string{"Array", "Binary Search"},
		Author:      "Eve Adams",
		Function: FunctionSignature{
			FunctionName: "findDuplicate",
			Parameters: []Parameter{
				{Name: "nums", Type: "int[]"},
			},
			ReturnType: "int",
		},
		TestCases: []TestCase{
			{
				Input: map[string]interface{}{
					"nums": []int{1, 3, 4, 2, 2},
				},
				Output: "2",
			},
			{
				Input: map[string]interface{}{
					"nums": []int{3, 1, 3, 4, 2},
				},
				Output: "3",
			},
			{
				Input: map[string]interface{}{
					"nums": []int{1, 1},
				},
				Output: "1",
			},
			{
				Input: map[string]interface{}{
					"nums": []int{1, 2, 3, 4, 5, 5},
				},
				Output: "5",
			},
		},
	}

	// Generate boilerplates
	// tsBoilerplate := GenerateTypescriptBoilerplate(problem)
	// fullTsBoilerplate := GenerateFullTypescriptBoilerplate(problem)

	javaBoilerplate := GenerateJavaBoilerplate(problem)
	fullJavaBoilerPlate := GenerateFullJavaBoilerplate(problem)

	// Output the boilerplates
	fmt.Println("Java Boilerplate:")
	fmt.Println(javaBoilerplate)

	fmt.Println("Full Java Boilerplate:")
	fmt.Println(fullJavaBoilerPlate)
}

func mapTypeToJava(t string) string {
	switch t {
	case "int":
		return "int"
	case "float":
		return "float"
	case "string":
		return "String"
	case "bool":
		return "boolean"
	case "int[]":
		return "int[]"
	case "string[]":
		return "String[]"
	case "float[]":
		return "float[]"
	case "bool[]":
		return "boolean[]"
	case "TreeNode":
		return "TreeNode"
	case "ListNode":
		return "ListNode"
	default:
		return t // Return as-is for other types
	}
}

func mapTypeToCpp(t string) string {
	switch t {
	case "int":
		return "int"
	case "float":
		return "float"
	case "string":
		return "std::string"
	case "bool":
		return "bool"
	case "int[]":
		return "std::vector<int>"
	case "string[]":
		return "std::vector<std::string>"
	case "float[]":
		return "std::vector<float>"
	case "bool[]":
		return "std::vector<bool>"
	case "TreeNode":
		return "TreeNode*"
	case "ListNode":
		return "ListNode*"
	default:
		return t // Return as-is for other types
	}
}

func mapTypeToTypescript(t string) string {
	switch t {
	case "int":
		return "number"
	case "float":
		return "number" // In TypeScript, both int and float are represented as number
	case "string":
		return "string"
	case "bool":
		return "boolean"
	case "int[]":
		return "number[]"
	case "string[]":
		return "string[]"
	case "float[]":
		return "number[]" // In TypeScript, both int[] and float[] are represented as number[]
	case "bool[]":
		return "boolean[]"
	case "TreeNode":
		return "TreeNode"
	case "ListNode":
		return "ListNode"
	default:
		return t // Return as-is for other types
	}
}
