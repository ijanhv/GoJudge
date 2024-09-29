import { useMutation, useQuery } from "@tanstack/react-query";
import axios from "axios";
import { useRouter } from "next/router";
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
  const router = useRouter();
  return useMutation({
    mutationFn: loginUser,
    onSuccess: (data) => {
      toast.success("Logged in successfully", {
        position: "top-center",
      });
      Cookies.set("token", data.data.access_token);
      router.push("/");
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
        Cookies.set("token", data.data.access_token);
        router.push("/");
      },
      onError: (error: any) => {
        toast.error("Error, logging in!", {
          position: "top-center",
        });
      },
    });
  };

const getAdminDetails = async () => {
  const res = await axios.get(`${apiUrl}/user/token`, {
    headers: {
      Authorization: `Bearer ${Cookies.get("token")}`,
    },
  });

  return res.data.user;
};

export const useGetAdminDetailsQuery = () => {
  return useQuery({
    queryKey: ["admin"],
    queryFn: getAdminDetails,
    staleTime: Infinity,
  });
};
