import { useGetCurrentSubmissionQuery } from "@/hooks/use-submission-query";
import { CheckCircle, XCircle } from "lucide-react";
import React from "react";

const Submission = ({ id }: { id: number }) => {
  const { data, isPending, isError } = useGetCurrentSubmissionQuery(id);

  if(isPending) return (<div>Loading..</div>)
  if(isError) return (<div>Error</div>)
  return (
    <div>
      {data?.testResults?.map((testCase) => (
        <div
          key={testCase.id}
          className={`my-4 p-4 border rounded-md flex items-start justify-between ${
            testCase.status === "success"
              ? "border-green-500/40 "
              : "border-red-500/40 "
          }`}
        >
          <div className="w-4/5">
            {/* <div className="font-semibold mb-2">Input:</div> */}
            {/* <code className="bg-muted p-2 rounded-md block mb-2">
              {testCase.}
            </code> */}
            <div className="font-semibold mb-2">Expected Output:</div>
            <code className="bg-muted p-2 rounded-md block">
              {testCase.output}
            </code>
          </div>
          <div className="flex items-center justify-end">
            {testCase.status === "success" ? (
              <CheckCircle className="text-green-500 w-6 h-6" />
            ) : (
              <XCircle className="text-red-500 w-6 h-6" />
            )}
          </div>
        </div>
      ))}
    </div>
  );
};

export default Submission;
