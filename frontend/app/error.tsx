"use client";
import { Button } from "@/components/ui/button";
import { RefreshCcw, Home } from "lucide-react";

export default function Error({ reset }: { reset: () => void }) {
  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-gradient-to-br from-[#FFCCEA] via-white to-[#BFECFF] p-8">
      <div className="max-w-2xl w-full space-y-12 text-center">
        {/* 謝罪イラスト */}
        <div className="w-48 h-48 mx-auto relative">
          <div className="absolute inset-0 flex items-center justify-center">
            <svg
              className="w-full h-full"
              viewBox="0 0 100 100"
              fill="none"
              xmlns="http://www.w3.org/2000/svg"
            >
              {/* 申し訳なさそうな顔文字 */}
              <circle
                cx="50"
                cy="50"
                r="45"
                stroke="url(#errorGradient)"
                strokeWidth="2"
                fill="white"
              />
              <path
                d="M35 65 C 35 65, 50 75, 65 65"
                stroke="url(#errorGradient)"
                strokeWidth="2"
                fill="none"
              />
              <circle cx="35" cy="45" r="5" fill="url(#errorGradient)" />
              <circle cx="65" cy="45" r="5" fill="url(#errorGradient)" />
              <defs>
                <linearGradient
                  id="errorGradient"
                  x1="0"
                  y1="0"
                  x2="100"
                  y2="100"
                >
                  <stop offset="0%" stopColor="#FF69B4" />
                  <stop offset="100%" stopColor="#0080FF" />
                </linearGradient>
              </defs>
            </svg>
          </div>
        </div>

        {/* エラーメッセージ */}
        <div className="space-y-6">
          <h1 className="text-3xl md:text-4xl font-bold bg-gradient-to-r from-pink-500 to-blue-500 bg-clip-text text-transparent">
            申し訳ありません
          </h1>
          <p className="text-lg md:text-xl text-gray-600 leading-relaxed max-w-lg mx-auto">
            予期せぬエラーが発生しました。
            <br />
            しばらく時間をおいて、再度お試しください。
          </p>
        </div>

        {/* アクションボタン */}
        <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
          <Button
            onClick={reset}
            size="lg"
            className="bg-gradient-to-r from-pink-500 to-blue-500 text-white hover:opacity-90 transition-all duration-200 px-8 py-6 text-lg rounded-full shadow-lg w-full sm:w-auto"
          >
            <RefreshCcw className="mr-2 h-5 w-5" />
            もう一度試す
          </Button>
          <Button
            onClick={() => (window.location.href = "/")}
            size="lg"
            variant="outline"
            className="border-2 text-gray-800 hover:bg-white/50 transition-all duration-200 px-8 py-6 text-lg rounded-full w-full sm:w-auto"
          >
            <Home className="mr-2 h-5 w-5" />
            ホームに戻る
          </Button>
        </div>
      </div>
    </div>
  );
}
