import React from 'react'

const Test = () => {
  return (
    <div>
    <div className="flex flex-col items-start gap-4 max-w-2xl w-full mx-auto bg-gradient-to-b from-neutral-900 to-neutral-950 p-10 rounded-3xl relative overflow-hidden">
      <div className="pointer-events-none absolute left-1/2 top-0  -ml-20 -mt-2 h-full w-full [mask-image:linear-gradient(white,transparent)]">
        <div className="absolute inset-0 bg-gradient-to-r  [mask-image:radial-gradient(farthest-side_at_top,white,transparent)] from-zinc-900/30 to-zinc-900/30 opacity-100">
          <svg
            aria-hidden="true"
            className="absolute inset-0 h-full w-full  mix-blend-overlay fill-white/10 stroke-white/10"
          >
            <defs>
              <pattern
                id=":rp:"
                width="20"
                height="20"
                patternUnits="userSpaceOnUse"
                x="-12"
                y="4"
              >
                <path d="M.5 20V.5H20" fill="none"></path>
              </pattern>
            </defs>
            <rect
              width="100%"
              height="100%"
              strokeWidth="0"
              fill="url(#:rp:)"
            ></rect>
            <svg x="-12" y="4" className="overflow-visible">
              <rect
                strokeWidth="0"
                width="21"
                height="21"
                x="180"
                y="20"
              ></rect>
              <rect
                strokeWidth="0"
                width="21"
                height="21"
                x="160"
                y="120"
              ></rect>
              <rect
                strokeWidth="0"
                width="21"
                height="21"
                x="140"
                y="40"
              ></rect>
              <rect
                strokeWidth="0"
                width="21"
                height="21"
                x="140"
                y="120"
              ></rect>
              <rect
                strokeWidth="0"
                width="21"
                height="21"
                x="140"
                y="80"
              ></rect>
            </svg>
          </svg>
        </div>
      </div>
      <div className="mb-4 w-full relative z-20">
        <label
          className="text-neutral-300 text-sm font-medium mb-2 inline-block"
          htmlFor="name"
        >
          Full name
        </label>
        <input
          id="name"
          placeholder="Manu Arora"
          className="h-10 pl-4 w-full rounded-md text-sm bg-charcoal border border-neutral-800 text-white placeholder-neutral-500 outline-none focus:outline-none active:outline-none focus:ring-2 focus:ring-neutral-800"
          type="text"
        />
      </div>
      <div className="mb-4 w-full relative z-20">
        <label
          className="text-neutral-300 text-sm font-medium mb-2 inline-block"
          htmlFor="email"
        >
          Email Address
        </label>
        <input
          id="email"
          placeholder="contact@aceternity.com"
          className="h-10 pl-4 w-full rounded-md text-sm bg-charcoal border border-neutral-800 text-white placeholder-neutral-500 outline-none focus:outline-none active:outline-none focus:ring-2 focus:ring-neutral-800"
          type="email"
        />
      </div>
      <div className="mb-4 w-full relative z-20">
        <label
          className="text-neutral-300 text-sm font-medium mb-2 inline-block"
          htmlFor="company"
        >
          Company
        </label>
        <input
          id="company"
          placeholder="contact@aceternity.com"
          className="h-10 pl-4 w-full rounded-md text-sm bg-charcoal border border-neutral-800 text-white placeholder-neutral-500 outline-none focus:outline-none active:outline-none focus:ring-2 focus:ring-neutral-800"
          type="text"
        />
      </div>
      <div className="mb-4 w-full relative z-20">
        <label
          className="text-neutral-300 text-sm font-medium mb-2 inline-block"
          htmlFor="message"
        >
          Message
        </label>
        <textarea
          id="message"
          placeholder="Type your message here"
          className="pl-4 pt-4 w-full rounded-md text-sm bg-charcoal border border-neutral-800 text-white placeholder-neutral-500 outline-none focus:outline-none active:outline-none focus:ring-2 focus:ring-neutral-800"
        ></textarea>
      </div>
      <button className="group hover:-translate-y-0.5 active:scale-[0.98] bg-neutral-800 relative z-10 hover:bg-neutral-900 border border-transparent text-white text-sm md:text-sm transition font-medium duration-200 rounded-md px-4 py-2 flex items-center justify-center shadow-[0px_1px_0px_0px_#FFFFFF20_inset]">
        Submit
      </button>
    </div>
  </div>
  )
}

export default Test