"use strict";
var __create = Object.create;
var __defProp = Object.defineProperty;
var __getOwnPropDesc = Object.getOwnPropertyDescriptor;
var __getOwnPropNames = Object.getOwnPropertyNames;
var __getProtoOf = Object.getPrototypeOf;
var __hasOwnProp = Object.prototype.hasOwnProperty;
var __copyProps = (to, from, except, desc) => {
  if (from && typeof from === "object" || typeof from === "function") {
    for (let key of __getOwnPropNames(from))
      if (!__hasOwnProp.call(to, key) && key !== except)
        __defProp(to, key, { get: () => from[key], enumerable: !(desc = __getOwnPropDesc(from, key)) || desc.enumerable });
  }
  return to;
};
var __toESM = (mod, isNodeMode, target) => (target = mod != null ? __create(__getProtoOf(mod)) : {}, __copyProps(
  // If the importer is in node compatibility mode or this is not an ESM
  // file that has been converted to a CommonJS file using a Babel-
  // compatible transform (i.e. "__esModule" has not been set), then set
  // "default" to the CommonJS "module.exports" for node compatibility.
  isNodeMode || !mod || !mod.__esModule ? __defProp(target, "default", { value: mod, enumerable: true }) : target,
  mod
));

// src/index.ts
var import_fs = __toESM(require("fs"));
var import_path = __toESM(require("path"));

// src/ProblemDefinitionGenerator.ts
var ProblemDefinitionParser = class {
  constructor() {
    this.problemName = "";
    this.functionName = "";
    this.inputFields = [];
    this.outputFields = [];
  }
  parse(input) {
    const lines = input.split("\n").map((line) => line.trim());
    let currentSection = null;
    lines.forEach((line) => {
      if (line.startsWith("Problem Name:")) {
        this.problemName = this.extractQuotedValue(line);
      } else if (line.startsWith("Function Name:")) {
        this.functionName = this.extractValue(line);
      } else if (line.startsWith("Input Structure:")) {
        currentSection = "input";
      } else if (line.startsWith("Output Structure:")) {
        currentSection = "output";
      } else if (line.startsWith("Input Field:")) {
        if (currentSection === "input") {
          const field = this.extractField(line);
          if (field) this.inputFields.push(field);
        }
      } else if (line.startsWith("Output Field:")) {
        if (currentSection === "output") {
          const field = this.extractField(line);
          if (field) this.outputFields.push(field);
        }
      }
    });
  }
  extractQuotedValue(line) {
    const match = line.match(/: "(.*)"$/);
    return match ? match[1] : "";
  }
  extractValue(line) {
    const match = line.match(/: (\w+)$/);
    return match ? match[1] : "";
  }
  extractField(line) {
    const match = line.match(/Field: (\w+(?:<\w+>)?) (\w+)$/);
    return match ? { type: match[1], name: match[2] } : null;
  }
  generateCpp() {
    const inputs = this.inputFields.map((field) => `${this.mapTypeToCpp(field.type)} ${field.name}`).join(", ");
    return `${this.mapTypeToCpp(this.outputFields[0].type)} ${this.functionName}(${inputs}) {
    // Implementation goes here
    return result;
}`;
  }
  generateJs() {
    const inputs = this.inputFields.map((field) => field.name).join(", ");
    return `function ${this.functionName}(${inputs}) {
    // Implementation goes here
    return result;
}`;
  }
  generateRust() {
    const inputs = this.inputFields.map((field) => `${field.name}: ${this.mapTypeToRust(field.type)}`).join(", ");
    const outputType = this.mapTypeToRust(this.outputFields[0].type);
    return `fn ${this.functionName}(${inputs}) -> ${outputType} {
    // Implementation goes here
    result
}`;
  }
  generateJava() {
    const inputs = this.inputFields.map((field) => `${this.mapTypeToJava(field.type)} ${field.name}`).join(", ");
    return `public static ${this.mapTypeToJava(this.outputFields[0].type)} ${this.functionName}(${inputs}) {
    // Implementation goes here
    return result;
}`;
  }
  mapTypeToRust(type) {
    switch (type) {
      case "int":
        return "i32";
      case "float":
        return "f64";
      case "string":
        return "String";
      case "bool":
        return "bool";
      case "list<int>":
        return "Vec<i32>";
      case "list<float>":
        return "Vec<f64>";
      case "list<string>":
        return "Vec<String>";
      case "list<bool>":
        return "Vec<bool>";
      default:
        return "unknown";
    }
  }
  mapTypeToCpp(type) {
    switch (type) {
      case "int":
        return "int";
      case "float":
        return "float";
      case "string":
        return "std::string";
      case "bool":
        return "bool";
      case "list<int>":
        return "std::vector<int>";
      case "list<float>":
        return "std::vector<float>";
      case "list<string>":
        return "std::vector<std::string>";
      case "list<bool>":
        return "std::vector<bool>";
      default:
        return "unknown";
    }
  }
  mapTypeToJava(type) {
    switch (type) {
      case "int":
        return "int";
      case "float":
        return "float";
      case "string":
        return "String";
      case "bool":
        return "boolean";
      case "list<int>":
        return "List<Integer>";
      case "list<float>":
        return "List<Float>";
      case "list<string>":
        return "List<String>";
      case "list<bool>":
        return "List<Boolean>";
      default:
        return "unknown";
    }
  }
};

// src/FullProblemDefinitionGenerator.ts
var FullProblemDefinitionParser = class {
  constructor() {
    this.problemName = "";
    this.functionName = "";
    this.inputFields = [];
    this.outputFields = [];
  }
  parse(input) {
    const lines = input.split("\n").map((line) => line.trim());
    let currentSection = null;
    lines.forEach((line) => {
      if (line.startsWith("Problem Name:")) {
        this.problemName = this.extractQuotedValue(line);
      } else if (line.startsWith("Function Name:")) {
        this.functionName = this.extractValue(line);
      } else if (line.startsWith("Input Structure:")) {
        currentSection = "input";
      } else if (line.startsWith("Output Structure:")) {
        currentSection = "output";
      } else if (line.startsWith("Input Field:")) {
        if (currentSection === "input") {
          const field = this.extractField(line);
          if (field) this.inputFields.push(field);
        }
      } else if (line.startsWith("Output Field:")) {
        if (currentSection === "output") {
          const field = this.extractField(line);
          if (field) this.outputFields.push(field);
        }
      }
    });
  }
  extractQuotedValue(line) {
    const match = line.match(/: "(.*)"$/);
    return match ? match[1] : "";
  }
  extractValue(line) {
    const match = line.match(/: (\w+)$/);
    return match ? match[1] : "";
  }
  extractField(line) {
    const match = line.match(/Field: (\w+(?:<\w+>)?) (\w+)$/);
    return match ? { type: match[1], name: match[2] } : null;
  }
  generateCpp() {
    const inputs = this.inputFields.map((field) => `${this.mapTypeToCpp(field.type)} ${field.name}`).join(", ");
    const inputReads = this.inputFields.map((field, index) => {
      if (field.type.startsWith("list<")) {
        return `int size_${field.name};
  std::istringstream(lines[${index}]) >> size_${field.name};
  ${this.mapTypeToCpp(field.type)} ${field.name}(size_${field.name});
  if(!size_${field.name}==0) {
  	std::istringstream iss(lines[${index + 1}]);
  	for (int i=0; i < size_arr; i++) iss >> arr[i];
  }`;
      } else {
        return `${this.mapTypeToCpp(field.type)} ${field.name};
  std::istringstream(lines[${index}]) >> ${field.name};`;
      }
    }).join("\n  ");
    const outputType = this.outputFields[0].type;
    const functionCall = `${outputType} result = ${this.functionName}(${this.inputFields.map((field) => field.name).join(", ")});`;
    const outputWrite = `std::cout << result << std::endl;`;
    return `#include <iostream>
  #include <fstream>
  #include <vector>
  #include <string>
  #include <sstream>
  #include <climits>
  
  ##USER_CODE_HERE##
  
  int main() {
    std::ifstream file("/dev/problems/${this.problemName.toLowerCase().replace(" ", "-")}/tests/inputs/##INPUT_FILE_INDEX##.txt");
    std::vector<std::string> lines;
    std::string line;
    while (std::getline(file, line)) lines.push_back(line);
  
    file.close();
    ${inputReads}
    ${functionCall}
    ${outputWrite}
    return 0;
  }
  `;
  }
  generateJava() {
    let inputReadIndex = 0;
    const inputReads = this.inputFields.map((field, index) => {
      if (field.type.startsWith("list<")) {
        let javaType = this.mapTypeToJava(field.type);
        let inputType = javaType.match(/<(.*?)>/);
        javaType = inputType ? inputType[1] : "Integer";
        let parseToType = javaType === "Integer" ? "Int" : javaType;
        return `int size_${field.name} = Integer.parseInt(lines.get(${inputReadIndex++}).trim());

          ${this.mapTypeToJava(field.type)} ${field.name} = new ArrayList<>(size_${field.name});

          String[] inputStream = lines.get(${inputReadIndex++}).trim().split("\\s+");

          for (String inputChar : inputStream)  {

            ${field.name}.add(${javaType}.parse${parseToType}(inputChar));

          }
`;
      } else {
        let javaType = this.mapTypeToJava(field.type);
        if (javaType === "int") {
          javaType = "Integer";
        } else if (javaType === "float") {
          javaType = "Float";
        } else if (javaType === "boolean") {
          javaType = "Boolean";
        } else if (javaType === "String") {
          javaType = "String";
        }
        let parseToType = javaType === "Integer" ? "Int" : javaType;
        return `${this.mapTypeToJava(field.type)} ${field.name} = ${javaType}.parse${parseToType}(lines.get(${inputReadIndex++}).trim());`;
      }
    }).join("\n  ");
    const outputType = this.mapTypeToJava(this.outputFields[0].type);
    const functionCall = `${outputType} result = ${this.functionName}(${this.inputFields.map((field) => field.name).join(", ")});`;
    const outputWrite = `System.out.println(result);`;
    return `
  import java.io.*;
  import java.util.*;
  
  public class Main {
      
      ##USER_CODE_HERE##
  
      public static void main(String[] args) {
          String filePath = "/dev/problems/${this.problemName.toLowerCase().replace(" ", "-")}/tests/inputs/##INPUT_FILE_INDEX##.txt"; 
          List<String> lines = readLinesFromFile(filePath);
          ${inputReads}
          ${functionCall}
          ${outputWrite}
      }
      public static List<String> readLinesFromFile(String filePath) {
          List<String> lines = new ArrayList<>();
          try (BufferedReader br = new BufferedReader(new FileReader(filePath))) {
              String line;
              while ((line = br.readLine()) != null) {
                  lines.add(line);
              }
          } catch (IOException e) {
              e.printStackTrace();
          }
          return lines;
      }
  }`;
  }
  generateJs() {
    const inputs = this.inputFields.map((field) => field.name).join(", ");
    const inputReads = this.inputFields.map((field) => {
      if (field.type.startsWith("list<")) {
        return `const size_${field.name} = parseInt(input.shift());
const ${field.name} = input.splice(0, size_${field.name}).map(Number);`;
      } else {
        return `const ${field.name} = parseInt(input.shift());`;
      }
    }).join("\n  ");
    const outputType = this.outputFields[0].type;
    const functionCall = `const result = ${this.functionName}(${this.inputFields.map((field) => field.name).join(", ")});`;
    const outputWrite = `console.log(result);`;
    return `##USER_CODE_HERE##
  
  const input = require('fs').readFileSync('/dev/problems/${this.problemName.toLowerCase().replace(
      " ",
      "-"
    )}/tests/inputs/##INPUT_FILE_INDEX##.txt', 'utf8').trim().split('\\n').join(' ').split(' ');
  ${inputReads}
  ${functionCall}
  ${outputWrite}
      `;
  }
  generateRust() {
    const inputs = this.inputFields.map((field) => `${field.name}: ${this.mapTypeToRust(field.type)}`).join(", ");
    const inputReads = this.inputFields.map((field) => {
      if (field.type.startsWith("list<")) {
        return `let size_${field.name}: usize = lines.next().and_then(|line| line.parse().ok()).unwrap_or(0);
	let ${field.name}: ${this.mapTypeToRust(field.type)} = parse_input(lines, size_${field.name});`;
      } else {
        return `let ${field.name}: ${this.mapTypeToRust(
          field.type
        )} = lines.next().unwrap().parse().unwrap();`;
      }
    }).join("\n  ");
    const containsVector = this.inputFields.find(
      (field) => field.type.startsWith("list<")
    );
    const outputType = this.mapTypeToRust(this.outputFields[0].type);
    const functionCall = `let result = ${this.functionName}(${this.inputFields.map((field) => field.name).join(", ")});`;
    const outputWrite = `println!("{}", result);`;
    return `use std::fs::read_to_string;
  use std::io::{self};
  use std::str::Lines;
  
  ##USER_CODE_HERE##
  
  fn main() -> io::Result<()> {
    let input = read_to_string("/dev/problems/${this.problemName.toLowerCase().replace(" ", "-")}/tests/inputs/##INPUT_FILE_INDEX##.txt")?;
    let mut lines = input.lines();
    ${inputReads}
    ${functionCall}
    ${outputWrite}
    Ok(())
  }${containsVector ? `
fn parse_input(mut input: Lines, size_arr: usize) -> Vec<i32> {
      let arr: Vec<i32> = input
          .next()
          .unwrap_or_default()
          .split_whitespace()
          .filter_map(|x| x.parse().ok())
          .collect();
  
      if size_arr == 0 {
          Vec::new()
      } else {
          arr
      }
  }` : ""}
  `;
  }
  mapTypeToCpp(type) {
    switch (type) {
      case "int":
        return "int";
      case "float":
        return "float";
      case "string":
        return "std::string";
      case "bool":
        return "bool";
      case "list<int>":
        return "std::vector<int>";
      case "list<float>":
        return "std::vector<float>";
      case "list<string>":
        return "std::vector<std::string>";
      case "list<bool>":
        return "std::vector<bool>";
      default:
        return "unknown";
    }
  }
  mapTypeToRust(type) {
    switch (type) {
      case "int":
        return "i32";
      case "float":
        return "f64";
      case "string":
        return "String";
      case "bool":
        return "bool";
      case "list<int>":
        return "Vec<i32>";
      case "list<float>":
        return "Vec<f64>";
      case "list<string>":
        return "Vec<String>";
      case "list<bool>":
        return "Vec<bool>";
      default:
        return "unknown";
    }
  }
  mapTypeToJava(type) {
    switch (type) {
      case "int":
        return "int";
      case "float":
        return "float";
      case "string":
        return "String";
      case "bool":
        return "boolean";
      case "list<int>":
        return "List<Integer>";
      case "list<float>":
        return "List<Float>";
      case "list<string>":
        return "List<String>";
      case "list<bool>":
        return "List<Boolean>";
      default:
        return "unknown";
    }
  }
};

// src/index.ts
function generatePartialBoilerplate(generatorFilePath) {
  const inputFilePath = import_path.default.join(__dirname, generatorFilePath, "Structure.md");
  const boilerplatePath = import_path.default.join(
    __dirname,
    generatorFilePath,
    "boilerplate"
  );
  const input = import_fs.default.readFileSync(inputFilePath, "utf-8");
  const parser = new ProblemDefinitionParser();
  parser.parse(input);
  const cppCode = parser.generateCpp();
  const jsCode = parser.generateJs();
  const rustCode = parser.generateRust();
  if (!import_fs.default.existsSync(boilerplatePath)) {
    import_fs.default.mkdirSync(boilerplatePath, { recursive: true });
  }
  import_fs.default.writeFileSync(import_path.default.join(boilerplatePath, "function.cpp"), cppCode);
  import_fs.default.writeFileSync(import_path.default.join(boilerplatePath, "function.js"), jsCode);
  import_fs.default.writeFileSync(import_path.default.join(boilerplatePath, "function.rs"), rustCode);
  console.log("Boilerplate code generated successfully!");
}
function generateFullBoilerPLate(generatorFilePath) {
  const inputFilePath = import_path.default.join(__dirname, generatorFilePath, "Structure.md");
  const boilerplatePath = import_path.default.join(
    __dirname,
    generatorFilePath,
    "boilerplate-full"
  );
  const input = import_fs.default.readFileSync(inputFilePath, "utf-8");
  const parser = new FullProblemDefinitionParser();
  parser.parse(input);
  const cppCode = parser.generateCpp();
  const jsCode = parser.generateJs();
  const rustCode = parser.generateRust();
  if (!import_fs.default.existsSync(boilerplatePath)) {
    import_fs.default.mkdirSync(boilerplatePath, { recursive: true });
  }
  import_fs.default.writeFileSync(import_path.default.join(boilerplatePath, "function.cpp"), cppCode);
  import_fs.default.writeFileSync(import_path.default.join(boilerplatePath, "function.js"), jsCode);
  import_fs.default.writeFileSync(import_path.default.join(boilerplatePath, "function.rs"), rustCode);
  console.log("Boilerplate code generated successfully!");
}
generatePartialBoilerplate(process.env.GENERATOR_FILE_PATH ?? "");
generateFullBoilerPLate(process.env.GENERATOR_FILE_PATH ?? "");
