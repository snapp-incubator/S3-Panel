import { createKeyStore } from '@/api/query-keys'

export const billingKeys = createKeyStore('billing', {
  totalCost: (team: string, date: string, region: string) => [
    team,
    date,
    region
  ],
  teams: () => [],
  trends: (team: string, date: string, region: string) => [team, date, region],
  costBreakdown: (team: string, date: string, region: string) => [
    team,
    date,
    region
  ]
})
