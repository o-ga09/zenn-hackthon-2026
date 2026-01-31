'use client'
import React, { useState } from 'react'
import { Button } from '../ui/button'
import { Upload, MapPin, Zap, Sparkles } from 'lucide-react'
import Link from 'next/link'
import { motion } from 'framer-motion'
import Image from 'next/image'

export default function Hero() {
  const [destination, setDestination] = useState('')

  return (
    <section className="max-w-[1200px] mx-auto px-6 pt-16 pb-24">
      <div className="grid lg:grid-cols-2 gap-12 items-center">
        {/* 左側: テキストコンテンツ */}
        <motion.div
          className="flex flex-col gap-8"
          initial={{ opacity: 0, x: -30 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ duration: 0.8 }}
        >
          <motion.div
            className="inline-flex items-center gap-2 px-3 py-1 bg-secondary/15 text-secondary rounded-full w-fit"
            initial={{ opacity: 0, scale: 0.8 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ delay: 0.2 }}
          >
            <Sparkles className="w-4 h-4" />
            <span className="text-xs font-bold uppercase tracking-wider">AI-Powered Magic</span>
          </motion.div>

          <motion.h2
            className="text-5xl lg:text-6xl font-black leading-[1.15] tracking-tight"
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.3, duration: 0.8 }}
          >
            あなたの旅を、
            <br />
            <span className="text-primary italic">一生の思い出に。</span>
          </motion.h2>

          <motion.p
            className="text-lg text-muted-foreground max-w-[500px] leading-relaxed"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 0.5, duration: 0.8 }}
          >
            思い出の写真をアップロードするだけ。AIが最高のシーンを選び、音楽と合わせて15秒のVlogを瞬時に作成します。
          </motion.p>

          <motion.div
            className="flex flex-wrap gap-4 items-center"
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.7, duration: 0.6 }}
          >
            <div className="flex -space-x-3 overflow-hidden">
              <Image
                className="inline-block h-10 w-10 rounded-full ring-2 ring-white"
                src="https://lh3.googleusercontent.com/aida-public/AB6AXuCatuVps2kfYxPOV1f-V2cTcG3hJ4u5x2rOSigJWbfXGhKoSYcnGW_sVN4R5d5Z3OlIhMhbOxSX0hplu-GweuGForXWPkANtC3DPDtKOsQ9bSLS8ML3FT5CFU2RITYt_pNPr9DwkI8VFRFmS5LeWoRVQLLdXQtb-70eNq-wxt9cUh8pgtDtHt7S3Sc7BUz__govjQ6gHEFb9WU2NqGsF-nA3jKvbFGFVJ_nPNaD4hbN_4kTmdOulqaeL4x9WR48FmvEbJtu3RwPAgc"
                alt="User 1"
                width={40}
                height={40}
              />
              <Image
                className="inline-block h-10 w-10 rounded-full ring-2 ring-white"
                src="https://lh3.googleusercontent.com/aida-public/AB6AXuAj7gT9ey7kXnFKOYzQyjTcoeJ9eUzVPXVOKhMM3nfPLyYMmHnIWCiWR25Yrti5cV7SavbZGkwcQkRV-RI8McOtdJ6JlnLWKvsTxfrk4XCwptgeKTat2KttCOC4xUNHwTUdoNyd8mplsgnaKJnEZmyC13OAV5fJl_vmUn6PBTcDmPu-r0DMGlurqqjwGi75Gqv_dUHHQLyCHkR1BW0PlVmI4gn8p7X9eIhHYvllKLNaVEYySOu-_6mf6XRzU2K6oHrjiC4sluKDN0I"
                alt="User 2"
                width={40}
                height={40}
              />
              <Image
                className="inline-block h-10 w-10 rounded-full ring-2 ring-white"
                src="https://lh3.googleusercontent.com/aida-public/AB6AXuAmWUnDqW1w1X0vWWd3SOunMW4NG8J7mVSb8VpD-rCN5JKeIvPXq1C6oVx1FZWxqOM1DO-AQb_ZofHzZe-IN2XCrP8MZi_h-M-JqZoDdwzNjQJHw2zkzC6cbugAkKXq7UbLlSA4NCyO4pWipvapctKi3onhN0bbOvZZ2iHAti9lpHHZYICt_h8Tr6jr8q2vBEsU0Iz0dUkLdPXDk-Ypm4QTmARWtLMqNo5sL6eVQ7_givWpjdt6eMqef16LDXNRxOH1NCZPdZleX18"
                alt="User 3"
                width={40}
                height={40}
              />
            </div>
            <div className="flex flex-col">
              <span className="text-sm font-bold">10,000+ creators</span>
              <span className="text-xs text-muted-foreground">making memories every day</span>
            </div>
          </motion.div>
        </motion.div>

        {/* 右側: Trial Generation カード */}
        <motion.div
          className="bg-white dark:bg-gray-900 rounded-xl p-8 border border-primary/10 kawaii-shadow relative overflow-hidden"
          initial={{ opacity: 0, x: 30 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ duration: 0.8, delay: 0.2 }}
        >
          <div className="absolute top-0 right-0 p-4 opacity-10">
            <svg
              className="w-16 h-16 text-primary"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z"
              />
            </svg>
          </div>

          <h3 className="text-xl font-bold mb-6 flex items-center gap-2">
            <Sparkles className="w-5 h-5 text-primary" />
            Trial Generation
          </h3>

          <div className="space-y-6">
            {/* アップロードエリア */}
            <Link href="/upload">
              <motion.div
                className="group cursor-pointer border-2 border-dashed border-gray-200 dark:border-gray-700 rounded-lg p-10 flex flex-col items-center justify-center gap-4 hover:border-primary/50 hover:bg-primary/5 transition-all"
                whileHover={{ scale: 1.02 }}
                whileTap={{ scale: 0.98 }}
              >
                <div className="size-16 bg-accent rounded-full flex items-center justify-center text-primary">
                  <Upload className="w-8 h-8" />
                </div>
                <div className="text-center">
                  <p className="font-bold text-sm">3〜5枚の写真を選択</p>
                  <p className="text-xs text-muted-foreground mt-1">
                    または短い動画クリップをドロップ
                  </p>
                </div>
              </motion.div>
            </Link>

            {/* Destination入力 */}
            <div className="space-y-2">
              <label className="text-xs font-bold uppercase text-muted-foreground ml-1">
                Destination
              </label>
              <div className="relative">
                <MapPin className="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground w-5 h-5" />
                <input
                  type="text"
                  value={destination}
                  onChange={e => setDestination(e.target.value)}
                  className="w-full pl-10 pr-4 py-3 bg-muted dark:bg-gray-800 border-none rounded-lg focus:ring-2 focus:ring-primary transition-all"
                  placeholder="例: 京都, パリ, 沖縄..."
                />
              </div>
            </div>

            {/* 生成ボタン */}
            <Link href="/upload" className="block">
              <Button className="w-full bg-primary hover:bg-primary/90 text-white font-black py-4 rounded-lg flex items-center justify-center gap-2 group transition-all h-14">
                <span>15秒のプレビューを生成</span>
                <Zap className="w-5 h-5 group-hover:translate-x-1 transition-transform" />
              </Button>
            </Link>

            <p className="text-center text-[10px] text-muted-foreground">
              ログイン不要・無料で今すぐお試しいただけます
            </p>
          </div>
        </motion.div>
      </div>
    </section>
  )
}
