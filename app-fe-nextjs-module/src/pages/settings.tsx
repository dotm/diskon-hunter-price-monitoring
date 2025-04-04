import Settings from '@/components/Settings'
import Sidebar from '@/components/Sidebar'
import { refreshIfNewAppVersionAvailable } from '@/utils/appversionutil'
import { LocalStorageKey } from '@/utils/constants'
import { useRouter } from 'next/router'
import { useEffect } from 'react'
import useLocalStorage from 'use-local-storage'

export default function SettingsPage() {
  const router = useRouter()
  const [appVersion, setAppVersion] = useLocalStorage<string>(LocalStorageKey.appVersion, "")

  useEffect(() => {
    refreshIfNewAppVersionAvailable(appVersion, setAppVersion, router)
  })

  return (
    <Sidebar pageTitle='Pengaturan' navTitle='Pengaturan'>
      <Settings />
    </Sidebar>
  )
}
