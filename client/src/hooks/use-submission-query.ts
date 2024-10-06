import { apiUrl } from "@/utils/url";
import {
  useMutation,
  useQuery,
  useQueryClient,
  UseQueryResult,
} from "@tanstack/react-query";
import axios from "axios";
import Cookies from "js-cookie";
import { useRouter } from "next/navigation";
import { toast } from "sonner";

const submitCode = async (data: TCodeSubmission) => {
  const res = await axios.post(`${apiUrl}/api/submission`, data, {
    headers: {
      Authorization: `Bearer ${Cookies.get("token")}`,
    },
  });
  return res.data;
};

export const useSubmitCodequery = () => {
  const router = useRouter();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: submitCode,
    onSuccess: (data) => {
      queryClient.invalidateQueries({
        queryKey: [`submission_${data.submission.id}`],
      });

      toast.success("Code submitted!", {
        position: "top-center",
      });

      router.push(
        `/problems/${data.submission.problem.slug}?id=${data.submission.id}`
      );
    },
    onError: (error: any) => {
      toast.error("Error, logging in!", {
        position: "top-center",
      });
    },
  });
};

// Get current submission
const getCurrentSubmission = async (id: number) => {
  const res = await axios.get(`${apiUrl}/api/submission/${id}`, {
    headers: {
      Authorization: `Bearer ${Cookies.get("token")}`,
    },
  });
  return res.data.submission;
};

export const useGetCurrentSubmissionQuery = (
  id: number
): UseQueryResult<TSubmission> => {
  return useQuery({
    queryKey: [`submission_${id}`],
    queryFn: () => getCurrentSubmission(id),
    staleTime: Infinity,
    refetchInterval: (data) => {
      return data?.state?.data?.status === 'pending' ? 3000 : false;
    },

  });
};
