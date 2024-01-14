package flowUntilCompanyCreationAndSubscription

import (
	searchedItemAddMultiple "diskon-hunter/price-monitoring-e2e-test/libSingle/searchedItem/addMultiple"
	searchedItemEditMultiple "diskon-hunter/price-monitoring-e2e-test/libSingle/searchedItem/editMultiple"
	searchedItemList "diskon-hunter/price-monitoring-e2e-test/libSingle/searchedItem/list"
	userEdit "diskon-hunter/price-monitoring-e2e-test/libSingle/user/edit"
	userResetPassword "diskon-hunter/price-monitoring-e2e-test/libSingle/user/resetPassword"
	userSignIn "diskon-hunter/price-monitoring-e2e-test/libSingle/user/signin"
	userSignUp "diskon-hunter/price-monitoring-e2e-test/libSingle/user/signup"
	userValidateOtp "diskon-hunter/price-monitoring-e2e-test/libSingle/user/validateOtp"
	"diskon-hunter/price-monitoring-e2e-test/shared/currencyutil"
	"diskon-hunter/price-monitoring-e2e-test/shared/dynamodbhelper"
	"fmt"
)

func KeepInImportStatement() {
	//do nothing.
	//the purpose of this function is to keep import statement in executeMultiple.
}

func CheckDatabaseForNormalUserFlow() {
	continueTesting := false

	userEmail := "diskon.hunter.e2e.bot@yopmail.com"
	password := "Test123!"

	userSignUpRequestDTO := userSignUp.GenerateRequestObject(userSignUp.GenerateRequestObjectArgs{
		Email:    userEmail,
		Password: password,
	})
	fmt.Printf("__execute__ userSignUpRequestDTO: %v\n", userSignUpRequestDTO)
	userSignUpResult, err := userSignUp.Execute(userSignUpRequestDTO)
	if err != nil {
		err = fmt.Errorf("error userSignUp: %s", err)
		fmt.Println(err)
		continueTesting = false
		return
	}
	requesterUserId := userSignUpResult.ResponseData.HubUserId
	continueTesting = CheckDatabaseForUserSignUp(userSignUpRequestDTO, userSignUpResult)
	if !continueTesting {
		return
	}
	userEmailHasOtpDetailForSignUpFlow, _, err := dynamodbhelper.GetUserEmailHasOtpDetailByEmail(
		dynamodbhelper.CreateClientFromSession(),
		userEmail,
	)
	if err != nil {
		err = fmt.Errorf("error GetUserEmailHasOtpDetailByEmail: %s", err)
		fmt.Println(err)
		continueTesting = false
		return
	}
	otpForSignUpFlow := userEmailHasOtpDetailForSignUpFlow.OTP
	userValidateOtpRequestDTOForSignUpFlow := userValidateOtp.GenerateRequestObject(userValidateOtp.GenerateRequestObjectArgs{
		Email: userEmail,
		OTP:   otpForSignUpFlow,
	})
	fmt.Printf("__execute__ userValidateOtpRequestDTO: %v\n", userValidateOtpRequestDTOForSignUpFlow)
	userValidateOtpResultForSignUpFlow, err := userValidateOtp.Execute(userValidateOtpRequestDTOForSignUpFlow)
	if err != nil {
		err = fmt.Errorf("error userValidateOtp: %s", err)
		fmt.Println(err)
		continueTesting = false
		return
	}
	continueTesting = CheckDatabaseForUserValidateOtp(userValidateOtpRequestDTOForSignUpFlow, userValidateOtpResultForSignUpFlow, requesterUserId)
	if !continueTesting {
		return
	}

	defer func() {
		_, err := dynamodbhelper.DeleteUserListByFilter(
			dynamodbhelper.CreateClientFromSession(),
			[]string{userValidateOtpResultForSignUpFlow.ResponseData.Id},
		)
		if err == nil {
			fmt.Printf("__cleanup__ DeleteUserListByFilter successful\n\n")
		} else {
			fmt.Printf("error DeleteUserListByFilter: %v\n\n", err)
		}
	}()

	userSignInRequestDTO := userSignIn.GenerateRequestObject(userSignIn.GenerateRequestObjectArgs{
		Email:    userEmail,
		Password: password,
	})
	fmt.Printf("__execute__ userSignInRequestDTO: %v\n", userSignInRequestDTO)
	userSignInResult, err := userSignIn.Execute(userSignInRequestDTO)
	if err != nil {
		err = fmt.Errorf("error userSignIn: %s", err)
		fmt.Println(err)
		continueTesting = false
		return
	}

	jwtToken := userSignInResult.JWTCookieString

	userEditRequestDTO := userEdit.GenerateRequestObject(userEdit.GenerateRequestObjectArgs{
		Password:       password,
		WhatsAppNumber: "+6281",
	})
	fmt.Printf("__execute__ userEditRequestDTO: %v\n", userEditRequestDTO)
	userEditResult, err := userEdit.Execute(userEditRequestDTO, jwtToken)
	if err != nil {
		err = fmt.Errorf("error userEdit: %s", err)
		fmt.Println(err)
		continueTesting = false
		return
	}
	continueTesting = CheckDatabaseForUserEdit(userEditRequestDTO, userEditResult, requesterUserId)
	if !continueTesting {
		return
	}

	searchedItemAddMultipleRequestDTO := searchedItemAddMultiple.GenerateRequestObject(searchedItemAddMultiple.GenerateRequestObjectArgs{
		SearchedItemList: []searchedItemAddMultiple.SearchedItemRequestDTOV1{
			{Name: "Smartphone A", Description: "Released at 2023", AlertPrice: currencyutil.NewFromNumberString("30000", "IDR")},
			{Name: "Smartphone B", Description: "Released at 2024", AlertPrice: currencyutil.NewFromNumberString("40000", "IDR")},
		},
	})
	fmt.Printf("__execute__ searchedItemAddMultipleRequestDTO: %v\n", searchedItemAddMultipleRequestDTO)
	searchedItemAddMultipleResult, err := searchedItemAddMultiple.Execute(searchedItemAddMultipleRequestDTO, jwtToken)
	if err != nil {
		err = fmt.Errorf("error searchedItemAddMultiple: %s", err)
		fmt.Println(err)
		continueTesting = false
		return
	}
	continueTesting = CheckDatabaseForSearchedItemAddMultiple(searchedItemAddMultipleRequestDTO, searchedItemAddMultipleResult, requesterUserId)
	if !continueTesting {
		return
	}
	defer func() {
		_, err = dynamodbhelper.DeleteStlUserSearchesItemDetailByIdList(
			dynamodbhelper.CreateClientFromSession(),
			requesterUserId,
			searchedItemAddMultipleResult.ResponseData.SearchedItemIdList,
		)
		if err == nil {
			fmt.Printf("__cleanup__ DeleteStlUserSearchesItemDetailByUrlList successful\n\n")
		} else {
			fmt.Printf("error DeleteStlUserSearchesItemDetailByUrlList: %v\n\n", err)
		}
	}()

	//check searchedItem.list
	searchedItemListRequestDTO := searchedItemList.DefaultRequestObject
	fmt.Printf("__execute__ searchedItemListRequestDTO: %v\n", searchedItemListRequestDTO)
	searchedItemListResult, err := searchedItemList.Execute(searchedItemListRequestDTO, jwtToken)
	if err != nil {
		err = fmt.Errorf("error searchedItemList: %s", err)
		fmt.Println(err)
		continueTesting = false
		return
	}
	if len(searchedItemListResult.ResponseData) != len(searchedItemAddMultipleResult.ResponseData.SearchedItemIdList) {
		fmt.Printf("error searchedItemListResult length is: %v (expected is %v)\n\n", len(searchedItemListResult.ResponseData), len(searchedItemAddMultipleResult.ResponseData.SearchedItemIdList))
		continueTesting = false
		return
	}
	for _, userSearchesItem := range searchedItemListResult.ResponseData {
		if userSearchesItem.HubUserId != requesterUserId {
			fmt.Printf("error fetch urls from other userId: %v\n\n", userSearchesItem.HubUserId)
			continueTesting = false
			return
		}
	}
	fmt.Printf("__success__ testing searchedItemListResult\n\n")

	editedSearchedItemList := []searchedItemEditMultiple.SearchedItemRequestDTOV1{}
	for index, id := range searchedItemAddMultipleResult.ResponseData.SearchedItemIdList {
		editedSearchedItemList = append(editedSearchedItemList, searchedItemEditMultiple.SearchedItemRequestDTOV1{
			HubSearchedItemId: id,
			Name:              fmt.Sprintf("Smartphone Pro %d", index+1),
			Description:       "Pro Version",
			AlertPrice:        currencyutil.NewFromNumberString("100000", "IDR"),
		})
	}
	searchedItemEditMultipleRequestDTO := searchedItemEditMultiple.GenerateRequestObject(searchedItemEditMultiple.GenerateRequestObjectArgs{
		SearchedItemList: editedSearchedItemList,
	})
	fmt.Printf("__execute__ searchedItemEditMultipleRequestDTO: %v\n", searchedItemEditMultipleRequestDTO)
	searchedItemEditMultipleResult, err := searchedItemEditMultiple.Execute(searchedItemEditMultipleRequestDTO, jwtToken)
	if err != nil {
		err = fmt.Errorf("error searchedItemEditMultiple: %s", err)
		fmt.Println(err)
		continueTesting = false
		return
	}
	continueTesting = CheckDatabaseForSearchedItemEditMultiple(searchedItemEditMultipleRequestDTO, searchedItemEditMultipleResult, requesterUserId)
	if !continueTesting {
		return
	}
	fmt.Printf("__success__ testing searchedItemEditMultipleResult\n\n")

	fmt.Printf("__finish__ all test successfully: %v\n\n", continueTesting)

	//clean up here
	//print that __cleanup__ is succesful
}

func CheckDatabaseForUserEdit(
	request userEdit.RequestDTOV1,
	result userEdit.ExecuteResult,
	requesterUserId string,
) (
	continueTesting bool,
) {
	//uncomment to disable this function
	// fmt.Printf("__disabled_checking__ should be commented out once all tests run successfully\n\n")
	// return

	err := userEdit.CheckResultIsCorrect(request, result, requesterUserId) //error already printed in CheckResultIsCorrect
	continueTesting = err == nil
	return
}

func CheckDatabaseForSearchedItemAddMultiple(
	request searchedItemAddMultiple.RequestDTOV1,
	result searchedItemAddMultiple.ExecuteResult,
	requesterUserId string,
) (
	continueTesting bool,
) {
	//uncomment to disable this function
	// fmt.Printf("__disabled_checking__ should be commented out once all tests run successfully\n\n")
	// return

	err := searchedItemAddMultiple.CheckResultIsCorrect(request, result, requesterUserId) //error already printed in CheckResultIsCorrect
	continueTesting = err == nil
	return
}

func CheckDatabaseForSearchedItemEditMultiple(
	request searchedItemEditMultiple.RequestDTOV1,
	result searchedItemEditMultiple.ExecuteResult,
	requesterUserId string,
) (
	continueTesting bool,
) {
	//uncomment to disable this function
	// fmt.Printf("__disabled_checking__ should be commented out once all tests run successfully\n\n")
	// return

	err := searchedItemEditMultiple.CheckResultIsCorrect(request, result, requesterUserId) //error already printed in CheckResultIsCorrect
	continueTesting = err == nil
	return
}

func CheckDatabaseForUserSignUp(request userSignUp.RequestDTOV1, result userSignUp.ExecuteResult) (continueTesting bool) {
	//uncomment to disable this function
	// fmt.Printf("__disabled_checking__ should be commented out once all tests run successfully\n\n")
	// return

	err := userSignUp.CheckResultIsCorrect(request, result) //error already printed in CheckResultIsCorrect
	continueTesting = err == nil
	return
}

func CheckDatabaseForUserResetPassword(request userResetPassword.RequestDTOV1, result userResetPassword.ExecuteResult) (continueTesting bool) {
	//uncomment to disable this function
	// fmt.Printf("__disabled_checking__ should be commented out once all tests run successfully\n\n")
	// return

	err := userResetPassword.CheckResultIsCorrect(request, result) //error already printed in CheckResultIsCorrect
	continueTesting = err == nil
	return
}

func CheckDatabaseForUserValidateOtp(request userValidateOtp.RequestDTOV1, result userValidateOtp.ExecuteResult, requesterUserId string) (continueTesting bool) {
	//uncomment to disable this function
	// fmt.Printf("__disabled_checking__ should be commented out once all tests run successfully\n\n")
	// return

	err := userValidateOtp.CheckResultIsCorrect(request, result, requesterUserId) //error already printed in CheckResultIsCorrect
	continueTesting = err == nil
	return
}
