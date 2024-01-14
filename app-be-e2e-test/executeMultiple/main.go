package main

import (
	flowUntilUnsubscribeFromMonitoring "diskon-hunter/price-monitoring-e2e-test/libMultiple/1flowUntilUnsubscribeFromMonitoring"
	flowUntilUserSearchesItem "diskon-hunter/price-monitoring-e2e-test/libMultiple/2flowUntilUserSearchesItem"
)

func main() {
	executeMultipleRequest()

	flowUntilUnsubscribeFromMonitoring.KeepInImportStatement()
	flowUntilUserSearchesItem.KeepInImportStatement()
}

func executeMultipleRequest() {
	// flowUntilUnsubscribeFromMonitoring.Seeding()
	flowUntilUserSearchesItem.CheckDatabaseForNormalUserFlow()
}
