'use client'
import React from 'react'
import { motion } from 'framer-motion'
import { Check, Verified } from 'lucide-react'
import { Button } from '../ui/button'
import Link from 'next/link'

export default function Pricing() {
  const plans = [
    {
      name: 'Free Plan',
      subtitle: 'まずはここから',
      price: '¥0',
      period: '/month',
      badge: 'CURRENT',
      features: [
        { text: '3 Vlogs per month', active: true },
        { text: '720p Resolution', active: false },
        { text: 'Standard Library', active: false },
      ],
      buttonText: 'Get Started',
      buttonVariant: 'outline' as const,
      featured: false,
    },
    {
      name: 'Pro Plan',
      subtitle: 'もっとクリエイティブに',
      price: '¥980',
      period: '/month',
      badge: 'Recommended',
      features: [
        { text: 'Unlimited Vlogs', active: true },
        { text: '4K Resolution', active: true },
        { text: 'Premium Music & Fonts', active: true },
        { text: 'No Watermark', active: true },
      ],
      buttonText: 'Start Pro Trial',
      buttonVariant: 'default' as const,
      featured: true,
    },
  ]

  return (
    <section id="pricing" className="py-24 max-w-[1200px] mx-auto px-6">
      <motion.div
        className="text-center mb-16"
        initial={{ opacity: 0, y: 20 }}
        whileInView={{ opacity: 1, y: 0 }}
        viewport={{ once: true }}
        transition={{ duration: 0.7 }}
      >
        <h2 className="text-3xl md:text-4xl font-black mb-4">Pricing Plans</h2>
        <p className="text-muted-foreground">ライフスタイルに合わせて選べるシンプルなプラン。</p>
      </motion.div>

      <div className="grid md:grid-cols-2 gap-8 max-w-4xl mx-auto">
        {plans.map((plan, index) => (
          <motion.div
            key={plan.name}
            className={`border-2 ${plan.featured ? 'border-primary kawaii-shadow' : 'border-border hover:border-primary/20'} rounded-xl p-8 flex flex-col gap-6 bg-white dark:bg-gray-900 relative transition-all`}
            initial={{ opacity: 0, y: 30 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            transition={{ duration: 0.5, delay: index * 0.15 }}
            whileHover={{ y: -4 }}
          >
            {plan.featured && (
              <div className="absolute -top-4 left-1/2 -translate-x-1/2 bg-primary text-white text-[10px] font-black px-4 py-1 rounded-full uppercase tracking-widest">
                {plan.badge}
              </div>
            )}

            <div className="flex justify-between items-start">
              <div>
                <h3 className="text-xl font-black">{plan.name}</h3>
                <p className="text-sm text-muted-foreground">{plan.subtitle}</p>
              </div>
              {plan.featured ? (
                <Verified className="w-6 h-6 text-primary" />
              ) : (
                <span className="px-3 py-1 bg-muted text-[10px] font-bold rounded-full">
                  {plan.badge}
                </span>
              )}
            </div>

            <div className="flex items-baseline gap-1">
              <span className="text-4xl font-black tracking-tight">{plan.price}</span>
              <span className="text-muted-foreground font-bold">{plan.period}</span>
            </div>

            <ul className="space-y-3 mt-4">
              {plan.features.map((feature, featureIndex) => (
                <li
                  key={featureIndex}
                  className={`flex items-center gap-3 text-sm ${feature.active ? '' : 'text-muted-foreground'}`}
                >
                  <Check
                    className={`w-4 h-4 ${feature.active ? 'text-secondary' : 'text-muted-foreground'}`}
                  />
                  {feature.text}
                </li>
              ))}
            </ul>

            <Link href="/upload" className="mt-auto">
              <Button
                className={`w-full py-3 rounded-lg font-bold ${
                  plan.featured
                    ? 'bg-primary text-white hover:bg-primary/90'
                    : 'bg-muted hover:bg-gray-200 text-foreground'
                } transition-colors`}
              >
                {plan.buttonText}
              </Button>
            </Link>
          </motion.div>
        ))}
      </div>
    </section>
  )
}
