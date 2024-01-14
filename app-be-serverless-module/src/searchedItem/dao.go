package searchedItem

import (
	"diskon-hunter/price-monitoring/shared/currencyutil"
	"diskon-hunter/price-monitoring/shared/envhelper"
	"fmt"
	"time"
)

//DAO is used as interface between domain and database persistence and retrieval.

func GetStlUserSearchesItemDetailDynamoDBTableV1() string {
	return fmt.Sprintf(
		"%s-%s-StlUserSearchesItemDetail",
		envhelper.GetEnvVar("deployment_environment_name"),
		envhelper.GetEnvVar("project_name_short"),
	)
}

type StlUserSearchesItemDetailDAOV1 struct {
	HubUserId         string
	HubSearchedItemId string
	Name              string
	Description       string
	Status            string
	AlertPrice        currencyutil.Currency
	TimeExpired       time.Time //DynamoDB time-to-live; latest expired time of all user
}

func NewStlUserSearchesItemDetailDAOV1(
	HubUserId string,
	HubSearchedItemId string,
	Name string,
	Description string,
	Status string,
	AlertPrice currencyutil.Currency,
	TimeExpired time.Time,
) StlUserSearchesItemDetailDAOV1 {
	return StlUserSearchesItemDetailDAOV1{
		HubUserId:         HubUserId,
		HubSearchedItemId: HubSearchedItemId,
		Name:              Name,
		Description:       Description,
		Status:            Status,
		AlertPrice:        AlertPrice,
		TimeExpired:       TimeExpired,
	}
}
