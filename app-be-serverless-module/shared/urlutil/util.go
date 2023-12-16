package urlutil

func CleanUrl(rawUrl string) string {
	//remove irrelevant query params ~kodok
	return rawUrl
}

func CleanUrlList(rawUrlList []string) (cleanedUrlList []string) {
	for _, rawUrl := range rawUrlList {
		cleanedUrlList = append(cleanedUrlList, CleanUrl(rawUrl))
	}
	return cleanedUrlList
}
