package user

import (
	"diskon-hunter/price-monitoring-e2e-test/shared/envhelper"
	"fmt"
	"time"
)

//DAO is used as interface between domain and database persistence and retrieval.

func GetStlUserDetailDynamoDBTableV1() string {
	return fmt.Sprintf(
		"%s-%s-StlUserDetail",
		envhelper.GetEnvVar("deployment_environment_name"),
		envhelper.GetEnvVar("project_name_short"),
	)
}

func GetStlUserEmailAuthenticationDynamoDBTableV1() string {
	return fmt.Sprintf(
		"%s-%s-StlUserEmailAuthentication",
		envhelper.GetEnvVar("deployment_environment_name"),
		envhelper.GetEnvVar("project_name_short"),
	)
}

func GetStlUserEmailHasOtpDetailDynamoDBTableV1() string {
	return fmt.Sprintf(
		"%s-%s-StlUserEmailHasOtpDetail",
		envhelper.GetEnvVar("deployment_environment_name"),
		envhelper.GetEnvVar("project_name_short"),
	)
}

type StlUserDetailDAOV1 struct {
	HubUserId      string
	Email          string
	HashedPassword string
}

type StlUserEmailAuthenticationDAOV1 struct {
	Email          string
	HubUserId      string
	HashedPassword string
}

type StlUserEmailHasOtpDetailDAOV1 struct {
	Email          string
	HubUserId      string
	HashedPassword string
	OTP            string
	TimeExpired    time.Time //DynamoDB time-to-live
}
