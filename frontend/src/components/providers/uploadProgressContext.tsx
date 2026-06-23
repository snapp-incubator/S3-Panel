import {
  createContext,
  type MutableRefObject,
  type ReactNode,
  useCallback,
  useContext,
  useEffect,
  useRef,
  useState
} from 'react'

import { t } from '@/i18n'
import type { IUploadNames } from '@/types/s3/upload.types'

interface UploadProgressContextProps {
  uploadNames: IUploadNames[]
  setUploadNames: (
    names: IUploadNames[] | ((prev: IUploadNames[]) => IUploadNames[])
  ) => void
  showUploadModal: boolean
  setShowUploadModal: (arg0: boolean) => void
  abortController: MutableRefObject<AbortController | null>
  updateUploadNamesState: (
    fileName: string,
    newObj: Partial<IUploadNames>
  ) => void
  processing: boolean
  setProcessing: (value: boolean) => void
}

const UploadProgressContext = createContext<
  UploadProgressContextProps | undefined
>(undefined)

export const UploadProgressProvider = ({
  children
}: {
  children: ReactNode
}) => {
  const [uploadNames, setUploadNames] = useState<IUploadNames[]>([])
  const abortController = useRef<AbortController | null>(null)
  const [showUploadModal, setShowUploadModal] = useState(false)
  const [processing, setProcessing] = useState(false)

  const checkUploadProgress = useCallback(() => {
    return !!processing
  }, [processing])

  useEffect(() => {
    const handleBeforeUnload = (event: BeforeUnloadEvent) => {
      if (checkUploadProgress()) {
        event.preventDefault()
        event.returnValue = ''

        return ''
      }
    }

    window.addEventListener('beforeunload', handleBeforeUnload)

    return () => {
      window.removeEventListener('beforeunload', handleBeforeUnload)
    }
  }, [checkUploadProgress])

  useEffect(() => {
    return () => {
      if (abortController.current) {
        abortController.current.abort()
      }
    }
  }, [])

  const updateUploadNamesState = (
    fileName: string,
    newObj: Partial<IUploadNames>
  ) => {
    setUploadNames(prevItems => {
      if (newObj.canceled === true) {
        // Remove the item from the list when canceled
        return prevItems.filter(item => item.name !== fileName)
      }

      return prevItems.map(item => {
        if (item.name === fileName) {
          // If the item is already canceled and the update is not resetting cancel, ignore update.
          if (item.canceled && newObj.canceled !== false) return item

          return { ...item, ...newObj }
        }

        return item
      })
    })
  }

  return (
    <UploadProgressContext.Provider
      value={{
        uploadNames,
        setUploadNames,
        showUploadModal,
        setShowUploadModal,
        abortController,
        updateUploadNamesState,
        processing,
        setProcessing
      }}
    >
      {children}
    </UploadProgressContext.Provider>
  )
}

export const useUploadProgress = () => {
  const context = useContext(UploadProgressContext)
  const errorMessage = t('use_upload_context_error')

  if (context === undefined) {
    throw new Error(errorMessage)
  }

  return context
}
