import { PanelLeftClose, PanelLeftOpen } from 'lucide-react'

import { t } from '@/i18n'
import { cn } from '@/lib/utils'

type SidebarToggleFloatingProps = {
  isOpen: boolean
  onToggle: () => void
  className?: string
}

export const SidebarToggleFloating = ({
  isOpen,
  onToggle,
  className
}: SidebarToggleFloatingProps) => {
  return (
    <button
      type="button"
      className={cn(
        'fixed left-4 top-4 z-50 flex h-10 w-10 items-center justify-center rounded-lg',
        'bg-background/80 backdrop-blur-sm border border-border/50 shadow-lg',
        'transition-all duration-200 hover:scale-105 hover:shadow-xl',
        'hover:bg-background hover:border-border',
        'focus:outline-none focus:ring-2 focus:ring-primary/20',
        // Hide when sidebar is open and positioned normally
        isOpen && 'md:opacity-0 md:pointer-events-none',
        className
      )}
      onClick={onToggle}
      aria-label={isOpen ? t('close_sidebar') : t('open_sidebar')}
    >
      {isOpen ? (
        <PanelLeftClose className="h-5 w-5 text-foreground" />
      ) : (
        <PanelLeftOpen className="h-5 w-5 text-foreground" />
      )}
    </button>
  )
}
