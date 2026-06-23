import { createKeyStore } from '@/api/query-keys'

export const baremetalKeys = createKeyStore('baremetals', {
  tenants: () => [],
  baremetalQuota: (venture: string) => [venture]
})
