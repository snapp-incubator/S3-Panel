import dayjs, { extend } from 'dayjs'
import jalaliday from 'jalaliday'

extend(jalaliday)

type TMonths = {
  title: string
  monthDate: string
  default?: boolean
}

/**
 * Get last 6 months in Gregorian calendar
 * Example:
 * [
 *   { title: "September", monthDate: "2025-09", default: true },
 *   { title: "August", monthDate: "2025-08", default: false },
 *   ...
 * ]
 */
const getLastMonths = (countLastMonth = 6): TMonths[] => {
  const now = dayjs()

  return Array.from({ length: countLastMonth }, (_, index) => ({
    title: now.subtract(index, 'month').format('MMMM'),
    monthDate: now.subtract(index, 'month').format('YYYY-MM'),
    default: index === 0
  }))
}

export const dateUtils = {
  getLastMonths
}
