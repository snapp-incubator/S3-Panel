import { createContext, useContext, useState, PropsWithChildren } from 'react'
import { TitleContextType } from './provider.types'

const TitleContext = createContext<TitleContextType>({
  title: '',
  setTitle: () => {}
})

export const useTitle = () => useContext(TitleContext)

export function TitleProvider({ children }: PropsWithChildren) {
  const [title, setTitle] = useState<string>('')

  return (
    <TitleContext.Provider value={{ title, setTitle }}>
      {children}
    </TitleContext.Provider>
  )
}
