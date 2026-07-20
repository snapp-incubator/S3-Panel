import { TriangleAlert } from 'lucide-react'
import { cloneElement, type ReactElement } from 'react'

import { t } from '@/i18n'

const DEFAULT_ICON = <TriangleAlert color="red" />

interface ErrorStateProps {
  message?: string
  icon?: ReactElement<{ color?: string }>
  iconColor?: string
}

export default function ErrorState({
  message = t('data_empty_204'),
  icon = DEFAULT_ICON,
  iconColor = 'red'
}: ErrorStateProps) {
  return (
    <div className="mt-2 flex flex-col items-center justify-center px-2">
      {cloneElement(icon, { color: iconColor })}
      <h4 className="mt-1 text-center">{message}</h4>
    </div>
  )
}
