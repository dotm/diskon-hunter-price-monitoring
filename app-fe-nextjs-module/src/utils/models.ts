import { Currency } from "./currencyutil";

export interface LoggedInUserData {
  jwt: string,
  userId: string,
  email: string,
}

export const AlertMethod = {
  Email: {
    backendValue: "Email",
    frontendValue: "Email",
  },
  PushNotification: {
    backendValue: "PushNotification",
    frontendValue: "Push Notification",
  },
  SMS: {
    backendValue: "SMS",
    frontendValue: "SMS",
  },
  WhatsApp: {
    backendValue: "WhatsApp",
    frontendValue: "WhatsApp",
  },
}
export const AvailableAlertMethodList = [
  AlertMethod.Email,
  // AlertMethod.PushNotification,
  // AlertMethod.SMS,
  // AlertMethod.WhatsApp,
]

export function emptyUserMonitorsLinkListAddRequestDTO(): UserMonitorsLinkListAddRequestDTO{
  return {MonitoredLinkList: []}
}
export interface UserMonitorsLinkListAddRequestDTO {
  MonitoredLinkList: UserMonitorsLinkDetailAddRequestDTO[],
}
export interface UserMonitorsLinkDetailAddRequestDTO {
  FrontendID: string, //used for key in add form. edit form should use HubMonitoredLinkUrl
  HubMonitoredLinkUrl: string,
  AlertPrice: Currency,
  AlertPriceString: string, //should be processed into AlertPrice Currency before sending request
  AlertMethodList: string[],
}
export function emptyUserMonitorsLinkListEditRequestDTO(): UserMonitorsLinkListEditRequestDTO{
  return {MonitoredLinkList: []}
}
export interface UserMonitorsLinkListEditRequestDTO {
  MonitoredLinkList: UserMonitorsLinkDetailEditRequestDTO[],
}
export interface UserMonitorsLinkDetailEditRequestDTO {
  HubMonitoredLinkUrl: string,
  AlertPrice: Currency,
  AlertPriceString: string, //should be processed into AlertPrice Currency before sending request
  ActiveAlertMethodList: string[] | null,
  PaidAlertMethodList: string[] | null,
}
export interface UserLinkDetail {
  HubUserId: string,
  HubMonitoredLinkUrl: string,
  AlertPrice: Currency,
  ActiveAlertMethodList: string[] | null,
  PaidAlertMethodList: string[] | null,
  TimeExpired: Date,
  
  //data from StlMonitoredLinkDetailDAOV1
  LatestPrice: Currency | null,
  TimeLatestScrapped: Date | null,
}
