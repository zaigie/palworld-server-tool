package auth

import (
	"sync"
	"time"
)

type LoginAttempt struct {
	Count      int       // 失败次数
	LastFailed time.Time // 最近一次失败时间
}

type LoginGuard struct {
	mu       sync.Mutex
	attempts map[string]*LoginAttempt
	Window   time.Duration
	Limit    int
	BlockDur time.Duration
}

var DefaultLoginGuard = &LoginGuard{
	attempts: make(map[string]*LoginAttempt),
	Window:   time.Minute,     // 1分钟之内
	Limit:    5,               // 最多失败5次
	BlockDur: 5 * time.Minute, // 封锁5分钟
}

// 检查是否被封禁
func (g *LoginGuard) IsBlocked(ip string) (bool, time.Duration) {
	g.mu.Lock()
	defer g.mu.Unlock()

	att, exists := g.attempts[ip]
	if !exists || att.Count < g.Limit {
		return false, 0
	}

	now := time.Now()
	if now.Sub(att.LastFailed) < g.BlockDur {
		return true, g.BlockDur - now.Sub(att.LastFailed)
	}
	return false, 0
}

// 登录失败记录
func (g *LoginGuard) RecordFailure(ip string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	now := time.Now()
	att, exists := g.attempts[ip]
	if !exists {
		att = &LoginAttempt{}
		g.attempts[ip] = att
	}

	if now.Sub(att.LastFailed) > g.Window {
		att.Count = 1
	} else {
		att.Count++
	}
	att.LastFailed = now
}

// 登录成功重置
func (g *LoginGuard) Reset(ip string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	att, exists := g.attempts[ip]
	if !exists {
		return
	}
	att.Count = 0
	att.LastFailed = time.Time{}
}

// 清理超过指定时间未活动的 IP
func (g *LoginGuard) Cleanup(expire time.Duration) {
	g.mu.Lock()
	defer g.mu.Unlock()

	now := time.Now()
	for ip, att := range g.attempts {
		if att.LastFailed.IsZero() {
			continue
		}
		if now.Sub(att.LastFailed) > expire {
			delete(g.attempts, ip)
		}
	}
}
