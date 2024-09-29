package generators

// import (
// 	"bolierplate/utils"
// 	"fmt"
// 	"strings"
// )

// // mapTypeToCpp maps Go data types to C++ data types
// func mapTypeToCpp(goType string) string {
// 	switch goType {
// 	case "int":
// 		return "int"
// 	case "float64":
// 		return "double" // C++ uses double for floating-point numbers
// 	case "string":
// 		return "std::string"
// 	case "bool":
// 		return "bool"
// 	default:
// 		return goType // For custom types, we'll just use the same name
// 	}
// }

// // mapVectorToCpp handles list types and converts them to C++ vector types
// func mapVectorToCpp(goType string) string {
// 	baseType := strings.TrimPrefix(goType, "list<")
// 	baseType = strings.TrimSuffix(baseType, ">")
// 	return fmt.Sprintf("std::vector<%s>", mapTypeToCpp(baseType))
// }

// func GeneratePartialCpp(problem utils.ProblemDefinition) string {
// 	inputs := strings.Replace(problem.InputStructure, "list<", "vector<", -1)
// 	inputs = strings.Replace(inputs, "string", "std::string", -1)
// 	output := strings.Replace(problem.OutputStructure, "list<", "vector<", -1)
// 	output = strings.Replace(output, "string", "std::string", -1)

// 	return fmt.Sprintf("%s %s(%s) {\n    // Implementation goes here\n    return result;\n}",
// 		strings.Split(output, " ")[0],
// 		problem.FunctionName,
// 		inputs)
// }

// func GenerateFullCpp(problem utils.ProblemDefinition) string {
// 	output := strings.Replace(problem.OutputStructure, "list<", "std::vector<", -1)
// 	output = strings.Replace(output, "string", "std::string", -1)

// 	return fmt.Sprintf(`#include <iostream>
// #include <vector>
// #include <string>
// #include <sstream>

// using namespace std;

// ##USER_CODE_HERE##

// int main() {
//     // Read input

//     // Call function
//     %s result = %s(%s);

//     // Print output
//     cout << result << endl;

//     return 0;
// }`,
// 		// generateInputReading(problem.InputStructure),
// 		strings.Split(output, " ")[0],
// 		problem.FunctionName,
// 		strings.Join(utils.GetInputNames(problem.InputStructure), ", "))
// }

// func generateInputReading(inputStructure string) string {
// 	inputs := strings.Split(inputStructure, ", ")
// 	var result strings.Builder

// 	for _, input := range inputs {
// 		parts := strings.Split(input, " ")
// 		goType := strings.Trim(parts[0], "list<>") // Get the Go type
// 		if strings.HasPrefix(parts[0], "list<") {
// 			cppType := mapVectorToCpp(parts[0]) // Map to C++ vector type
// 			result.WriteString(fmt.Sprintf(`int %s_size;
//     std::cin >> %s_size;
//     %s %s(%s_size);
//     for (int i = 0; i < %s_size; ++i) {
//         std::cin >> %s[i];
//     }
// `, parts[1], parts[1], cppType, parts[1], parts[1], parts[1], parts[1]))
// 		} else {
// 			cppType := mapTypeToCpp(goType) // Map to C++ type
// 			result.WriteString(fmt.Sprintf("    %s %s;\n    std::cin >> %s;\n", cppType, parts[1], parts[1]))
// 		}
// 	}

// 	return result.String()
// }

import (
	"strings"
)

// GenerateFullCpp generates the full C++ boilerplate code without any test case logic.
func GenerateFullCpp(problem utils.ProblemDefinition) string {
	var sb strings.Builder

	// Include necessary headers
	sb.WriteString("#include <iostream>\n")
	sb.WriteString("#include <vector>\n")
	sb.WriteString("#include <sstream>\n")
	sb.WriteString("#include <string>\n")
	sb.WriteString("using namespace std;\n\n")

	// User code placeholder
	sb.WriteString("##USER_CODE_HERE##\n\n")

	// Main function
	sb.WriteString("int main() {\n")
	sb.WriteString("    // Your test cases will be injected here.\n")
	sb.WriteString("    // Example:\n")
	sb.WriteString("    // vector<pair<string, string>> test_cases = {{\"input1\", \"output1\"}, {\"input2\", \"output2\"}};\n")
	sb.WriteString("    // for (const auto& [input, expected] : test_cases) {\n")
	sb.WriteString("    //     // Execute your function and compare output.\n")
	sb.WriteString("    // }\n")
	sb.WriteString("    return 0;\n")
	sb.WriteString("}\n")

	return sb.String()
}
