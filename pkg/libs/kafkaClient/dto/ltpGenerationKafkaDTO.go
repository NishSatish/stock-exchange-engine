package dto

import "time"

const (
	LtpGenerationTopic = "topics.trades.ltpGeneration"
)

// TODO: Only adding support to stream the LTP for now, not all the data for a full candle
type LtpGenerationKafkaDTO struct {
	StockID   string    `json:"stock_id"`
	Price     float64   `json:"price"`
	Timestamp time.Time `json:"timestamp"`
}
