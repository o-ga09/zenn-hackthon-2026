'use client'

import { Button } from '@/components/ui/button'
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Camera, Upload } from 'lucide-react'
import React, { useState, useRef } from 'react'
import { X } from 'lucide-react'
import { useUploadForm } from './UploadFormContext'

export default function PhotoUpload() {
  const { uploadedFiles, addFiles, removeFile } = useUploadForm()
  const [isDragging, setIsDragging] = useState(false)
  const fileInputRef = useRef<HTMLInputElement>(null)

  const handleFileUpload = (event: React.ChangeEvent<HTMLInputElement>) => {
    const files = Array.from(event.target.files || [])
    if (files.length > 0) {
      addFiles(files)
    }
  }

  const handleDragEnter = (e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault()
    e.stopPropagation()
    setIsDragging(true)
  }

  const handleDragLeave = (e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault()
    e.stopPropagation()
    setIsDragging(false)
  }

  const handleDragOver = (e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault()
    e.stopPropagation()
    setIsDragging(true)
  }

  const handleDrop = (e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault()
    e.stopPropagation()
    setIsDragging(false)

    const files = Array.from(e.dataTransfer.files)
    if (files.length > 0) {
      const imageFiles = files.filter(file => file.type.startsWith('image/'))
      if (imageFiles.length > 0) {
        addFiles(imageFiles)
      }
    }
  }

  const handleButtonClick = () => {
    fileInputRef.current?.click()
  }

  return (
    <div className="space-y-6">
      <Card className="border-0 shadow-lg bg-card/50 backdrop-blur-sm">
        <CardHeader>
          <CardTitle className="flex items-center space-x-2 text-base md:text-lg lg:text-xl">
            <Camera className="w-4 h-4 md:w-5 md:h-5" />
            <span>写真をアップロード</span>
          </CardTitle>
          <CardDescription className="text-xs md:text-sm">
            旅行の思い出の写真を選択してください（複数選択可能）
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            <div
              className={`border-2 border-dashed rounded-lg p-4 md:p-8 text-center transition-colors ${
                isDragging ? 'border-primary bg-primary/5' : 'border-border hover:border-primary/50'
              }`}
              onDragEnter={handleDragEnter}
              onDragOver={handleDragOver}
              onDragLeave={handleDragLeave}
              onDrop={handleDrop}
              onClick={handleButtonClick}
            >
              <Upload
                className={`w-10 h-10 md:w-12 md:h-12 mx-auto mb-3 md:mb-4 ${
                  isDragging ? 'text-primary' : 'text-muted-foreground'
                }`}
              />
              <p className="text-sm md:text-base text-muted-foreground mb-3 md:mb-4">
                ここに写真をドラッグ&ドロップするか、クリックして選択
              </p>
              <Input
                ref={fileInputRef}
                type="file"
                multiple
                accept="image/*"
                onChange={handleFileUpload}
                className="hidden"
                id="file-upload"
              />
              <Button
                variant="outline"
                className="cursor-pointer bg-transparent text-sm md:text-base"
              >
                写真を選択
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Uploaded Files Preview */}
      {uploadedFiles.length > 0 && (
        <div className="space-y-1 md:space-y-2">
          <h4 className="font-semibold text-sm md:text-base">
            アップロード済み写真 ({uploadedFiles.length}枚)
          </h4>
          <div className="grid grid-cols-3 sm:grid-cols-4 md:grid-cols-5 gap-1 md:gap-2 max-h-36 md:max-h-48 overflow-y-auto p-1 md:p-2 border rounded-lg">
            {uploadedFiles.map((file, index) => (
              <div key={index} className="relative group">
                <img
                  src={URL.createObjectURL(file)}
                  alt={`Upload ${index + 1}`}
                  className="w-full h-16 md:h-20 object-cover rounded-lg"
                />
                <Button
                  size="sm"
                  variant="destructive"
                  className="absolute top-1 right-1 w-5 h-5 md:w-6 md:h-6 p-0 opacity-0 group-hover:opacity-100 transition-opacity"
                  onClick={() => removeFile(index)}
                >
                  <X className="w-2 h-2 md:w-3 md:h-3" />
                </Button>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  )
}
