"use client";
import React, { Suspense, useState } from "react";

import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from "@/components/ui/sheet";

import dynamic from "next/dynamic";
import { buttonVariants } from "../ui/button";
import { Menu } from "lucide-react";
import { LogoIcon } from "../home/Icons";
import Link from "next/link";
import ThemeToggle from "../mode-toggle";

const Check = dynamic(() => import("@/components/globals/check"), {
  loading: () => <p>Loading...</p>,
});

interface RouteProps {
  href: string;
  label: string;
}

const routeList: RouteProps[] = [
  {
    href: "#features",
    label: "Features",
  },
  {
    href: "#testimonials",
    label: "Testimonials",
  },
  {
    href: "#pricing",
    label: "Pricing",
  },
  {
    href: "#faq",
    label: "FAQ",
  },
];

export const Navbar = () => {
  const [isOpen, setIsOpen] = useState<boolean>(false);
  return (
    <div className="sticky border-b-[1px] top-0 z-40 w-full bg-white dark:border-b-slate-700 dark:bg-background">
      <div className="max-w-7xl mx-auto h-14 px-4 py-3 w-screen flex justify-between ">
        <div className="font-bold flex">
          <Link href="/" className="ml-2 font-bold text-xl flex">
            <LogoIcon />
            GoJudge
          </Link>
        </div>

        {/* mobile */}
        <div className="flex md:hidden"></div>

        {/* desktop */}
        {/* <div className="hidden md:flex gap-2">
          {routeList.map((route: RouteProps, i) => (
            <Link
              href={route.href}
              key={i}
              className={`text-[17px] ${buttonVariants({
                variant: "ghost",
              })}`}
            >
              {route.label}
            </Link>
          ))}
        </div> */}

        <div className="flex gap-2 items-center ">
          <Sheet open={isOpen} onOpenChange={setIsOpen}>
            <SheetTrigger className="px-2">
              <Menu
                className="flex md:hidden h-5 w-5"
                onClick={() => setIsOpen(true)}
              ></Menu>
            </SheetTrigger>

            <SheetContent side={"left"}>
              <SheetHeader>
                <SheetTitle className="font-bold text-xl">
                  Shadcn/React
                </SheetTitle>
              </SheetHeader>
              <div className="flex flex-col justify-center items-center gap-2 mt-4">
                {routeList.map(({ href, label }: RouteProps) => (
                  <a
                    rel="noreferrer noopener"
                    key={label}
                    href={href}
                    onClick={() => setIsOpen(false)}
                    className={buttonVariants({ variant: "ghost" })}
                  >
                    {label}
                  </a>
                ))}
              </div>
            </SheetContent>
          </Sheet>
          <Suspense fallback={<div>Loading..</div>}>
            <Check />
          </Suspense>

          <ThemeToggle />
        </div>
      </div>
    </div>
  );
};
