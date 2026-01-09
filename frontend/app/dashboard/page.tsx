import MainLayout from '@/components/layout/MainLayout'
import QuickAction from './_components/QuickAction'
import RecentVideo from './_components/RecentVideo'

export const metadata = {
  title: 'ダッシュボード | Tavinikkiy',
  description: 'あなたの旅行動画を管理しましょう',
}

export default function DashboardPage() {
  return (
    <MainLayout description="あなたの旅行動画を管理しましょう">
      <QuickAction />
      <RecentVideo />
    </MainLayout>
  )
}
