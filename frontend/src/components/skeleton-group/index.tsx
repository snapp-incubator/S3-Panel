import { Skeleton } from '@/components/shadcn/skeleton'

type SkeletonGroupProps = {
  count?: number
  height?: number | string
  orientation?: 'horizontal' | 'vertical'
  className?: string
}

export default function SkeletonGroup({
  count = 1,
  height = 130,
  orientation = 'horizontal',
  className = ''
}: SkeletonGroupProps) {
  const heightValue = typeof height === 'number' ? `${height}px` : height

  const skeletons = Array.from({ length: count }, (_, index) => (
    <Skeleton
      key={index}
      className={`w-full ${className}`}
      style={{ height: heightValue }}
    />
  ))

  const containerClass =
    orientation === 'horizontal'
      ? `grid gap-4 md:grid-cols-${count <= 6 ? count : 6}`
      : 'flex flex-col gap-4'

  return (
    <div className={`my-3 w-full overflow-hidden ${containerClass}`}>
      {skeletons}
    </div>
  )
}
