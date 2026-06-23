import { AlertCircle } from 'lucide-react'

import { Alert, AlertDescription, AlertTitle } from '@/components/shadcn/alert'

import type { MessageAlertProps } from './message.types'

export const AlertMessage = ({
  title,
  message,
  variant
}: MessageAlertProps) => {
  return (
    <Alert variant={variant}>
      <AlertCircle className="size-4" />
      <AlertTitle>{title}</AlertTitle>
      <AlertDescription>{message}</AlertDescription>
    </Alert>
  )
}
