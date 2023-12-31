package flowUntilCompanyCreationAndSubscription

import (
	monitoredLinkAddMultiple "diskon-hunter/price-monitoring-e2e-test/libSingle/monitoredLink/addMultiple"
	monitoredLinkEditMultiple "diskon-hunter/price-monitoring-e2e-test/libSingle/monitoredLink/editMultiple"
	monitoredLinkList "diskon-hunter/price-monitoring-e2e-test/libSingle/monitoredLink/list"
	userEdit "diskon-hunter/price-monitoring-e2e-test/libSingle/user/edit"
	userResetPassword "diskon-hunter/price-monitoring-e2e-test/libSingle/user/resetPassword"
	userSignIn "diskon-hunter/price-monitoring-e2e-test/libSingle/user/signin"
	userSignUp "diskon-hunter/price-monitoring-e2e-test/libSingle/user/signup"
	userValidateOtp "diskon-hunter/price-monitoring-e2e-test/libSingle/user/validateOtp"
	"diskon-hunter/price-monitoring-e2e-test/shared/constenum"
	"diskon-hunter/price-monitoring-e2e-test/shared/currencyutil"
	"diskon-hunter/price-monitoring-e2e-test/shared/delivery/monitoredLink"
	"diskon-hunter/price-monitoring-e2e-test/shared/dynamodbhelper"
	"fmt"
	"time"
)

func KeepInImportStatement() {
	//do nothing.
	//the purpose of this function is to keep import statement in executeMultiple.
}

func CheckDatabaseForNormalUserFlow() {
	continueTesting := false

	userEmail := "diskon.hunter.e2e@yopmail.com"
	firstPassword := "Test123!"
	secondPassword := "Test12!"
	password := "Test123!"

	userSignUpRequestDTO := userSignUp.GenerateRequestObject(userSignUp.GenerateRequestObjectArgs{
		Email:    userEmail,
		Password: firstPassword,
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

	userResetPasswordRequestDTO := userResetPassword.GenerateRequestObject(userResetPassword.GenerateRequestObjectArgs{
		Email:    userEmail,
		Password: secondPassword,
	})
	fmt.Printf("__execute__ userResetPasswordRequestDTO: %v\n", userResetPasswordRequestDTO)
	userResetPasswordResult, err := userResetPassword.Execute(userResetPasswordRequestDTO)
	if err != nil {
		err = fmt.Errorf("error userResetPassword: %s", err)
		fmt.Println(err)
		continueTesting = false
		return
	}
	continueTesting = CheckDatabaseForUserResetPassword(userResetPasswordRequestDTO, userResetPasswordResult)
	if !continueTesting {
		return
	}

	userEmailHasOtpDetailForResetPasswordFlow, _, err := dynamodbhelper.GetUserEmailHasOtpDetailByEmail(
		dynamodbhelper.CreateClientFromSession(),
		userEmail,
	)
	if err != nil {
		err = fmt.Errorf("error GetUserEmailHasOtpDetailByEmail forResetPasswordFlow: %s", err)
		fmt.Println(err)
		continueTesting = false
		return
	}
	otpForResetPasswordFlow := userEmailHasOtpDetailForResetPasswordFlow.OTP
	userValidateOtpRequestDTOForResetPasswordFlow := userValidateOtp.GenerateRequestObject(userValidateOtp.GenerateRequestObjectArgs{
		Email: userEmail,
		OTP:   otpForResetPasswordFlow,
	})
	fmt.Printf("__execute__ userValidateOtpRequestDTO forResetPasswordFlow: %v\n", userValidateOtpRequestDTOForResetPasswordFlow)
	userValidateOtpResultForResetPasswordFlow, err := userValidateOtp.Execute(userValidateOtpRequestDTOForResetPasswordFlow)
	if err != nil {
		err = fmt.Errorf("error userValidateOtp forResetPasswordFlow: %s", err)
		fmt.Println(err)
		continueTesting = false
		return
	}
	continueTesting = CheckDatabaseForUserValidateOtp(userValidateOtpRequestDTOForResetPasswordFlow, userValidateOtpResultForResetPasswordFlow, requesterUserId)
	if !continueTesting {
		return
	}
	defer func() {
		_, err := dynamodbhelper.DeleteUserListByFilter(
			dynamodbhelper.CreateClientFromSession(),
			[]string{userValidateOtpResultForResetPasswordFlow.ResponseData.Id},
		)
		if err == nil {
			fmt.Printf("__cleanup__ DeleteUserListByFilter successful\n\n")
		} else {
			fmt.Printf("error DeleteUserListByFilter: %v\n\n", err)
		}
	}()

	userSignInRequestDTO := userSignIn.GenerateRequestObject(userSignIn.GenerateRequestObjectArgs{
		Email:    userEmail,
		Password: secondPassword,
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
		Password: password,
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

	firstProductUrl := "https://mock.com/product/1"
	twiceInputProductUrl := "https://mock.com/product/2"
	thirdProductUrl := "https://mock.com/product/3"
	monitoredLinkAddMultipleRequestDTO := monitoredLinkAddMultiple.GenerateRequestObject(monitoredLinkAddMultiple.GenerateRequestObjectArgs{
		MonitoredLinkList: []monitoredLinkAddMultiple.MonitoredLinkRequestDTOV1{
			{HubMonitoredLinkUrl: firstProductUrl, AlertPrice: currencyutil.NewFromNumberString("50000", "IDR"), AlertMethodList: []constenum.AlertMethod{constenum.AlertMethodEmail, constenum.AlertMethodPushNotification}},
			{HubMonitoredLinkUrl: twiceInputProductUrl, AlertPrice: currencyutil.NewFromNumberString("40000", "IDR"), AlertMethodList: []constenum.AlertMethod{}},
		},
	})
	fmt.Printf("__execute__ monitoredLinkAddMultipleRequestDTO: %v\n", monitoredLinkAddMultipleRequestDTO)
	monitoredLinkAddMultipleResult, err := monitoredLinkAddMultiple.Execute(monitoredLinkAddMultipleRequestDTO, jwtToken)
	if err != nil {
		err = fmt.Errorf("error monitoredLinkAddMultiple: %s", err)
		fmt.Println(err)
		continueTesting = false
		return
	}
	continueTesting, monitoredLinkDAOList := CheckDatabaseForMonitoredLinkAddMultiple(monitoredLinkAddMultipleRequestDTO, monitoredLinkAddMultipleResult, requesterUserId)
	if !continueTesting {
		return
	}

	//check twice input the same url
	twiceInputProductCleanedUrl := monitoredLinkAddMultipleResult.ResponseData.MonitoredLinkRawToCleanedMap[twiceInputProductUrl]
	twiceInputProductUrlFirstExpireTime := time.Time{}
	for _, monitoredLinkDAO := range monitoredLinkDAOList {
		if monitoredLinkDAO.HubMonitoredLinkUrl == twiceInputProductCleanedUrl {
			twiceInputProductUrlFirstExpireTime = monitoredLinkDAO.TimeExpired
		}
	}
	if twiceInputProductUrlFirstExpireTime.IsZero() {
		err = fmt.Errorf("error twiceInputProductUrlFirstExpireTime isZero")
		fmt.Println(err)
		continueTesting = false
		return
	}
	monitoredLinkAddMultipleRequestDTO = monitoredLinkAddMultiple.GenerateRequestObject(monitoredLinkAddMultiple.GenerateRequestObjectArgs{
		MonitoredLinkList: []monitoredLinkAddMultiple.MonitoredLinkRequestDTOV1{
			{HubMonitoredLinkUrl: thirdProductUrl, AlertPrice: currencyutil.NewFromNumberString("50000", "IDR"), AlertMethodList: []constenum.AlertMethod{constenum.AlertMethodEmail, constenum.AlertMethodPushNotification}},
			{HubMonitoredLinkUrl: twiceInputProductUrl, AlertPrice: currencyutil.NewFromNumberString("40000", "IDR"), AlertMethodList: []constenum.AlertMethod{}},
		},
	})
	fmt.Printf("__execute__ monitoredLinkAddMultipleRequestDTO: %v\n", monitoredLinkAddMultipleRequestDTO)
	monitoredLinkAddMultipleResult, err = monitoredLinkAddMultiple.Execute(monitoredLinkAddMultipleRequestDTO, jwtToken)
	if err != nil {
		err = fmt.Errorf("error monitoredLinkAddMultiple: %s", err)
		fmt.Println(err)
		continueTesting = false
		return
	}
	continueTesting, monitoredLinkDAOList = CheckDatabaseForMonitoredLinkAddMultiple(monitoredLinkAddMultipleRequestDTO, monitoredLinkAddMultipleResult, requesterUserId)
	if !continueTesting {
		return
	}
	twiceInputProductUrlSecondExpireTime := time.Time{}
	for _, monitoredLinkDAO := range monitoredLinkDAOList {
		if monitoredLinkDAO.HubMonitoredLinkUrl == twiceInputProductCleanedUrl {
			twiceInputProductUrlSecondExpireTime = monitoredLinkDAO.TimeExpired
		}
	}
	if twiceInputProductUrlSecondExpireTime.IsZero() {
		err = fmt.Errorf("error twiceInputProductUrlSecondExpireTime isZero")
		fmt.Println(err)
		continueTesting = false
		return
	}
	if twiceInputProductUrlFirstExpireTime.Before(twiceInputProductUrlSecondExpireTime) {
		fmt.Printf("__success__ TimeExpired renewed in StlMonitoredLinkDetail\n\n")
	} else {
		err = fmt.Errorf("error TimeExpired in StlMonitoredLinkDetail is not renewed")
		fmt.Println(err)
		continueTesting = false
		return
	}

	defer func() {
		cleanedMonitoredLinkUrlList := []string{firstProductUrl, twiceInputProductUrl, thirdProductUrl}
		_, err := dynamodbhelper.DeleteStlMonitoredLinkDetailByUrlList(
			dynamodbhelper.CreateClientFromSession(),
			cleanedMonitoredLinkUrlList,
		)
		if err == nil {
			fmt.Printf("__cleanup__ DeleteStlMonitoredLinkDetailByUrlList successful\n\n")
		} else {
			fmt.Printf("error DeleteStlMonitoredLinkDetailByUrlList: %v\n\n", err)
		}
		_, err = dynamodbhelper.DeleteStlUserMonitorsLinkDetailByUrlList(
			dynamodbhelper.CreateClientFromSession(),
			requesterUserId,
			cleanedMonitoredLinkUrlList,
		)
		if err == nil {
			fmt.Printf("__cleanup__ DeleteStlUserMonitorsLinkDetailByUrlList successful\n\n")
		} else {
			fmt.Printf("error DeleteStlUserMonitorsLinkDetailByUrlList: %v\n\n", err)
		}
	}()

	//check monitoredLink.list
	monitoredLinkListRequestDTO := monitoredLinkList.DefaultRequestObject
	fmt.Printf("__execute__ monitoredLinkListRequestDTO: %v\n", monitoredLinkListRequestDTO)
	monitoredLinkListResult, err := monitoredLinkList.Execute(monitoredLinkListRequestDTO, jwtToken)
	if err != nil {
		err = fmt.Errorf("error monitoredLinkList: %s", err)
		fmt.Println(err)
		continueTesting = false
		return
	}
	if len(monitoredLinkListResult.ResponseData) != 3 {
		fmt.Printf("error monitoredLinkListResult length is: %v\n\n", len(monitoredLinkListResult.ResponseData))
		continueTesting = false
		return
	}
	monitoredLinkUrlSet := map[string]bool{
		firstProductUrl:      true,
		twiceInputProductUrl: true,
		thirdProductUrl:      true,
	}
	for _, userMonitorsLink := range monitoredLinkListResult.ResponseData {
		if _, ok := monitoredLinkUrlSet[userMonitorsLink.HubMonitoredLinkUrl]; !ok {
			fmt.Printf("error url not in input set: %v\n\n", userMonitorsLink.HubMonitoredLinkUrl)
			continueTesting = false
			return
		}
		if userMonitorsLink.HubUserId != requesterUserId {
			fmt.Printf("error fetch urls from other userId: %v\n\n", userMonitorsLink.HubUserId)
			continueTesting = false
			return
		}
	}
	fmt.Printf("__success__ testing monitoredLinkListResult\n\n")

	monitoredLinkEditMultipleRequestDTO := monitoredLinkEditMultiple.GenerateRequestObject(monitoredLinkEditMultiple.GenerateRequestObjectArgs{
		MonitoredLinkList: []monitoredLinkEditMultiple.MonitoredLinkRequestDTOV1{
			{HubMonitoredLinkUrl: firstProductUrl, AlertPrice: currencyutil.NewFromNumberString("150000", "IDR"), ActiveAlertMethodList: []constenum.AlertMethod{constenum.AlertMethodEmail}},
			{HubMonitoredLinkUrl: thirdProductUrl, AlertPrice: currencyutil.NewFromNumberString("150000", "IDR"), ActiveAlertMethodList: []constenum.AlertMethod{}},
		},
	})
	fmt.Printf("__execute__ monitoredLinkEditMultipleRequestDTO: %v\n", monitoredLinkEditMultipleRequestDTO)
	monitoredLinkEditMultipleResult, err := monitoredLinkEditMultiple.Execute(monitoredLinkEditMultipleRequestDTO, jwtToken)
	if err != nil {
		err = fmt.Errorf("error monitoredLinkEditMultiple: %s", err)
		fmt.Println(err)
		continueTesting = false
		return
	}
	continueTesting = CheckDatabaseForMonitoredLinkEditMultiple(monitoredLinkEditMultipleRequestDTO, monitoredLinkEditMultipleResult, requesterUserId)
	if !continueTesting {
		return
	}
	fmt.Printf("__success__ testing monitoredLinkEditMultipleResult\n\n")

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

func CheckDatabaseForMonitoredLinkAddMultiple(
	request monitoredLinkAddMultiple.RequestDTOV1,
	result monitoredLinkAddMultiple.ExecuteResult,
	requesterUserId string,
) (
	continueTesting bool,
	monitoredLinkDAOList []monitoredLink.StlMonitoredLinkDetailDAOV1,
) {
	//uncomment to disable this function
	// fmt.Printf("__disabled_checking__ should be commented out once all tests run successfully\n\n")
	// return

	monitoredLinkDAOList, err := monitoredLinkAddMultiple.CheckResultIsCorrect(request, result, requesterUserId) //error already printed in CheckResultIsCorrect
	continueTesting = err == nil
	return
}

func CheckDatabaseForMonitoredLinkEditMultiple(
	request monitoredLinkEditMultiple.RequestDTOV1,
	result monitoredLinkEditMultiple.ExecuteResult,
	requesterUserId string,
) (
	continueTesting bool,
) {
	//uncomment to disable this function
	// fmt.Printf("__disabled_checking__ should be commented out once all tests run successfully\n\n")
	// return

	err := monitoredLinkEditMultiple.CheckResultIsCorrect(request, result, requesterUserId) //error already printed in CheckResultIsCorrect
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
