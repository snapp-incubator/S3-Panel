import { Info } from 'lucide-react'

import {
  Card,
  CardHeader,
  CardTitle,
  CardContent
} from '@/components/shadcn/card'
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger
} from '@/components/shadcn/tooltip'

import ErrorState from '../error-state'

import { TotalCostCardProps } from './cost.types'

const TotalCostCard = ({
  title,
  unit = '',
  amount,
  tooltipContent,
  hasError
}: TotalCostCardProps) => {
  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="flex items-center gap-2 text-sm font-medium">
          {title}
          {tooltipContent ? (
            <TooltipProvider delayDuration={50}>
              <Tooltip delayDuration={50}>
                <TooltipTrigger asChild>
                  <Info className="size-4 cursor-help text-muted-foreground" />
                </TooltipTrigger>
                <TooltipContent>
                  <p className="max-w-prose p-0 text-base leading-relaxed text-gray-700">
                    {tooltipContent}
                  </p>
                </TooltipContent>
              </Tooltip>
            </TooltipProvider>
          ) : (
            ''
          )}
        </CardTitle>
      </CardHeader>
      <CardContent className="flex items-center justify-center">
        {hasError ? (
          <ErrorState />
        ) : (
          <div className="text-2xl font-bold">{`${amount} ${unit}`}</div>
        )}
      </CardContent>
    </Card>
  )
}

export default TotalCostCard
