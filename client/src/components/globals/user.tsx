import React, { useEffect } from "react";
import { Avatar, AvatarFallback, AvatarImage } from "../ui/avatar";
import { useGetProfileDetailsQuery } from "@/hooks/use-auth-query";
import Link from "next/link";
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
import { deleteCookie, getCookie } from "cookies-next";

const User = () => {
  const router = useRouter();
  const { data, isPending, isError } = useGetProfileDetailsQuery();

  if (isPending) return <></>;
  if (isError) return <Login />;

  if (!data)
    return (
      <Link
        href="/auth"
        className="border py-1 hover:border hover:border-primary flex items-center justify-center rounded-full px-5"
      >
        Login
      </Link>
    );
  if (data && getCookie("token"))
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
              deleteCookie("token");
              router.refresh();
            }}
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
    deleteCookie("token");
  }, []);
  return (
    <Link
      href="/auth"
      className="border py-1 hover:border hover:border-primary flex items-center justify-center rounded-full px-5"
    >
      Login
    </Link>
  );
};
