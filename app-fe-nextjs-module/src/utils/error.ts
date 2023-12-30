export function handleErrorInFrontend(error: any){
  if(error === null || error === undefined){
    return
  }
  const errorMessage = convertErrorCodeToErrorMessage((error.message ?? "") as string)
  let displayedError = errorMessage
  if(errorMessage === "TypeError: Failed to fetch"){
    displayedError = `Gagal mengambil data. Pastikan anda memiliki akses ke internet.`
  }
  alert(displayedError)
}

function convertErrorCodeToErrorMessage(errorMessage: string): string {
  switch (errorMessage) {
    case "user/email-not-registered": return "Email atau password salah."
    case "user/credential-incorrect": return "Email atau password salah."
    default: return errorMessage
  }
}