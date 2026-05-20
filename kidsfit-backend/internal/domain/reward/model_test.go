package reward

import (
	"testing"
	"time"
)

// TestChallenge_IsExpired 测试挑战是否过期
func TestChallenge_IsExpired(t *testing.T) {
	t.Run("无过期时间不过期", func(t *testing.T) {
		c := &Challenge{ExpiresAt: nil}
		if c.IsExpired() {
			t.Error("无过期时间的挑战不应判断为过期")
		}
	})

	t.Run("过期时间在未来不过期", func(t *testing.T) {
		future := time.Now().Add(24 * time.Hour)
		c := &Challenge{ExpiresAt: &future}
		if c.IsExpired() {
			t.Error("过期时间在未来的挑战不应判断为过期")
		}
	})

	t.Run("过期时间在过去已过期", func(t *testing.T) {
		past := time.Now().Add(-1 * time.Hour)
		c := &Challenge{ExpiresAt: &past}
		if !c.IsExpired() {
			t.Error("过期时间在过去的挑战应判断为已过期")
		}
	})
}

// TestChallenge_CanAccept 测试挑战是否可接受（pending状态且未过期）
func TestChallenge_CanAccept(t *testing.T) {
	t.Run("pending状态且未过期可接受", func(t *testing.T) {
		future := time.Now().Add(24 * time.Hour)
		c := &Challenge{Status: ChallengeStatusPending, ExpiresAt: &future}
		if !c.CanAccept() {
			t.Error("pending状态且未过期的挑战应可接受")
		}
	})

	t.Run("pending状态但已过期不可接受", func(t *testing.T) {
		past := time.Now().Add(-1 * time.Hour)
		c := &Challenge{Status: ChallengeStatusPending, ExpiresAt: &past}
		if c.CanAccept() {
			t.Error("已过期的挑战不应可接受")
		}
	})

	t.Run("accepted状态不可接受", func(t *testing.T) {
		future := time.Now().Add(24 * time.Hour)
		c := &Challenge{Status: ChallengeStatusAccepted, ExpiresAt: &future}
		if c.CanAccept() {
			t.Error("已接受状态的挑战不应可接受")
		}
	})

	t.Run("completed状态不可接受", func(t *testing.T) {
		future := time.Now().Add(24 * time.Hour)
		c := &Challenge{Status: ChallengeStatusCompleted, ExpiresAt: &future}
		if c.CanAccept() {
			t.Error("已完成状态的挑战不应可接受")
		}
	})
}

// TestChallenge_CanSubmit 测试挑战是否可提交成绩（accepted状态且未过期）
func TestChallenge_CanSubmit(t *testing.T) {
	t.Run("accepted状态且未过期可提交", func(t *testing.T) {
		future := time.Now().Add(24 * time.Hour)
		c := &Challenge{Status: ChallengeStatusAccepted, ExpiresAt: &future}
		if !c.CanSubmit() {
			t.Error("accepted状态且未过期的挑战应可提交")
		}
	})

	t.Run("accepted状态但已过期不可提交", func(t *testing.T) {
		past := time.Now().Add(-1 * time.Hour)
		c := &Challenge{Status: ChallengeStatusAccepted, ExpiresAt: &past}
		if c.CanSubmit() {
			t.Error("已过期的挑战不应可提交")
		}
	})

	t.Run("pending状态不可提交", func(t *testing.T) {
		future := time.Now().Add(24 * time.Hour)
		c := &Challenge{Status: ChallengeStatusPending, ExpiresAt: &future}
		if c.CanSubmit() {
			t.Error("pending状态的挑战不应可提交")
		}
	})

	t.Run("completed状态不可提交", func(t *testing.T) {
		future := time.Now().Add(24 * time.Hour)
		c := &Challenge{Status: ChallengeStatusCompleted, ExpiresAt: &future}
		if c.CanSubmit() {
			t.Error("已完成状态的挑战不应可提交")
		}
	})
}
