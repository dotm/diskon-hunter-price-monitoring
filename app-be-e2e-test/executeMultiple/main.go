package main

import (
	flowUntilUnsubscribeFromMonitoring "diskon-hunter/price-monitoring-e2e-test/libMultiple/1flowUntilUnsubscribeFromMonitoring"
)

func main() {
	executeMultipleRequest()

	flowUntilUnsubscribeFromMonitoring.KeepInImportStatement()
}

func executeMultipleRequest() {
	// flowUntilUnsubscribeFromMonitoring.Seeding()
	flowUntilUnsubscribeFromMonitoring.CheckDatabaseForNormalUserFlow()
}
