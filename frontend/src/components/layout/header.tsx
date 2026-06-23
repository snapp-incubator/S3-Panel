import { Link, useMatchRoute } from '@tanstack/react-router'
import { Boxes } from 'lucide-react'

import { ThemeToggle } from '@/components/layout/theme-toggle'
import { useTitle } from '@/components/providers/titleProvider'
import useS3Credentials from '@/hooks/useS3Credentials'

import UserDetails from './user-details'

const Header = () => {
  const { title } = useTitle()
  const { isLogin } = useS3Credentials()

  const matchRoute = useMatchRoute()

  const isObjectStorageRoute = matchRoute({
    to: '/object-storage/s3-bucket',
    fuzzy: true
  })

  return (
    <header className="supports-backdrop-blur:bg-background/60 fixed inset-x-0 top-0 z-20 border-b bg-background/95 backdrop-blur">
      <nav className="flex h-16 items-center justify-between px-4">
        <Link
          to="/"
          className="hidden items-center justify-between gap-2 md:flex"
        >
          <Boxes className="size-6" />
          <div className="flex items-center justify-center gap-2">
            <h1 className="text-lg font-semibold">{title}</h1>
          </div>
        </Link>
        <div className="flex items-center gap-2">
          {isLogin() && isObjectStorageRoute ? <UserDetails /> : null}
          <ThemeToggle />
        </div>
      </nav>
    </header>
  )
}

export default Header
