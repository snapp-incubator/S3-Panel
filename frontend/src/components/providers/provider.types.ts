import { Dispatch, SetStateAction } from 'react'
export type Theme = 'dark' | 'light'

export type ThemeProviderProps = {
  children: React.ReactNode
  defaultTheme?: Theme
  storageKey?: string
}

export type ThemeProviderState = {
  theme: Theme
  setTheme: (theme: Theme) => void
}

export type TitleContextType = {
  title: string
  setTitle: Dispatch<SetStateAction<string>>
}
