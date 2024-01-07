import MonitorPriceAddLinks from '@/components/MonitorPriceAddLinks'
import Sidebar from '@/components/Sidebar'
import { LocalStorageKey } from '@/utils/constants'
import { useRouter } from 'next/router'
import useLocalStorage from 'use-local-storage'

export default function MonitorPriceAddLinksPage() {
  const router = useRouter()
  const [appVersion, setAppVersion] = useLocalStorage<string>(LocalStorageKey.appVersion, "")

  //refreshIfNewAppVersionAvailable useEffect is done inside the component

  return (
    <Sidebar pageTitle='Tambah Link Monitor Harga' navTitle='Monitor Harga'>
      <MonitorPriceAddLinks appVersion={appVersion} setAppVersion={setAppVersion} />
    </Sidebar>
  )
}
