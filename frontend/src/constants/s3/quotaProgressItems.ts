// import { changeByteToMegabyte } from '@/lib/utils'
import { t } from '@/i18n'
import type {
  TProgressHardKeys,
  TProgressItemsKeys,
  TProgressUsedKeys
} from '@/types/s3/userQuota.types'

type TProgressItem = {
  key: TProgressItemsKeys
  title: string
  hardKey: TProgressHardKeys
  usedKey: TProgressUsedKeys
}

export const progressItems: TProgressItem[] = [
  {
    title: t('bucket_usage'),
    key: 'buckets',
    hardKey: 'hard_buckets',
    usedKey: 'used_buckets'
  },
  {
    title: t('object_usage'),
    key: 'objects',
    hardKey: 'hard_objects',
    usedKey: 'used_objects'
  }
]
