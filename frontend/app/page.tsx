'use client'
import LPFooter from '@/components/footer/LP-Footer'
import LPHeader from '@/components/header/LP-Header'
import Benefit from '@/components/lp/Benefit'
import CTA from '@/components/lp/CTA'
import Feature from '@/components/lp/Feature'
import Hero from '@/components/lp/Hero'
import { motion } from 'framer-motion'

export default function HomePage() {
  return (
    <motion.div
      className="min-h-screen bg-gradient-to-br from-pink-100 via-purple-50 to-blue-100"
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      transition={{ duration: 0.5 }}
    >
      <motion.div
        initial={{ opacity: 0, y: -20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.7 }}
      >
        <LPHeader />
      </motion.div>

      <Hero />
      <Feature />
      <Benefit />
      <CTA />

      <motion.div
        initial={{ opacity: 0, y: 20 }}
        whileInView={{ opacity: 1, y: 0 }}
        viewport={{ once: true }}
        transition={{ duration: 0.7 }}
      >
        <LPFooter />
      </motion.div>
    </motion.div>
  )
}
