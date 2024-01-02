import { LocalStorageKey } from "@/utils/constants"
import { LoggedInUserData } from "@/utils/models"
import useLocalStorage from "use-local-storage"

export default function SignOutForm() {
  const [loggedInUserData, setLoggedInUserData] =
    useLocalStorage<LoggedInUserData | undefined>(LocalStorageKey.loggedInUser, undefined)
  
  return (
    <div className="flex min-h-full flex-1 flex-col justify-center">
      <div className="sm:mx-auto sm:w-full sm:max-w-[480px]">
        <div className="bg-gray-800 p-8 shadow sm:rounded-lg">
          <div className="space-y-6">
            <p className="text-center text-white">
              Halo, {loggedInUserData?.email ?? "Diskon Hunter"}!
              <br/>
              <br/>
              Silahkan mulai menggunakan aplikasi ini dengan meng-klik salah satu menu di samping kiri.
              <br/>
              <br/>
              Jika anda menggunakan handphone, silahkan meng-klik ☰ di kiri atas layar anda untuk melihat menu yang tersedia.
            </p>
            <div>
              <button
                onClick={event=>setLoggedInUserData(undefined)}
                className="disabled:bg-slate-600 flex w-full justify-center rounded-md bg-indigo-600 px-3 py-1.5 text-sm font-semibold leading-6 text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
              >
                Sign out
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
