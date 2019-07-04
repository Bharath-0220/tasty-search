package review

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var ReviewList ReviewCollection

var WordIndexes map[string][]*Review

func (a ListReviews) Len() int           { return len(a) }
func (a ListReviews) Less(i, j int) bool { return a[i].Score > a[j].Score }
func (a ListReviews) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func ParseReviews() {

	var line string
	var err error

	WordIndexes = make(map[string][]*Review)

	f, errs := os.Open("finefoods.txt")

	if errs != nil {
		fmt.Println("error", errs.Error())
	}

	reader := bufio.NewReader(f)

	review := new(Review)

	count := 0
	for err == nil || err != io.EOF {
		count++
		line, err = reader.ReadString('\n')

		if err != nil {
			continue
		}

		if len(line) == 1 {
			ReviewList.ReviewList = append(ReviewList.ReviewList, review)
			review = new(Review)
			count = 0
			continue
		}

		temp := strings.TrimSuffix(strings.Split(line, ": ")[1], "\n")

		switch count {
		case 1:
			review.ProductId = temp
		case 2:
			review.UserId = temp
		case 3:
			review.ProfileName = temp
		case 4:
			review.HelpFullNess = temp
		case 5:
			review.Score, _ = strconv.ParseFloat(temp, 64)
		case 6:
			temp, _ := strconv.Atoi(temp)
			review.TimeStamp = int64(temp)
		case 7:
			review.Summary = temp
			review.buildIndexesForSummary()
		case 8:
			review.Text = temp
			review.buildIndexesForText()
		}

	}
}

func (review *Review) buildIndexesForText() {

	wordCheck := make(map[string]bool)

	reg, _ := regexp.Compile("[^a-zA-Z]+")

	wordsSplit := strings.Split(review.Text, " ")

	for _, word := range wordsSplit {

		word = strings.ToLower(word)
		word = reg.ReplaceAllString(word, "")

		if word == "" {
			continue
		}

		if _, ok := WordIndexes[word]; !ok {
			WordIndexes[word] = make([]*Review, 0)
		}

		if _, ok := wordCheck[word]; ok {
			continue
		}
		WordIndexes[word] = append(WordIndexes[word], review)

		wordCheck[word] = true
	}
}

func (review *Review) buildIndexesForSummary() {

	wordCheck := make(map[string]bool)

	reg, _ := regexp.Compile("[^a-zA-Z]+")

	wordsSplit := strings.Split(review.Text, " ")

	for _, word := range wordsSplit {

		word = strings.ToLower(word)

		word = reg.ReplaceAllString(word, "")

		if word == "" {
			continue
		}

		if _, ok := WordIndexes[word]; !ok {
			WordIndexes[word] = make([]*Review, 0)
		}

		if _, ok := wordCheck[word]; ok {
			continue
		}
		WordIndexes[word] = append(WordIndexes[word], review)

		wordCheck[word] = true
	}
}

func SortWordIndexes() {

	var ListReview ListReviews
	for _, wordList := range WordIndexes {

		ListReview = wordList
		sort.Sort(ListReview)
	}
}

func SearchTopReviews(c *gin.Context) ListReviews {

	req := c.Query("s")

	r := &Request{}
	_ = json.Unmarshal([]byte(req), r)

	var list ListReviews
	for _, s := range r.Tokens {

		s = strings.ToLower(s)
		if val, ok := WordIndexes[s]; ok {
			list = append(list, val...)
		}
	}

	reviewCheck := make(map[string]int)
	reviewToPointer := make(map[string]*Review)

	for _, l := range list {

		if _, ok := reviewCheck[l.Summary]; !ok {
			reviewCheck[l.Summary] = 0
			reviewToPointer[l.Summary] = l
		}

		reviewCheck[l.Summary]++
	}

	var a []int
	countToReviewMap := make(map[int][]string)

	for k, v := range reviewCheck {

		if _, ok := countToReviewMap[v]; !ok {
			countToReviewMap[v] = []string{k}
			a = append(a, v)
		} else {
			countToReviewMap[v] = append(countToReviewMap[v], k)
		}
	}

	sort.Sort(sort.Reverse(sort.IntSlice(a)))

	result := []string{}
	for _, v := range a {
		val, _ := countToReviewMap[v]
		result = append(result, val...)
	}

	if len(result) > 20 {
		result = result[:20]
	}

	var finalResult ListReviews

	for _, s := range result {

		finalResult = append(finalResult, reviewToPointer[s])
	}

	return finalResult
}
