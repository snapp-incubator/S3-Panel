import { ThemeProvider } from '@/components/providers/themeProvider'
import type { ErrorPageProps } from './error.type'

export default function ErrorPage({
  title,
  subtitle,
  buttonText,
  buttonLink
}: ErrorPageProps) {
  return (
    <ThemeProvider defaultTheme="light" storageKey="vite-ui-theme">
      {title}
      <p>{subtitle}</p>
      <a href={buttonLink}>{buttonText}</a>
    </ThemeProvider>
  )
}
