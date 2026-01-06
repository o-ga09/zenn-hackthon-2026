import MainLayout from '@/components/layout/MainLayout'
import StepperWrapper from './_components/StepperWrapper'

export const metadata = {
  title: 'アップロード | Tavinikkiy',
  description: '写真と旅行情報をアップロードして、AIが素敵な動画を作成します',
}

export default function UploadPage() {
  return (
    <MainLayout
      title="動画生成 - アップロード"
      description="写真と旅行情報をアップロードして、AIが素敵な動画を作成します"
    >
      <StepperWrapper />
    </MainLayout>
  )
}
