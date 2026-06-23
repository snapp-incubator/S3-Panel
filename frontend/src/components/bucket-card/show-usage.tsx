import { Progress } from '@/components/shadcn/progress'
import { calculateValue } from '@/lib/utils'

import type { TShowUsageProps } from './showUsage.types'

export default function ShowUsage({
  quotaEnabled,
  hardData,
  usedData,
  hardUnit = '',
  usedUnit = '',
  hardRaw,
  usedRaw,
  withUnit
}: TShowUsageProps) {
  const progressValue = withUnit
    ? calculateValue(usedRaw!, hardRaw!)
    : calculateValue(usedData, hardData)

  return (
    <>
      {quotaEnabled ? (
        <>
          <Progress
            value={progressValue}
            className="max-h-3 w-full text-green-700"
          />
          <span className="ml-1 text-sm">
            {`${usedData} ${usedUnit} of ${hardData} ${hardUnit}`}
          </span>
        </>
      ) : (
        <span className="text-base font-bold">{`${usedData} ${usedUnit}`}</span>
      )}
    </>
  )
}
