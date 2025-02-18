package models

import (
	"time"
)

type ExpectedTransaction struct {
	ID           string    `json:"id" gorm:"primaryKey"`
	Revenue      float64   `json:"revenue" binding:"required"`
	Transactions int       `json:"transactions" binding:"required"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
}

type TelecomReport struct {
	ID           string    `gorm:"primaryKey;type:uuid" json:"id"`
	Telecom      string    `json:"telecom"`
	Revenue      float64   `gorm:"default:0" json:"revenue"`
	Transactions int       `gorm:"default:0" json:"transactions"`
	UniqueHits   int       `gorm:"column:uniqueHits;default:0" json:"uniqueHits"`
	FromDate     time.Time `json:"fromDate"`
	ToDate       time.Time `json:"toDate"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
type ReportItemDetails struct {
	Revenue      float64 `json:"revenue"`
	Transactions int     `json:"transactions"`
}

type ExpectedTransactionRequest struct {
	Revenue      float64 `json:"revenue" binding:"required"`
	Transactions int     `json:"transactions" binding:"required"`
}

type TopTenTransactionResponse struct {
	PhoneNumber string  `json:"phone_number"`
	Amount      float64 `json:"amount"`
	Count       int     `json:"count"`
}

type MonthlyReportItem struct {
	Title    string            `json:"title"` // Month name (e.g., "Jan")
	Actual   ReportItemDetails `json:"actual"`
	Expected ReportItemDetails `json:"expected"`
}

type WeeklyReportItem struct {
	Title    string            `json:"title"` // Week description (e.g., "Week Jan 1")
	Actual   ReportItemDetails `json:"actual"`
	Expected ReportItemDetails `json:"expected"`
}

func (ExpectedTransaction) TableName() string {
	return "expected_transaction"
}

func (TelecomReport) TableName() string {
	return "transaction_report_telecom"
}

// Response DTOs
type CumulativeResponse struct {
	Revenue      float64 `json:"revenue"`
	Transactions int     `json:"transactions"`
}

type GeneralReportResponse struct {
	Cumulative struct {
		Actual   CumulativeResponse `json:"actual"`
		Expected CumulativeResponse `json:"expected"`
	} `json:"cumulative"`
	Current struct {
		Actual CumulativeResponse `json:"actual"`
	} `json:"current"`
	Previous struct {
		Total CumulativeResponse `json:"total"`
	} `json:"previous"`
	Telecoms []map[string]interface{} `json:"telecoms"`
}

type TransactionReportResponse struct {
	TotalTransaction struct {
		Cumulative struct {
			Revenue     float64 `json:"revenue"`
			Transaction int     `json:"transaction"`
		} `json:"cumulative"`
		Current struct {
			Revenue     float64 `json:"revenue"`
			Transaction int     `json:"transaction"`
		} `json:"current"`
	} `json:"total_transaction"`
	Unique struct {
		Cumulative struct {
			UniqueHits int `json:"uniqueHits"`
		} `json:"cumulative"`
		Current struct {
			UniqueHits int `json:"uniqueHits"`
		} `json:"current"`
	} `json:"unique"`
}

type RecentTransaction struct {
	PhoneNumber string    `json:"phone_number"`
	Amount      float64   `json:"amount"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

// DailyTransactionStats represents the structure for daily transaction statistics
type DailyTransactionStats struct {
	TotalAmount           float64 `json:"total_amount"`
	TotalSuccessfulAmount float64 `json:"total_successful_amount"`
	TotalFailedAmount     float64 `json:"total_failed_amount"`
}

type Report struct {
	Draw                  interface{} `json:"draw"`
	Margin                float64     `json:"margin"`
	ProductExpectedAmount float64     `json:"product_expected_amount"`
	IncomingAmount        float64     `json:"incoming_amount"`
}

type DailyTransactionReport struct {
	TrxDate                     time.Time `json:"trx_date"`
	TotSuccessfulTrxAmount      float64   `json:"tot_successful_trx_amount"`
	TotSuccessfulTrxCount       int64     `json:"tot_successful_trx_count"`
	TotSuccessfulTrxUniqueCount int64     `json:"tot_successful_trx_unique_count"`
	TotPendingTrxAmount         float64   `json:"tot_pending_trx_amount"`
	TotPendingTrxCount          int64     `json:"tot_pending_trx_count"`
	TotPendingTrxUniqueCount    int64     `json:"tot_pending_trx_unique_count"`
	TotFailedTrxAmount          float64   `json:"tot_failed_trx_amount"`
	TotFailedTrxCount           int64     `json:"tot_failed_trx_count"`
	TotFailedTrxUniqueCount     int64     `json:"tot_failed_trx_unique_count"`
}

type MonthlyTransactionReport struct {
	TrxYear                     int     `json:"trx_year"`
	TrxMonth                    string  `json:"trx_month"`
	TotSuccessfulTrxAmount      float64 `json:"tot_successful_trx_amount"`
	TotSuccessfulTrxCount       int64   `json:"tot_successful_trx_count"`
	TotSuccessfulTrxUniqueCount int64   `json:"tot_successful_trx_unique_count"`
	TotPendingTrxAmount         float64 `json:"tot_pending_trx_amount"`
	TotPendingTrxCount          int64   `json:"tot_pending_trx_count"`
	TotPendingTrxUniqueCount    int64   `json:"tot_pending_trx_unique_count"`
	TotFailedTrxAmount          float64 `json:"tot_failed_trx_amount"`
	TotFailedTrxCount           int64   `json:"tot_failed_trx_count"`
	TotFailedTrxUniqueCount     int64   `json:"tot_failed_trx_unique_count"`
}

type HourlyTransactionReport struct {
	TrxDate                     time.Time `json:"trx_date"`
	Hour                        string    `json:"hour"`
	TotSuccessfulTrxAmount      float64   `json:"tot_successful_trx_amount"`
	TotSuccessfulTrxCount       int64     `json:"tot_successful_trx_count"`
	TotSuccessfulTrxUniqueCount int64     `json:"tot_successful_trx_unique_count"`
	TotPendingTrxAmount         float64   `json:"tot_pending_trx_amount"`
	TotPendingTrxCount          int64     `json:"tot_pending_trx_count"`
	TotPendingTrxUniqueCount    int64     `json:"tot_pending_trx_unique_count"`
	TotFailedTrxAmount          float64   `json:"tot_failed_trx_amount"`
	TotFailedTrxCount           int64     `json:"tot_failed_trx_count"`
	TotFailedTrxUniqueCount     int64     `json:"tot_failed_trx_unique_count"`
}
type ProductTransactionReport struct {
	TrxDate                     string  `json:"trx_date,omitempty"`
	TrxYear                     int     `json:"trx_year,omitempty"`
	TrxMonth                    string  `json:"trx_month,omitempty"`
	ProductName                 string  `json:"product_name"`
	ProductPlayedAmount         float64 `json:"product_played_amount"`
	Hour                        string  `json:"hour,omitempty"`
	TotSuccessfulTrxAmount      float64 `json:"tot_successful_trx_amount"`
	TotSuccessfulTrxCount       int64   `json:"tot_successful_trx_count"`
	TotSuccessfulTrxUniqueCount int64   `json:"tot_successful_trx_unique_count"`
	TotPendingTrxAmount         float64 `json:"tot_pending_trx_amount"`
	TotPendingTrxCount          int64   `json:"tot_pending_trx_count"`
	TotPendingTrxUniqueCount    int64   `json:"tot_pending_trx_unique_count"`
	TotFailedTrxAmount          float64 `json:"tot_failed_trx_amount"`
	TotFailedTrxCount           int64   `json:"tot_failed_trx_count"`
	TotFailedTrxUniqueCount     int64   `json:"tot_failed_trx_unique_count"`
}

type PaginationResponse struct {
	Status       string      `json:"status"`
	List         interface{} `json:"list"`
	Total        int         `json:"total"`
	PreviousPage *int        `json:"previous_page"`
	NextPage     *int        `json:"next_page"`
	LastPage     int         `json:"last_page"`
	CurrentPage  int         `json:"current_page"`
}
