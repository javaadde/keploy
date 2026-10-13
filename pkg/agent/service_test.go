package agent

import (
	"sync"
	"testing"

	"go.keploy.io/server/v3/pkg/models"
)

func TestNewSessions(t *testing.T) {
	sessions := NewSessions()
	if sessions == nil {
		t.Fatal("expected NewSessions() to return non-nil *Sessions")
	}

	// Verify Get on empty Sessions
	sess, ok := sessions.Get(1)
	if ok || sess != nil {
		t.Errorf("expected Get(1) on empty Sessions to return (nil, false), got (%v, %v)", sess, ok)
	}
}

func TestSessions_SetAndGet(t *testing.T) {
	sessions := NewSessions()

	session1 := &Session{
		ID:   1,
		Mode: models.MODE_RECORD,
	}

	session2 := &Session{
		ID:   2,
		Mode: models.MODE_TEST,
	}

	sessions.Set(session1.ID, session1)
	sessions.Set(session2.ID, session2)

	// Get session 1
	got1, ok := sessions.Get(1)
	if !ok || got1 == nil {
		t.Fatalf("expected Get(1) to return (session1, true), got (%v, %v)", got1, ok)
	}
	if got1.ID != session1.ID || got1.Mode != session1.Mode {
		t.Errorf("expected session ID %d and Mode %v, got ID %d and Mode %v", session1.ID, session1.Mode, got1.ID, got1.Mode)
	}

	// Get session 2
	got2, ok := sessions.Get(2)
	if !ok || got2 == nil {
		t.Fatalf("expected Get(2) to return (session2, true), got (%v, %v)", got2, ok)
	}
	if got2.ID != session2.ID || got2.Mode != session2.Mode {
		t.Errorf("expected session ID %d and Mode %v, got ID %d and Mode %v", session2.ID, session2.Mode, got2.ID, got2.Mode)
	}

	// Get non-existent session 3
	got3, ok := sessions.Get(3)
	if ok || got3 != nil {
		t.Errorf("expected Get(3) to return (nil, false), got (%v, %v)", got3, ok)
	}
}

func TestSessions_Delete(t *testing.T) {
	sessions := NewSessions()

	session := &Session{
		ID:   100,
		Mode: models.MODE_RECORD,
	}

	sessions.Set(session.ID, session)

	got, ok := sessions.Get(100)
	if !ok || got == nil {
		t.Fatalf("expected Get(100) to return (session, true), got (%v, %v)", got, ok)
	}

	sessions.Delete(100)

	gotAfterDelete, okAfterDelete := sessions.Get(100)
	if okAfterDelete || gotAfterDelete != nil {
		t.Errorf("expected Get(100) after Delete to return (nil, false), got (%v, %v)", gotAfterDelete, okAfterDelete)
	}
}

func TestSessions_Concurrency(t *testing.T) {
	sessions := NewSessions()
	const numGoroutines = 100
	const numOpsPerGoroutine = 50

	var wg sync.WaitGroup
	wg.Add(numGoroutines * 3)

	// Goroutines performing Set
	for i := 0; i < numGoroutines; i++ {
		go func(gID int) {
			defer wg.Done()
			for op := 0; op < numOpsPerGoroutine; op++ {
				id := uint64(gID*numOpsPerGoroutine + op)
				sessions.Set(id, &Session{ID: id, Mode: models.MODE_RECORD})
			}
		}(i)
	}

	// Goroutines performing Get
	for i := 0; i < numGoroutines; i++ {
		go func(gID int) {
			defer wg.Done()
			for op := 0; op < numOpsPerGoroutine; op++ {
				id := uint64(gID*numOpsPerGoroutine + op)
				sessions.Get(id)
			}
		}(i)
	}

	// Goroutines performing Delete
	for i := 0; i < numGoroutines; i++ {
		go func(gID int) {
			defer wg.Done()
			for op := 0; op < numOpsPerGoroutine; op++ {
				id := uint64(gID*numOpsPerGoroutine + op)
				sessions.Delete(id)
			}
		}(i)
	}

	wg.Wait()
}
