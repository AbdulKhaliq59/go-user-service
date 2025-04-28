package utils

import (
	"go-user-service/models"
	"math"
	"sort"
	"time"
)

type PaginationQuery struct {
	PageNumber int    `form:"pageNumber"`
	PageSize   int    `form:"pageSize"`
	From       string `form:"from"`
	To         string `form:"to"`
	Status     string `form:"status"`
	Telecom    string `form:"telecom"`
	ProductID  string `form:"productId"`
	AgentCode  string `form:"agentCode"`
	Search     string `form:"search"`
}

type PaginationResponse struct {
	Status       string      `json:"status"`
	Data         interface{} `json:"data"`
	Total        int64       `json:"total"`
	CurrentPage  int         `json:"currentPage"`
	PreviousPage *int        `json:"previousPage"`
	NextPage     *int        `json:"nextPage"`
	LastPage     int         `json:"lastPage"`
	PageSize     int         `json:"pageSize"`
}

func CreatePaginationResponse(data interface{}, total int64, query PaginationQuery) PaginationResponse {
	// Set default values if not provided
	if query.PageSize == 0 {
		query.PageSize = 10
	}
	if query.PageNumber == 0 {
		query.PageNumber = 1
	}

	lastPage := int(math.Ceil(float64(total) / float64(query.PageSize)))

	var prevPage *int
	if query.PageNumber > 1 {
		prev := query.PageNumber - 1
		prevPage = &prev
	}

	var nextPage *int
	if query.PageNumber < lastPage {
		next := query.PageNumber + 1
		nextPage = &next
	}

	return PaginationResponse{
		Status:       "success",
		Data:         data,
		Total:        total,
		CurrentPage:  query.PageNumber,
		PageSize:     query.PageSize,
		PreviousPage: prevPage,
		NextPage:     nextPage,
		LastPage:     lastPage,
	}
}

func SortProductStatsByPercentage(stats []models.ProductStats) {
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Percentage > stats[j].Percentage
	})
}

func GetPaginationOffset(pageNumber, pageSize int) int {
	return (pageNumber - 1) * pageSize
}

// GetPaginationLimit calculates the ending index for pagination
func GetPaginationLimit(pageNumber, pageSize, total int) int {
	end := GetPaginationOffset(pageNumber, pageSize) + pageSize
	if end > total {
		return total
	}
	return end
}

// PaginationHelper provides pagination functionality
type PaginationHelper struct{}

func NewPaginationHelper() *PaginationHelper {
	return &PaginationHelper{}
}

// Paginate applies pagination to a slice of data
func (p *PaginationHelper) Paginate(data interface{}, pageNumber, pageSize int, total int64) PaginationResponse {
	query := PaginationQuery{
		PageNumber: pageNumber,
		PageSize:   pageSize,
	}
	return CreatePaginationResponse(data, total, query)
}

func GetFormattedTimestamp() string {
	return time.Now().Format("2006-01-02T15:04:05Z07:00")
}
