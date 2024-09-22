import React from "react";
import { Hero } from "@/components/Hero";
import { Navbar } from "@/components/Navbar";
import { Sponsors } from "@/components/Sponsors";
import { About } from "@/components/About";
import { Features } from "@/components/Features";
import { HowItWorks } from "@/components/HowItWorks";
import { Services } from "@/components/Services";
import { Cta } from "@/components/Cta";
import { Newsletter } from "@/components/Newsletter";
import { FAQ } from "@/components/FAQ";
import { ScrollToTop } from "@/components/ScrollToTop";


export default function Home() {
  return (
    <>
      <Navbar />
      <Hero />
      <Sponsors />
      <About />
      <HowItWorks />
      <Features />
      <Services />
      <Cta />
      <Newsletter />
      <FAQ />

      <ScrollToTop />
    </>
  );
}

