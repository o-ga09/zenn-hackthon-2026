'use client'
import React from 'react'
import { Button } from '../ui/button'
import Link from 'next/link'
import { motion } from 'framer-motion'

export default function CTA() {
  return (
    <motion.section
      className="py-24 px-6 bg-primary text-white text-center"
      initial={{ opacity: 0 }}
      whileInView={{ opacity: 1 }}
      viewport={{ once: true, margin: '-100px' }}
      transition={{ duration: 0.8 }}
    >
      <motion.h2
        className="text-4xl font-black mb-6 italic"
        initial={{ opacity: 0, y: 20 }}
        whileInView={{ opacity: 1, y: 0 }}
        viewport={{ once: true }}
        transition={{ duration: 0.7, delay: 0.2 }}
      >
        旅はまだ、終わらない。
      </motion.h2>

      <motion.p
        className="text-white/80 max-w-xl mx-auto mb-10"
        initial={{ opacity: 0 }}
        whileInView={{ opacity: 1 }}
        viewport={{ once: true }}
        transition={{ duration: 0.7, delay: 0.4 }}
      >
        今すぐあなたのiPhoneから最高の一枚を選んで、最初のVlogを作ってみませんか？
      </motion.p>

      <motion.div
        initial={{ opacity: 0, y: 20 }}
        whileInView={{ opacity: 1, y: 0 }}
        viewport={{ once: true }}
        transition={{ duration: 0.7, delay: 0.6 }}
        whileHover={{ scale: 1.05 }}
        whileTap={{ scale: 0.95 }}
        className="inline-block"
      >
        <Link href="/upload">
          <Button
            size="lg"
            className="bg-white text-primary font-black px-10 py-6 rounded-full text-lg shadow-xl hover:bg-white/90 transition-all"
          >
            今すぐ無料で始める
          </Button>
        </Link>
      </motion.div>
    </motion.section>
  )
}
