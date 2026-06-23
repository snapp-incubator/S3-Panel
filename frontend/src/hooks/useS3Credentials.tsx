import { create } from 'zustand'
import { createJSONStorage, persist } from 'zustand/middleware'

import { s3Keys } from '@/api/s3Keys'
import { S3StorageKeys } from '@/constants/s3/storeKeys'
import { queryClient } from '@/services/http/query-client'
import type { TRegions } from '@/types/regions.types'

interface ICredential {
  access_key: string | null
  secret_key: string | null
}

type TFillCredentials = {
  fillCredentials: (credentials: ICredential, region?: TRegions) => void
}

type TFillRegion = {
  fillRegion: (region: TRegions) => void
}

type TResetData = {
  resetData: () => void
}

type TIsLogin = {
  isLogin: () => boolean
}

type TLogout = {
  logout: () => Promise<void>
}

type IUseS3Credentials = ICredential & {
  region: TRegions | null
} & TResetData &
  TIsLogin &
  TLogout &
  TFillCredentials &
  TFillRegion

const useS3Credentials = create(
  persist<IUseS3Credentials>(
    (set, get) => ({
      access_key: null,
      secret_key: null,
      region: null,
      fillRegion: region =>
        set(() => ({
          region
        })),
      fillCredentials: (credentials, region) => {
        set(() => ({
          access_key: credentials.access_key,
          secret_key: credentials.secret_key
        }))

        if (region) {
          get().fillRegion(region)
        }
      },
      isLogin: (): boolean => {
        const { access_key, secret_key } = get()

        if (access_key && secret_key) return true

        return false
      },
      resetData: () =>
        set(() => ({
          access_key: null,
          secret_key: null
        })),
      logout: () => {
        return new Promise<void>(resolve => {
          get().resetData()
          sessionStorage.removeItem(S3StorageKeys.s3_credentials)

          queryClient.removeQueries({
            queryKey: [s3Keys.user]
          })
          queryClient.removeQueries({
            queryKey: [s3Keys.objects]
          })
          queryClient.removeQueries({
            queryKey: [s3Keys.buckets]
          })

          resolve()
        })
      }
    }),
    {
      name: S3StorageKeys.s3_credentials,
      storage: createJSONStorage(() => sessionStorage)
    }
  )
)

export default useS3Credentials
