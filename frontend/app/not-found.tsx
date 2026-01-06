import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Home } from "lucide-react";

export default function NotFound() {
  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-gradient-to-br from-[#FFCCEA] via-white to-[#BFECFF] p-8">
      <div className="max-w-2xl w-full space-y-12 text-center">
        {/* エラーコード */}
        <h1 className="text-7xl md:text-8xl font-bold bg-gradient-to-r from-pink-500 to-blue-500 bg-clip-text text-transparent">
          404
        </h1>

        {/* メインメッセージ */}
        <div className="space-y-6">
          <h2 className="text-2xl md:text-3xl font-semibold text-gray-800">
            ページが見つかりません
          </h2>
          <p className="text-lg md:text-xl text-gray-600 leading-relaxed max-w-lg mx-auto">
            お探しのページは存在しないか、移動した可能性があります。
            <br />
            URLをご確認いただくか、以下のボタンからホームページへお戻りください。
          </p>
        </div>

        {/* アクション */}
        <div className="pt-4">
          <Link href="/">
            <Button
              size="lg"
              className="bg-gradient-to-r from-pink-500 to-blue-500 text-white hover:opacity-90 transition-all duration-200 px-8 py-6 text-lg rounded-full shadow-lg"
            >
              <Home className="mr-2 h-5 w-5" />
              ホームに戻る
            </Button>
          </Link>
        </div>
      </div>
    </div>
  );
}
