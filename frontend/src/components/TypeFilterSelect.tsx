const TYPES = [
  'normal', 'fighting', 'flying', 'poison', 'ground', 'rock', 'bug', 'ghost',
  'steel', 'fire', 'water', 'grass', 'electric', 'psychic', 'ice', 'dragon',
  'dark', 'fairy',
];

interface TypeFilterSelectProps {
  value: string;
  onChange: (value: string) => void;
}

export function TypeFilterSelect({ value, onChange }: TypeFilterSelectProps) {
  return (
    <select className="filter-select" value={value} onChange={(e) => onChange(e.target.value)}>
      <option value="">All types</option>
      {TYPES.map((t) => (
        <option key={t} value={t}>
          {t}
        </option>
      ))}
    </select>
  );
}
