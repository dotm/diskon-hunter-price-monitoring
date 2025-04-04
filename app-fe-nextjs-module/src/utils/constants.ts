export const supportEmailAddress = "diskon.hunter.official@gmail.com"

export const backendBaseUrl = "https://ueb8gkzdtj.execute-api.ap-southeast-1.amazonaws.com/prod"
export function backendHeadersForPostRequest(jwtToken: string | null = null){
  const headers = {
    'Accept': '*/*',
    'Accept-Encoding': 'gzip, deflate, br',
    'Connection': 'keep-alive',
    'Content-Type': 'application/json',
    'Authorization': '',
  }
  if(jwtToken){
    headers['Authorization'] = 'Bearer ' + jwtToken
  }
  return headers
}

export const LocalStorageKey = {
  lastRefreshDate: "dhpm-lastRefreshDate",
  appVersion: "dhpm-appVersion",
  loggedInUser: "dhpm-loggedInUser",
  UserMonitorsLinkListAddRequestDTO: "dhpm-UserMonitorsLinkListAddRequestDTO",
  UserMonitorsLinkListEditRequestDTO: "dhpm-UserMonitorsLinkListEditRequestDTO",
  UserSearchesItemListAddRequestDTO: "dhpm-UserSearchesItemListAddRequestDTO",
  UserSearchesItemListEditRequestDTO: "dhpm-UserSearchesItemListEditRequestDTO",
}
