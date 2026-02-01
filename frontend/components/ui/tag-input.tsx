'use client'

import { useState, KeyboardEvent } from 'react'
import { Badge } from '@/components/ui/badge'
import { Input } from '@/components/ui/input'
import { X } from 'lucide-react'

interface TagInputProps {
  tags: string[]
  onChange: (tags: string[]) => void
  placeholder?: string
  className?: string
}

export function TagInput({
  tags,
  onChange,
  placeholder = 'タグを追加...',
  className = '',
}: TagInputProps) {
  const [inputValue, setInputValue] = useState('')

  const addTag = (tag: string) => {
    const trimmed = tag.trim()
    if (trimmed && !tags.includes(trimmed)) {
      onChange([...tags, trimmed])
    }
    setInputValue('')
  }

  const removeTag = (index: number) => {
    onChange(tags.filter((_, i) => i !== index))
  }

  const handleKeyDown = (e: KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter') {
      e.preventDefault()
      addTag(inputValue)
    } else if (e.key === 'Backspace' && !inputValue && tags.length > 0) {
      // Backspaceキーで最後のタグを削除
      onChange(tags.slice(0, -1))
    }
  }

  return (
    <div className={`flex flex-wrap gap-2 p-2 border rounded-md min-h-[44px] ${className}`}>
      {tags.map((tag, index) => (
        <Badge key={index} variant="secondary" className="gap-1 h-7 px-2 flex items-center">
          <span className="text-sm">{tag}</span>
          <button
            type="button"
            onClick={() => removeTag(index)}
            className="hover:bg-destructive/20 rounded-full p-0.5 transition-colors"
            aria-label={`${tag}を削除`}
          >
            <X className="w-3 h-3" />
          </button>
        </Badge>
      ))}
      <Input
        value={inputValue}
        onChange={e => setInputValue(e.target.value)}
        onKeyDown={handleKeyDown}
        placeholder={placeholder}
        className="flex-1 min-w-[120px] border-0 focus-visible:ring-0 shadow-none h-auto p-0"
      />
    </div>
  )
}
