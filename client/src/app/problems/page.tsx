import ProblemsList from "@/components/problems/list";
import React, { Suspense } from "react";

const ProblemsPage = () => {
  return (
    <div>
      <Suspense>
      <ProblemsList />
      </Suspense>

    </div>
  );
};

export default ProblemsPage;
