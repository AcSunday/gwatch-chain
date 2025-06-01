package loadbalance

// services list
// var services = []string{
// 	"https://eth.llamarpc.com",
// 	"https://rpc.ankr.com/eth",
// 	"https://eth.api.onfinality.io/public",
// 	"https://eth-pokt.nodies.app",
// 	"https://xlayertestrpc.okx.com",
// 	"http://failed.example.com",
// }

// func Test_SimpleCheckHealthy(t *testing.T) {
// 	lb := New(services, evmClientFactory)

// 	for i := 0; i < 50; i++ {
// 		cli := lb.NextClient()
// 		if cli != nil {
// 			t.Logf("check healthy, cli url: %s", cli.GetRawUrl())
// 		}

// 		time.Sleep(1 * time.Second)
// 	}
// }
