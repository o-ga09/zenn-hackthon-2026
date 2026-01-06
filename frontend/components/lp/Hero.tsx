'use client'
import React from 'react'
import { Button } from '../ui/button'
import { Badge } from '../ui/badge'
import { Play, Sparkles, Video } from 'lucide-react'
import Link from 'next/link'
import { motion } from 'framer-motion'

export default function Hero() {
  return (
    <section className="container mx-auto px-4 min-h-screen flex items-center justify-center text-center pt-10 md:pt-0 pb-10 md:pb-20 -mt-16">
      <div className="max-w-4xl mx-auto">
        <motion.div
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.6 }}
          className="mb-10 md:mb-16"
        >
          <Badge className="py-2 px-3 md:px-4 text-sm md:text-lg bg-primary text-white font-medium border-primary hover:bg-primary/90 transition-colors shadow-md">
            <Sparkles className="w-4 h-4 md:w-5 md:h-5 mr-1 md:mr-2 text-white" />
            AI搭載の新しい旅行体験
          </Badge>
        </motion.div>

        <motion.h1
          className="text-4xl sm:text-5xl md:text-7xl font-bold text-foreground mb-6 text-balance relative"
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8, delay: 0.2 }}
        >
          <motion.div
            className="absolute -top-5 sm:-top-6 md:-top-8 right-[calc(50%-80px)] sm:right-[calc(50%-100px)] md:right-[calc(50%-140px)] bg-primary text-white px-3 sm:px-5 py-1 sm:py-2 rounded-xl text-base sm:text-lg md:text-xl font-bold shadow-lg"
            initial={{ opacity: 0, scale: 0.5, rotate: -5 }}
            animate={{ opacity: 1, scale: 1, rotate: 5 }}
            transition={{ duration: 0.5, delay: 1.2, type: 'spring', bounce: 0.4 }}
          >
            <div className="absolute -bottom-2 sm:-bottom-3 right-5 h-3 sm:h-4 w-3 sm:w-4 bg-primary rotate-45"></div>
            動画で
          </motion.div>
          旅の記録を、
          <br />
          <motion.span
            className="text-primary"
            initial={{ opacity: 0, scale: 0.8 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ duration: 0.8, delay: 0.6, type: 'spring', stiffness: 200 }}
          >
            もっと簡単に
          </motion.span>
        </motion.h1>

        <motion.p
          className="text-base sm:text-lg md:text-xl text-muted-foreground mb-6 sm:mb-8 max-w-2xl mx-auto text-pretty px-2"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ duration: 1, delay: 0.8 }}
        >
          旅行の写真をアップロードするだけで、AIが自動的にエモーショナルな縦画面動画を生成。 YouTube
          ShortsやInstagramストーリーズで簡単にシェアできます。
        </motion.p>

        <motion.div
          className="flex flex-col sm:flex-row gap-4 justify-center items-center"
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8, delay: 1 }}
        >
          <motion.div
            whileHover={{ scale: 1.05 }}
            whileTap={{ scale: 0.95 }}
            className="w-full sm:w-auto"
          >
            <Button
              size="lg"
              className="bg-primary hover:bg-primary/90 text-base sm:text-lg px-6 sm:px-8 py-5 sm:py-6 w-full sm:w-auto"
            >
              <Play className="w-4 h-4 sm:w-5 sm:h-5 mr-1 sm:mr-2" />
              <Link href="/videos">無料で動画を作成</Link>
            </Button>
          </motion.div>
          <motion.div
            whileHover={{ scale: 1.05 }}
            whileTap={{ scale: 0.95 }}
            className="w-full sm:w-auto"
          >
            <Button
              size="lg"
              variant="outline"
              className="text-base sm:text-lg px-6 sm:px-8 py-5 sm:py-6 bg-transparent w-full sm:w-auto"
            >
              <Video className="w-4 h-4 sm:w-5 sm:h-5 mr-1 sm:mr-2" />
              <Link href="/any-one">デモを見る</Link>
            </Button>
          </motion.div>
        </motion.div>
      </div>
    </section>
  )
}
