package api

import (
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// rateLimiter 基于令牌桶算法的 per-IP 限流器。
// 每个 IP 独立维护一个桶，按时间自动补充令牌。
// 选择内存实现而非 Redis：单实例部署场景下够用，且零外部依赖。
type rateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*bucket
	rate     int           // 每秒补充的令牌数
	burst    int           // 桶的最大容量（允许瞬时峰值）
	cleanup  time.Duration // 过期 visitor 清理间隔
}

type bucket struct {
	tokens    int
	lastCheck time.Time
}

func newRateLimiter(ratePerSec, burst int) *rateLimiter {
	rl := &rateLimiter{
		visitors: make(map[string]*bucket),
		rate:     ratePerSec,
		burst:    burst,
		cleanup:  5 * time.Minute,
	}
	go rl.cleanupLoop()
	return rl
}

// allow 判断该 IP 是否还有可用令牌。首次访问直接放行并初始化桶。
func (rl *rateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	b, exists := rl.visitors[ip]
	if !exists {
		rl.visitors[ip] = &bucket{tokens: rl.burst - 1, lastCheck: time.Now()}
		return true
	}

	// 按经过时间补充令牌，上限为 burst
	elapsed := time.Since(b.lastCheck)
	b.lastCheck = time.Now()
	b.tokens += int(elapsed.Seconds()) * rl.rate
	if b.tokens > rl.burst {
		b.tokens = rl.burst
	}

	if b.tokens <= 0 {
		return false
	}
	b.tokens--
	return true
}

// cleanupLoop 定期清理长时间未访问的 visitor，防止 map 无限增长。
func (rl *rateLimiter) cleanupLoop() {
	ticker := time.NewTicker(rl.cleanup)
	defer ticker.Stop()
	for range ticker.C {
		rl.mu.Lock()
		for ip, b := range rl.visitors {
			if time.Since(b.lastCheck) > rl.cleanup {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// rateLimitMiddleware 将限流器包装为 chi 中间件，被限流时返回 429 + Retry-After。
func rateLimitMiddleware(rl *rateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := extractIP(r)
			if !rl.allow(ip) {
				w.Header().Set("Retry-After", "1")
				respondErr(w, http.StatusTooManyRequests, "rate limit exceeded")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// extractIP 从请求中提取客户端真实 IP。
// 优先级：X-Forwarded-For 第一段 > X-Real-IP > RemoteAddr（适配反向代理场景）。
func extractIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.SplitN(xff, ",", 2)
		return strings.TrimSpace(parts[0])
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	parts := strings.SplitN(r.RemoteAddr, ":", 2)
	return parts[0]
}

// corsMiddleware 处理跨域请求。
// 生产环境应通过 CORS_ALLOWED_ORIGINS 环境变量指定具体域名，避免使用通配符 *。
func corsMiddleware(next http.Handler) http.Handler {
	allowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "*"
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if allowedOrigins == "*" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		} else {
			// 逐一匹配逗号分隔的白名单
			for _, allowed := range strings.Split(allowedOrigins, ",") {
				if strings.TrimSpace(allowed) == origin {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					break
				}
			}
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID")
		w.Header().Set("Access-Control-Max-Age", "86400") // 预检缓存 24h
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// securityHeaders 注入通用安全响应头，防范 MIME 嗅探、点击劫持、XSS 等常见攻击。
func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		next.ServeHTTP(w, r)
	})
}