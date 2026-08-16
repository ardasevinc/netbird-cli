package ledger

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ardasevinc/netbird-cli/internal/mutation"
	_ "modernc.org/sqlite"
)

type Store struct {
	db *sql.DB
}

type StageInput struct {
	Profile        string
	ServerIdentity string
	AccountID      string
	Operation      string
	Request        json.RawMessage
	Before         json.RawMessage
	IntendedAfter  json.RawMessage
	Impact         json.RawMessage
	Findings       []Finding
}

type Finding struct {
	Code     string `json:"code"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
}

type Stage struct {
	ID             string
	Revision       int
	Profile        string
	ServerIdentity string
	AccountID      string
	Operation      string
	Request        json.RawMessage
	Before         json.RawMessage
	IntendedAfter  json.RawMessage
	Impact         json.RawMessage
	Digest         string
	Findings       []Finding
	Cancelled      bool
	CreatedAt      time.Time
}

type Attempt struct {
	ID        string
	StageID   string
	Revision  int
	State     string
	CreatedAt time.Time
	IntentAt  time.Time
}

type Receipt struct {
	ID        string
	AttemptID string
	StageID   string
	Revision  int
	State     string
	Result    json.RawMessage
	CreatedAt time.Time
}

func Open(path string) (*Store, error) {
	if path == "" {
		return nil, errors.New("ledger path is empty")
	}
	if dir := filepath.Dir(path); dir != "." {
		if err := os.MkdirAll(dir, 0o700); err != nil {
			return nil, fmt.Errorf("create ledger directory: %w", err)
		}
		if err := os.Chmod(dir, 0o700); err != nil { // #nosec G302 -- ledger directories must be owner-only and need execute permission.
			return nil, fmt.Errorf("protect ledger directory: %w", err)
		}
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open ledger: %w", err)
	}
	db.SetMaxOpenConns(1)
	store := &Store{db: db}
	if err := store.migrate(context.Background()); err != nil {
		_ = db.Close()
		return nil, err
	}
	if err := os.Chmod(path, 0o600); err != nil { // #nosec G302 -- the ledger is machine-local owner-only state.
		_ = db.Close()
		return nil, fmt.Errorf("protect ledger file: %w", err)
	}
	return store, nil
}

func (s *Store) Close() error { return s.db.Close() }

func (s *Store) migrate(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS schema_version (version INTEGER NOT NULL);
INSERT INTO schema_version(version)
SELECT 3 WHERE NOT EXISTS (SELECT 1 FROM schema_version);
CREATE TABLE IF NOT EXISTS stages (
  id TEXT NOT NULL,
  revision INTEGER NOT NULL,
  profile TEXT NOT NULL,
  server_identity TEXT NOT NULL,
  account_id TEXT NOT NULL,
  operation TEXT NOT NULL,
  request_json BLOB NOT NULL,
  before_json BLOB NOT NULL,
  intended_after_json BLOB NOT NULL,
  impact_json BLOB NOT NULL DEFAULT '{}',
  digest TEXT NOT NULL,
  findings_json BLOB NOT NULL,
  cancelled INTEGER NOT NULL DEFAULT 0,
  created_at TEXT NOT NULL,
  PRIMARY KEY (id, revision)
);
CREATE TABLE IF NOT EXISTS attempts (
  id TEXT PRIMARY KEY,
  stage_id TEXT NOT NULL,
  revision INTEGER NOT NULL,
  state TEXT NOT NULL,
  created_at TEXT NOT NULL,
  intent_at TEXT NOT NULL,
  FOREIGN KEY (stage_id, revision) REFERENCES stages(id, revision)

);
CREATE TABLE IF NOT EXISTS receipts (
  id TEXT PRIMARY KEY,
  attempt_id TEXT NOT NULL UNIQUE,
  stage_id TEXT NOT NULL,
  revision INTEGER NOT NULL,
  state TEXT NOT NULL,
  result_json BLOB NOT NULL,
  created_at TEXT NOT NULL,
  FOREIGN KEY (attempt_id) REFERENCES attempts(id),
  FOREIGN KEY (stage_id, revision) REFERENCES stages(id, revision)
);
	`)
	if err != nil {
		return fmt.Errorf("migrate ledger: %w", err)
	}
	var version int
	if err := s.db.QueryRowContext(ctx, `SELECT version FROM schema_version LIMIT 1`).Scan(&version); err != nil {
		return fmt.Errorf("read ledger schema version: %w", err)
	}
	if version < 3 {
		if _, err := s.db.ExecContext(ctx, `ALTER TABLE stages ADD COLUMN impact_json BLOB NOT NULL DEFAULT '{}'`); err != nil {
			return fmt.Errorf("migrate stage impact: %w", err)
		}
		if _, err := s.db.ExecContext(ctx, `UPDATE schema_version SET version = 3 WHERE version < 3`); err != nil {
			return fmt.Errorf("update ledger schema version: %w", err)
		}
	}
	return nil
}

// BeginAttempt journals the exact revision before a dispatcher may send any
// request bytes to a remote service.
func (s *Store) BeginAttempt(ctx context.Context, stageID string, revision int) (Attempt, error) {
	if _, err := s.Get(ctx, stageID, revision); err != nil {
		return Attempt{}, err
	}
	id, err := newID()
	if err != nil {
		return Attempt{}, err
	}
	now := time.Now().UTC().Truncate(time.Microsecond)
	_, err = s.db.ExecContext(ctx, `INSERT INTO attempts (id, stage_id, revision, state, created_at, intent_at) VALUES (?, ?, ?, 'not_dispatched', ?, ?)`, id, stageID, revision, now.Format(time.RFC3339Nano), now.Format(time.RFC3339Nano))
	if err != nil {
		return Attempt{}, fmt.Errorf("journal dispatch intent: %w", err)
	}
	return Attempt{ID: id, StageID: stageID, Revision: revision, State: "not_dispatched", CreatedAt: now, IntentAt: now}, nil
}

func (s *Store) GetAttempt(ctx context.Context, id string) (Attempt, error) {
	var attempt Attempt
	var created, intent string
	err := s.db.QueryRowContext(ctx, `SELECT id, stage_id, revision, state, created_at, intent_at FROM attempts WHERE id = ?`, id).Scan(&attempt.ID, &attempt.StageID, &attempt.Revision, &attempt.State, &created, &intent)
	if errors.Is(err, sql.ErrNoRows) {
		return Attempt{}, fmt.Errorf("attempt %s not found", id)
	}
	if err != nil {
		return Attempt{}, fmt.Errorf("read attempt: %w", err)
	}
	attempt.CreatedAt, err = time.Parse(time.RFC3339Nano, created)
	if err != nil {
		return Attempt{}, fmt.Errorf("decode attempt timestamp: %w", err)
	}
	attempt.IntentAt, err = time.Parse(time.RFC3339Nano, intent)
	if err != nil {
		return Attempt{}, fmt.Errorf("decode intent timestamp: %w", err)
	}
	return attempt, nil
}

func (s *Store) SetAttemptState(ctx context.Context, id, state string) error {
	attempt, err := s.GetAttempt(ctx, id)
	if err != nil {
		return err
	}
	if err := mutation.ValidateDispatchTransition(mutation.DispatchState(attempt.State), mutation.DispatchState(state)); err != nil {
		return fmt.Errorf("attempt %s state transition: %w", id, err)
	}
	result, err := s.db.ExecContext(ctx, `UPDATE attempts SET state = ? WHERE id = ? AND state = ?`, state, id, attempt.State)
	if err != nil {
		return fmt.Errorf("update attempt state: %w", err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check attempt state update: %w", err)
	}
	if count != 1 {
		return fmt.Errorf("attempt %s state changed concurrently", id)
	}
	return nil
}

func (s *Store) RecordReceipt(ctx context.Context, receipt Receipt) error {
	if receipt.AttemptID == "" || receipt.StageID == "" || receipt.State == "" {
		return errors.New("receipt identity and state are required")
	}
	result, err := mutation.CanonicalJSON(receipt.Result)
	if err != nil {
		return fmt.Errorf("canonical receipt: %w", err)
	}
	if receipt.ID == "" {
		receipt.ID, err = newID()
		if err != nil {
			return err
		}
	}
	created := receipt.CreatedAt.UTC().Truncate(time.Microsecond)
	if created.IsZero() {
		created = time.Now().UTC().Truncate(time.Microsecond)
	}
	_, err = s.db.ExecContext(ctx, `INSERT INTO receipts (id, attempt_id, stage_id, revision, state, result_json, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)`, receipt.ID, receipt.AttemptID, receipt.StageID, receipt.Revision, receipt.State, result, created.Format(time.RFC3339Nano))
	if err != nil {
		return fmt.Errorf("persist receipt: %w", err)
	}
	return nil
}

func (s *Store) GetReceipt(ctx context.Context, attemptID string) (Receipt, error) {
	var receipt Receipt
	var created string
	err := s.db.QueryRowContext(ctx, `SELECT id, attempt_id, stage_id, revision, state, result_json, created_at FROM receipts WHERE attempt_id = ?`, attemptID).Scan(&receipt.ID, &receipt.AttemptID, &receipt.StageID, &receipt.Revision, &receipt.State, &receipt.Result, &created)
	if errors.Is(err, sql.ErrNoRows) {
		return Receipt{}, fmt.Errorf("receipt for attempt %s not found", attemptID)
	}
	if err != nil {
		return Receipt{}, fmt.Errorf("read receipt: %w", err)
	}
	receipt.CreatedAt, err = time.Parse(time.RFC3339Nano, created)
	if err != nil {
		return Receipt{}, fmt.Errorf("decode receipt timestamp: %w", err)
	}
	return receipt, nil
}

func (s *Store) Create(ctx context.Context, input StageInput) (Stage, error) {
	if input.Profile == "" || input.Operation == "" {
		return Stage{}, errors.New("profile and operation are required")
	}
	request, err := mutation.CanonicalJSON(input.Request)
	if err != nil {
		return Stage{}, fmt.Errorf("canonical request: %w", err)
	}
	before, err := mutation.CanonicalJSON(input.Before)
	if err != nil {
		return Stage{}, fmt.Errorf("canonical preimage: %w", err)
	}
	after, err := mutation.CanonicalJSON(input.IntendedAfter)
	if err != nil {
		return Stage{}, fmt.Errorf("canonical intended state: %w", err)
	}
	impact := input.Impact
	if len(impact) == 0 {
		impact = json.RawMessage(`{}`)
	}
	impact, err = mutation.CanonicalJSON(impact)
	if err != nil {
		return Stage{}, fmt.Errorf("canonical impact evidence: %w", err)
	}
	findings, err := json.Marshal(input.Findings)
	if err != nil {
		return Stage{}, fmt.Errorf("encode findings: %w", err)
	}
	digest, err := mutation.Digest(request, before, after, impact, findings)
	if err != nil {
		return Stage{}, fmt.Errorf("digest stage: %w", err)
	}
	id, err := newID()
	if err != nil {
		return Stage{}, err
	}
	created := time.Now().UTC().Truncate(time.Microsecond)
	_, err = s.db.ExecContext(ctx, `INSERT INTO stages (id, revision, profile, server_identity, account_id, operation, request_json, before_json, intended_after_json, impact_json, digest, findings_json, created_at) VALUES (?, 1, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, id, input.Profile, input.ServerIdentity, input.AccountID, input.Operation, request, before, after, impact, digest, findings, created.Format(time.RFC3339Nano))
	if err != nil {
		return Stage{}, fmt.Errorf("persist stage: %w", err)
	}
	return Stage{ID: id, Revision: 1, Profile: input.Profile, ServerIdentity: input.ServerIdentity, AccountID: input.AccountID, Operation: input.Operation, Request: request, Before: before, IntendedAfter: after, Impact: impact, Digest: digest, Findings: input.Findings, CreatedAt: created}, nil
}

func (s *Store) Get(ctx context.Context, id string, revision int) (Stage, error) {
	var stage Stage
	var findings []byte
	var created string
	var cancelled int
	err := s.db.QueryRowContext(ctx, `SELECT id, revision, profile, server_identity, account_id, operation, request_json, before_json, intended_after_json, impact_json, digest, findings_json, cancelled, created_at FROM stages WHERE id = ? AND revision = ?`, id, revision).Scan(&stage.ID, &stage.Revision, &stage.Profile, &stage.ServerIdentity, &stage.AccountID, &stage.Operation, &stage.Request, &stage.Before, &stage.IntendedAfter, &stage.Impact, &stage.Digest, &findings, &cancelled, &created)
	if errors.Is(err, sql.ErrNoRows) {
		return Stage{}, fmt.Errorf("stage %s@%d not found", id, revision)
	}
	if err != nil {
		return Stage{}, fmt.Errorf("read stage: %w", err)
	}
	if err := json.Unmarshal(findings, &stage.Findings); err != nil {
		return Stage{}, fmt.Errorf("decode stage findings: %w", err)
	}
	stage.Cancelled = cancelled != 0
	stage.CreatedAt, err = time.Parse(time.RFC3339Nano, created)
	if err != nil {
		return Stage{}, fmt.Errorf("decode stage timestamp: %w", err)
	}
	return stage, nil
}

func (s *Store) Cancel(ctx context.Context, id string, revision int) error {
	result, err := s.db.ExecContext(ctx, `UPDATE stages SET cancelled = 1 WHERE id = ? AND revision = ? AND cancelled = 0`, id, revision)
	if err != nil {
		return fmt.Errorf("cancel stage: %w", err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check stage cancellation: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("stage %s@%d not found or already cancelled", id, revision)
	}
	return nil
}

func newID() (string, error) {
	var raw [16]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "", fmt.Errorf("generate stage id: %w", err)
	}
	return hex.EncodeToString(raw[:]), nil
}
