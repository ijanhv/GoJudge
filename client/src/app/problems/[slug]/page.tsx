import ProblemSolving from "@/components/problems/solving";
import React from "react";

const ProblemPage = ({
  params,
}: {
  params: {
    slug: string;
  };
}) => {
  return <ProblemSolving slug={params.slug} />
};

export default ProblemPage;
