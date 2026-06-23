import { useState } from 'react'

import { NavItems } from '@/constants/navItems'
import { useSidebarStore } from '@/hooks/useSidebarStore'
import { cn } from '@/lib/utils'

import { SideNav } from './navbar-items'
import { SidebarToggle } from './sidebar-toggle'

type SidebarProps = {
  className?: string
}

export default function Sidebar({ className }: SidebarProps) {
  const { isOpen, toggle } = useSidebarStore()
  const [status, setStatus] = useState(false)

  const handleToggle = () => {
    setStatus(true)
    toggle()
    setTimeout(() => setStatus(false), 500)
  }

  return (
    <nav
      className={cn(
        `relative hidden h-screen bg-background text-3xl text-foreground border-r pt-0 md:block transition-all duration-300`,
        status && 'duration-500',
        isOpen ? 'w-72' : 'w-[78px]',
        className
      )}
    >
      <div className="relative mt-2 flex h-5 items-center justify-between  px-4">
        <SidebarToggle isOpen={isOpen} onToggle={handleToggle} />
      </div>
      <div className="flex-1 overflow-y-auto py-4">
        <div
          className={cn(
            'transition-all duration-300',
            isOpen ? 'px-3' : 'px-1'
          )}
        >
          <SideNav items={NavItems} />
        </div>
      </div>

      <div
        className={cn(
          'border-t p-4 transition-opacity duration-200',
          isOpen && 'opacity-0 pointer-events-none'
        )}
      >
        <div className="flex justify-center">
          <div className="size-2 rounded-full bg-muted-foreground/30" />
        </div>
      </div>
    </nav>
  )
}
