'use client'
import React from 'react'
import Link from 'next/link'

export default function LPFooter() {
  return (
    <footer className="bg-white dark:bg-gray-950 border-t border-border py-12">
      <div className="max-w-[1200px] mx-auto px-6 grid grid-cols-2 md:grid-cols-4 gap-8">
        {/* ブランドセクション */}
        <div className="col-span-2">
          <div className="flex items-center gap-2 mb-4">
            <div className="size-6 bg-primary rounded flex items-center justify-center text-white">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="14"
                height="14"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
                strokeLinejoin="round"
              >
                <path d="M12 3v3m6.366-.366-2.12 2.12M21 12h-3m.366 6.366-2.12-2.12M12 18v3m-4.246-4.246-2.12 2.12M6 12H3m4.246-4.246-2.12-2.12" />
                <circle cx="12" cy="12" r="4" />
              </svg>
            </div>
            <h2 className="text-md font-extrabold tracking-tight">Tavinikkiy Agent</h2>
          </div>
          <p className="text-xs text-muted-foreground max-w-xs">
            AIを活用した次世代のトラベルVlogエージェント。あなたの思い出を美しく編集し、記録として残します。
          </p>
          <div className="flex gap-4 mt-6">
            <a href="#" className="text-muted-foreground hover:text-primary transition-colors">
              <svg
                className="w-5 h-5"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                strokeWidth={1.5}
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="M6.827 6.175A2.31 2.31 0 015.186 7.23c-.38.054-.757.112-1.134.175C2.999 7.58 2.25 8.507 2.25 9.574V18a2.25 2.25 0 002.25 2.25h15A2.25 2.25 0 0021.75 18V9.574c0-1.067-.75-1.994-1.802-2.169a47.865 47.865 0 00-1.134-.175 2.31 2.31 0 01-1.64-1.055l-.822-1.316a2.192 2.192 0 00-1.736-1.039 48.774 48.774 0 00-5.232 0 2.192 2.192 0 00-1.736 1.039l-.821 1.316z"
                />
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="M16.5 12.75a4.5 4.5 0 11-9 0 4.5 4.5 0 019 0zM18.75 10.5h.008v.008h-.008V10.5z"
                />
              </svg>
            </a>
            <a href="#" className="text-muted-foreground hover:text-primary transition-colors">
              <svg
                className="w-5 h-5"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                strokeWidth={1.5}
              >
                <path
                  strokeLinecap="round"
                  d="M15.75 10.5l4.72-4.72a.75.75 0 011.28.53v11.38a.75.75 0 01-1.28.53l-4.72-4.72M4.5 18.75h9a2.25 2.25 0 002.25-2.25v-9a2.25 2.25 0 00-2.25-2.25h-9A2.25 2.25 0 002.25 7.5v9a2.25 2.25 0 002.25 2.25z"
                />
              </svg>
            </a>
            <a href="#" className="text-muted-foreground hover:text-primary transition-colors">
              <svg
                className="w-5 h-5"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                strokeWidth={1.5}
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="M12 21a9.004 9.004 0 008.716-6.747M12 21a9.004 9.004 0 01-8.716-6.747M12 21c2.485 0 4.5-4.03 4.5-9S14.485 3 12 3m0 18c-2.485 0-4.5-4.03-4.5-9S9.515 3 12 3m0 0a8.997 8.997 0 017.843 4.582M12 3a8.997 8.997 0 00-7.843 4.582m15.686 0A11.953 11.953 0 0112 10.5c-2.998 0-5.74-1.1-7.843-2.918m15.686 0A8.959 8.959 0 0121 12c0 .778-.099 1.533-.284 2.253m0 0A17.919 17.919 0 0112 16.5c-3.162 0-6.133-.815-8.716-2.247m0 0A9.015 9.015 0 013 12c0-1.605.42-3.113 1.157-4.418"
                />
              </svg>
            </a>
          </div>
        </div>

        {/* Serviceリンク */}
        <div>
          <h4 className="font-bold text-sm mb-4">Service</h4>
          <ul className="text-xs text-muted-foreground space-y-2">
            <li>
              <a href="#features" className="hover:text-primary transition-colors">
                Features
              </a>
            </li>
            <li>
              <a href="#" className="hover:text-primary transition-colors">
                Showcase
              </a>
            </li>
            <li>
              <a href="#pricing" className="hover:text-primary transition-colors">
                Pricing
              </a>
            </li>
          </ul>
        </div>

        {/* Companyリンク */}
        <div id="contact">
          <h4 className="font-bold text-sm mb-4">Company</h4>
          <ul className="text-xs text-muted-foreground space-y-2">
            <li>
              <a href="#" className="hover:text-primary transition-colors">
                About Us
              </a>
            </li>
            <li>
              <Link href="/terms-of-service" className="hover:text-primary transition-colors">
                Terms
              </Link>
            </li>
            <li>
              <Link href="/privacy-policy" className="hover:text-primary transition-colors">
                Privacy
              </Link>
            </li>
          </ul>
        </div>
      </div>

      {/* 下部バー */}
      <div className="max-w-[1200px] mx-auto px-6 mt-12 pt-8 border-t border-border flex justify-between items-center">
        <p className="text-[10px] text-muted-foreground">
          © 2024 Tavinikkiy Agent. All rights reserved.
        </p>
        <div className="flex gap-4">
          <button className="text-[10px] font-bold text-muted-foreground uppercase tracking-widest hover:text-foreground transition-colors">
            English
          </button>
          <button className="text-[10px] font-bold text-primary uppercase tracking-widest underline decoration-2 underline-offset-4">
            日本語
          </button>
        </div>
      </div>
    </footer>
  )
}
