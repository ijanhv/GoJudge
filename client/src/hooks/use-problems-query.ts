import { apiUrl } from "@/utils/url";
import { useQuery, UseQueryResult } from "@tanstack/react-query";
import axios from "axios";

const getAllProblems = async () => {
  const res = await axios.get(`${apiUrl}/api/problems`);
  return res.data.problems;
};

export const useGetAllProblemsQuery = (): UseQueryResult<TProblem[]> => {
  return useQuery({
    queryKey: ["problems"],
    queryFn: () => getAllProblems(),
    staleTime: Infinity,
  });
};

const getProblem = async (slug: string) => {
  const res = await axios.get(`${apiUrl}/api/problems/${slug}`);
  const problem = {
    ...res.data.problem,
   boilerplate:  res.data.boilerplate
  }
  return problem
};

export const useGetProblemQuery = (slug:string): UseQueryResult<TProblem> => {
  return useQuery({
    queryKey: [`problem_${slug}`],
    queryFn: () => getProblem(slug),
    staleTime: Infinity,
  });
};
