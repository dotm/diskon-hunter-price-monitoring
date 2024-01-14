import { refreshIfNewAppVersionAvailable } from "@/utils/appversionutil";
import { LocalStorageKey, backendBaseUrl, backendHeadersForPostRequest } from "@/utils/constants";
import { convertCurrencyToNumber, convertNumberStringToCurrency, createZeroCurrency, displayCurrencyInUI } from "@/utils/currencyutil";
import { handleErrorInFrontend } from "@/utils/error";
import { disableChangingNumberValueOnScroll } from "@/utils/eventhandler";
import { LoggedInUserData, UserSearchesItemDetailAddRequestDTO, UserSearchesItemListAddRequestDTO, emptyUserSearchesItemListAddRequestDTO } from "@/utils/models";
import Link from "next/link";
import { useRouter } from 'next/router';
import { MouseEvent, useEffect, useState } from "react";
import useLocalStorage from "use-local-storage";
import { v4 } from 'uuid';

export default function SearchItemAddItems({
  appVersion,
  setAppVersion,
}:{
  appVersion: string,
  setAppVersion: (value: string)=>void,
}) {
  const router = useRouter()
  const [loggedInUserData, setLoggedInUserData] =
    useLocalStorage<LoggedInUserData | undefined>(LocalStorageKey.loggedInUser, undefined)
  const [requestDTO, setRequestDTO] = useLocalStorage<UserSearchesItemListAddRequestDTO>(
    LocalStorageKey.UserSearchesItemListAddRequestDTO,
    emptyUserSearchesItemListAddRequestDTO(),
  )
  const [loading, setLoading] = useState(false)

  function updateSearchedItemList(searchedItem: UserSearchesItemDetailAddRequestDTO){
    setRequestDTO({
      SearchedItemList: requestDTO.SearchedItemList.map(o=>{
        if(o.FrontendID !== searchedItem.FrontendID){
          return o
        }
        return searchedItem
      })
    })
  }
  function resetRequestDTO(){
    setRequestDTO(emptyUserSearchesItemListAddRequestDTO())
  }
  useEffect(() => {
    refreshIfNewAppVersionAvailable(appVersion, setAppVersion, router, function(){resetRequestDTO()})
  })
  async function interactor_searchedItemAddMultiple(event: MouseEvent<HTMLButtonElement, globalThis.MouseEvent>){
    event.preventDefault()
    try {
      setLoading(true)
      if(loggedInUserData === undefined){
        throw new Error("Mohon sign in terlebih dahulu")
      }
      
      if(requestDTO.SearchedItemList.find(o=>!o.Name) !== undefined){
        throw new Error("Mohon isi semua Nama Barang")
      }
      const processedSearchedItemList = requestDTO.SearchedItemList.map(o=>{
        return {
          ...o,
          AlertPrice: convertNumberStringToCurrency(o.AlertPriceString, "IDR")
        }
      })
      if(processedSearchedItemList.find(o=>convertCurrencyToNumber(o.AlertPrice) <= 0) !== undefined){
        throw new Error("Mohon isi semua Harga")
      }
      const validatedRequestDTO = {
        SearchedItemList: processedSearchedItemList
      }
      const searchedItemAddRespJson = await fetch(`${backendBaseUrl}/v1/searchedItem.addMultiple`, {
        method: 'POST',
        headers: backendHeadersForPostRequest(loggedInUserData.jwt),
        body: JSON.stringify(validatedRequestDTO),
      })
      .then(response => response.json())
      if(!searchedItemAddRespJson.ok || !searchedItemAddRespJson.data){
        throw new Error(searchedItemAddRespJson.err?.code ?? "error searchedItemAddRespJson")
      }
      //do something with searchedItemAddRespJson.data
      resetRequestDTO()
      alert(`Berhasil menambahkan barang yang dicari`)
      router.replace('/searchItem')
    } catch (error) {
      handleErrorInFrontend(error)
    } finally {
      setLoading(false)
    }
  }
  
  return (
    <div className="space-y-3">
      {
        loggedInUserData === undefined
        ?
        <Link href="/settings" className="bg-gray-800 hover:bg-gray-700 outline-green-500 text-white block w-[100%] outline pb-2 pt-1 px-3 rounded-xl mx-auto">
          <p className="text-sm">
          Mohon sign in ke akun anda terlebih dahulu.
          <br />
          Dan pastikan anda telah mengisi nomor WhatsApp untuk akun anda.
          </p>
        </Link>
        :
        <></>
      }
      {
        loggedInUserData !== undefined && !loggedInUserData?.whatsAppNumber
        ?
        <Link href="/settings" className="bg-gray-800 hover:bg-gray-700 outline-green-500 text-white block w-[100%] outline pb-2 pt-1 px-3 rounded-xl mx-auto">
          <p className="text-sm">
          Mohon isi nomor WhatsApp untuk akun anda terlebih dahulu.
          </p>
        </Link>
        :
        <></>
      }
      { loggedInUserData !== undefined && !!loggedInUserData?.whatsAppNumber
        ?
        <>
          <div className="divide-y divide-gray-400 bg-gray-800 text-white block w-[100%] px-3 pb-3 rounded-xl mx-auto">
            <h2 className="text-center font-bold p-2">Daftar Barang yang Anda Cari</h2>
            <p className="text-center text-sm pt-1 pb-2">
              Klik Tambah Barang untuk mulai menambahkan barang yang anda ingin cari.
              <br />
              Mohon klik Simpan untuk mencatat barang yang anda telah Tambahkan.
              <br />
              Kami akan membantu anda mencari barang yang anda Tambahkan selama satu tahun kedepan.
            </p>
            <div className='flex justify-center flex-wrap text-center space-x-2'>
              <button
                type="button"
                className="my-2 rounded-md w-20 disabled:bg-slate-600 bg-blue-600 py-1 text-sm font-semibold text-white shadow-sm hover:bg-blue-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600"
                onClick={interactor_searchedItemAddMultiple}
                disabled={loading}
              >
                Simpan
              </button>
              <button
                type="button"
                className="my-2 rounded-md w-20 disabled:bg-slate-600 bg-red-600 py-1 text-sm font-semibold text-white shadow-sm hover:bg-red-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-red-600"
                onClick={() => {
                  router.back()
                }}
                disabled={loading}
              >
                Kembali
              </button>
            </div>
            <div className="bg-gray-800 text-white block w-[100%] pb-3 px-3 mx-auto space-y-4 divide-y divide-gray-400">
              {
                requestDTO.SearchedItemList.map((itemDetail)=>{
                  return (
                    <div
                      key={itemDetail.FrontendID} //edit form should use HubSearchedItemId
                      className="pt-4"
                    >
                      <div>
                        <label
                          htmlFor={`searchedItemAdd-name-${itemDetail.FrontendID}`}
                          className="sr-only block text-sm font-medium leading-6 text-gray-900"
                        >
                          Nama Barang
                        </label>
                        <div>
                          <input
                            id={`searchedItemAdd-name-${itemDetail.FrontendID}`}
                            name={`searchedItemAdd-name-${itemDetail.FrontendID}`}
                            type="text"
                            placeholder="Nama Barang"
                            value={itemDetail.Name}
                            onChange={event=>updateSearchedItemList({
                              ...itemDetail,
                              Name: event.target.value,
                            })}
                            className="placeholder:text-gray-400 block w-full rounded-md border-0 bg-white/5 py-1.5 px-2.5 text-white shadow-sm ring-1 ring-inset ring-white/10 focus:ring-2 focus:ring-inset focus:ring-indigo-500 sm:text-sm sm:leading-6"
                          />
                        </div>
                      </div>
                      <div>
                        <label
                          htmlFor={`searchedItemAdd-name-${itemDetail.FrontendID}`}
                          className="block text-sm font-medium leading-6 text-white"
                        >
                          Deskripsi Barang (max 2000 karakter)
                        </label>
                        <div>
                          <textarea
                            id={`searchedItemAdd-description-${itemDetail.FrontendID}`}
                            name={`searchedItemAdd-description-${itemDetail.FrontendID}`}
                            placeholder="misal: contoh link barang, link untuk gambar barang, spesifikasi teknis, tahun keluaran, jenis varian, warna, ukuran, dsb."
                            value={itemDetail.Description}
                            onChange={event=>updateSearchedItemList({
                              ...itemDetail,
                              Description: event.target.value,
                            })}
                            className="h-44 placeholder:text-gray-400 block w-full rounded-md border-0 bg-white/5 py-1.5 px-2.5 text-white shadow-sm ring-1 ring-inset ring-white/10 focus:ring-2 focus:ring-inset focus:ring-indigo-500 sm:text-sm sm:leading-6"
                          />
                        </div>
                      </div>
                      <div>
                        <label
                          htmlFor={`searchedItemAdd-alertPriceString-${itemDetail.FrontendID}`}
                          className="block text-sm font-medium leading-6 text-gray-100"
                        >
                          Ingatkan saya pada harga: {displayCurrencyInUI(
                            convertNumberStringToCurrency(itemDetail.AlertPriceString,"IDR"))
                          }
                        </label>
                        <div>
                          <input
                            id={`searchedItemAdd-alertPriceString-${itemDetail.FrontendID}`}
                            name={`searchedItemAdd-alertPriceString-${itemDetail.FrontendID}`}
                            type="number"
                            placeholder="Masukkan harga (angka tanpa titik atau koma)"
                            value={itemDetail.AlertPriceString}
                            onWheel={disableChangingNumberValueOnScroll}
                            onKeyDown={event=>{
                              if(["e",".",",","-"].includes(event.key)){ //disable keys
                                event.preventDefault()
                              }
                            }}
                            onChange={event=>updateSearchedItemList({
                              ...itemDetail,
                              AlertPriceString: event.target.value,
                            })}
                            className="placeholder:text-gray-400 block w-full rounded-md border-0 bg-white/5 py-1.5 px-2.5 text-white shadow-sm ring-1 ring-inset ring-white/10 focus:ring-2 focus:ring-inset focus:ring-indigo-500 sm:text-sm sm:leading-6"
                          />
                        </div>
                      </div>
                      <div className='flex justify-end flex-wrap text-center space-x-2'>
                        <button
                          type="button"
                          className="mt-3 rounded-md w-40 disabled:bg-slate-600 bg-red-600 py-1 text-sm font-semibold text-white shadow-sm hover:bg-red-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-red-600"
                          onClick={()=>{
                            setRequestDTO({
                              SearchedItemList: [...requestDTO.SearchedItemList.filter(o=>o.FrontendID!==itemDetail.FrontendID)]
                            })
                          }}
                          disabled={loading}
                        >
                          Hapus Barang
                        </button>
                      </div>
                    </div>
                  )
                })
              }
              <div className='flex justify-center flex-wrap text-center space-x-2 mt-5'>
                <button
                  type="button"
                  className="my-2 rounded-md w-40 disabled:bg-slate-600 bg-green-600 py-1 text-sm font-semibold text-white shadow-sm hover:bg-green-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-green-600"
                  onClick={()=>{
                    setRequestDTO({
                      SearchedItemList: [...requestDTO.SearchedItemList, {
                        FrontendID: v4(),
                        Name: "",
                        Description: "",
                        AlertPrice: createZeroCurrency("IDR"),
                        AlertPriceString: "",
                      }]
                    })
                  }}
                  disabled={loading}
                >
                  Tambah Barang
                </button>
              </div>
            </div>
          </div>
        </>
        :
        <></>
      }
    </div>
  )
}
