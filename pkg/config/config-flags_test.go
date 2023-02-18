package config

// TODO - REWRITE THIS

// import (
// 	"flag"
// 	"reflect"
// 	"testing"
// )

// var testNumber = flag.Uint("test-number", 1, "")

// func TestFlags(t *testing.T) {
// 	switch *testNumber {
// 	case 1:
// 		cfg, err := GetConfigOrDie()
// 		if err != nil {
// 			t.Errorf("Error get config: %s", err.Error())
// 		}
// 		if !reflect.DeepEqual(cfg, Config{
// 			VaultToken: "my-token",
// 			VaultUrl:   "https://test-vault.url",
// 		}) {
// 			t.Errorf("Failed creating configuration from flags")
// 		}
// 	default:
// 		t.Errorf("Invalid test number")
// 	}
// }
