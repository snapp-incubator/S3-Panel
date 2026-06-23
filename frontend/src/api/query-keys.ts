import { QueryKey } from '@tanstack/react-query'

type FunctionQueryKey = (...args: any[]) => QueryKey

type KeysObject = Record<string, FunctionQueryKey>

export const createKeyStore = <B extends string, K extends KeysObject>(
  baseName: B,
  keysObject: K
) => {
  const newKeysObject: Partial<K> = {}

  Object.keys(keysObject).forEach(key => {
    newKeysObject[key as keyof K] = (...args: any[]) => {
      let filteredKeys: QueryKey = []

      if (keysObject[key]) {
        const keys = keysObject[key](...args)

        filteredKeys = keys.filter(k => k !== undefined && k !== null)
      }

      return [baseName, key, ...filteredKeys]
    }
  })

  return newKeysObject as K
}
