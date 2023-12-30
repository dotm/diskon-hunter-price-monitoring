import { backendBaseUrl, backendHeadersForPostRequest } from "@/utils/constants"
import { handleErrorInFrontend } from "@/utils/error"
import { FormEvent, useState } from "react"

export default function ResetPasswordForm() {
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [otp, setOtp] = useState("")
  const [resetPasswordStep, setResetPasswordStep] =
    useState<"input_new_credential"|"input_otp"|"finish">("input_new_credential")
  const [loading, setLoading] = useState(false)
  async function interactor_userResetPassword(event: FormEvent<HTMLFormElement>){
    event.preventDefault()
    
    try {
      setLoading(true)
      const resetPasswordRespJson = await fetch(`${backendBaseUrl}/v1/user.resetPassword`, {
        method: 'POST',
        headers: backendHeadersForPostRequest(),
        body: JSON.stringify({
          email: email,
          password: password,
        }),
      })
      .then(response => response.json())
      if(!resetPasswordRespJson.ok || resetPasswordRespJson.data === undefined){
        throw new Error(resetPasswordRespJson.err?.code ?? "error resetPasswordRespJson")
      }
      setResetPasswordStep("input_otp")
    } catch (error) {
      handleErrorInFrontend(error)
    } finally {
      setLoading(false)
    }
  }

  async function interactor_userValidateOtp(event: FormEvent<HTMLFormElement>){
    event.preventDefault()
    
    try {
      setLoading(true)
      const validateOtpRespJson = await fetch(`${backendBaseUrl}/v1/user.validateOtp`, {
        method: 'POST',
        headers: backendHeadersForPostRequest(),
        body: JSON.stringify({
          email: email,
          otp: otp,
        }),
      })
      .then(response => response.json())
      if(!validateOtpRespJson.ok || validateOtpRespJson.data === undefined){
        throw new Error(validateOtpRespJson.err?.code ?? "error resetPasswordRespJson")
      }
      setResetPasswordStep("finish")
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
          {
            resetPasswordStep === "input_new_credential"
            ?
            <form className="space-y-6" onSubmit={interactor_userResetPassword}>
              <p className="text-white">Silahkan isi email dan password baru anda. Jika memang email anda ada di database kami, kami akan mengirimkan OTP untuk mengganti password anda.</p>
              <div>
                <label htmlFor="resetPassword-email" className="sr-only block text-sm font-medium leading-6 text-gray-900">
                  Alamat Email
                </label>
                <div>
                  <input
                    id="resetPassword-email"
                    name="resetPassword-email"
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
                <label htmlFor="resetPassword-password" className="sr-only block text-sm font-medium leading-6 text-gray-900">
                  Password
                </label>
                <div className="">
                  <input
                    id="resetPassword-password"
                    name="resetPassword-password"
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
                  Reset Password
                </button>
              </div>
            </form>
            : <></>
          }
          {
            resetPasswordStep === "input_otp"
            ?
            <form className="space-y-6" onSubmit={interactor_userValidateOtp}>
              <p className="text-white">
                Silahkan cek email anda untuk mendapatkan OTP. Jangan bagikan OTP kepada siapapun termasuk admin Diskon Hunter.
                <br/>
                Jika email tidak terkirim, silahkan klik Kembali dan lakukan Reset Password lagi.
              </p>
              <div>
                <label htmlFor="resetPassword-otp" className="sr-only block text-sm font-medium leading-6 text-gray-900">
                  OTP
                </label>
                <div>
                  <input
                    id="resetPassword-otp"
                    name="resetPassword-otp"
                    type="text"
                    placeholder="OTP"
                    required
                    value={otp}
                    onChange={event=>setOtp(event.target.value)}
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
                  Masukkan OTP
                </button>
              </div>
              <div>
                <button
                  onClick={event=>setResetPasswordStep('input_new_credential')}
                  disabled={loading}
                  className="disabled:bg-slate-600 flex w-full justify-center rounded-md bg-red-600 px-3 py-1.5 text-sm font-semibold leading-6 text-white shadow-sm hover:bg-red-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-red-600"
                >
                  Kembali
                </button>
              </div>
            </form>
            : <></>
          }
          {
            resetPasswordStep === "finish"
            ?
            <div className="space-y-6">
              <p className="text-white">Password anda berhasil diganti!<br/>Silahkan melakukan Sign In.</p>
              <div>
                <button
                  onClick={event=>setResetPasswordStep('input_new_credential')}
                  disabled={loading}
                  className="disabled:bg-slate-600 flex w-full justify-center rounded-md bg-indigo-600 px-3 py-1.5 text-sm font-semibold leading-6 text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
                >
                  Kembali
                </button>
              </div>
            </div>
            : <></>
          }
        </div>
      </div>
    </div>
  )
}