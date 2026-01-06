import MainLayout from '@/components/layout/MainLayout'
import ShortsViewer from './_components/ShortsViewer'

export const metadata = {
  title: 'リール | Tavinikkiy',
  description: 'みんなの旅行動画をショート形式で楽しもう',
}

export default function AnyOnePage() {
  return (
    <MainLayout
      title="リール"
      description="みんなの旅行動画を楽しもう"
      className="p-0 max-w-full"
      showTitle={false}
      showFooter={false}
      bgClassName="bg-white"
      disableBodyScroll={true}
    >
      <ShortsViewer />
    </MainLayout>
  )
}
