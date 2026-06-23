import { PropsWithChildren } from 'react'

import { QueryClientProvider as ReactQueryClientProvider } from '@tanstack/react-query'

import { queryClient } from '@/services/http/query-client'

export default function QueryClientProvider({ children }: PropsWithChildren) {
  return (
    <ReactQueryClientProvider client={queryClient}>
      {children}
    </ReactQueryClientProvider>
  )
}
