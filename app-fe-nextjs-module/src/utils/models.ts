export interface LoggedInUserData {
  jwt: string,
  userId: string,
  email: string,
}

export interface UserLinkDetail {
  HubUserId: string,
  HubMonitoredLinkUrl: string,
  AlertPrice: Price,
  ActiveAlertMethodList: string[],
  PaidAlertMethodList: string[],
  TimeExpired: Date,
  
  //data from StlMonitoredLinkDetailDAOV1
  LatestPrice: Price,
  TimeLatestScrapped: Date,
}

export interface Price {
  amount: number,
  currency: "IDR",
}
