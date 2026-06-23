import { useQuery } from '@tanstack/react-query'
import { useNavigate } from '@tanstack/react-router'
import { CircleUserRound, ChevronDown, Loader2 } from 'lucide-react'

import { fetchUserDetails } from '@/api/s3'
import { userKeys } from '@/api/s3Keys'
import { Button } from '@/components/shadcn/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from '@/components/shadcn/dropdown-menu'
import useS3Credentials from '@/hooks/useS3Credentials'

export default function UserDetails() {
  const { logout } = useS3Credentials()
  const navigate = useNavigate()

  const { data: user, isPending } = useQuery({
    queryKey: userKeys.details(),
    queryFn: fetchUserDetails
  })

  const logoutHandler = () => {
    logout().then(() => {
      navigate({
        to: '/object-storage/s3-bucket'
      })
    })
  }

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button
          variant="outline"
          disabled={isPending}
          className="flex items-center gap-2"
        >
          {isPending ? (
            <Loader2 className="animate-spin" />
          ) : (
            <>
              <CircleUserRound /> {user?.display_name} <ChevronDown />
            </>
          )}
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent>
        <DropdownMenuItem asChild>
          <Button
            variant="ghost"
            size="sm"
            className="w-full cursor-pointer text-red-600"
            onClick={logoutHandler}
          >
            Logout
          </Button>
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
