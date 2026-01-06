'use client'
import { UploadCloud, FileText, WandSparkles } from 'lucide-react'
import React from 'react'
import { Card, CardHeader, CardTitle, CardContent, CardDescription } from '../ui/card'
import { motion } from 'framer-motion'

export default function Feature() {
  // アニメーション変数
  const containerVariants = {
    hidden: { opacity: 0 },
    visible: {
      opacity: 1,
      transition: {
        staggerChildren: 0.3,
      },
    },
  }

  const cardVariants = {
    hidden: { opacity: 0, y: 50 },
    visible: {
      opacity: 1,
      y: 0,
      transition: {
        type: 'spring' as const,
        stiffness: 100,
        damping: 12,
      },
    },
  }

  const iconVariants = {
    hidden: { scale: 0, rotate: -180 },
    visible: {
      scale: 1,
      rotate: 0,
      transition: {
        type: 'spring' as const,
        stiffness: 260,
        damping: 20,
        delay: 0.2,
      },
    },
  }

  return (
    <section className="container mx-auto px-4 min-h-screen flex flex-col justify-center py-10 sm:py-16">
      <motion.div
        className="text-center mb-10 sm:mb-16"
        initial={{ opacity: 0, y: 20 }}
        whileInView={{ opacity: 1, y: 0 }}
        viewport={{ once: true, margin: '-100px' }}
        transition={{ duration: 0.7 }}
      >
        <h2 className="text-3xl sm:text-4xl font-bold text-foreground mb-3 sm:mb-4">
          3つのステップで完成
        </h2>
        <p className="text-base sm:text-xl text-muted-foreground max-w-2xl mx-auto px-2">
          複雑な編集作業は不要。シンプルな操作で、プロ品質の動画が完成します。
        </p>
      </motion.div>

      <motion.div
        className="grid sm:grid-cols-2 md:grid-cols-3 gap-6 sm:gap-8 max-w-5xl mx-auto"
        variants={containerVariants}
        initial="hidden"
        whileInView="visible"
        viewport={{ once: true, margin: '-50px' }}
      >
        <motion.div variants={cardVariants} whileHover={{ y: -10, transition: { duration: 0.3 } }}>
          <Card className="text-center border-0 shadow-lg bg-card/50 backdrop-blur-sm">
            <CardHeader>
              <motion.div
                className="w-16 h-16 bg-orange-100 rounded-full flex items-center justify-center mx-auto mb-4 shadow-md"
                variants={iconVariants}
              >
                <div className="text-orange-500">
                  <UploadCloud size={36} strokeWidth={2} />
                </div>
              </motion.div>
              <CardTitle className="text-2xl">写真をアップロード</CardTitle>
            </CardHeader>
            <CardContent>
              <CardDescription className="text-base">
                旅行の思い出の写真を選んでアップロード。複数枚でも一度に処理できます。
              </CardDescription>
            </CardContent>
          </Card>
        </motion.div>

        <motion.div variants={cardVariants} whileHover={{ y: -10, transition: { duration: 0.3 } }}>
          <Card className="text-center border-0 shadow-lg bg-card/50 backdrop-blur-sm">
            <CardHeader>
              <motion.div
                className="w-16 h-16 bg-blue-100 rounded-full flex items-center justify-center mx-auto mb-4 shadow-md"
                variants={iconVariants}
              >
                <div className="text-blue-500">
                  <FileText size={36} strokeWidth={2} />
                </div>
              </motion.div>
              <CardTitle className="text-2xl">旅行情報を入力</CardTitle>
            </CardHeader>
            <CardContent>
              <CardDescription className="text-base">
                旅行のタイトルと日付を入力。AIがこの情報を使って動画をパーソナライズします。
              </CardDescription>
            </CardContent>
          </Card>
        </motion.div>

        <motion.div variants={cardVariants} whileHover={{ y: -10, transition: { duration: 0.3 } }}>
          <Card className="text-center border-0 shadow-lg bg-card/50 backdrop-blur-sm">
            <CardHeader>
              <motion.div
                className="w-16 h-16 bg-purple-100 rounded-full flex items-center justify-center mx-auto mb-4 shadow-md"
                variants={iconVariants}
              >
                <div className="text-purple-500">
                  <WandSparkles size={36} strokeWidth={2} />
                </div>
              </motion.div>
              <CardTitle className="text-2xl">AI動画生成</CardTitle>
            </CardHeader>
            <CardContent>
              <CardDescription className="text-base">
                AIが自動的にエモーショナルな動画を生成。音楽やエフェクトも自動で追加されます。
              </CardDescription>
            </CardContent>
          </Card>
        </motion.div>
      </motion.div>
    </section>
  )
}
