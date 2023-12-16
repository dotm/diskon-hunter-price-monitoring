package monitoredLink

import (
	"diskon-hunter/price-monitoring/shared/constenum"
	"diskon-hunter/price-monitoring/shared/currencyutil"
	"diskon-hunter/price-monitoring/shared/envhelper"
	"fmt"
	"time"
)

//DAO is used as interface between domain and database persistence and retrieval.

func GetStlMonitoredLinkDetailDynamoDBTableV1() string {
	return fmt.Sprintf(
		"%s-%s-StlMonitoredLinkDetail",
		envhelper.GetEnvVar("deployment_environment_name"),
		envhelper.GetEnvVar("project_name_short"),
	)
}

func GetStlUserMonitorsLinkDetailDynamoDBTableV1() string {
	return fmt.Sprintf(
		"%s-%s-StlUserMonitorsLinkDetail",
		envhelper.GetEnvVar("deployment_environment_name"),
		envhelper.GetEnvVar("project_name_short"),
	)
}

// Use New function to avoid unintentional missing field.
// The only exception is if you want to initialize empty struct with zero values.
type StlMonitoredLinkDetailDAOV1 struct {
	HubMonitoredLinkUrl string                 //unnecessary query params must be cleaned to avoid duplicate scrapping
	LatestPrice         *currencyutil.Currency //can be nil
	TimeLatestScrapped  *time.Time
	TimeExpired         time.Time //DynamoDB time-to-live; latest expired time of all user
}

func NewStlMonitoredLinkDetailDAOV1(
	HubMonitoredLinkUrl string,
	LatestPrice *currencyutil.Currency,
	TimeLatestScrapped *time.Time,
	TimeExpired time.Time,
) StlMonitoredLinkDetailDAOV1 {
	return StlMonitoredLinkDetailDAOV1{
		HubMonitoredLinkUrl: HubMonitoredLinkUrl,
		LatestPrice:         LatestPrice,
		TimeLatestScrapped:  TimeLatestScrapped,
		TimeExpired:         TimeExpired,
	}
}

type StlUserMonitorsLinkDetailDAOV1 struct {
	HubUserId             string
	HubMonitoredLinkUrl   string
	AlertPrice            currencyutil.Currency
	ActiveAlertMethodList []constenum.AlertMethod
	PaidAlertMethodList   []constenum.AlertMethod
	TimeExpired           time.Time //DynamoDB time-to-live; latest expired time per user
}

func NewStlUserMonitorsLinkDetailDAOV1(
	HubUserId string,
	HubMonitoredLinkUrl string,
	AlertPrice currencyutil.Currency,
	ActiveAlertMethodList []constenum.AlertMethod,
	PaidAlertMethodList []constenum.AlertMethod,
	TimeExpired time.Time,
) StlUserMonitorsLinkDetailDAOV1 {
	return StlUserMonitorsLinkDetailDAOV1{
		HubUserId:             HubUserId,
		HubMonitoredLinkUrl:   HubMonitoredLinkUrl,
		AlertPrice:            AlertPrice,
		ActiveAlertMethodList: ActiveAlertMethodList,
		PaidAlertMethodList:   PaidAlertMethodList,
		TimeExpired:           TimeExpired,
	}
}
