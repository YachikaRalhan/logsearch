package model

import "time"

type LogSearchRequest struct {
	SearchKeyword string
	From          time.Time
	To            time.Time
}
