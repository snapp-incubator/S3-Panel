import type { PropsWithChildren } from 'react'

import Header from '@/components/layout/header'
import { Toaster } from '@/components/shadcn/toaster'
import Sidebar from '@/components/side-navbar'

export default function Layout({ children }: PropsWithChildren) {
  return (
    <>
      <Header />
      <div className="flex h-screen border-collapse overflow-hidden pt-10">
        <Sidebar />
        <main className="flex-1 overflow-y-auto overflow-x-hidden bg-secondary/10 px-5 pb-1 pt-10">
          {children}
          <Toaster />
        </main>
      </div>
    </>
  )
}
