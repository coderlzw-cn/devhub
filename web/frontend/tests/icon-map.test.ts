import { describe, expect, it } from 'vitest'
import * as lucide from 'lucide-vue-next'
import { resolveIcon } from '../src/utils/icon-map'

const databaseIconNames = [
  'Activity',
  'Binary',
  'Bot',
  'BookOpen',
  'ChartNoAxesCombined',
  'Clapperboard',
  'Cloud',
  'CloudUpload',
  'Code2',
  'CodeXml',
  'Cog',
  'Database',
  'GitBranch',
  'Globe',
  'GraduationCap',
  'Hammer',
  'Images',
  'KeyRound',
  'Mail',
  'MonitorSmartphone',
  'Network',
  'Palette',
  'Paintbrush',
  'PanelsTopLeft',
  'Presentation',
  'Server',
  'ServerCog',
  'Shapes',
  'ShieldCheck',
  'TestTube',
  'WandSparkles',
  'Webhook',
  'Workflow',
]

describe('database Lucide icon names', () => {
  it('resolves each configured name to its matching Lucide Vue component', () => {
    for (const name of databaseIconNames) {
      expect(resolveIcon(name), name).toBe(lucide[name as keyof typeof lucide])
    }
  })

  it('uses the category fallback for unknown names and URLs', () => {
    expect(resolveIcon('https://lucide.dev/icons/')).toBe(lucide.FolderTree)
    expect(resolveIcon('unknown-icon')).toBe(lucide.FolderTree)
  })
})
