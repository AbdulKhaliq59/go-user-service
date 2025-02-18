package services

import (
	"fmt"
	"go-user-service/models"
	"go-user-service/repository"
	"time"
)

type ReportService struct {
	ReportRepo *repository.ReportRepository
}

func NewReportService(reportRepo *repository.ReportRepository) *ReportService {
	return &ReportService{ReportRepo: reportRepo}
}

func (s *ReportService) GenerateReport(fromDate, toDate time.Time) (models.GeneralReportResponse, error) {
	// Calculate days difference
	daysDifference := int(toDate.Sub(fromDate).Hours() / 24)

	// Calculate previous date range
	previousToDate := fromDate
	previousFromDate := fromDate.AddDate(0, 0, -daysDifference)

	// Get all required data
	cumulativeActual, err := s.ReportRepo.CalculateCumulativeActual()
	if err != nil {
		return models.GeneralReportResponse{}, fmt.Errorf("error calculating cumulative actual: %v", err)
	}

	currentActual, err := s.ReportRepo.CalculateCurrentActual(fromDate, toDate)
	if err != nil {
		return models.GeneralReportResponse{}, fmt.Errorf("error calculating current actual: %v", err)
	}

	cumulativeExpected, err := s.ReportRepo.CalculateCumulativeExpected(fromDate, toDate)
	if err != nil {
		return models.GeneralReportResponse{}, fmt.Errorf("error calculating cumulative expected: %v", err)
	}

	previousTotal, err := s.ReportRepo.CalculatePreviousTotal(previousFromDate, previousToDate)
	if err != nil {
		return models.GeneralReportResponse{}, fmt.Errorf("error calculating previous total: %v", err)
	}

	telecomData, err := s.ReportRepo.CalculateTelecomData(fromDate, toDate)
	if err != nil {
		return models.GeneralReportResponse{}, fmt.Errorf("error calculating telecom data: %v", err)
	}

	// Build and return response
	response := models.GeneralReportResponse{}
	response.Cumulative.Actual = cumulativeActual
	response.Cumulative.Expected = cumulativeExpected
	response.Current.Actual = currentActual
	response.Previous.Total = previousTotal
	response.Telecoms = telecomData

	return response, nil
}

func (s *ReportService) TransactionReport(fromDate, toDate time.Time, productId, status string) (models.TransactionReportResponse, error) {
	return s.ReportRepo.TransactionReport(fromDate, toDate, productId, status)
}

func (s *ReportService) GenerateMonthlyReport() ([]models.MonthlyReportItem, error) {
	// Get monthly data from repository
	monthlyData, err := s.ReportRepo.GetMonthlyReport()
	if err != nil {
		return nil, fmt.Errorf("error fetching monthly report data: %w", err)
	}

	// Get latest expected values
	expectedValues, err := s.ReportRepo.GetLatestExpectedTransaction()
	if err != nil {
		return nil, fmt.Errorf("error fetching expected transaction values: %w", err)
	}

	// If no data, return empty array with a message
	if len(monthlyData) == 0 {
		return []models.MonthlyReportItem{}, nil
	}

	report := make([]models.MonthlyReportItem, 0, len(monthlyData))

	for _, row := range monthlyData {
		// Calculate days in month for expected values
		daysInMonth := time.Date(row.Month.Year(), row.Month.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()

		expectedRevenue := 0.0
		expectedTransactions := 0

		if expectedValues != nil {
			expectedRevenue = expectedValues.Revenue * float64(daysInMonth)
			expectedTransactions = expectedValues.Transactions * daysInMonth
		}

		monthDetails := models.MonthlyReportItem{
			Title: row.Month.Format("Jan"), // Format as short month name
			Actual: models.ReportItemDetails{
				Revenue:      row.Revenue,
				Transactions: row.Transactions,
			},
			Expected: models.ReportItemDetails{
				Revenue:      expectedRevenue,
				Transactions: expectedTransactions,
			},
		}

		report = append(report, monthDetails)
	}

	return report, nil
}

// GenerateWeeklyReport implements the weekly report logic
func (s *ReportService) GenerateWeeklyReport() ([]models.WeeklyReportItem, error) {
	// Get weekly data from repository
	weeklyData, err := s.ReportRepo.GetWeeklyReport()
	if err != nil {
		return nil, fmt.Errorf("error fetching weekly report data: %w", err)
	}

	// Get latest expected values
	expectedValues, err := s.ReportRepo.GetLatestExpectedTransaction()
	if err != nil {
		return nil, fmt.Errorf("error fetching expected transaction values: %w", err)
	}

	// If no data, return empty array
	if len(weeklyData) == 0 {
		return []models.WeeklyReportItem{}, nil
	}

	report := make([]models.WeeklyReportItem, 0, len(weeklyData))

	for _, row := range weeklyData {
		// For weekly reports, we use 7 days as multiplier
		const daysInWeek = 7

		expectedRevenue := 0.0
		expectedTransactions := 0

		if expectedValues != nil {
			expectedRevenue = expectedValues.Revenue * daysInWeek
			expectedTransactions = expectedValues.Transactions * daysInWeek
		}

		// Format the week title
		weekTitle := fmt.Sprintf("Week %s", row.Week.Format("Jan 2"))

		weekDetails := models.WeeklyReportItem{
			Title: weekTitle,
			Actual: models.ReportItemDetails{
				Revenue:      row.Revenue,
				Transactions: row.Transactions,
			},
			Expected: models.ReportItemDetails{
				Revenue:      expectedRevenue,
				Transactions: expectedTransactions,
			},
		}

		report = append(report, weekDetails)
	}

	return report, nil
}

func (s *ReportService) AddExpectedTransaction(expected models.ExpectedTransaction) error {
	return s.ReportRepo.AddExpectedTransaction(expected)
}

func (s *ReportService) GetTopTenTransactions() ([]models.TopTenTransactionResponse, error) {
	return s.ReportRepo.GetTopTenTransactions()
}

func (s *ReportService) GetTopTenRecentTransactions() ([]models.RecentTransaction, error) {
	return s.ReportRepo.GetTopTenRecentTransactions()
}

// GetDailyTransactionStatistics fetches daily transaction statistics
func (s *ReportService) GetDailyTransactionStatistics() (models.DailyTransactionStats, error) {
	return s.ReportRepo.GetDailyTransactionStatistics()
}

func (s *ReportService) GetDailyTransactionReport() ([]models.DailyTransactionReport, error) {
	return s.ReportRepo.GetDailyTransactionReport()
}

func (s *ReportService) GetMonthlyTransactionReport() ([]models.MonthlyTransactionReport, error) {
	return s.ReportRepo.GetMonthlyTransactionReport()
}

func (s *ReportService) GetHourlyTransactionReport(date string) ([]models.HourlyTransactionReport, error) {
	return s.ReportRepo.GetHourlyTransactionReport(date)
}
func (s *ReportService) GetProductMonthlyReport(productID string) ([]models.ProductTransactionReport, error) {
	return s.ReportRepo.GetProductMonthlyReport(productID)
}

func (s *ReportService) GetProductDailyReport(productID, date string) ([]models.ProductTransactionReport, error) {
	return s.ReportRepo.GetProductDailyReport(productID, date)
}

func (s *ReportService) GetProductHourlyReport(productID, date, fromDate, toDate, orderBy string) ([]models.ProductTransactionReport, error) {
	return s.ReportRepo.GetProductHourlyReport(productID, date, fromDate, toDate, orderBy)
}

func Paginate[T any](data []T, page, perPage int) models.PaginationResponse {
	total := len(data)
	lastPage := (total + perPage - 1) / perPage
	currentPage := min(max(1, page), lastPage)
	previousPage := currentPage - 1
	nextPage := currentPage + 1
	if nextPage > lastPage {
		nextPage = 0
	}

	startIndex := (currentPage - 1) * perPage
	endIndex := min(startIndex+perPage, total)

	return models.PaginationResponse{
		Status:       "success",
		List:         data[startIndex:endIndex],
		Total:        total,
		PreviousPage: &previousPage,
		NextPage:     &nextPage,
		LastPage:     lastPage,
		CurrentPage:  currentPage,
	}
}
