"use client"
import React from "react";
import dynamic from "next/dynamic";
import Link from "next/link";
import { getCookie } from "cookies-next";

const User = dynamic(() => import("@/components/globals/user"), {
  loading: () => <p>Loading...</p>,
});

const Check = () => {
  const token = getCookie("token")
  return (
    <>

    

      {token ? (

          <User />

      ) : (
        <Link
          href="/auth"
          className="border py-1 hover:border hover:border-primary flex items-center justify-center rounded-full px-5"
        >
          Login
        </Link>
      )}
    </>
  );
};

export default Check;

