// SELECT *
// FROM pg_catalog.pg_tables
// WHERE schemaname != 'pg_catalog' AND
//    schemaname != 'information_schema';

// select * from anomalies;

package report

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Anomaly struct {
	SessionId    string
	Frequency    float64
	UtcTimestamp int64
}

func Report(ID string, Frequency float64, Timestamp int64) error {
	dsn := "host=localhost user=pitermar password=1243 dbname=anomaly port=5051 sslmode=disable TimeZone=Europe/Moscow"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	// Миграция схемы
	db.AutoMigrate(&Anomaly{})

	// Создание новой записи
	anomaly := Anomaly{SessionId: ID, Frequency: Frequency, UtcTimestamp: Timestamp}
	result := db.Create(&anomaly)
	if result.Error != nil {
		return err
	}

	return nil
}

func ReadReport() error {
	dsn := "host=localhost user=pitermar password=1243 dbname=anomaly port=5051 sslmode=disable TimeZone=Europe/Moscow"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	var anomalies []Anomaly
	result := db.Find(&anomalies)
	if result.Error != nil {
		return err
	}
	for _, anomaly := range anomalies {
		log.Println(anomaly)
	}

	return nil
}
