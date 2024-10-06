"use client";
import React, { useEffect, useState } from "react";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Button } from "@/components/ui/button";
import { ChevronLeft, PlayCircle } from "lucide-react";
import Link from "next/link";
import { useGetProblemQuery } from "@/hooks/use-problems-query";
import Loader from "../common/loader";
import Error from "../common/error";
import { useTheme } from "next-themes";
import dynamic from "next/dynamic";
import { useSearchParams } from "next/navigation";
import Submission from "./submission";
import { useSubmitCodequery } from "@/hooks/use-submission-query";

const Editor = dynamic(() => import("@monaco-editor/react"), {
  loading: () => <p>Loading...</p>,
});

const ProblemSolving = ({ slug }: { slug: string }) => {
  const searchParams = useSearchParams();

  const submisssionId = searchParams.get("id")
  const { data, isPending, isError } = useGetProblemQuery(slug);

  const { mutate, isPending: isSubmitting } = useSubmitCodequery();
  const { theme } = useTheme();

  const [selectedLanguage, setSelectedLanguage] = useState("problem.js");
  const [code, setCode] = useState<string | undefined>(
    data?.boilerplate["problem.js"] ?? ""
  );

  useEffect(() => {
    if (data) {
      setCode(data.boilerplate[selectedLanguage]);
    }
  }, [data, selectedLanguage]);

  if (isPending) return <Loader />;
  if (isError) return <Error />;

  const handleRunCode = () => {
    if (code && data) {
      mutate({
        code,
        problemId: data?.id as number,
        language: getMonacoLanguage(selectedLanguage),
      });
    }
  };

  const getMonacoLanguage = (lang: string) => {
    switch (lang) {
      case "problem.js":
        return "javascript";
      case "problem.java":
        return "java";
      case "problem.cpp":
        return "cpp";
      default:
        return "javascript";
    }
  };

  return (
    <div className="flex lg:flex-row flex-col gap-10 w-full  space-y-6">
      <div className="flex flex-col gap-3 lg:w-1/3 w-full">
        <div className="flex justify-between items-center">
          <Link
            href="/problems"
            className="flex items-center text-blue-500 hover:text-blue-700"
          >
            <ChevronLeft className="w-4 h-4 mr-1" />
            Back to Problems
          </Link>
          <Badge
            variant={data.difficulty === "Easy" ? "secondary" : "destructive"}
          >
            {data.difficulty}
          </Badge>
        </div>

        <Card>
          <CardHeader>
            <CardTitle className="text-2xl">{data.title}</CardTitle>
            <CardDescription>{data.description}</CardDescription>
          </CardHeader>
          <CardContent>
            <h3 className="font-semibold mb-2">Function Signature:</h3>
            <code className="bg-muted p-2 rounded-md block">
              {`${data.function.returnType} ${data.function.functionName}(${data?.function?.parameters?.map((p) => `${p?.type} ${p?.name}`).join(", ")})`}
            </code>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Test Cases</CardTitle>
          </CardHeader>
          <CardContent>
            {data?.testCases?.map((testCase) => (
              <div key={testCase.id} className="mb-4 p-4 border rounded-md">
                <div className="font-semibold mb-2">Input:</div>
                <code className="bg-muted p-2 rounded-md block mb-2">
                  {testCase.input}
                </code>
                <div className="font-semibold mb-2">Expected Output:</div>
                <code className="bg-muted p-2 rounded-md block">
                  {testCase.output}
                </code>
              </div>
            ))}
          </CardContent>
        </Card>
      </div>

      <div className=" lg:w-2/3 w-full">
        <Tabs
          value={selectedLanguage}
          onValueChange={(value) => {
            setSelectedLanguage(value);
            setCode(data?.boilerplate[value]);
          }}
        >
          <TabsList>
            <TabsTrigger value="problem.js">JavaScript</TabsTrigger>
            <TabsTrigger value="problem.java">Java</TabsTrigger>
            <TabsTrigger value="problem.cpp">C++</TabsTrigger>
          </TabsList>
          <TabsContent value={selectedLanguage} className="mt-4">
            <Editor
              height={"60vh"}
              value={code}
              theme={theme === "light" ? "" : "vs-dark"}
              onMount={() => {}}
              options={{
                fontSize: 14,
                scrollBeyondLastLine: false,
              }}
              language={getMonacoLanguage(selectedLanguage)}
              onChange={(value) => setCode(value)}
              defaultLanguage="javascript"
            />

            {submisssionId && <Submission id={Number(submisssionId)} />}
          </TabsContent>
        </Tabs>
        <Button
          disabled={isPending || isSubmitting}
          className="mt-4"
          onClick={() => handleRunCode()}
        >
          <PlayCircle className="w-4 h-4 mr-2" />
          Run Code
        </Button>
      </div>
    </div>
  );
};

export default ProblemSolving;
