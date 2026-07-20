import type { LucideIcon } from 'lucide-react'

// Feature flag that gates a nav item. Currently only object storage ('s3').
type TFeature = 's3'

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
