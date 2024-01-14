package searchedItemEditMultiple

import (
	"context"
	"diskon-hunter/price-monitoring/shared/constenum"
	"diskon-hunter/price-monitoring/shared/createerror"
	"diskon-hunter/price-monitoring/shared/currencyutil"
	"diskon-hunter/price-monitoring/shared/dynamodbhelper"
	"diskon-hunter/price-monitoring/shared/lazylogger"
	"diskon-hunter/price-monitoring/shared/serverresponse"
	"diskon-hunter/price-monitoring/src/searchedItem"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

/*
Commands represent input from client through API requests.
Addition, change, or removal of struct fields might cause version increment
*/
type CommandV1 struct {
	Version          string
	RequesterUserId  string
	SearchedItemList []SearchedItemDetailCommandV1
}

func NewCommandV1(
	Version string, //should follow the struct name suffix
	RequesterUserId string,
	SearchedItemList []SearchedItemDetailCommandV1,
) CommandV1 {
	return CommandV1{
		Version:          Version,
		RequesterUserId:  RequesterUserId,
		SearchedItemList: SearchedItemList,
	}
}

type SearchedItemDetailCommandV1 struct {
	HubSearchedItemId string
	Name              string
	Description       string
	AlertPrice        currencyutil.Currency
}

func NewSearchedItemDetailCommandV1(
	HubSearchedItemId string,
	Name string,
	Description string,
	AlertPrice currencyutil.Currency,
) SearchedItemDetailCommandV1 {
	return SearchedItemDetailCommandV1{
		HubSearchedItemId: HubSearchedItemId,
		Name:              Name,
		Description:       Description,
		AlertPrice:        AlertPrice,
	}
}

func (x CommandV1) createLoggableString() (string, error) {
	//strip any sensitive information.
	//strip any fields that are too large to be printed (e.g. image blob).
	loggableCommand := x //no sensitive info and no large fields so we'll just use x
	byteSlice, err := json.Marshal(loggableCommand)
	if err != nil {
		return "", err
	} else {
		return string(byteSlice), nil
	}
}

type CommandV1Dependencies struct {
	Logger         *lazylogger.Instance
	DynamoDBClient *dynamodb.DynamoDB
}

type CommandV1DataResponse = []searchedItem.StlUserSearchesItemDetailDAOV1

func CommandV1Handler(
	ctx context.Context,
	dependencies CommandV1Dependencies,
	command CommandV1,
) (CommandV1DataResponse, *serverresponse.ErrorObj) {
	//don't mutate this. emptyResponse should be used when returning error.
	emptyResponse := CommandV1DataResponse{}

	//log the command
	loggableCommand, err := command.createLoggableString()
	if err != nil {
		err = fmt.Errorf("error creating loggable string: %v", err)
		dependencies.Logger.EnqueueErrorLog(err, true)
		return emptyResponse, createerror.InternalException(err)
	}
	dependencies.Logger.EnqueueCommandLog(loggableCommand, true)

	/* Validations
	Validations from auth, write model,
	or domain model's business logic (from projections or from events replay).
	*/

	for i := 0; i < len(command.SearchedItemList); i++ {
		if len(command.SearchedItemList[i].Description) > constenum.MaxCharForSearchedItemDescription {
			err = fmt.Errorf("error description length of %v exceed max length of %v", len(command.SearchedItemList[i].Description), constenum.MaxCharForSearchedItemDescription)
			return emptyResponse, createerror.ClientBadRequest(err)
		}
	}

	/* Business Logic
	Perform business logic preferably through domain model's methods.
	*/

	searchedItemIdMapToNewDetail := map[string]SearchedItemDetailCommandV1{}
	idList := []string{}
	for _, searchedItemData := range command.SearchedItemList {
		searchedItemIdMapToNewDetail[searchedItemData.HubSearchedItemId] = searchedItemData
		idList = append(idList, searchedItemData.HubSearchedItemId)
	}
	existingSearchedItemList, errObj, err := dynamodbhelper.GetStlUserSearchesItemDetailList(
		dependencies.DynamoDBClient,
		command.RequesterUserId,
		idList,
	)
	if errObj != nil {
		//error already well described on the calling method
		dependencies.Logger.EnqueueErrorLog(err, true)
		return emptyResponse, errObj
	}

	transactItems := []*dynamodb.TransactWriteItem{}
	editedSearchedItemDAOList := []searchedItem.StlUserSearchesItemDetailDAOV1{}
	for i := 0; i < len(existingSearchedItemList); i++ {
		id := existingSearchedItemList[i].HubSearchedItemId

		//use existingDAO[i].field for fields that can't be changed
		//use idMapToNewDetail[id].field for fields that can be changed
		editedSearchedItemDAO := searchedItem.NewStlUserSearchesItemDetailDAOV1(
			existingSearchedItemList[i].HubUserId,         //HubUserId
			existingSearchedItemList[i].HubSearchedItemId, //HubSearchedItemId
			searchedItemIdMapToNewDetail[id].Name,         //Name
			searchedItemIdMapToNewDetail[id].Description,  //Description
			existingSearchedItemList[i].Status,            //Status
			searchedItemIdMapToNewDetail[id].AlertPrice,   //AlertPrice
			existingSearchedItemList[i].TimeExpired,       //TimeExpired
		)
		editedSearchedItemDAOList = append(editedSearchedItemDAOList, editedSearchedItemDAO)
		editedSearchedItemDAOItem, err := dynamodbattribute.MarshalMap(editedSearchedItemDAO)
		if err != nil {
			err = fmt.Errorf("error marshaling editedSearchedItemDAO: %v", err)
			dependencies.Logger.EnqueueErrorLog(err, true)
			return emptyResponse, createerror.InternalException(err)
		}
		transactItems = append(transactItems, &dynamodb.TransactWriteItem{
			Put: &dynamodb.Put{
				Item:      editedSearchedItemDAOItem,
				TableName: aws.String(searchedItem.GetStlUserSearchesItemDetailDynamoDBTableV1()),
			},
		})
	}

	/* Persisting Data
	Persist event to event store.
	If write model is used, also persist write model with atomic transaction.
	*/

	errObj, err = dynamodbhelper.TransactWriteItemsInWaves(dependencies.DynamoDBClient, transactItems)
	if errObj != nil {
		//error already well described on the calling method
		dependencies.Logger.EnqueueErrorLog(err, true)
		return emptyResponse, errObj
	}

	return editedSearchedItemDAOList, nil
}
