import { supportEmailAddress } from "@/utils/constants"
import { Transition } from "@headlessui/react"
import { CheckCircleIcon, XMarkIcon } from "@heroicons/react/24/outline"
import React, { useState } from "react"

export default function ContactUs() {
  const [show, setShow] = useState(false)

  return (
    <>
      <button
        className="bg-gray-800 hover:bg-gray-700 outline-red-500 text-white block w-[100%] outline pb-2 pt-1 px-3 rounded-xl mx-auto"
        onClick={function(){
          navigator.clipboard.writeText(supportEmailAddress)
          setShow(true)
          setTimeout(() => {
            setShow(false)
          }, 5000);
        }}
      >
        <p className="text-center">
        Ada masalah atau kesulitan dalam menggunakan aplikasi ini?
        <br/>
        Punya pertanyaan atau kritik? Punya saran fitur untuk kami?
        <br/>
        <br/>
        Mohon kirimkan email ke: {supportEmailAddress}
        <br/>
        Silahkan klik disini untuk meng-copy email tersebut.
        </p>
      </button>

      {/* Global notification live region, render this permanently at the end of the document */}
      <div
        aria-live="assertive"
        className="pointer-events-none fixed inset-0 flex items-end px-4 py-6 sm:items-start sm:p-6"
      >
        <div className="flex w-full flex-col items-center space-y-4 sm:items-end">
          {/* Notification panel, dynamically insert this into the live region when it needs to be displayed */}
          <Transition
            show={show}
            as={React.Fragment}
            enter="transform ease-out duration-300 transition"
            enterFrom="translate-y-2 opacity-0 sm:translate-y-0 sm:translate-x-2"
            enterTo="translate-y-0 opacity-100 sm:translate-x-0"
            leave="transition ease-in duration-100"
            leaveFrom="opacity-100"
            leaveTo="opacity-0"
          >
            <div className="pointer-events-auto w-full max-w-sm overflow-hidden rounded-lg bg-gray-800 shadow-sm shadow-white outline outline-green-400">
              <div className="p-4">
                <div className="flex items-start">
                  <div className="flex-shrink-0">
                    <CheckCircleIcon className="h-6 w-6 text-green-400" aria-hidden="true" />
                  </div>
                  <div className="ml-3 w-0 flex-1 pt-0.5">
                    <p className="text-sm font-medium text-gray-100">Alamat email berhasil di-copy</p>
                    <p className="mt-1 text-sm text-gray-200">Kami tunggu email dari anda</p>
                  </div>
                  <div className="ml-4 flex flex-shrink-0">
                    <button
                      type="button"
                      className="inline-flex rounded-md bg-gray-800 text-gray-100 hover:text-gray-300 focus:outline-none focus:ring-2 focus:ring-green-500 focus:ring-offset-2"
                      onClick={() => {
                        setShow(false)
                      }}
                    >
                      <span className="sr-only">Close</span>
                      <XMarkIcon className="h-5 w-5" aria-hidden="true" />
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </Transition>
        </div>
      </div>
    </>
  )
}
