package review

type Review struct {
	ProductId    string
	UserId       string
	ProfileName  string
	HelpFullNess string
	Score        float64
	TimeStamp    int64
	Summary      string
	Text         string
}

type ReviewCollection struct {
	ReviewList []*Review
}

type ListReviews []*Review

type Request struct {
	Tokens []string `json:"tokens"`
}


