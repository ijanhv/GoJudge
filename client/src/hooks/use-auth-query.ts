import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import axios from "axios";
import { useRouter } from "next/navigation";
import { toast } from "sonner";
import Cookies from "js-cookie";
import { apiUrl } from "@/utils/url";

const loginUser = async (data: TLogin) => {
  const res = await axios.post(`${apiUrl}/api/auth/login`, {
    email: data.email,
    password: data.password,
  });
  return res.data;
};

export const useLoginUserQuery = () => {
  const queryClient = useQueryClient();
  const router = useRouter();
  return useMutation({
    mutationFn: loginUser,
    onSuccess: (data) => {
      router.refresh();
      if (data.token) {
        queryClient.invalidateQueries({ queryKey: ["profile"] });

        router.push("/problems");
        toast.success("Logged in successfully", {
          position: "top-center",
        });
        Cookies.set("token", data.token);
      }
    },
    onError: (error: any) => {
      toast.error("Error, logging in!", {
        position: "top-center",
      });
    },
  });
};

const registerUser = async (data: TRegister) => {
  const res = await axios.post(`${apiUrl}/api/auth/register`, data);
  return res.data;
};

export const useRegisterQuery = () => {
  const router = useRouter();
  return useMutation({
    mutationFn: registerUser,
    onSuccess: (data) => {
      toast.success("Registered successfully", {
        position: "top-center",
      });

      router.push("/auth");
    },
    onError: (error: any) => {
      toast.error("Error, logging in!", {
        position: "top-center",
      });
    },
  });
};

const getProfileDetails = async () => {
  const res = await axios.get(`${apiUrl}/api/user/profile`, {
    headers: {
      Authorization: `Bearer ${Cookies.get("token")}`,
    },
  });

  return res.data.user;
};

export const useGetProfileDetailsQuery = () => {
  return useQuery({
    queryKey: ["profile"],
    queryFn: getProfileDetails,
    staleTime: Infinity,
  });
};
