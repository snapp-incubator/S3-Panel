import { PackageOpen } from 'lucide-react'

import type { NavItem } from './navItems.types'

export const NavItems: NavItem[] = [
  {
    title: 'Object Storage',
    shortName: 'S3',
    icon: PackageOpen,
    href: '/object-storage/s3-bucket',
    color: 'text-green-500',
    feature: 's3'
  }
]
