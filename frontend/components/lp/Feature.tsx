'use client'
import { CloudUpload, Wand2, Share2 } from 'lucide-react'
import React from 'react'
import { motion } from 'framer-motion'

export default function Feature() {
  const containerVariants = {
    hidden: { opacity: 0 },
    visible: {
      opacity: 1,
      transition: {
        staggerChildren: 0.2,
      },
    },
  }

  const cardVariants = {
    hidden: { opacity: 0, y: 30 },
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

  const features = [
    {
      icon: CloudUpload,
      title: 'Quick Upload',
      description:
        'お気に入りの旅行写真やクリップをドラッグ＆ドロップするだけ。手間のかかる整理は不要です。',
      bgColor: 'bg-primary/10',
      iconColor: 'text-primary',
    },
    {
      icon: Wand2,
      title: 'AI Magic Editing',
      description:
        'AIが内容を分析し、リズムに合わせてカット割り。最適なエフェクトと字幕を自動で追加します。',
      bgColor: 'bg-secondary/15',
      iconColor: 'text-secondary',
    },
    {
      icon: Share2,
      title: 'Instant Sharing',
      description:
        '出来上がったVlogをワンタップで書き出し。InstagramやTikTokへすぐにシェアできます。',
      bgColor: 'bg-primary/10',
      iconColor: 'text-primary',
    },
  ]

  return (
    <section id="features" className="bg-accent dark:bg-gray-950 py-24">
      <div className="max-w-[1200px] mx-auto px-6">
        <motion.div
          className="text-center mb-16"
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true, margin: '-100px' }}
          transition={{ duration: 0.7 }}
        >
          <h2 className="text-3xl md:text-4xl font-black mb-4">Easy as 1-2-3</h2>
          <p className="text-muted-foreground">
            最新のAI技術が、あなたの代わりにストーリーを紡ぎます。
          </p>
        </motion.div>

        <motion.div
          className="grid md:grid-cols-3 gap-8"
          variants={containerVariants}
          initial="hidden"
          whileInView="visible"
          viewport={{ once: true, margin: '-50px' }}
        >
          {features.map((feature, index) => (
            <motion.div
              key={index}
              variants={cardVariants}
              whileHover={{ y: -8, transition: { duration: 0.3 } }}
              className="bg-white dark:bg-gray-900 p-8 rounded-lg flex flex-col items-center text-center gap-4 kawaii-shadow"
            >
              <div
                className={`size-16 ${feature.bgColor} ${feature.iconColor} rounded-full flex items-center justify-center`}
              >
                <feature.icon className="w-8 h-8" strokeWidth={1.5} />
              </div>
              <h3 className="text-xl font-bold">{feature.title}</h3>
              <p className="text-sm text-muted-foreground leading-relaxed">{feature.description}</p>
            </motion.div>
          ))}
        </motion.div>
      </div>
    </section>
  )
}
