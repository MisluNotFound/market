package lib

import (
	"fmt"
	"os"
	"testing"
)

func TestGenerate(t *testing.T) {
	os.Setenv("m_market_config", "../../../config.yaml")
	accessToken, err := GenerateAccessToken("7bed6ab8-378f-11f0-939b-0242ac150003")
	if err != nil {
		t.Fatal(err)
		return
	}

	fmt.Println(accessToken)
}
