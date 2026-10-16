const SORT_OPTIONS: { label: string; value: string }[] = [
  { label: 'Dex order', value: '' },
  { label: 'Name (A-Z)', value: 'name' },
  { label: 'Name (Z-A)', value: '-name' },
  { label: 'Attack (high-low)', value: '-attack' },
  { label: 'Attack (low-high)', value: 'attack' },
  { label: 'Defense (high-low)', value: '-defense' },
  { label: 'Speed (high-low)', value: '-speed' },
  { label: 'HP (high-low)', value: '-hp' },
];

interface SortSelectProps {
  value: string;
  onChange: (value: string) => void;
}

export function SortSelect({ value, onChange }: SortSelectProps) {
  return (
    <select className="filter-select" value={value} onChange={(e) => onChange(e.target.value)}>
      {SORT_OPTIONS.map((opt) => (
        <option key={opt.value} value={opt.value}>
          {opt.label}
        </option>
      ))}
    </select>
  );
}
