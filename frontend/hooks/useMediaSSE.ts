import { useEffect, useRef, useState } from 'react'
import { Media } from '@/api/types'
import { useQueryClient } from '@tanstack/react-query'
import { MEDIA_QUERY_KEYS } from '@/api/mediaApi'

interface MediaStatusUpdate {
  medias: Media[]
  totalItems: number
  completedItems: number
  failedItems: number
  allCompleted: boolean
  // バックエンドからのスネークケースもサポート
  all_completed?: boolean
  total_items?: number
  completed_items?: number
  failed_items?: number
}

interface UseMediaSSEOptions {
  mediaIds: string[]
  onComplete?: (medias: Media[]) => void
  onError?: (error: Error) => void
  onProgress?: (update: MediaStatusUpdate) => void
  onFetchNotifications?: () => Promise<void>
  enabled?: boolean
}

/**
 * メディア分析のSSEストリームを監視するカスタムフック
 */
export function useMediaSSE({
  mediaIds,
  onComplete,
  onError,
  onProgress,
  onFetchNotifications,
  enabled = true,
}: UseMediaSSEOptions) {
  const [status, setStatus] = useState<MediaStatusUpdate | null>(null)
  const [isConnected, setIsConnected] = useState(false)
  const eventSourceRef = useRef<EventSource | null>(null)
  const isCompletedRef = useRef(false)
  const hasReceivedMessageRef = useRef(false)
  const initialMediaIdsRef = useRef<string[]>([])
  const lastReceivedStatusRef = useRef<MediaStatusUpdate | null>(null) // refに変更
  const queryClient = useQueryClient()

  // [修正1] コールバック関数を ref に安定させる
  // → 毎レンダレンタルで新しいインスタンスになるコールバックを
  //    useEffect の依存配列から外し、再接続を防止する
  const onCompleteRef = useRef(onComplete)
  const onErrorRef = useRef(onError)
  const onProgressRef = useRef(onProgress)
  const onFetchNotificationsRef = useRef(onFetchNotifications)
  useEffect(() => {
    onCompleteRef.current = onComplete
  }, [onComplete])
  useEffect(() => {
    onErrorRef.current = onError
  }, [onError])
  useEffect(() => {
    onProgressRef.current = onProgress
  }, [onProgress])
  useEffect(() => {
    onFetchNotificationsRef.current = onFetchNotifications
  }, [onFetchNotifications])

  // [修正2] mediaIds を文字列化して依存配列に含める
  // → 配列参照の変化による不必要な再実行を防つつ、
  //    内容が変わった時だけ再接続する
  const mediaIdsKey = mediaIds.join(',')

  // [修正3] enabled の現在の値を ref で保持する
  // → クリーンアップ関数がクロージャで古い値を参照するのを修正
  const enabledRef = useRef(enabled)
  useEffect(() => {
    enabledRef.current = enabled
  }, [enabled])

  useEffect(() => {
    if (!enabled) {
      console.log('[SSE] スキップ: enabled=false')
      return
    }

    // 初回接続時に mediaIds を保存
    if (initialMediaIdsRef.current.length === 0 && mediaIds.length > 0) {
      initialMediaIdsRef.current = mediaIds
      console.log('[SSE] 初回接続用のmediaIDsを保存:', mediaIds)
    }

    const targetMediaIds =
      initialMediaIdsRef.current.length > 0 ? initialMediaIdsRef.current : mediaIds

    if (targetMediaIds.length === 0) {
      console.log('[SSE] スキップ: mediaIds.length=0')
      return
    }

    // フラグをリセット
    isCompletedRef.current = false
    hasReceivedMessageRef.current = false
    lastReceivedStatusRef.current = null // refをリセット

    // SSE接続を確立
    const idsParam = targetMediaIds.join(',')
    const baseURL = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080'
    const url = `${baseURL}/api/agent/analyze-media/stream?ids=${idsParam}`

    console.log('[SSE] 接続開始:', url)
    console.log('[SSE] 使用するMedia IDs:', targetMediaIds)

    const eventSource = new EventSource(url, {
      withCredentials: true,
    })

    eventSourceRef.current = eventSource

    eventSource.onopen = () => {
      setIsConnected(true)
      console.log('[SSE] 接続成功 (readyState:', eventSource.readyState, ')')
    }

    eventSource.onmessage = event => {
      hasReceivedMessageRef.current = true
      console.log('[SSE] メッセージ受信:', event.data)
      try {
        const rawData: any = JSON.parse(event.data)
        console.log('[SSE] パース済みデータ（raw）:', rawData)

        // スネークケースからキャメルケースに正規化
        const data: MediaStatusUpdate = {
          medias: rawData.medias || [],
          totalItems: rawData.total_items || rawData.totalItems || 0,
          completedItems: rawData.completed_items || rawData.completedItems || 0,
          failedItems: rawData.failed_items || rawData.failedItems || 0,
          allCompleted: rawData.all_completed ?? rawData.allCompleted ?? false,
          // 元のフィールドも保持
          all_completed: rawData.all_completed,
          total_items: rawData.total_items,
          completed_items: rawData.completed_items,
          failed_items: rawData.failed_items,
        }
        console.log('[SSE] 正規化済みデータ:', data)
        setStatus(data)
        lastReceivedStatusRef.current = data // refに保存

        // 進捗コールバック
        if (onProgressRef.current) {
          onProgressRef.current(data)
        }

        // 全て完了した場合のみキャッシュを無効化する
        // [修正4] 毎イベントでの invalidateQueries を廃止
        // → リフェッチの連鎖を防止し、完了時にまとめて更新する
        if (data.allCompleted) {
          console.log('[SSE] 全ての分析が完了しました（onmessage）')
          console.log('[SSE] onComplete実行前 - medias:', data.medias)

          isCompletedRef.current = true

          queryClient.invalidateQueries({ queryKey: MEDIA_QUERY_KEYS.images() })

          // 通知をフェッチ
          if (onFetchNotificationsRef.current) {
            onFetchNotificationsRef
              .current()
              .catch(err => console.error('[SSE] Failed to fetch notifications:', err))
          }

          // コールバックを実行
          if (onCompleteRef.current) {
            console.log('[SSE] onCompleteコールバックを実行します')
            onCompleteRef.current(data.medias)
          } else {
            console.warn('[SSE] onCompleteコールバックがありません')
          }

          // [修正5] setTimeout を削除し直接閉じる
          // → isCompletedRef が既に true なので onerror側で正しくハンドリングされる
          eventSource.close()
          setIsConnected(false)
        }
      } catch (error) {
        console.error('[SSE] データパースエラー:', error, 'rawData:', event.data)
        if (onErrorRef.current) {
          onErrorRef.current(error as Error)
        }
      }
    }

    // [追加] 明示的な完了イベントをリッスン
    // バックエンドから "event: complete" が送信された場合に確実に完了を検知
    eventSource.addEventListener('complete', event => {
      console.log('[SSE] 完了イベント受信:', event)
      console.log('[SSE] 完了イベント時の状態:', {
        isCompleted: isCompletedRef.current,
        lastReceivedStatus: lastReceivedStatusRef.current,
        onCompleteExists: !!onCompleteRef.current,
      })

      isCompletedRef.current = true

      // onmessage で既に完了処理されていない場合のフォールバック
      const lastStatus = lastReceivedStatusRef.current
      const isAllCompleted = lastStatus?.allCompleted || lastStatus?.all_completed

      console.log('[SSE] 完了イベント詳細チェック:', {
        lastStatus: !!lastStatus,
        allCompleted: lastStatus?.allCompleted,
        all_completed: lastStatus?.all_completed,
        isAllCompleted,
        onCompleteExists: !!onCompleteRef.current,
      })

      if (lastStatus && isAllCompleted && onCompleteRef.current) {
        console.log('[SSE] 完了イベントでフォールバック処理を実行')
        console.log('[SSE] フォールバック - medias:', lastStatus.medias)
        queryClient.invalidateQueries({ queryKey: MEDIA_QUERY_KEYS.images() })
        onCompleteRef.current(lastStatus.medias)
      } else {
        console.warn('[SSE] 完了イベント - フォールバック実行条件を満たしていません:', {
          lastStatus: !!lastStatus,
          allCompleted: lastStatus?.allCompleted,
          all_completed: lastStatus?.all_completed,
          isAllCompleted,
          onComplete: !!onCompleteRef.current,
        })
      }

      eventSource.close()
      setIsConnected(false)
    })

    eventSource.onerror = error => {
      console.log('[SSE] onerror 判定変数:', {
        isCompletedRef: isCompletedRef.current,
        hasReceivedMessageRef: hasReceivedMessageRef.current,
        readyState: eventSource.readyState,
        readyStateName:
          eventSource.readyState === EventSource.CONNECTING
            ? 'CONNECTING'
            : eventSource.readyState === EventSource.OPEN
              ? 'OPEN'
              : 'CLOSED',
      })
      setIsConnected(false)

      // 正常完了後のエラーイベントは無視する。
      // 判定に isCompletedRef だけを使うと、close() がフラグ設定より先に
      // onerror を発火させるタイミング問題で漏れることがある。
      // そのため「CLOSED かつ既にメッセージを受信した」も正常終了と判定する。
      // [強化] allCompleted を受信済みの場合も正常終了とみなす
      const lastStatus = lastReceivedStatusRef.current
      const lastStatusAllCompleted = lastStatus?.allCompleted || lastStatus?.all_completed || false

      console.log('[SSE] onerror - フォールバック判定詳細:', {
        lastStatus: !!lastStatus,
        allCompleted: lastStatus?.allCompleted,
        all_completed: lastStatus?.all_completed,
        lastStatusAllCompleted,
        isCompleted: isCompletedRef.current,
        onCompleteExists: !!onCompleteRef.current,
      })

      const isNormalClose =
        isCompletedRef.current ||
        lastStatusAllCompleted ||
        (hasReceivedMessageRef.current && eventSource.readyState === EventSource.CLOSED)

      if (isNormalClose) {
        console.log('[SSE] 正常完了後の接続クローズを検出、エラーコールバックはスキップ')

        // フォールバック: onComplete が未呼び出しの場合、ここで呼ぶ
        const isAllCompleted = lastStatus?.allCompleted || lastStatus?.all_completed
        if (!isCompletedRef.current && isAllCompleted && onCompleteRef.current) {
          console.log('[SSE] onerror内でフォールバック完了処理を実行')
          console.log('[SSE] onerrorフォールバック - medias:', lastStatus.medias)
          isCompletedRef.current = true
          queryClient.invalidateQueries({ queryKey: MEDIA_QUERY_KEYS.images() })
          onCompleteRef.current(lastStatus.medias)
        }
      } else if (!hasReceivedMessageRef.current && eventSource.readyState === EventSource.CLOSED) {
        console.error('[SSE] 接続エラー:', error)
        console.error('[SSE] readyState:', eventSource.readyState)
        console.error('[SSE] hasReceivedMessage:', hasReceivedMessageRef.current)
        console.error('[SSE] isCompleted:', isCompletedRef.current)
        if (onErrorRef.current) {
          onErrorRef.current(
            new Error('SSE初期接続失敗: サーバーに接続できないか、認証エラーの可能性があります')
          )
        }
      } else if (onErrorRef.current) {
        onErrorRef.current(
          new Error(`SSE connection error (readyState: ${eventSource.readyState})`)
        )
      }

      eventSource.close()
      setIsConnected(false)
    }

    // クリーンアップ
    return () => {
      console.log('[SSE] クリーンアップ (readyState:', eventSource.readyState, ')')
      if (eventSource.readyState !== EventSource.CLOSED) {
        eventSource.close()
        setIsConnected(false)
      }
      // [修正3] enabledRef で現在の値を参照する
      // → クロージャで古い値を参照するバグを修正
      if (!enabledRef.current) {
        initialMediaIdsRef.current = []
        console.log('[SSE] 初回接続IDをリセット')
      }
    }
    // [修正1, 修正2] コールバックを依存配列から外し、mediaIdsKey を追加
  }, [enabled, mediaIdsKey, queryClient])

  return {
    status,
    isConnected,
    close: () => {
      if (eventSourceRef.current) {
        eventSourceRef.current.close()
        setIsConnected(false)
      }
    },
  }
}
