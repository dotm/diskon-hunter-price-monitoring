import SearchItemEditItems from '@/components/SearchItemEditItems'
import Sidebar from '@/components/Sidebar'
import { LocalStorageKey } from '@/utils/constants'
import useLocalStorage from 'use-local-storage'

export default function SearchItemEditItemsPage() {
  const [appVersion, setAppVersion] = useLocalStorage<string>(LocalStorageKey.appVersion, "")

  //refreshIfNewAppVersionAvailable useEffect is done inside the component

  return (
    <Sidebar pageTitle='Ubah Barang yang Dicari' navTitle='Cari Barang'>
      <SearchItemEditItems appVersion={appVersion} setAppVersion={setAppVersion} />
    </Sidebar>
  )
}
