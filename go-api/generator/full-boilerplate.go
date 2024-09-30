package generator

import (
	"fmt"
	"strings"
)

// C++ Boilerplate code
func GenerateFullCPlusPlusBoilerplate() string {
	cppBoilerplate := `
#include <iostream>
using namespace std;
#USER CODE HERE#

int main() {
    #TEST CASE INPUT#
    return 0;
}
`
	return cppBoilerplate
}

// Java Boilerplate code
func GenerateFullJavaBoilerplate() string {
	javaBoilerplate := `
public class Main {
    #USER CODE HERE#
    
    public static void main(String[] args) {
        #TEST CASE INPUT#
    }
}
`
	return javaBoilerplate
}

// JavaScript Boilerplate code
func GenerateFullJavaScriptBoilerplate() string {
	jsBoilerplate := `
// JavaScript solution function
#USER CODE HERE#

// #TEST CASE INPUT#
`
	return jsBoilerplate
}

// Generate Boilerplate for any specified language
func GenerateBoilerplate(language string) (string, error) {
	switch strings.ToLower(language) {
	case "cpp", "c++":
		return GenerateFullCPlusPlusBoilerplate(), nil
	case "java":
		return GenerateFullJavaBoilerplate(), nil
	case "js", "javascript":
		return GenerateFullJavaScriptBoilerplate(), nil
	default:
		return "", fmt.Errorf("unsupported language: %s", language)
	}
}
