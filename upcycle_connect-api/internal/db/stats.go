package db

import (
	"upcycle_connect-api/internal/config"
	"upcycle_connect-api/internal/models"
)

func GetStats() (models.Stats, error) {
    var stats models.Stats

    err := config.Conn.QueryRow("SELECT COUNT(id_user) FROM USER").Scan(&stats.UserCount)
    if err != nil {
        return stats, err
    }

    err = config.Conn.QueryRow("SELECT COUNT(id_event) FROM EVENT").Scan(&stats.EventCount)
    if err != nil {
        return stats, err
    }

	err = config.Conn.QueryRow("SELECT COUNT(id_announcement) FROM ANNOUNCEMENT").Scan(&stats.AnnouncementCount)
	if err != nil {
		return stats, err
	}

	err = config.Conn.QueryRow("SELECT COUNT(id_report) FROM REPORT WHERE status = 'pending'").Scan(&stats.ReportCount)
	if err != nil {
		return stats, err
	}

	return stats, nil
}
