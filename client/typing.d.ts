/* eslint-disable no-unused-vars */
interface TLogin {
  email: string;
  password: string;
}

interface TRegister {
  username: string;
  email: string;
  password: string;
}

interface FunctionDetails {
  id: number;
  createdAt: string;
  updatedAt: string;
  problemId: number;
  functionName: string;
  parameters: {
    id: number;
    signatureId: number;
    name: string;
    type: string;
  }[];
  returnType: string;
}

interface TProblem {
  id: number;
  createdAt: string;
  updatedAt: string;
  title: string;
  slug: string;
  description: string;
  boilerplate: Boilerplate;

  difficulty: string;
  function: FunctionDetails;
  testCases: TestCase[] | null; // Adjust type as needed if test cases have a different structure
}

interface TestCase {
  id: number;
  input: string;
  output: string;
}

interface Boilerplate {
  [key: string]: string;
}

interface TProblemData {
  id: number;
  createdAt: string;
  updatedAt: string;
  title: string;
  slug: string;
  description: string;
  boilerplate: Boilerplate[];
  difficulty: string;
  function: FunctionDetails;
  testCases: TestCase[];
}
