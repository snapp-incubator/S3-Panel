import { Skeleton } from '@/components/shadcn/skeleton'

interface IBucketCardSkeletonProps {
  count: number
}

export default function BucketCardSkeleton({
  count
}: IBucketCardSkeletonProps) {
  return new Array(count).fill('').map((_, index) => (
    <div key={index} className="flex h-[210px] w-full min-w-[400px]">
      <Skeleton className="w-[45%] rounded-lg" />
      <div className="flex w-[55%] flex-col gap-3 p-2">
        <Skeleton className="h-4 w-full rounded-lg" />
        <Skeleton className="h-4 w-1/2 rounded-lg" />
        <Skeleton className="h-4 w-full rounded-lg" />
        <Skeleton className="h-4 w-1/2 rounded-lg" />
        <Skeleton className="mt-auto h-12 w-full rounded-lg" />
      </div>
    </div>
  ))
}
