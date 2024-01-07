import { LocalStorageKey, backendBaseUrl, backendHeadersForPostRequest } from "@/utils/constants"
import { handleErrorInFrontend } from "@/utils/error"
import { cleanUpJWT } from "@/utils/jwt"
import { LoggedInUserData } from "@/utils/models"
import { FormEvent, useState } from "react"
import useLocalStorage from "use-local-storage"

export default function SignInForm() {
  const [loggedInUserData, setLoggedInUserData] =
    useLocalStorage<LoggedInUserData | undefined>(LocalStorageKey.loggedInUser, undefined)
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [loading, setLoading] = useState(false)
  async function interactor_userSignIn(event: FormEvent<HTMLFormElement>){
    event.preventDefault()
    
    try {
      setLoading(true)
      const signInRespJson = await fetch(`${backendBaseUrl}/v1/user.signin`, {
        method: 'POST',
        headers: backendHeadersForPostRequest(),
        body: JSON.stringify({
          email: email,
          password: password
        }),
      })
      .then(response => response.json())
      if(!signInRespJson.ok || !signInRespJson.data){
        throw new Error(signInRespJson.err?.code ?? "error signInRespJson")
      }
      const jwt = cleanUpJWT(signInRespJson.data.JwtCookieString)

      setLoggedInUserData({
        jwt: jwt,
        userId: signInRespJson.data.HubUserId,
        email: signInRespJson.data.Email,
      })
    } catch (error) {
      handleErrorInFrontend(error)
    } finally {
      setLoading(false)
    }
  }
  async function interactor_userMe(event: FormEvent<HTMLFormElement>){
    event.preventDefault()
    try {
      setLoading(true)
      if(loggedInUserData === undefined){
        throw new Error("Mohon sign in terlebih dahulu")
      }
      const userMeRespJson = await fetch(`${backendBaseUrl}/v1/user.me`, {
        method: 'POST',
        headers: backendHeadersForPostRequest(loggedInUserData.jwt),
        body: JSON.stringify({}),
      })
      .then(response => response.json())
      if(!userMeRespJson.ok || !userMeRespJson.data){
        throw new Error(userMeRespJson.err?.code ?? "error userMeRespJson")
      }
      //do something with userMeRespJson.data
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
          <form className="space-y-6" onSubmit={interactor_userSignIn}>
            <div>
              <label htmlFor="signIn-email" className="sr-only block text-sm font-medium leading-6 text-gray-900">
                Alamat Email
              </label>
              <div>
                <input
                  id="signIn-email"
                  name="signIn-email"
                  type="email"
                  autoComplete="email"
                  placeholder="Alamat Email"
                  required
                  value={email}
                  onChange={event=>setEmail(event.target.value)}
                  className="block w-full rounded-md border-0 bg-white/5 py-1.5 px-2.5 text-white shadow-sm ring-1 ring-inset ring-white/10 focus:ring-2 focus:ring-inset focus:ring-indigo-500 sm:text-sm sm:leading-6"
                />
              </div>
            </div>

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
                Sign in
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  )
}
