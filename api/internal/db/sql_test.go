package db

import (
	"os"
	"testing"

	"github.com/mislu/market-api/internal/utils/log"
)

func TestFirstOrCreate(t *testing.T) {
	os.Setenv("m_market_config", "../../config.yaml")
	logger := log.NewLogger()
	Init(logger)

	type MockTable struct {
		ID int `gorm:"column:id"`
	}

	DB.AutoMigrate(&MockTable{})

	var mocks = []MockTable{
		{
			ID: 1,
		},
		{
			ID: 1,
		},
	}

	for _, mock := range mocks {
		err := FirstOrCreate(&mock)
		if err != nil {
			t.Fatal(err)
		}
	}

	count, err := GetCount[MockTable](
		Equal("id", 1),
	)
	if err != nil {
		t.Fatal(err)
	}

	if count != 1 {
		t.Fail()
	}
}
