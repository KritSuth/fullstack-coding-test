import { useState, useEffect, useCallback, useRef } from "react";

type ItemType = "Fruit" | "Vegetable";

interface TodoItem {
  id: number;
  type: ItemType;
  name: string;
}

const INITIAL_ITEMS: Omit<TodoItem, "id">[] = [
  { type: "Fruit", name: "Apple" },
  { type: "Vegetable", name: "Broccoli" },
  { type: "Vegetable", name: "Mushroom" },
  { type: "Fruit", name: "Banana" },
  { type: "Vegetable", name: "Tomato" },
  { type: "Fruit", name: "Orange" },
  { type: "Fruit", name: "Mango" },
  { type: "Fruit", name: "Pineapple" },
  { type: "Vegetable", name: "Cucumber" },
  { type: "Fruit", name: "Watermelon" },
  { type: "Vegetable", name: "Carrot" },
];

const AUTO_RETURN_MS = 5000;

const ITEMS: TodoItem[] = INITIAL_ITEMS.map((item, index) => ({ ...item, id: index }));

// ─── Styles ────────────────────────────────────────────────────────────────────

const styles = {
  app: {
    padding: "1.5rem",
    fontFamily: "-apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif",
    background: "#f9f9f7",
    minHeight: "100vh",
  } as React.CSSProperties,

  title: {
    fontSize: "20px",
    fontWeight: 600,
    color: "#1a1a1a",
    marginBottom: "1.5rem",
  } as React.CSSProperties,

  grid: {
    display: "grid",
    gridTemplateColumns: "1fr 1fr 1fr",
    gap: "2rem",
  } as React.CSSProperties,

  colHeader: {
    display: "flex",
    alignItems: "center",
    gap: "6px",
    fontSize: "12px",
    fontWeight: 600,
    letterSpacing: "0.07em",
    textTransform: "uppercase" as const,
    marginBottom: "0.75rem",
    paddingBottom: "0.75rem",
    borderBottom: "1px solid #e5e5e5",
  } as React.CSSProperties,

  dot: (color: string) => ({
    width: "10px",
    height: "10px",
    borderRadius: "50%",
    background: color,
    flexShrink: 0,
  } as React.CSSProperties),

  emptyText: {
    color: "#aaa",
    fontStyle: "italic" as const,
    fontSize: "14px",
    padding: "4px 0",
  } as React.CSSProperties,

  itemBtn: {
    display: "block",
    width: "100%",
    textAlign: "left" as const,
    padding: "10px 14px",
    marginBottom: "6px",
    background: "#ffffff",
    border: "1px solid #e5e5e5",
    borderRadius: "8px",
    cursor: "pointer",
    fontSize: "14px",
    color: "#1a1a1a",
    position: "relative" as const,
    overflow: "hidden" as const,
    transition: "border-color 0.15s, background 0.15s",
  } as React.CSSProperties,
};

// ─── App ───────────────────────────────────────────────────────────────────────

export default function App() {
  const [mainList, setMainList] = useState<TodoItem[]>(ITEMS);
  const [fruitColumn, setFruitColumn] = useState<TodoItem[]>([]);
  const [vegColumn, setVegColumn] = useState<TodoItem[]>([]);
  const timersRef = useRef<Record<number, ReturnType<typeof setTimeout>>>({});

  const returnToMain = useCallback((item: TodoItem) => {
    clearTimeout(timersRef.current[item.id]);
    delete timersRef.current[item.id];
    setFruitColumn((prev) => prev.filter((i) => i.id !== item.id));
    setVegColumn((prev) => prev.filter((i) => i.id !== item.id));
    setMainList((prev) => [...prev, item]);
  }, []);

  const moveToColumn = useCallback(
    (item: TodoItem) => {
      setMainList((prev) => prev.filter((i) => i.id !== item.id));
      if (item.type === "Fruit") setFruitColumn((prev) => [...prev, item]);
      else setVegColumn((prev) => [...prev, item]);
      timersRef.current[item.id] = setTimeout(() => returnToMain(item), AUTO_RETURN_MS);
    },
    [returnToMain]
  );

  useEffect(() => {
    const timers = timersRef.current;
    return () => Object.values(timers).forEach(clearTimeout);
  }, []);

  return (
    <div style={styles.app}>
      <h1 style={styles.title}>Todo List</h1>
      <div style={styles.grid}>
        <Column
          label="All Items"
          dotColor={undefined}
          items={mainList}
          onItemClick={moveToColumn}
        />
        <Column
          label="Fruit"
          dotColor="#4caf82"
          items={fruitColumn}
          onItemClick={returnToMain}
          autoReturnMs={AUTO_RETURN_MS}
        />
        <Column
          label="Vegetable"
          dotColor="#5b8dee"
          items={vegColumn}
          onItemClick={returnToMain}
          autoReturnMs={AUTO_RETURN_MS}
        />
      </div>
    </div>
  );
}

// ─── Column ────────────────────────────────────────────────────────────────────

interface ColumnProps {
  label: string;
  dotColor?: string;
  items: TodoItem[];
  onItemClick: (item: TodoItem) => void;
  autoReturnMs?: number;
}

function Column({ label, dotColor, items, onItemClick, autoReturnMs }: ColumnProps) {
  return (
    <div>
      <div style={styles.colHeader}>
        {dotColor && <span style={styles.dot(dotColor)} />}
        <span style={{ color: dotColor ?? "#888" }}>{label}</span>
      </div>

      {items.length === 0 ? (
        <span style={styles.emptyText}>Empty</span>
      ) : (
        items.map((item) => (
          <ItemButton
            key={item.id}
            item={item}
            onClick={onItemClick}
            autoReturnMs={autoReturnMs}
          />
        ))
      )}
    </div>
  );
}

// ─── Item Button ───────────────────────────────────────────────────────────────

interface ItemButtonProps {
  item: TodoItem;
  onClick: (item: TodoItem) => void;
  autoReturnMs?: number;
}

function ItemButton({ item, onClick, autoReturnMs }: ItemButtonProps) {
  const addedAt = useRef(Date.now());

  return (
    <button style={styles.itemBtn} onClick={() => onClick(item)}>
      {item.name}
      {autoReturnMs && (
        <TimerBar durationMs={autoReturnMs} startedAt={addedAt.current} />
      )}
    </button>
  );
}

// ─── Timer Bar ─────────────────────────────────────────────────────────────────

function TimerBar({ durationMs, startedAt }: { durationMs: number; startedAt: number }) {
  const [width, setWidth] = useState(100);

  useEffect(() => {
    const tick = () => {
      const elapsed = Date.now() - startedAt;
      setWidth(Math.max(0, (1 - elapsed / durationMs) * 100));
    };
    const interval = setInterval(tick, 50);
    return () => clearInterval(interval);
  }, [durationMs, startedAt]);

  return (
    <div
      style={{
        position: "absolute",
        bottom: 0,
        left: 0,
        height: "3px",
        width: `${width}%`,
        background: "#e07b54",
        transition: "width 50ms linear",
        borderRadius: "0 0 0 8px",
      }}
    />
  );
}
