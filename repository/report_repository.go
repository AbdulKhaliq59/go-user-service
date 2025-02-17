package repository

import (
	"encoding/json"
	"go-user-service/models"
	"time"

	"gorm.io/gorm"
)

type ReportRepository struct {
	DB *gorm.DB
}

func NewReportRepository(db *gorm.DB) *ReportRepository {
	return &ReportRepository{DB: db}
}

func (r *ReportRepository) CalculateCumulativeActual() (models.CumulativeResponse, error) {
	var result struct {
		Revenue      float64
		Transactions int
	}

	if err := r.DB.Table("transaction_report_telecom").
		Select("SUM(revenue) as revenue, SUM(transactions) as transactions").
		Scan(&result).Error; err != nil {
		return models.CumulativeResponse{}, err
	}

	return models.CumulativeResponse{
		Revenue:      result.Revenue,
		Transactions: result.Transactions,
	}, nil
}

func (r *ReportRepository) CalculateCurrentActual(fromDate, toDate time.Time) (models.CumulativeResponse, error) {
	var result struct {
		Revenue      float64
		Transactions int
	}

	if err := r.DB.Table("transaction_report_telecom").
		Select("SUM(revenue) as revenue, SUM(transactions) as transactions").
		Where("created_at BETWEEN ? AND ?", fromDate, toDate).
		Scan(&result).Error; err != nil {
		return models.CumulativeResponse{}, err
	}

	return models.CumulativeResponse{
		Revenue:      result.Revenue,
		Transactions: result.Transactions,
	}, nil
}

func (r *ReportRepository) CalculateCumulativeExpected(fromDate, toDate time.Time) (models.CumulativeResponse, error) {
	var earliestTransaction models.TelecomReport
	if err := r.DB.Order("created_at ASC").First(&earliestTransaction).Error; err != nil {
		return models.CumulativeResponse{}, err
	}

	var latestTransaction models.TelecomReport
	if err := r.DB.Order("created_at DESC").First(&latestTransaction).Error; err != nil {
		return models.CumulativeResponse{}, err
	}

	daysDifference := int(latestTransaction.CreatedAt.Sub(earliestTransaction.CreatedAt).Hours() / 24)

	var expectedValues models.ExpectedTransaction
	if err := r.DB.Order("created_at DESC").First(&expectedValues).Error; err != nil {
		return models.CumulativeResponse{}, err
	}

	expectedRevenue := expectedValues.Revenue * float64(daysDifference)
	expectedTransactions := expectedValues.Transactions * daysDifference

	return models.CumulativeResponse{
		Revenue:      expectedRevenue,
		Transactions: expectedTransactions,
	}, nil
}

func (r *ReportRepository) CalculatePreviousTotal(fromDate, toDate time.Time) (models.CumulativeResponse, error) {
	var result struct {
		Revenue      float64
		Transactions int
	}

	if err := r.DB.Table("transaction_report_telecom").
		Select("SUM(revenue) as revenue, SUM(transactions) as transactions").
		Where("created_at BETWEEN ? AND ?", fromDate, toDate).
		Scan(&result).Error; err != nil {
		return models.CumulativeResponse{}, err
	}

	return models.CumulativeResponse{
		Revenue:      result.Revenue,
		Transactions: result.Transactions,
	}, nil
}

func (r *ReportRepository) CalculateTelecomData(fromDate, toDate time.Time) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	rows, err := r.DB.Table("transaction_report_telecom").
		Select("telecom, SUM(revenue) as revenue, SUM(transactions) as transactions").
		Where("created_at BETWEEN ? AND ?", fromDate, toDate).
		Group("telecom").
		Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var telecom string
		var revenue float64
		var transactions int

		if err := rows.Scan(&telecom, &revenue, &transactions); err != nil {
			return nil, err
		}

		results = append(results, map[string]interface{}{
			"telecom":      telecom,
			"revenue":      revenue,
			"transactions": transactions,
		})
	}

	return results, nil
}

func (r *ReportRepository) TransactionReport(fromDate, toDate time.Time, productId, status string) (models.TransactionReportResponse, error) {
	// Initialize the response
	response := models.TransactionReportResponse{}

	// Build the base query for revenue and transaction count
	baseQuery := "SELECT SUM(final_amount) as revenue, COUNT(id) as transactions FROM transactions_new WHERE 1=1"
	baseQueryParams := []interface{}{}

	// Build the unique hits query
	uniqueQuery := "SELECT COUNT(DISTINCT phone_number) as unique_hits FROM transactions_new WHERE 1=1"
	uniqueQueryParams := []interface{}{}

	// Add product filter if provided
	if productId != "" {
		baseQuery += " AND product_id IN (SELECT id FROM products WHERE product_id = ?)"
		baseQueryParams = append(baseQueryParams, productId)

		uniqueQuery += " AND product_id IN (SELECT id FROM products WHERE product_id = ?)"
		uniqueQueryParams = append(uniqueQueryParams, productId)
	}

	// Add status filter if provided
	if status != "" {
		baseQuery += " AND status = ?"
		baseQueryParams = append(baseQueryParams, status)

		uniqueQuery += " AND status = ?"
		uniqueQueryParams = append(uniqueQueryParams, status)
	}

	// Execute base query for cumulative results
	var cumulativeResult struct {
		Revenue      float64
		Transactions int
	}
	if err := r.DB.Raw(baseQuery, baseQueryParams...).Scan(&cumulativeResult).Error; err != nil {
		return response, err
	}

	// Execute base query with date filter for current results
	currentBaseQuery := baseQuery + " AND created_at BETWEEN ? AND ?"
	currentBaseQueryParams := append(baseQueryParams, fromDate, toDate)

	var currentResult struct {
		Revenue      float64
		Transactions int
	}
	if err := r.DB.Raw(currentBaseQuery, currentBaseQueryParams...).Scan(&currentResult).Error; err != nil {
		return response, err
	}

	// Execute unique query for cumulative unique hits
	var cumulativeUniqueResult struct {
		UniqueHits int
	}
	if err := r.DB.Raw(uniqueQuery, uniqueQueryParams...).Scan(&cumulativeUniqueResult).Error; err != nil {
		return response, err
	}

	// Execute unique query with date filter for current unique hits
	currentUniqueQuery := uniqueQuery + " AND created_at BETWEEN ? AND ?"
	currentUniqueQueryParams := append(uniqueQueryParams, fromDate, toDate)

	var currentUniqueResult struct {
		UniqueHits int
	}
	if err := r.DB.Raw(currentUniqueQuery, currentUniqueQueryParams...).Scan(&currentUniqueResult).Error; err != nil {
		return response, err
	}

	// Construct the final response
	result := map[string]interface{}{
		"total_transaction": map[string]interface{}{
			"cumulative": map[string]interface{}{
				"revenue":     cumulativeResult.Revenue,
				"transaction": cumulativeResult.Transactions,
			},
			"current": map[string]interface{}{
				"revenue":     currentResult.Revenue,
				"transaction": currentResult.Transactions,
			},
		},
		"unique": map[string]interface{}{
			"cumulative": map[string]interface{}{
				"uniqueHits": cumulativeUniqueResult.UniqueHits,
			},
			"current": map[string]interface{}{
				"uniqueHits": currentUniqueResult.UniqueHits,
			},
		},
	}

	// Convert to JSON
	jsonData, err := json.Marshal(result)
	if err != nil {
		return response, err
	}

	// Parse the JSON back into our response structure
	if err := json.Unmarshal(jsonData, &response); err != nil {
		return response, err
	}

	return response, nil
}

func (r *ReportRepository) GetLatestExpectedTransaction() (*models.ExpectedTransaction, error) {
	var expected models.ExpectedTransaction
	result := r.DB.Order("created_at DESC").First(&expected)
	if result.Error != nil {
		// If no record is found, return nil instead of error
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &expected, nil
}

// Monthly report implementation
type MonthlyReportResult struct {
	Month        time.Time `json:"month"`
	Revenue      float64   `json:"revenue"`
	Transactions int       `json:"transactions"`
}

func (r *ReportRepository) GetMonthlyReport() ([]MonthlyReportResult, error) {
	var results []MonthlyReportResult

	// Corrected SQL query
	query := `
		SELECT 
			DATE_TRUNC('month', created_at) AS month,
			COALESCE(SUM(final_amount), 0) AS revenue,
			COUNT(*) AS transactions
		FROM
			transactions_new
		WHERE
			status = 'SUCCESS'
		GROUP BY
			month
		ORDER BY
			month
	`

	err := r.DB.Raw(query).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	// If no results, return an empty array instead of nil
	if len(results) == 0 {
		return []MonthlyReportResult{}, nil
	}

	return results, nil
}

// Weekly report implementation
type WeeklyReportResult struct {
	Week         time.Time `json:"week"`
	Revenue      float64   `json:"revenue"`
	Transactions int       `json:"transactions"`
}

func (r *ReportRepository) GetWeeklyReport() ([]WeeklyReportResult, error) {
	var results []WeeklyReportResult

	// Corrected SQL query
	query := `
		SELECT 
			DATE_TRUNC('week', created_at) AS week,
			COALESCE(SUM(final_amount), 0) AS revenue,
			COUNT(*) AS transactions
		FROM
			transactions_new
		WHERE
			status = 'SUCCESS'
		GROUP BY
			week
		ORDER BY
			week
	`

	err := r.DB.Raw(query).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	// If no results, return an empty array instead of nil
	if len(results) == 0 {
		return []WeeklyReportResult{}, nil
	}

	return results, nil
}

func (r *ReportRepository) AddExpectedTransaction(expected models.ExpectedTransaction) error {
	return r.DB.Create(&expected).Error
}

func (r *ReportRepository) GetTopTenTransactions() ([]models.TopTenTransactionResponse, error) {
	var topTransactions []models.TopTenTransactionResponse
	err := r.DB.Table("transactions_new").
		Select("phone_number, SUM(final_amount) as amount, COUNT(phone_number) as count").
		Group("phone_number").
		Order("amount DESC").
		Limit(10).
		Scan(&topTransactions).Error

	return topTransactions, err
}

// GetTopTenRecentTransactions fetches the most recent 10 successful transactions
func (r *ReportRepository) GetTopTenRecentTransactions() ([]models.RecentTransaction, error) {
	var recentTransactions []models.RecentTransaction

	err := r.DB.Table("transactions_new").
		Select("phone_number, COALESCE(final_amount, 0) as amount, COALESCE(status, '') as status, created_at").
		Where("status = ?", "SUCCESS").
		Order("created_at DESC").
		Limit(10).
		Scan(&recentTransactions).Error

	if err != nil {
		return nil, err
	}

	return recentTransactions, nil
}

// GetDailyTransactionStatistics calculates daily transaction statistics
func (r *ReportRepository) GetDailyTransactionStatistics() (models.DailyTransactionStats, error) {
	var stats models.DailyTransactionStats

	startDate := time.Now().Truncate(24 * time.Hour)
	endDate := startDate.Add(24 * time.Hour)

	// Query total transactions (handle NULL with COALESCE)
	err := r.DB.Table("transactions_new").
		Select("COALESCE(SUM(final_amount), 0) as total_amount").
		Where("created_at >= ? AND created_at < ?", startDate, endDate).
		Scan(&stats.TotalAmount).Error

	if err != nil {
		return stats, err
	}

	// Query successful transactions (handle NULL with COALESCE)
	err = r.DB.Table("transactions_new").
		Select("COALESCE(SUM(final_amount), 0) as total_successful_amount").
		Where("created_at >= ? AND created_at < ? AND status = ?", startDate, endDate, "SUCCESS").
		Scan(&stats.TotalSuccessfulAmount).Error

	if err != nil {
		return stats, err
	}

	// Query failed transactions (handle NULL with COALESCE)
	err = r.DB.Table("transactions_new").
		Select("COALESCE(SUM(final_amount), 0) as total_failed_amount").
		Where("created_at >= ? AND created_at < ? AND status != ?", startDate, endDate, "SUCCESS").
		Scan(&stats.TotalFailedAmount).Error

	return stats, err
}

func (r *ReportRepository) GetDailyTransactionReport() ([]models.DailyTransactionReport, error) {
	var report []models.DailyTransactionReport
	query := `
        SELECT DATE(created_at) AS trx_date,
               SUM(CASE WHEN status = 'SUCCESS' THEN final_amount ELSE 0 END) AS tot_successful_trx_amount,
               SUM(CASE WHEN status = 'PENDING' THEN final_amount ELSE 0 END) AS tot_pending_trx_amount,
               SUM(CASE WHEN status = 'FAILED' THEN final_amount ELSE 0 END) AS tot_failed_trx_amount
        FROM transactions_new
        GROUP BY DATE(created_at)
        ORDER BY DATE(created_at) DESC;
    `
	err := r.DB.Raw(query).Scan(&report).Error
	return report, err
}

// Get Monthly Transaction Report
func (r *ReportRepository) GetMonthlyTransactionReport() ([]models.MonthlyTransactionReport, error) {
	var report []models.MonthlyTransactionReport
	query := `
        SELECT EXTRACT(YEAR FROM created_at) AS year,
               TO_CHAR(created_at, 'FMMonth') AS month,
               SUM(CASE WHEN status = 'SUCCESS' THEN final_amount ELSE 0 END) AS tot_successful_trx_amount
        FROM transactions_new
        GROUP BY year, month
        ORDER BY year DESC, month DESC;
    `
	err := r.DB.Raw(query).Scan(&report).Error
	return report, err
}

// Get Hourly Transaction Report
func (r *ReportRepository) GetHourlyTransactionReport() ([]models.HourlyTransactionReport, error) {
	var report []models.HourlyTransactionReport
	query := `
        SELECT DATE(created_at) AS trx_date,
               TO_CHAR(EXTRACT(HOUR FROM created_at), 'FM00') || ':00 - ' ||
               TO_CHAR(EXTRACT(HOUR FROM created_at) + 1, 'FM00') || ':00' AS hour,
               SUM(CASE WHEN status = 'SUCCESS' THEN final_amount ELSE 0 END) AS tot_successful_trx_amount,
               SUM(CASE WHEN status = 'PENDING' THEN final_amount ELSE 0 END) AS tot_pending_trx_amount,
               SUM(CASE WHEN status = 'FAILED' THEN final_amount ELSE 0 END) AS tot_failed_trx_amount
        FROM transactions_new
        GROUP BY trx_date, hour
        ORDER BY trx_date DESC, hour DESC;
    `
	err := r.DB.Raw(query).Scan(&report).Error
	return report, err
}
