import React from "react";
import { Hero } from "@/components/home/Hero";
import { Sponsors } from "@/components/home/Sponsors";
import { About } from "@/components/home/About";
import { Features } from "@/components/home/Features";
import { HowItWorks } from "@/components/home/HowItWorks";
import { Services } from "@/components/home/Services";
import { Cta } from "@/components/home/Cta";

import { FAQ } from "@/components/home/FAQ";
import { ScrollToTop } from "@/components/home/ScrollToTop";


export default function Home() {
  return (
    <>
      <Hero />
      <Sponsors />
      <About />
      <HowItWorks />
      <Features />
      <Services />
      <Cta />

      <FAQ />

      <ScrollToTop />
    </>
  );
}

