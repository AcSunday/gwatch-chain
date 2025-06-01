package loadbalance

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// MockRPCClient 是一个模拟的RPC客户端实现
type MockRPCClient struct {
	url     string
	chainId uint64
}

func (m MockRPCClient) Close() {}

func (m MockRPCClient) GetRawUrl() string {
	return m.url
}

func (m MockRPCClient) GetChainId() uint64 {
	return m.chainId
}

// mockClientFactory 创建MockRPCClient的工厂函数
func mockClientFactory(url string) (MockRPCClient, error) {
	return MockRPCClient{
		url:     url,
		chainId: 1, // 模拟chainId
	}, nil
}

// TestGenericLoadBalance_Basic 测试基本功能
func TestGenericLoadBalance_Basic(t *testing.T) {
	// 创建测试服务器
	servers := make([]*httptest.Server, 3)
	urls := make([]string, 3)

	for i := range servers {
		servers[i] = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		urls[i] = servers[i].URL
	}

	// 清理测试服务器
	defer func() {
		for _, server := range servers {
			server.Close()
		}
	}()

	// 创建负载均衡器
	lb := New[MockRPCClient](urls, mockClientFactory)
	if lb == nil {
		t.Fatal("failed to create load balancer")
	}
	defer lb.Close()

	// 测试客户端轮询
	for i := 0; i < 6; i++ {
		client := lb.NextClient()
		if client == (MockRPCClient{}) {
			t.Fatal("got zero client")
		}
		t.Logf("got client with url: %s", client.GetRawUrl())
	}
}

// BenchmarkGenericLoadBalance_NextClient 基准测试
func BenchmarkGenericLoadBalance_NextClient(b *testing.B) {
	// 创建测试服务器
	servers := make([]*httptest.Server, 3)
	urls := make([]string, 3)

	for i := range servers {
		servers[i] = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		urls[i] = servers[i].URL
	}

	// 清理测试服务器
	defer func() {
		for _, server := range servers {
			server.Close()
		}
	}()

	// 创建负载均衡器
	lb := New[MockRPCClient](urls, mockClientFactory)
	if lb == nil {
		b.Fatal("failed to create load balancer")
	}
	defer lb.Close()

	// 等待健康检查完成
	time.Sleep(100 * time.Millisecond)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			client := lb.NextClient()
			if client == (MockRPCClient{}) {
				b.Fatal("got zero client")
			}
		}
	})
}

// Example 展示如何使用泛型负载均衡器
func Example() {
	// 定义URLs
	urls := []string{
		"http://localhost:8545",
		"http://localhost:8546",
		"http://localhost:8547",
	}

	// 创建负载均衡器
	lb := New[MockRPCClient](urls, mockClientFactory)
	if lb == nil {
		return
	}
	defer lb.Close()

	// 获取客户端并使用
	client := lb.NextClient()
	_ = client.GetRawUrl() // 实际使用中调用客户端方法
}
