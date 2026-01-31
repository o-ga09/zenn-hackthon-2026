'use client'
import React, { useState, useRef, useCallback } from 'react'
import { Button } from '../ui/button'
import { Upload, MapPin, Zap, Sparkles, X, CheckCircle, Loader2, AlertCircle, Video } from 'lucide-react'
import { motion, AnimatePresence } from 'framer-motion'
import Image from 'next/image'
import { createTrialVlog, TrialVlogResponse } from '@/api/trialApi'

// 生成結果のモーダル
function GenerationResultModal({
  isOpen,
  onClose,
  result,
  isGenerating,
  error,
}: {
  isOpen: boolean
  onClose: () => void
  result: TrialVlogResponse | null
  isGenerating: boolean
  error: string | null
}) {
  if (!isOpen) return null

  return (
    <AnimatePresence>
      <motion.div
        className="fixed inset-0 bg-black/50 backdrop-blur-sm z-50 flex items-center justify-center p-4"
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        exit={{ opacity: 0 }}
        onClick={onClose}
      >
        <motion.div
          className="bg-white dark:bg-gray-900 rounded-2xl shadow-2xl max-w-md w-full p-6"
          initial={{ scale: 0.9, opacity: 0 }}
          animate={{ scale: 1, opacity: 1 }}
          exit={{ scale: 0.9, opacity: 0 }}
          onClick={e => e.stopPropagation()}
        >
          {isGenerating ? (
            <div className="text-center py-8">
              <Loader2 className="w-16 h-16 text-primary mx-auto mb-4 animate-spin" />
              <h3 className="text-xl font-bold mb-2">VLogを生成中...</h3>
              <p className="text-muted-foreground text-sm">
                AIが素敵な動画を作成しています。
                <br />
                しばらくお待ちください。
              </p>
            </div>
          ) : error ? (
            <div className="text-center py-8">
              <AlertCircle className="w-16 h-16 text-destructive mx-auto mb-4" />
              <h3 className="text-xl font-bold mb-2">エラーが発生しました</h3>
              <p className="text-muted-foreground text-sm mb-4">{error}</p>
              <Button onClick={onClose} variant="outline">
                閉じる
              </Button>
            </div>
          ) : result ? (
            <div className="text-center py-4">
              <CheckCircle className="w-16 h-16 text-green-500 mx-auto mb-4" />
              <h3 className="text-xl font-bold mb-2">VLog生成完了！</h3>
              <p className="text-muted-foreground text-sm mb-4">
                {result.message || 'VLogが正常に生成されました。'}
              </p>
              {result.videoUrl && (
                <div className="mb-4">
                  <video
                    src={result.videoUrl}
                    controls
                    className="w-full rounded-lg shadow-lg"
                  />
                </div>
              )}
              <div className="flex gap-2 justify-center">
                {result.videoUrl && (
                  <Button
                    asChild
                    className="bg-primary hover:bg-primary/90"
                  >
                    <a href={result.videoUrl} download>
                      ダウンロード
                    </a>
                  </Button>
                )}
                <Button onClick={onClose} variant="outline">
                  閉じる
                </Button>
              </div>
            </div>
          ) : null}
        </motion.div>
      </motion.div>
    </AnimatePresence>
  )
}

export default function Hero() {
  const [destination, setDestination] = useState('')
  const [uploadedFiles, setUploadedFiles] = useState<File[]>([])
  const [isDragging, setIsDragging] = useState(false)
  const [isGenerating, setIsGenerating] = useState(false)
  const [showModal, setShowModal] = useState(false)
  const [generationResult, setGenerationResult] = useState<TrialVlogResponse | null>(null)
  const [generationError, setGenerationError] = useState<string | null>(null)
  const fileInputRef = useRef<HTMLInputElement>(null)

  const handleFileUpload = useCallback((files: FileList | null) => {
    if (!files) return
    const imageFiles = Array.from(files).filter(file => 
      file.type.startsWith('image/') || file.type.startsWith('video/')
    )
    if (imageFiles.length > 0) {
      setUploadedFiles(prev => [...prev, ...imageFiles].slice(0, 10)) // 最大10ファイル
    }
  }, [])

  const handleDragEnter = (e: React.DragEvent) => {
    e.preventDefault()
    e.stopPropagation()
    setIsDragging(true)
  }

  const handleDragLeave = (e: React.DragEvent) => {
    e.preventDefault()
    e.stopPropagation()
    setIsDragging(false)
  }

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault()
    e.stopPropagation()
    setIsDragging(true)
  }

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault()
    e.stopPropagation()
    setIsDragging(false)
    handleFileUpload(e.dataTransfer.files)
  }

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    handleFileUpload(e.target.files)
  }

  const removeFile = (index: number) => {
    setUploadedFiles(prev => prev.filter((_, i) => i !== index))
  }

  const handleGenerateVlog = async () => {
    if (uploadedFiles.length === 0) {
      setGenerationError('ファイルがアップロードされていません。')
      setShowModal(true)
      return
    }

    setIsGenerating(true)
    setGenerationError(null)
    setGenerationResult(null)
    setShowModal(true)

    try {
      const result = await createTrialVlog({
        files: uploadedFiles,
        destination: destination || undefined,
        title: `${destination || 'My'} Travel Vlog`,
        theme: 'adventure',
        duration: 15,
      })
      setGenerationResult(result)
    } catch {
      setGenerationError(
        '生成中にエラーが発生しました。サーバーが起動していない可能性があります。しばらく後でお試しください。'
      )
    } finally {
      setIsGenerating(false)
    }
  }

  const closeModal = () => {
    setShowModal(false)
    if (generationResult) {
      // 成功時はファイルをクリア
      setUploadedFiles([])
      setDestination('')
    }
  }

  return (
    <>
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
              <Video className="w-16 h-16 text-primary" />
            </div>

            <h3 className="text-xl font-bold mb-6 flex items-center gap-2">
              <Sparkles className="w-5 h-5 text-primary" />
              Trial Generation
            </h3>

            <div className="space-y-6">
              {/* アップロードエリア */}
              <motion.div
                className={`cursor-pointer border-2 border-dashed rounded-lg p-6 flex flex-col items-center justify-center gap-3 transition-all ${
                  isDragging
                    ? 'border-primary bg-primary/5 scale-[1.02]'
                    : 'border-gray-200 dark:border-gray-700 hover:border-primary/50 hover:bg-primary/5'
                }`}
                whileTap={{ scale: 0.98 }}
                onDragEnter={handleDragEnter}
                onDragOver={handleDragOver}
                onDragLeave={handleDragLeave}
                onDrop={handleDrop}
                onClick={() => fileInputRef.current?.click()}
              >
                <input
                  ref={fileInputRef}
                  type="file"
                  multiple
                  accept="image/*,video/*"
                  onChange={handleInputChange}
                  className="hidden"
                />
                <div className="size-12 bg-accent rounded-full flex items-center justify-center text-primary">
                  <Upload className="w-6 h-6" />
                </div>
                <div className="text-center">
                  <p className="font-bold text-sm">3〜5枚の写真を選択</p>
                  <p className="text-xs text-muted-foreground mt-1">
                    または短い動画クリップをドロップ
                  </p>
                </div>
              </motion.div>

              {/* アップロード済みファイルプレビュー */}
              <AnimatePresence>
                {uploadedFiles.length > 0 && (
                  <motion.div
                    className="space-y-2"
                    initial={{ opacity: 0, height: 0 }}
                    animate={{ opacity: 1, height: 'auto' }}
                    exit={{ opacity: 0, height: 0 }}
                  >
                    <div className="flex items-center justify-between">
                      <span className="text-xs font-bold text-muted-foreground">
                        アップロード済み ({uploadedFiles.length}枚)
                      </span>
                      <button
                        onClick={e => {
                          e.stopPropagation()
                          setUploadedFiles([])
                        }}
                        className="text-xs text-destructive hover:underline"
                      >
                        すべて削除
                      </button>
                    </div>
                    <div className="grid grid-cols-5 gap-2">
                      {uploadedFiles.slice(0, 5).map((file, index) => (
                        <motion.div
                          key={index}
                          className="relative group aspect-square"
                          initial={{ opacity: 0, scale: 0.8 }}
                          animate={{ opacity: 1, scale: 1 }}
                          exit={{ opacity: 0, scale: 0.8 }}
                        >
                          <img
                            src={URL.createObjectURL(file)}
                            alt={`Upload ${index + 1}`}
                            className="w-full h-full object-cover rounded-lg"
                          />
                          <button
                            onClick={e => {
                              e.stopPropagation()
                              removeFile(index)
                            }}
                            className="absolute -top-1 -right-1 w-5 h-5 bg-destructive text-white rounded-full flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity"
                          >
                            <X className="w-3 h-3" />
                          </button>
                        </motion.div>
                      ))}
                      {uploadedFiles.length > 5 && (
                        <div className="aspect-square bg-muted rounded-lg flex items-center justify-center text-sm font-bold text-muted-foreground">
                          +{uploadedFiles.length - 5}
                        </div>
                      )}
                    </div>
                  </motion.div>
                )}
              </AnimatePresence>

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
              <Button
                className="w-full bg-primary hover:bg-primary/90 text-white font-black py-4 rounded-lg flex items-center justify-center gap-2 group transition-all h-14 disabled:opacity-50 disabled:cursor-not-allowed"
                onClick={handleGenerateVlog}
                disabled={uploadedFiles.length === 0 || isGenerating}
              >
                {isGenerating ? (
                  <>
                    <Loader2 className="w-5 h-5 animate-spin" />
                    <span>生成中...</span>
                  </>
                ) : (
                  <>
                    <span>15秒のプレビューを生成</span>
                    <Zap className="w-5 h-5 group-hover:translate-x-1 transition-transform" />
                  </>
                )}
              </Button>

              <p className="text-center text-[10px] text-muted-foreground">
                ログイン不要・無料で今すぐお試しいただけます
              </p>
            </div>
          </motion.div>
        </div>
      </section>

      {/* 生成結果モーダル */}
      <GenerationResultModal
        isOpen={showModal}
        onClose={closeModal}
        result={generationResult}
        isGenerating={isGenerating}
        error={generationError}
      />
    </>
  )
}
