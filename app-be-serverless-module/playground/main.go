//Run with:
//	go run playground/*.go
//In windows, run with:
//  go run .\playground\
//In bash you can calculate execution time with:
//  time go run playground/*.go

//Use this playground to quickly prototype your functions or check functionalities.
//What happens in playground, stays in playground.
//  DO NOT REFERENCE ANYTHING IN THIS DIRECTORY OUTSIDE OF PLAYGROUND.
//  DO NOT COMMIT ANYTHING IN THIS DIRECTORY INTO GIT. Unless it's just a comment update.
//Feel free to import anything (standard libraries, github modules, types from this project, etc.).

package main

import (
	"diskon-hunter/price-monitoring/shared/constenum"
	"diskon-hunter/price-monitoring/shared/currencyutil"
	"diskon-hunter/price-monitoring/src/monitoredLink"
	"fmt"
	"time"

	handlebars "github.com/aymerick/raymond"
)

func main() {
	//The only thing that should exist in this function after you're done experimenting is this comment.
	//Any merge request where there is other things aside from this comment in this function should be rejected.

	//Code experiment goes here...
	// envhelper.SetLocalEnvVar()
	// a, b, c := dynamodbhelper.GetLatestCompanySubscriptionOfCompanyId(
	// 	dynamodbhelper.CreateClientFromSession(),
	// 	"83fac6fc-7a13-42da-8742-17e4e9cabcb6",
	// )
	// fmt.Println(c)
	// fmt.Println(b)
	// fmt.Println(a)
	linksAlertedByEmail := []monitoredLink.CombinedUserMonitoredLinkDataV1{
		{
			StlUserMonitorsLinkDetailDAOV1: monitoredLink.StlUserMonitorsLinkDetailDAOV1{
				HubUserId:             "id",
				HubMonitoredLinkUrl:   "url",
				AlertPrice:            currencyutil.Currency{Significand: "2", Exponent: "4", CurrencyUnit: "IDR"},
				ActiveAlertMethodList: []constenum.AlertMethod{},
				PaidAlertMethodList:   []constenum.AlertMethod{},
				TimeExpired:           time.Time{},
			},
			LatestPrice:        &currencyutil.Currency{Significand: "1", Exponent: "4", CurrencyUnit: "IDR"},
			TimeLatestScrapped: &time.Time{},
		},
		{
			StlUserMonitorsLinkDetailDAOV1: monitoredLink.StlUserMonitorsLinkDetailDAOV1{
				HubUserId:             "id",
				HubMonitoredLinkUrl:   "url",
				AlertPrice:            currencyutil.Currency{Significand: "4", Exponent: "4", CurrencyUnit: "IDR"},
				ActiveAlertMethodList: []constenum.AlertMethod{},
				PaidAlertMethodList:   []constenum.AlertMethod{},
				TimeExpired:           time.Time{},
			},
			LatestPrice:        &currencyutil.Currency{Significand: "3", Exponent: "4", CurrencyUnit: "IDR"},
			TimeLatestScrapped: &time.Time{},
		},
	}

	type linkData struct {
		HubMonitoredLinkUrl string
		LatestPriceString   string
		TimeLatestScrapped  string
	}
	type templateDataContext struct {
		LinkDataList []linkData
	}
	templateData := templateDataContext{
		LinkDataList: []linkData{},
	}
	for _, alertedLink := range linksAlertedByEmail {
		templateData.LinkDataList = append(templateData.LinkDataList, linkData{
			HubMonitoredLinkUrl: alertedLink.HubMonitoredLinkUrl,
			LatestPriceString:   alertedLink.LatestPrice.ToDisplayString(),
			TimeLatestScrapped:  alertedLink.TimeLatestScrapped.Format("2006-01-02"),
		})
	}

	htmlTemplate := `<div>
		<p>
			Halo Diskon Hunter! Kami mendeteksi link yang kamu monitor ada dibawah harga yang kamu input:
		</p>
		<ol>
			{{#each linkDataList}}
			<li>
				<ul>
					<li>{{HubMonitoredLinkUrl}}</li>
					<li>
						Harga terakhir: {{LatestPriceString}} (pada {{TimeLatestScrapped}})
					</li>
				</ul>
			</li>
			{{/each}}
		</ol>
		<p>
			Jika kamu tidak ingin menerima email ini lagi, silahkan matikan notifikasi email untuk semua link yang kamu monitor dari website kami.
		</p>
		<p>
			Jika kamu perlu bantuan, silahkan Kontak Kami.
		</p>
	</div>`
	textTemplate := `Halo Diskon Hunter! Kami mendeteksi link yang kamu monitor ada dibawah harga yang kamu input:
		
		{{#each linkDataList}}
			- {{HubMonitoredLinkUrl}}
				- Harga terakhir: {{LatestPriceString}} (pada {{TimeLatestScrapped}})
		{{/each}}

		Jika kamu tidak ingin menerima email ini lagi, silahkan matikan notifikasi email untuk semua link yang kamu monitor dari website kami.
		Jika kamu perlu bantuan, silahkan Kontak Kami.
	`
	htmlBody, err := handlebars.Render(htmlTemplate, templateData)
	textBody, err := handlebars.Render(textTemplate, templateData)
	fmt.Println(err, htmlBody, textBody)
}
