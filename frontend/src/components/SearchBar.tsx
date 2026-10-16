interface SearchBarProps {
  value: string;
  onChange: (value: string) => void;
}

export function SearchBar({ value, onChange }: SearchBarProps) {
  return (
    <input
      type="text"
      className="search-bar"
      placeholder="Search Pokemon by name..."
      value={value}
      onChange={(e) => onChange(e.target.value)}
    />
  );
}
