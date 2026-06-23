import { Menu, X } from 'lucide-react'

import { t } from '@/i18n'
import { cn } from '@/lib/utils'

type SidebarToggleHoverProps = {
  isOpen: boolean
  onToggle: () => void
  className?: string
}

export const SidebarToggleHover = ({
  isOpen,
  onToggle,
  className
}: SidebarToggleHoverProps) => {
  return (
    <div className="group relative">
      <div
        className={cn(
          'absolute -right-4 top-0 h-16 w-8 cursor-pointer',
          'hover:bg-accent/20 transition-colors duration-200'
        )}
        onClick={onToggle}
      />
      <button
        type="button"
        className={cn(
          'absolute -right-3 top-6 z-50 flex h-6 w-6 items-center justify-center rounded-full',
          'bg-background border border-border/50 shadow-sm',
          'opacity-60 group-hover:opacity-100 transition-all duration-200',
          'hover:scale-110 hover:shadow-md hover:border-border',
          'focus:opacity-100 focus:outline-none focus:ring-2 focus:ring-primary/20',
          className
        )}
        onClick={onToggle}
        aria-label={isOpen ? t('close_sidebar') : t('open_sidebar')}
      >
        {isOpen ? (
          <X className="size-3 text-muted-foreground" />
        ) : (
          <Menu className="size-3 text-muted-foreground" />
        )}
      </button>
    </div>
  )
}
