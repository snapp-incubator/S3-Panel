import { create } from 'zustand'
import { persist } from 'zustand/middleware'

interface SidebarState {
  isOpen: boolean
  toggle: () => void
}

export const useSidebarStore = create<SidebarState>()(
  persist(
    set => ({
      isOpen: false,
      toggle: () => set(state => ({ isOpen: !state.isOpen }))
    }),
    {
      name: 'sidebar-storage'
    }
  )
)
