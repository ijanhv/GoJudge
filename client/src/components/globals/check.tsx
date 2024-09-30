import React from "react";
import Cookies from "js-cookie";
import User from "./user";
import Link from "next/link";

const Check = () => {
  return (
    <>
      {Cookies.get("token") ? (
        <User />
      ) : (
        <Link
          href="/auth"
          className="border py-1 hover:border hover:border-primary flex items-center justify-center rounded-full px-5"
        >
          Register
        </Link>
      )}
    </>
  );
};

export default Check;
