import { Suspense } from 'react'
import StaticLayout from '@/components/layout/StaticLayout'

export default async function Layout({ children }: { children: React.ReactNode }) {
  return <StaticLayout>{children}</StaticLayout>
}
