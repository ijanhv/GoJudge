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


// {
//   "problemId": 31,
//   "userId": 8,
//   "language": "cpp",
//   "code": "int sum(int a, int b) {\n    return a + b;\n    \n}"

// }


interface TCodeSubmission {
  problemId: number,
  language: string;
  code: string
}

interface TSubmission {
  id: number;
  createdAt: string;
  updatedAt: string;
  problemId: number;
  problem: Problem;
  userId: number;
  testResults: TestResult[];
  submissionTime: string;
  status: 'pending' | 'success' | 'failure';
  errorMessage: string;
  language: string;
  code: string;
}


interface TestResult {
  id: number;
  createdAt: string;
  updatedAt: string;
  submissionId: number;
  testCaseId: number;
  status: 'success' | 'failure';
  output: string;
  errorMessage: string;
}