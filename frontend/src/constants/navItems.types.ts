import type { LucideIcon } from 'lucide-react'

import type { TFeature } from '@/services/http'

export type NavItem = {
  title: string
  href: string
  icon: LucideIcon
  feature: TFeature
  shortName: string
  color?: string
  isChidren?: boolean
  children?: NavItem[]
}
