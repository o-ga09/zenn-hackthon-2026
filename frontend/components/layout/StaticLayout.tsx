'use client'

import React from 'react'
import LPHeader from '../header/LP-Header'
import LPFooter from '../footer/LP-Footer'

export default function StaticLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className="flex flex-col h-screen">
      <LPHeader />
      <main className="flex-1 overflow-auto">{children}</main>
      <LPFooter />
    </div>
  )
}
