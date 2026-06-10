# Lottery Search System — Design Proposal

## Overview

Design a system to search 1 million 6-digit lottery tickets using wildcard patterns, while ensuring no two users receive the same ticket simultaneously.

---

## 1. Data Storage

### Recommended Database: PostgreSQL

PostgreSQL is chosen for the following reasons:

- **ACID transactions** — critical for atomic ticket reservation to prevent duplicate allocation
- **Row-level locking** — `SELECT ... FOR UPDATE SKIP LOCKED` allows concurrent reservation without conflicts
- **Partial indexes** — can index specific digit positions for fast wildcard lookups
- **Mature and production-proven** — reliable for high-concurrency workloads

### Schema

```sql
CREATE TABLE lottery_tickets (
    id        SERIAL PRIMARY KEY,
    number    CHAR(6)     NOT NULL UNIQUE,  -- e.g. '123456'
    d1        CHAR(1)     NOT NULL,          -- digit 1
    d2        CHAR(1)     NOT NULL,          -- digit 2
    d3        CHAR(1)     NOT NULL,          -- digit 3
    d4        CHAR(1)     NOT NULL,          -- digit 4
    d5        CHAR(1)     NOT NULL,          -- digit 5
    d6        CHAR(1)     NOT NULL,          -- digit 6
    status    VARCHAR(20) NOT NULL DEFAULT 'available',  -- available | reserved | sold
    reserved_at  TIMESTAMPTZ,
    reserved_by  TEXT
);

-- Index each digit column for fast wildcard matching
CREATE INDEX idx_d1 ON lottery_tickets (d1);
CREATE INDEX idx_d2 ON lottery_tickets (d2);
CREATE INDEX idx_d3 ON lottery_tickets (d3);
CREATE INDEX idx_d4 ON lottery_tickets (d4);
CREATE INDEX idx_d5 ON lottery_tickets (d5);
CREATE INDEX idx_d6 ON lottery_tickets (d6);

-- Composite index for common multi-digit patterns
CREATE INDEX idx_d5_d6 ON lottery_tickets (d5, d6);
CREATE INDEX idx_d1_d6 ON lottery_tickets (d1, d6);
```

**Why store each digit separately?**

A pattern like `****23` maps directly to `WHERE d5 = '2' AND d6 = '3'` — which uses the index instead of a full table scan. This is significantly faster than `WHERE number LIKE '____23'` on 1M rows.

---

## 2. Wildcard Search Algorithm

### Pattern Parsing

A 6-character pattern is parsed into fixed digits and wildcard positions:

```
Pattern: ****23
→ d5 = '2', d6 = '3'    (fixed)
→ d1, d2, d3, d4         (wildcards — ignored in WHERE clause)

Pattern: 1****5
→ d1 = '1', d6 = '5'    (fixed)

Pattern: 123***
→ d1 = '1', d2 = '2', d3 = '3'  (fixed)
```

### Query Example

For pattern `****23`:

```sql
SELECT id, number
FROM lottery_tickets
WHERE status = 'available'
  AND d5 = '2'
  AND d6 = '3'
LIMIT 50;
```

Since `d5` and `d6` are indexed, PostgreSQL uses an index scan — O(log n) lookup instead of O(n) full scan.

---

## 3. Preventing Duplicate Simultaneous Results

### The Problem

If two users search `****23` at the same time, both queries might return the same ticket before either reserves it.

### Solution: Atomic Reservation with `SELECT FOR UPDATE SKIP LOCKED`

```sql
BEGIN;

SELECT id, number
FROM lottery_tickets
WHERE status = 'available'
  AND d5 = '2'
  AND d6 = '3'
ORDER BY id
LIMIT 1
FOR UPDATE SKIP LOCKED;

-- If a row is returned, reserve it immediately
UPDATE lottery_tickets
SET status = 'reserved',
    reserved_at = NOW(),
    reserved_by = $user_id
WHERE id = $ticket_id;

COMMIT;
```

**How it works:**
- `FOR UPDATE` — locks the selected row so no other transaction can modify it
- `SKIP LOCKED` — other concurrent queries automatically skip rows that are already locked, so they move on to the next available ticket
- The lock is held only for the duration of the transaction (milliseconds), so throughput remains high

### Reservation Expiry

Reservations expire after a fixed window (e.g. 10 minutes) to prevent tickets from being held indefinitely:

```sql
-- Release expired reservations (run periodically)
UPDATE lottery_tickets
SET status = 'available',
    reserved_at = NULL,
    reserved_by = NULL
WHERE status = 'reserved'
  AND reserved_at < NOW() - INTERVAL '10 minutes';
```

---

## 4. Performance Analysis

| Scenario | Approach | Estimated Complexity |
|----------|----------|----------------------|
| Search by 1 fixed digit | Single index scan | O(log n) |
| Search by 2 fixed digits | Composite index scan | O(log n) |
| Reservation | Row-level lock + update | O(1) per ticket |
| Expiry cleanup | Periodic background job | O(expired rows) |

### Estimated throughput

- 1M rows with indexed digit columns → search returns results in **< 5ms**
- `SKIP LOCKED` ensures concurrent reservations do not block each other
- PostgreSQL connection pooling (e.g. PgBouncer) handles high concurrency efficiently

### Tradeoffs

| Decision | Benefit | Tradeoff |
|----------|---------|----------|
| Store digits separately | Fast index-based search | Slightly more storage |
| PostgreSQL over Redis | ACID + persistent storage | Slower than in-memory |
| Row-level locking | No duplicates | Slight overhead per reservation |
| Reservation expiry | Tickets not held forever | Requires background job |

---

## 5. Architecture Summary

```
User Request
     │
     ▼
 API Server
     │
     ├── Parse pattern → extract fixed digits
     │
     ├── Query PostgreSQL
     │   WHERE status = 'available'
     │   AND d{n} = '{digit}' ...
     │   FOR UPDATE SKIP LOCKED
     │
     ├── Reserve ticket atomically in same transaction
     │
     └── Return ticket to user
```

---

## 6. Future Improvements

- **Caching** — cache pattern → available ticket IDs in Redis, invalidate on reservation. Reduces DB load for popular patterns.
- **Read replica** — use a read replica for search queries; write primary only for reservations.
- **Partitioning** — partition the table by ticket number range if dataset grows beyond 10M rows.