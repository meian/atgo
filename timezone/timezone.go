package timezone

import (
	"log/slog"
	"time"
)

var Tokyo *time.Location

func init() {
	var err error
	Tokyo, err = time.LoadLocation("Asia/Tokyo")
	if err != nil {
		slog.With("err", err).Error("Failed to load Asia/Tokyo timezone")
		Tokyo = time.FixedZone("Asia/Tokyo", 9*60*60)
		slog.Info("Using fixed timezone for Asia/Tokyo")
	}
}
