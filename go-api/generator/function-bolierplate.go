package generator

import (
	"fmt"
	"gojudge/models"
	"strings"
)

// GenerateCPlusPlusBoilerplate generates C++ boilerplate code for the solution.
func GenerateCPlusPlusBoilerplate(problem models.Problem) string {
	funcSig := fmt.Sprintf("%s %s(", mapTypeToCpp(problem.Function.ReturnType), problem.Function.FunctionName)
	params := []string{}

	for _, param := range problem.Function.Parameters {
		params = append(params, fmt.Sprintf("%s %s", mapTypeToCpp(param.Type), param.Name))
	}
	funcSig += strings.Join(params, ", ") + ") {\n\t// Your code here\n\treturn result ;\n}\n"

	return fmt.Sprintf("// Problem: %s\n// Difficulty: %s\n\n%s\n",
		problem.Title, problem.Difficulty, funcSig)
}

// GenerateJavaScriptBoilerplate generates JavaScript boilerplate code for the solution.
func GenerateJavaScriptBoilerplate(problem models.Problem) string {
	funcSig := fmt.Sprintf("function %s(", problem.Function.FunctionName)
	params := []string{}

	for _, param := range problem.Function.Parameters {
		params = append(params, param.Name) // JavaScript parameters don't require types
	}
	funcSig += strings.Join(params, ", ") + ") {\n\t// Your code here\n\treturn result;\n}\n"

	return fmt.Sprintf("// Problem: %s\n// Difficulty: %s\n// \n%s\n",
		problem.Title, problem.Difficulty, funcSig)
}

// GenerateJavaBoilerplate generates Java boilerplate code for the solution.
func GenerateJavaBoilerplate(problem models.Problem) string {
	funcSig := fmt.Sprintf("public %s %s(", mapTypeToJava(problem.Function.ReturnType), problem.Function.FunctionName)
	params := []string{}

	for _, param := range problem.Function.Parameters {
		params = append(params, fmt.Sprintf("%s %s", mapTypeToJava(param.Type), param.Name))
	}
	funcSig += strings.Join(params, ", ") + ") {\n\t// Your code here\n\treturn result;\n}\n"

	// Constructing the Java class boilerplate
	boilerplate := fmt.Sprintf(`// Problem: %s
// Difficulty: %s

    %s

`, problem.Title, problem.Difficulty, funcSig)

	return boilerplate
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
