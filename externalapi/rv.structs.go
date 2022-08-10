package externalapi

type RVT_CreateTestCase struct {
	Url string `json:"url"`
}

type RVT_Inst struct {
	Id  string `json:"id"`
	Url string `json:"url"`
}

type RVT_ListRvts struct {
	Rvts   []RVT_Inst       `json:"rvts"`
	Status FdoConfApiStatus `json:"status"`
}
