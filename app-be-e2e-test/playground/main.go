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
	"diskon-hunter/price-monitoring-e2e-test/shared/dynamodbhelper"
	"diskon-hunter/price-monitoring-e2e-test/shared/envhelper"
	"fmt"
)

func main() {
	//The only thing that should exist in this function after you're done experimenting is this comment.
	//Any merge request where there is other things aside from this comment in this function should be rejected.

	//Code experiment goes here...
	envhelper.SetLocalEnvVar()
	errObj, err := dynamodbhelper.DeleteUserListByFilter(
		dynamodbhelper.CreateClientFromSession(), []string{"f1ffd659-bcb7-4bac-add2-1a1ab6d62727"})
	fmt.Println(err)
	fmt.Println(errObj)
}
