import { ChevronLeft, ChevronRight } from 'lucide-react'

import { t } from '@/i18n'
import { cn } from '@/lib/utils'

type SidebarToggleProps = {
  isOpen: boolean
  onToggle: () => void
  className?: string
}

export const SidebarToggle = ({
  isOpen,
  onToggle,
  className
}: SidebarToggleProps) => {
  return (
    <button
      type="button"
      className={cn(
        'absolute -right-3 top-6 z-50 flex h-6 w-6 items-center justify-center rounded-full border bg-background shadow-md transition-all duration-200 hover:scale-110 hover:shadow-lg',
        'border-border/50 hover:border-border',
        className
      )}
      onClick={onToggle}
      aria-label={isOpen ? t('close_sidebar') : t('open_sidebar')}
    >
      {isOpen ? (
        <ChevronLeft className="size-3 text-muted-foreground" />
      ) : (
        <ChevronRight className="size-3 text-muted-foreground" />
      )}
    </button>
  )
}
