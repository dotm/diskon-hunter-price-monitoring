import { LocalStorageKey, backendBaseUrl, backendHeadersForPostRequest } from "@/utils/constants"
import { handleErrorInFrontend } from "@/utils/error"
import { LoggedInUserData } from "@/utils/models"
import { FormEvent, useState } from "react"
import useLocalStorage from "use-local-storage"

export default function EditUserDataForm() {
  const [loggedInUserData, setLoggedInUserData] =
    useLocalStorage<LoggedInUserData | undefined>(LocalStorageKey.loggedInUser, undefined)
  //value other than password should use existing value from loggedInUserData
  const [password, setPassword] = useState("")
  const [loading, setLoading] = useState(false)
  async function interactor_userEdit(event: FormEvent<HTMLFormElement>){
    event.preventDefault()
    
    try {
      setLoading(true)
      if(loggedInUserData === undefined){
        throw new Error("Mohon sign in terlebih dahulu")
      }
      const editUserRespJson = await fetch(`${backendBaseUrl}/v1/user.edit`, {
        method: 'POST',
        headers: backendHeadersForPostRequest(loggedInUserData.jwt),
        body: JSON.stringify({
          password: password
        }),
      })
      .then(response => response.json())
      if(!editUserRespJson.ok || editUserRespJson.data === undefined){
        throw new Error(editUserRespJson.err?.code ?? "error editUserRespJson")
      }
      alert("Data user berhasil diubah.")
    } catch (error) {
      handleErrorInFrontend(error)
    } finally {
      setLoading(false)
    }
  }
  return (
    <div className="flex min-h-full flex-1 flex-col justify-center">
      <div className="sm:mx-auto sm:w-full sm:max-w-[480px]">
        <div className="bg-gray-800 p-8 shadow sm:rounded-lg">
          <form className="space-y-6" onSubmit={interactor_userEdit}>
            <p className="text-center text-white">{loggedInUserData?.email ?? "Email N/A"}</p>

            <div>
              <label htmlFor="signIn-password" className="sr-only block text-sm font-medium leading-6 text-gray-900">
                Password
              </label>
              <div className="">
                <input
                  id="signIn-password"
                  name="signIn-password"
                  type="password"
                  autoComplete="current-password"
                  placeholder="Password"
                  required
                  value={password}
                  onChange={event=>setPassword(event.target.value)}
                  className="block w-full rounded-md border-0 bg-white/5 py-1.5 px-2.5 text-white shadow-sm ring-1 ring-inset ring-white/10 focus:ring-2 focus:ring-inset focus:ring-indigo-500 sm:text-sm sm:leading-6"
                />
              </div>
            </div>

            <div>
              <button
                type="submit"
                disabled={loading}
                className="disabled:bg-slate-600 flex w-full justify-center rounded-md bg-indigo-600 px-3 py-1.5 text-sm font-semibold leading-6 text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
              >
                Edit User
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  )
}
