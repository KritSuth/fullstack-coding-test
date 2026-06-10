# User Department Grouping API (TypeScript)

A Node.js + TypeScript HTTP server that fetches users from [dummyjson.com](https://dummyjson.com/users) and transforms the data grouped by department.

> **Note:** This is an optional part of the assignment.

## Tech Stack

- **TypeScript** — strict typing throughout
- **Node.js** — native `fetch` (Node 18+)
- **Express** — HTTP framework
- **Jest + ts-jest** — unit testing

## Project Structure

```
backend-ts/
├── src/
│   ├── __tests__/
│   │   └── userService.test.ts  # Unit tests for groupByDepartment
│   ├── index.ts                 # Express server entry point
│   ├── types.ts                 # TypeScript interfaces
│   └── userService.ts           # Fetch + transform logic
└── package.json
└── tsconfig.json                # TypeScript configuration
```

### Key design decisions

- **Single-pass grouping** — `groupByDepartment` iterates users once O(n), using a `Map` to accumulate per-department data before summarizing.
- **Separation of concerns** — fetching (`fetchAndGroupUsers`) is decoupled from transformation (`groupByDepartment`), making the transform logic independently testable without hitting the network.

---

## Getting Started

### Prerequisites

- Node.js >= 18
- npm

### Installation

```bash
git clone https://github.com/KritSuth/fullstack-coding-test.git
cd fullstack-coding-test/backend-ts
npm install
```

### Run locally

```bash
npm start
```

Server runs at `http://localhost:3000`

### Build

```bash
npm run build
```

---

## API Endpoint

### `GET /users/grouped`

Fetches all users from dummyjson and returns them grouped by department.

```bash
curl http://localhost:3000/users/grouped
```

**Sample response:**

```json
{
  "Engineering": {
    "male": 3,
    "female": 2,
    "ageRange": "25-45",
    "hair": {
      "Black": 2,
      "Blond": 1,
      "Brown": 2
    },
    "addressUser": {
      "TerryMedhurst": "10001",
      "AliceSmith": "10002"
    }
  },
  "Marketing": {
    "male": 1,
    "female": 3,
    "ageRange": "28-50",
    "hair": {
      "Chestnut": 2,
      "Black": 2
    },
    "addressUser": {
      "BobJones": "20001"
    }
  }
}
```

**Response fields:**

| Field | Description |
|-------|-------------|
| `male` | Count of male employees |
| `female` | Count of female employees |
| `ageRange` | `"min-max"` age in the department |
| `hair` | Hair color counts `{ color: count }` |
| `addressUser` | `{ firstNamelastName: postalCode }` |

---

## Running Tests

```bash
npm test
```

Tests use mock user data — no network calls required.

Test cases cover:

- Correct department keys are produced
- Gender counts are accurate
- Age range min/max is computed correctly
- Hair colors are aggregated
- Address map uses `firstNamelastName` as key