import { Link } from '@tanstack/react-router'
import { ChevronRight } from 'lucide-react'
import { useRef, useState } from 'react'

import { buttonVariants } from '@/components/shadcn/button'
import type { NavItem } from '@/constants/navItems.types'
import { useSidebarStore } from '@/hooks/useSidebarStore'
import { cn } from '@/lib/utils'

import NavbarItem from './navbar-item'
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger
} from './subnav-accordion'

type SideNavProps = {
  items: NavItem[]
  setOpen?: (open: boolean) => void
  className?: string
}

export function SideNav({ items }: SideNavProps) {
  const { isOpen } = useSidebarStore()
  const [openItems, setOpenItems] = useState<string[]>([])
  const prevIsOpenRef = useRef(isOpen)

  if (isOpen !== prevIsOpenRef.current && isOpen) {
    prevIsOpenRef.current = isOpen
    setOpenItems([])
  }

  const handleAccordionChange = (values: string[]) => {
    setOpenItems(values)
  }

  return (
    <nav
      className={cn(
        'transition-all duration-200',
        isOpen ? 'space-y-2' : 'space-y-1'
      )}
    >
      {items.map(item =>
        item.isChidren ? (
          <Accordion
            type="multiple"
            className="space-y-2"
            key={item.title}
            value={openItems}
            onValueChange={handleAccordionChange}
          >
            <AccordionItem value={item.title} className="border-none">
              <AccordionTrigger
                className={cn(
                  buttonVariants({ variant: 'ghost' }),
                  'group relative flex justify-between py-2 text-base duration-200 hover:bg-muted hover:no-underline transition-all',
                  openItems.includes(item.title) ? 'bg-muted' : null,
                  isOpen ? 'h-12 px-4' : 'h-16 px-2 flex-col gap-1'
                )}
              >
                <div
                  className={cn(
                    'flex items-center transition-all duration-200',
                    isOpen ? 'gap-3' : 'flex-col gap-1'
                  )}
                >
                  <item.icon
                    className={cn('h-5 w-5 flex-shrink-0', item.color)}
                  />
                  <div
                    className={cn(
                      'transition-all duration-200',
                      isOpen
                        ? 'text-base'
                        : 'text-xs text-center leading-tight max-w-[60px] truncate'
                    )}
                  >
                    {isOpen ? item.title : item.shortName}
                  </div>
                </div>
                <div
                  className={cn(
                    'duration-200',
                    !isOpen && 'hidden',
                    openItems.includes(item.title) ? 'rotate-90' : 'rotate-0'
                  )}
                >
                  <ChevronRight className="size-4" />
                </div>
              </AccordionTrigger>
              <AccordionContent
                className={cn('mt-2 space-y-2 pb-1', isOpen ? 'pl-8' : 'pl-0')}
              >
                {item.children?.map(child => (
                  <Link
                    key={child.title}
                    to={child.href}
                    className={cn(
                      buttonVariants({ variant: 'ghost' }),
                      'group relative flex justify-start items-center transition-all duration-200',
                      isOpen ? 'h-12 gap-3 px-4' : 'h-14 flex-col gap-1 px-2'
                    )}
                  >
                    <child.icon
                      className={cn('h-4 w-4 flex-shrink-0', child.color)}
                    />
                    <div
                      className={cn(
                        'transition-all duration-200',
                        isOpen
                          ? 'absolute left-12 text-base'
                          : 'text-xs text-center leading-tight max-w-[60px] truncate'
                      )}
                    >
                      {isOpen ? child.title : child.shortName}
                    </div>
                  </Link>
                ))}
              </AccordionContent>
            </AccordionItem>
          </Accordion>
        ) : (
          <NavbarItem key={item.title} item={item} isOpen={isOpen} />
        )
      )}
    </nav>
  )
}
