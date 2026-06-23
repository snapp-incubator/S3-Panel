import { Link } from '@tanstack/react-router'

import { buttonVariants } from '@/components/shadcn/button'
import type { NavItem } from '@/constants/navItems.types'
import { cn } from '@/lib/utils'

type NavbarItemProps = {
  item: NavItem
  isOpen: boolean
}

export default function NavbarItem({ item, isOpen }: NavbarItemProps) {
  return (
    <Link
      to={item.href}
      className={cn(
        buttonVariants({ variant: 'ghost' }),
        'group relative flex justify-start items-center transition-all duration-200',
        isOpen ? 'h-12 gap-4 px-4' : 'h-16 flex-col gap-1 px-2'
      )}
    >
      <item.icon className={cn('h-5 w-5 flex-shrink-0', item.color)} />
      <div
        className={cn(
          'transition-all duration-200',
          isOpen
            ? 'absolute left-12 text-base'
            : 'text-xs text-center leading-tight max-w-[60px] truncate'
        )}
      >
        {isOpen ? item.title : item.shortName}
      </div>
    </Link>
  )
}
