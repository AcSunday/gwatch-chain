package loadbalance

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AcSunday/gwatch-chain/rpcclient"
)

// evmClientFactory 创建EvmClient的工厂函数
func evmClientFactory(url string) (*rpcclient.EvmClient, error) {
	return rpcclient.NewEvmRpcClient(url)
}

// 创建模拟的 HTTP 服务器
func createMockServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"jsonrpc":"2.0","id":1,"result":"0x1"}`))
	}))
}

func setupTestServers() ([]string, []*httptest.Server) {
	var urls []string
	var servers []*httptest.Server

	// 创建5个模拟服务器
	for i := 0; i < 5; i++ {
		server := createMockServer()
		urls = append(urls, server.URL)
		servers = append(servers, server)
	}

	return urls, servers
}

func cleanupTestServers(servers []*httptest.Server) {
	for _, server := range servers {
		server.Close()
	}
}

func TestLB_CheckHealthy(t *testing.T) {
	urls, servers := setupTestServers()
	defer cleanupTestServers(servers)

	lb := New(urls, evmClientFactory)

	for i := 0; i < 5; i++ {
		cli := lb.NextClient()
		if cli != nil {
			t.Logf("check healthy, cli url: %s", cli.GetRawUrl())
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// 基准测试：测试获取下一个客户端的性能
func BenchmarkLoadBalance_NextClient(b *testing.B) {
	urls, servers := setupTestServers()
	defer cleanupTestServers(servers)

	lb := New(urls, evmClientFactory)
	// 等待健康检查完成
	time.Sleep(100 * time.Millisecond)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cli := lb.NextClient()
		if cli == nil {
			b.Fatal("no client available")
		}
	}
}

// 基准测试：测试并发获取客户端的性能
func BenchmarkLoadBalance_NextClient_Parallel(b *testing.B) {
	urls, servers := setupTestServers()
	defer cleanupTestServers(servers)

	lb := New(urls, evmClientFactory)
	// 等待健康检查完成
	time.Sleep(100 * time.Millisecond)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			cli := lb.NextClient()
			if cli == nil {
				b.Fatal("no client available")
			}
		}
	})
}

// 基准测试：测试健康检查的性能
func BenchmarkLoadBalance_HealthCheck(b *testing.B) {
	urls, servers := setupTestServers()
	defer cleanupTestServers(servers)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lb := New(urls, evmClientFactory)
		// 等待一个健康检查周期的一小部分时间
		time.Sleep(100 * time.Millisecond)
		lb.Close()
	}
}

// 基准测试：测试在高负载下的性能
func BenchmarkLoadBalance_HighLoad(b *testing.B) {
	urls, servers := setupTestServers()
	defer cleanupTestServers(servers)

	lb := New(urls, evmClientFactory)
	// 等待健康检查完成
	time.Sleep(100 * time.Millisecond)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for j := 0; j < 100; j++ { // 模拟每个goroutine进行100次请求
				cli := lb.NextClient()
				if cli == nil {
					b.Fatal("no client available")
				}
			}
		}
	})
}
