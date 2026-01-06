'use client'
import { Sparkles } from 'lucide-react'
import React from 'react'
import { Button } from '../ui/button'
import Link from 'next/link'
import { motion } from 'framer-motion'

export default function CTA() {
  return (
    <motion.section
      className="container mx-auto px-4 min-h-screen flex items-center justify-center py-10 w-full"
      initial={{ opacity: 0 }}
      whileInView={{ opacity: 1 }}
      viewport={{ once: true, margin: '-100px' }}
      transition={{ duration: 0.8 }}
    >
      <motion.div
        className="bg-gradient-to-br from-primary/90 via-primary/70 to-primary/40 rounded-2xl sm:rounded-3xl p-6 sm:p-8 md:p-12 text-center text-white relative overflow-hidden w-full"
        initial={{ scale: 0.95, y: 30 }}
        whileInView={{ scale: 1, y: 0 }}
        viewport={{ once: true, margin: '-100px' }}
        transition={{
          type: 'spring',
          stiffness: 100,
          damping: 15,
          delay: 0.2,
        }}
      >
        <motion.div
          className="absolute inset-0 bg-white/10 rounded-2xl sm:rounded-3xl backdrop-blur-[2px]"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ duration: 1, delay: 0.5 }}
        ></motion.div>
        <div className="relative z-10">
          <motion.h3
            className="text-2xl sm:text-3xl md:text-4xl font-bold mb-3 sm:mb-4 text-white drop-shadow-lg"
            initial={{ opacity: 0, y: -20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            transition={{ duration: 0.7, delay: 0.3 }}
          >
            今すぐ始めよう
          </motion.h3>
          <motion.p
            className="text-base sm:text-lg md:text-xl mb-5 sm:mb-8 text-white/95 max-w-2xl mx-auto drop-shadow-md"
            initial={{ opacity: 0 }}
            whileInView={{ opacity: 1 }}
            viewport={{ once: true }}
            transition={{ duration: 0.7, delay: 0.5 }}
          >
            あなたの旅行の思い出を、AIの力で特別な動画に変えませんか？ 無料で始められます。
          </motion.p>
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            transition={{ duration: 0.7, delay: 0.7 }}
            whileHover={{ scale: 1.05 }}
            whileTap={{ scale: 0.95 }}
            className="w-full sm:w-auto inline-block"
          >
            <Button
              size="lg"
              variant="secondary"
              className="text-base sm:text-lg px-6 sm:px-8 py-5 sm:py-6 bg-white text-primary hover:bg-white/90 shadow-lg w-full sm:w-auto"
            >
              <motion.div
                initial={{ rotate: -10 }}
                animate={{ rotate: [0, 15, 0] }}
                transition={{
                  repeat: Infinity,
                  repeatDelay: 2,
                  duration: 1,
                }}
              >
                <Sparkles className="w-5 h-5 mr-2" />
              </motion.div>
              <Link href="/videos">無料で動画を作成する</Link>
            </Button>
          </motion.div>
        </div>
      </motion.div>
    </motion.section>
  )
}
