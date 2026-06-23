import { RotateCcw } from 'lucide-react'

import { Alert, AlertDescription, AlertTitle } from '@/components/shadcn/alert'
import { Button } from '@/components/shadcn/button'
import { t } from '@/i18n'

interface IFailedAlertProps {
  refetch: () => void
}

export default function FailedAlert({ refetch }: IFailedAlertProps) {
  return (
    <Alert variant="destructive" className="flex items-center justify-between">
      <div className="flex flex-col gap-1">
        <AlertTitle>{t('error')}</AlertTitle>
        <AlertDescription>{t('failed_load_total_costs')}</AlertDescription>
      </div>
      <Button
        variant="outline"
        className="flex gap-2"
        onClick={() => refetch()}
      >
        <RotateCcw size={20} />
        {t('try_again')}
      </Button>
    </Alert>
  )
}
