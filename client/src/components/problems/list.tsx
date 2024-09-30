"use client";
import { useGetAllProblemsQuery } from "@/hooks/use-problems-query";
import React from "react";
import Loader from "../common/loader";
import Error from "../common/error";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "../ui/card";
import { Badge } from "../ui/badge";
import { ChevronRight, Code2, FileText } from "lucide-react";
import Link from "next/link";

const ProblemsList = () => {
  const { data, isPending, isError } = useGetAllProblemsQuery();

  if (isPending) return <Loader />;
  if (isError) return <Error />;
  return <div className="space-y-6">
    <h1 className="text-3xl font-bold text-center mb-8">Coding Problems</h1>
    {data.map((problem) => (
      <Card key={problem.id} className="hover:shadow-md transition-shadow">
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-2xl font-bold">{problem.title}</CardTitle>
          <Badge variant={problem.difficulty === 'Easy' ? 'secondary' : 'destructive'}>
            {problem.difficulty}
          </Badge>
        </CardHeader>
        <CardContent>
          <CardDescription className="mt-2">{problem.description}</CardDescription>
          <div className="flex justify-between items-center mt-4">
            <div className="flex space-x-2">
              <Badge variant="outline" className="flex items-center space-x-1">
                <Code2 className="w-3 h-3" />
                <span>Code</span>
              </Badge>
              <Badge variant="outline" className="flex items-center space-x-1">
                <FileText className="w-3 h-3" />
                <span>Description</span>
              </Badge>
            </div>
            <Link
              href={`/problems/${problem.slug}`}
              className="text-sm text-blue-500 hover:text-blue-700 flex items-center"
            >
              Solve Challenge
              <ChevronRight className="w-4 h-4 ml-1" />
            </Link>
          </div>
        </CardContent>
      </Card>
    ))}
  </div>

};

export default ProblemsList;
