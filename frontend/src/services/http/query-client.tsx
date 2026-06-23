import * as Sentry from '@sentry/react'
import { MutationCache, QueryCache, QueryClient } from '@tanstack/react-query'

function toastErrorHandler(error: Error) {
  Sentry.captureException(error)

  return Error
}

export const queryClient = new QueryClient({
  queryCache: new QueryCache({
    onError: toastErrorHandler
  }),
  mutationCache: new MutationCache({
    onError: toastErrorHandler
  }),
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      retry: 1
    }
  }
})

export const buildQueryString = (
  params: Record<string, string | undefined | number>
) =>
  Object.entries(params)
    .filter(([, v]) => v !== undefined)
    .map(([k, v]) => `${k}=${encodeURIComponent(v!)}`)
    .join('&')
