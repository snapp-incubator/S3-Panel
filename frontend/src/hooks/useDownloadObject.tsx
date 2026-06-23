import { useMutation } from '@tanstack/react-query'

import { downloadObject } from '@/api/s3'
import { downloadUrl } from '@/lib/utils'
import { HTTPClientError } from '@/services/http/interceptorsConfig'

import { useToast } from './use-toast'

const useDownloadObject = (downloadLink: string) => {
  const { toast } = useToast()

  return useMutation({
    mutationFn: () => downloadObject(downloadLink),
    onSuccess: res => {
      const url = res.url

      downloadUrl(url)
    },
    onError: (err: HTTPClientError<{ message?: string }>) => {
      toast({
        variant: 'destructive',
        title: err.response?.data.message || err.message
      })
    }
  })
}

export default useDownloadObject
