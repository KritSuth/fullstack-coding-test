# Auto Delete Todo List
 
A React + TypeScript todo list where items are categorized into **Fruit** and **Vegetable** columns, auto-return to the main list after 5 seconds, and can be manually sent back by clicking them in their column.
 
🔗 **Live Demo:** [https://fullstack-coding-test-tawny.vercel.app](https://fullstack-coding-test-tawny.vercel.app)
 
---
 
## Features
 
- Click any item in the main list → moves to its type column (Fruit / Vegetable)
- A countdown timer bar shows remaining time before auto-return (5 seconds)
- Click an item inside its column → returns to the bottom of the main list immediately
- Empty state shown when a column has no items
## Tech Stack
 
- React 18
- TypeScript
- Vite
---
 
## Getting Started
 
### Prerequisites
 
- Node.js >= 18
- npm or yarn
### Installation
 
```bash
git clone https://github.com/KritSuth/fullstack-coding-test.git
cd fullstack-coding-test/frontend
npm install
```
 
### Run locally
 
```bash
npm run dev
```
 
Open [http://localhost:5173](http://localhost:5173) in your browser.
 
### Build for production
 
```bash
npm run build
```
 
---
 
## Project Structure
 
```
src/
└── App.tsx       # All components and logic (Column, ItemButton, TimerBar)
```
 
### Component overview
 
| Component | Responsibility |
|-----------|---------------|
| `App` | State management — main list, fruit/vegetable columns, timers |
| `Column` | Renders a labeled column with items or an empty state |
| `ItemButton` | Single clickable item, hosts the `TimerBar` when in a column |
| `TimerBar` | Animated progress bar that depletes over `AUTO_RETURN_MS` (5 s) |
 
### Key design decisions
 
- **Timers live in a `useRef`** — avoids re-renders when a timer is set/cleared
- **`returnToMain` is `useCallback`** — stable reference so `setTimeout` closures always call the latest version
- **Single `AUTO_RETURN_MS` constant** — change the value in one place to adjust the delay globally
- **Styles co-located as a `styles` object** — keeps JSX readable without adding a CSS file or extra dependency
---
 
## Deployment (Vercel)
 
This project is deployed on [Vercel](https://vercel.com) with zero configuration.