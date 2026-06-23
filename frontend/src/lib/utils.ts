import { redirect } from '@tanstack/react-router'
import { type ClassValue, clsx } from 'clsx'
import { twMerge } from 'tailwind-merge'

import useS3Credentials from '@/hooks/useS3Credentials'

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function calculateValue(usedValue: number, hardValue: number) {
  return (usedValue * 100) / hardValue
}

export function copyString(key: string) {
  try {
    navigator.clipboard.writeText(key)

    return Promise.resolve(key)
  } catch {
    return Promise.reject()
  }
}

/**
 * Formats a Date object into a string with format "YYYY/MM/DD HH:MM:SS" in UTC timezone
 * @param date - The Date object to format
 * @returns A formatted date-time string
 */
export function formatFullDateTime(date: Date): string {
  return new Intl.DateTimeFormat('en-CA', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: false,
    timeZone: 'UTC'
  })
    .format(date)
    .replace(/-/g, '/')
}

export function dateFormat(dateStr: string) {
  try {
    const normalized = dateStr.replace(' ', 'T').replace(' +0000 UTC', 'Z')
    const date = new Date(normalized)

    if (isNaN(date.getTime())) throw new Error('Invalid date')

    return formatFullDateTime(date)
  } catch (err) {
    return 'Invalid Date'
  }
}

export function handleAuthRedirect(to: string, shouldBeAuthenticated = true) {
  const isLogin = useS3Credentials.getState().isLogin()

  if (shouldBeAuthenticated && !isLogin) {
    throw redirect({ to }) // Redirect if user is not logged in but should be
  }

  if (!shouldBeAuthenticated && isLogin) {
    throw redirect({ to }) // Redirect if user is logged in but shouldn't be
  }
}

export function changeByteToMegabyte(byte: number) {
  return Number((byte / 1000000000).toFixed(2))
}

export function downloadUrl(url: string) {
  const link = document.createElement('a')

  link.href = url
  link.target = '_blank'
  link.click()
  window.URL.revokeObjectURL(url)
}

/**
 * Format text and show first 100 chars + ... for long texts
 * @param text - The Date object to format
 * @returns Short text plus ...
 */
export function shortText(text: string) {
  if (text.length > 100) {
    text = text.slice(0, 100)

    return `${text}...`
  }
}

/**
 * Format Toman
 * @param number - 1000000
 * @returns 1,000,000 T
 */
export const formatNumber = (value: number, prefix?: string) =>
  `${value.toLocaleString()} ${prefix || ''}`

/**
 * Scroll to section
 * @param id - 'cost-breakdown'
 * @returns void
 */
export const scrollToSection = (id: string) => {
  const element = document.getElementById(id)

  if (element) {
    element.scrollIntoView({ behavior: 'smooth', block: 'start' })
  }
}
