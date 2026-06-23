import { Info } from 'lucide-react'

import {
  Card,
  CardContent,
  CardHeader,
  CardTitle
} from '@/components/shadcn/card'
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger
} from '@/components/shadcn/tooltip'

import ErrorState from '../error-state'

import { TotalCostCardProps } from './cost.types'

const CPUCostCard = ({
  title,
  tooltipContent,
  unit = '',
  amount,
  hasError
}: TotalCostCardProps) => {
  return (
    <Card>
      <CardHeader className="mb-3 flex flex-row items-center justify-between space-y-0 p-0 pb-2 pl-5 pt-4">
        <CardTitle className="flex items-center gap-1 p-0 text-sm font-medium">
          {tooltipContent ? (
            <TooltipProvider>
              <Tooltip>
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
          <span className="text-1xl mr-1">{title}</span>
        </CardTitle>
      </CardHeader>
      <CardContent className="flex justify-around p-0">
        {hasError ? (
          <ErrorState />
        ) : (
          <>
            <div className="text-1xl">
              <div className="flex items-center justify-start">
                <div className="mr-1 font-bold">{amount.request.title}</div>
                {amount.request.help ? (
                  <TooltipProvider delayDuration={50}>
                    <Tooltip delayDuration={50}>
                      <TooltipTrigger asChild>
                        <Info className="size-4 cursor-help text-muted-foreground" />
                      </TooltipTrigger>
                      <TooltipContent>
                        <p>{amount.request.help}</p>
                      </TooltipContent>
                    </Tooltip>
                  </TooltipProvider>
                ) : (
                  ''
                )}
              </div>
              <p className="font-bold">{`${amount.request.value} ${unit}`}</p>
            </div>
            <div className="text-1xl">
              <div className="flex items-center">
                <div className="mr-1 font-bold">{amount.limit.title}</div>
                {amount.limit.help ? (
                  <TooltipProvider>
                    <Tooltip>
                      <TooltipTrigger asChild>
                        <Info className="size-4 cursor-help text-muted-foreground" />
                      </TooltipTrigger>
                      <TooltipContent>
                        <p>{amount.limit.help}</p>
                      </TooltipContent>
                    </Tooltip>
                  </TooltipProvider>
                ) : (
                  ''
                )}
              </div>
              <p className="font-bold">{`${amount.limit.value} ${unit}`}</p>
            </div>
          </>
        )}
      </CardContent>
    </Card>
  )
}

export default CPUCostCard
