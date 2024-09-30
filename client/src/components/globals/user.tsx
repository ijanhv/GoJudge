import React, { useEffect } from "react";
import { Avatar, AvatarFallback, AvatarImage } from "../ui/avatar";
import { useGetProfileDetailsQuery } from "@/hooks/use-auth-query";
import Link from "next/link";
import Cookies from "js-cookie";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { LogOutIcon } from "lucide-react";
import { useRouter } from "next/navigation";

const User = () => {
    const router = useRouter()
  const { data, isPending, isError } = useGetProfileDetailsQuery();

  if (isPending) return <></>;
  if (isError) return <Login />;
  if (data && Cookies.get("token"))
    return (
      <DropdownMenu>
        <DropdownMenuTrigger>
          <Avatar className="border cursor-pointer">
            <AvatarImage src="https://avatars1.githubusercontent.com/u/687449?v=4" />
            <AvatarFallback>CN</AvatarFallback>
          </Avatar>
        </DropdownMenuTrigger>
        <DropdownMenuContent>
          <DropdownMenuLabel>{data.username}</DropdownMenuLabel>
          <DropdownMenuSeparator />
          <DropdownMenuItem>{data.email}</DropdownMenuItem>
          <DropdownMenuItem
            onClick={() => {
                Cookies.remove("token")
                router.refresh()
            }
           
            }
            className="flex items-center gap-3"
          >
            Logout
            <LogOutIcon size={18} />
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    );
};

export default User;

const Login = () => {
  useEffect(() => {
    Cookies.remove("token");
  }, []);
  return (
    <Link
      href="/auth"
      className="border py-1 hover:border hover:border-primary flex items-center justify-center rounded-full px-5"
    >
      Register
    </Link>
  );
};
