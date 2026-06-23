import clsx from 'clsx'
import { type LucideIcon } from 'lucide-react'

import { Button } from '@/components/shadcn/button'
import { Card, CardContent } from '@/components/shadcn/card'
import { t } from '@/i18n'
import { formatNumber, scrollToSection } from '@/lib/utils'

interface ITotalCostCardItem {
  title: string
  total: number
  icon: LucideIcon
  default?: boolean
  rawTotal?: number
  maintenance?: number
}

interface ITotalCostCardProps {
  item: ITotalCostCardItem
}

export default function TotalCostCard({ item }: ITotalCostCardProps) {
  return (
    <Card
      className={clsx(
        'rounded-xl border-l-4',
        item.default
          ? 'border-l-green-500'
          : 'border-l-black dark:border-l-slate-700'
      )}
    >
      <CardContent className="py-7">
        <div className="flex items-start justify-between">
          <div className="flex flex-col">
            <h3 className="text-xl font-semibold text-gray-500">
              {item.title}
            </h3>
            <span className="text-2xl font-bold">
              {formatNumber(item.total, 'T')}
            </span>
          </div>
          <div
            className={clsx(
              'rounded-full  p-3',
              item.default
                ? 'bg-gray-100 text-green-500 dark:bg-green-500 dark:text-white'
                : 'bg-gray-100 text-current dark:bg-slate-700 dark:text-white'
            )}
          >
            <item.icon />
          </div>
        </div>
        <div>
          {item.default ? (
            <div className="mt-6 flex items-center gap-4">
              <div className="flex flex-col gap-1">
                <span className="text-sm font-semibold text-gray-600 dark:text-gray-400">
                  {t('total_raw')}
                </span>
                <span className="text-base font-semibold text-gray-700 dark:text-gray-300">
                  {formatNumber(item.rawTotal!, 'T')}
                </span>
              </div>
              <div className="flex flex-col gap-1">
                <span className="text-sm font-semibold text-gray-600 dark:text-gray-400">
                  {t('maintenance')}
                </span>
                <span className="text-base font-semibold text-gray-700 dark:text-gray-300">
                  {formatNumber(item.maintenance!, 'T')}
                </span>
              </div>
            </div>
          ) : (
            <div className="mt-6 w-full">
              <Button
                size="lg"
                className="w-full bg-green-600 text-white dark:bg-slate-700"
                onClick={() => scrollToSection('cost-breakdown')}
              >
                {t('more_details')}
              </Button>
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  )
}
