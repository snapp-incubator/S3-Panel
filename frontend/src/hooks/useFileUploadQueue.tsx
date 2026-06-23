import { useEffect, useCallback } from 'react'

import type { UseMutateAsyncFunction } from '@tanstack/react-query'

import { useUploadProgress } from '@/components/providers/uploadProgressContext'
import { HTTPClientError } from '@/services/http/interceptorsConfig'
import { IUploadNames } from '@/types/s3/upload.types'

interface IUseFileUploadQueueProps {
  mutateAsync: UseMutateAsyncFunction<
    {
      created: boolean
    },
    HTTPClientError<{
      message?: string
    }>,
    File,
    unknown
  >
}

export default function useFileUploadQueue({
  mutateAsync
}: IUseFileUploadQueueProps) {
  const {
    uploadNames,
    setShowUploadModal,
    setUploadNames,
    processing,
    setProcessing
  } = useUploadProgress()

  const runQueue = useCallback(
    async (pending: IUploadNames[]) => {
      setProcessing(true)
      setShowUploadModal(true)

      let index = 0
      const concurrency = 1
      const runNext = async () => {
        if (index >= pending.length) return

        const currentItem = pending[index++]

        try {
          await mutateAsync(currentItem.file!)
          setUploadNames(prev =>
            prev.map(item =>
              item.name === currentItem.name
                ? { ...item, completed: true }
                : item
            )
          )
        } catch (err) {
          setUploadNames(prev =>
            prev.map(item =>
              item.name === currentItem.name ? { ...item, failed: true } : item
            )
          )
        } finally {
          await runNext()
        }
      }
      const workers = []

      for (let i = 0; i < concurrency; i++) {
        workers.push(runNext())
      }

      await Promise.all(workers)
      setProcessing(false)
      setShowUploadModal(false)
    },
    [mutateAsync, setShowUploadModal, setUploadNames, setProcessing]
  )

  useEffect(() => {
    // Clear if all uploads are finished
    if (
      uploadNames.length &&
      uploadNames.every(item => item.canceled || item.completed || item.failed)
    ) {
      setUploadNames([])

      return
    }

    if (processing || uploadNames.length === 0) return

    const pending = uploadNames.filter(
      item => !item.canceled && !item.completed && !item.failed
    )

    if (pending.length === 0) return

    runQueue(pending)
  }, [
    uploadNames,
    processing,
    mutateAsync,
    setUploadNames,
    setShowUploadModal,
    runQueue
  ])
}
