package generators

import (
	"bolierplate/utils"
	"fmt"
	"strings"
)

func GeneratePartialCpp(problem utils.ProblemDefinition) string {
	inputs := strings.Replace(problem.InputStructure, "list<", "vector<", -1)
	inputs = strings.Replace(inputs, "string", "string", -1)
	output := strings.Replace(problem.OutputStructure, "list<", "vector<", -1)
	output = strings.Replace(output, "string", "string", -1)

	return fmt.Sprintf("%s %s(%s) {\n    // Implementation goes here\n    return result;\n}", 
		strings.Split(output, " ")[0], 
		problem.FunctionName, 
		inputs)
}

func GenerateFullCpp(problem utils.ProblemDefinition) string {
	output := strings.Replace(problem.OutputStructure, "list<", "std::vector<", -1)
	output = strings.Replace(output, "string", "string", -1)

	return fmt.Sprintf(`#include <iostream>
#include <vector>
#include <string>
#include <sstream>

using namespace std;

##USER_CODE_HERE##

int main() {
    // Read input
    %s

    // Call function
    %s result = %s(%s);

    // Print output
    cout << result << endl;

    return 0;
}`, 
		generateInputReading(problem.InputStructure),
		strings.Split(output, " ")[0],
		problem.FunctionName,
		strings.Join(utils.GetInputNames(problem.InputStructure), ", "))
}
func generateInputReading(inputStructure string) string {
	inputs := strings.Split(inputStructure, ", ")
	var result strings.Builder

	for _, input := range inputs {
		parts := strings.Split(input, " ")
		if strings.HasPrefix(parts[0], "list<") {
			result.WriteString(fmt.Sprintf(`int %s_size;
    std::cin >> %s_size;
    std::vector<%s> %s(%s_size);
    for (int i = 0; i < %s_size; ++i) {
        std::cin >> %s[i];
    }
`, parts[1], parts[1], strings.Trim(parts[0], "list<>"), parts[1], parts[1], parts[1], parts[1]))
		} else {
			result.WriteString(fmt.Sprintf("    %s %s;\n    std::cin >> %s;\n", parts[0], parts[1], parts[1]))
		}
	}

	return result.String()
}