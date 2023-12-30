import { NextRouter } from "next/router";

type AppVersionResponseData = {
  version: string
}

export async function refreshIfNewAppVersionAvailable(
  appVersion: string,
  setAppVersion: (value: string)=>void,
  router: NextRouter,
) {
  //implement /v1/appVersion in backend ~kodok
  // const response = await fetch(
  //   `${backendBaseUrl}/v1/appVersion`,
  //   {
  //     method: "POST",
  //   }
  // );
  // const respJson: AppVersionResponseData = await response.json();
  // const currentVersion = respJson.version
  // if(appVersion !== currentVersion){
  //   alert("Kami akan memuat ulang halaman ini agar anda dapat mendapatkan versi terbaru aplikasi kami.")
  //   setAppVersion(currentVersion)
  //   router.reload()
  // }
}
