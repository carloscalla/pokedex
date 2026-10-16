import { usePokemonDetail } from '../api/queries';

interface PokemonDetailModalProps {
  name: string;
  onClose: () => void;
}

const STAT_LABELS: Record<string, string> = {
  hp: 'HP',
  attack: 'Attack',
  defense: 'Defense',
  'special-attack': 'Sp. Attack',
  'special-defense': 'Sp. Defense',
  speed: 'Speed',
};

export function PokemonDetailModal({ name, onClose }: PokemonDetailModalProps) {
  const { data, isLoading, isError, error } = usePokemonDetail(name);

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content" onClick={(e) => e.stopPropagation()}>
        <button className="modal-close" onClick={onClose} aria-label="Close">
          &times;
        </button>

        {isLoading && <p>Loading...</p>}
        {isError && <p className="error-text">{(error as Error).message}</p>}

        {data && (
          <>
            <h2>{data.name}</h2>
            {data.spriteUrl && <img src={data.spriteUrl} alt={data.name} className="modal-sprite" />}
            <div className="type-badge-row">
              {data.types.map((t) => (
                <span key={t} className={`type-badge type-${t}`}>
                  {t}
                </span>
              ))}
            </div>
            <p>
              Height: {data.height / 10} m · Weight: {data.weight / 10} kg
            </p>
            <h3>Stats</h3>
            <ul className="stat-list">
              {Object.entries(data.stats).map(([stat, value]) => (
                <li key={stat}>
                  <span>{STAT_LABELS[stat] ?? stat}</span>
                  <span>{value}</span>
                </li>
              ))}
            </ul>
            <h3>Abilities</h3>
            <p>{data.abilities.join(', ')}</p>
          </>
        )}
      </div>
    </div>
  );
}
