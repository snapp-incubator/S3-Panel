import {
  BrainCircuit,
  CloudCog,
  FolderGit2,
  Info,
  MonitorCog,
  SquareTerminal
} from 'lucide-react'

export const serviceIconsMap = {
  Info: {
    icon: Info,
    color: 'text-blue-800',
    bgColor: 'bg-blue-200'
  },
  Observability: {
    icon: MonitorCog,
    color: 'text-green-800',
    bgColor: 'bg-green-200'
  },
  'Developer Experience': {
    icon: SquareTerminal,
    color: 'text-purple-800',
    bgColor: 'bg-purple-200'
  },
  'Version Control System': {
    icon: FolderGit2,
    color: 'text-orange-800',
    bgColor: 'bg-orange-200'
  },
  'SnappCloud AI': {
    icon: BrainCircuit,
    color: 'text-indigo-800',
    bgColor: 'bg-indigo-200'
  },
  Infrastructure: {
    icon: CloudCog,
    color: 'text-red-800',
    bgColor: 'bg-red-200'
  }
}
