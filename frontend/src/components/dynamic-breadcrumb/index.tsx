import { Fragment } from 'react'

import { Link } from '@tanstack/react-router'

import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator
} from '@/components/shadcn/breadcrumb'

type CrumbItem = {
  type: 'link' | 'page'
  label: string
  to?: string
}

type DynamicBreadcrumbProps = {
  items: CrumbItem[]
  className?: string
}

export const DynamicBreadcrumb = ({
  items,
  className = ''
}: DynamicBreadcrumbProps) => {
  return (
    <Breadcrumb className={className}>
      <BreadcrumbList>
        {items.map((item, index) => {
          const key = `${item.to ?? ''}-${item.label}-${index}`

          return (
            <Fragment key={key}>
              <BreadcrumbItem>
                {item.type === 'link' && item.to ? (
                  <BreadcrumbLink>
                    <Link to={item.to}>{item.label}</Link>
                  </BreadcrumbLink>
                ) : (
                  <BreadcrumbPage>{item.label}</BreadcrumbPage>
                )}
              </BreadcrumbItem>
              {index < items.length - 1 && <BreadcrumbSeparator />}
            </Fragment>
          )
        })}
      </BreadcrumbList>
    </Breadcrumb>
  )
}
