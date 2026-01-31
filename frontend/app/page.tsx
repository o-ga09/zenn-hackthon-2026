'use client'
import LPFooter from '@/components/footer/LP-Footer'
import LPHeader from '@/components/header/LP-Header'
import Feature from '@/components/lp/Feature'
import Hero from '@/components/lp/Hero'

import CTA from '@/components/lp/CTA'
import { motion } from 'framer-motion'
import Pricing from '@/components/lp/Pricing'

export default function HomePage() {
  return (
    <motion.div
      className="min-h-screen bg-background dark:bg-[#101922] relative overflow-hidden transition-colors duration-300"
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      transition={{ duration: 0.5 }}
    >
      {/* Blob backgrounds */}
      <div className="blob-bg top-[-10%] left-[-5%] w-[400px] h-[400px] bg-primary fixed" />
      <div className="blob-bg bottom-[20%] right-[-5%] w-[300px] h-[300px] bg-secondary fixed" />

      <motion.div
        initial={{ opacity: 0, y: -20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.7 }}
        className="relative z-10"
      >
        <LPHeader />
      </motion.div>

      <main className="relative z-10">
        <Hero />
        <Feature />
        <Pricing />
        <CTA />
      </main>

      <motion.div
        initial={{ opacity: 0, y: 20 }}
        whileInView={{ opacity: 1, y: 0 }}
        viewport={{ once: true }}
        transition={{ duration: 0.7 }}
        className="relative z-10"
      >
        <LPFooter />
      </motion.div>
    </motion.div>
  )
}
