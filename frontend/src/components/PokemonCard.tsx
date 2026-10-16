import type { PokemonSummary } from '../api/types';
import { useAuth } from '../context/AuthContext';
import { useTeam, useAddToTeam } from '../api/queries';

interface PokemonCardProps {
  pokemon: PokemonSummary;
  onSelect: (name: string) => void;
}

export function PokemonCard({ pokemon, onSelect }: PokemonCardProps) {
  const { isAuthenticated } = useAuth();
  const teamQuery = useTeam(isAuthenticated);
  const addToTeam = useAddToTeam();

  const team = teamQuery.data;
  const alreadyOnTeam = team?.members.some((m) => m.pokemonId === pokemon.id) ?? false;
  const teamFull = (team?.size ?? 0) >= (team?.max ?? 6);

  // Mirrors the backend's own validation order (not-authenticated, then
  // already-a-member, then full) purely for UX; the mutation's response is
  // still the actual source of truth if these ever disagree.
  let addLabel = 'Add to team';
  let addDisabled = false;
  if (!isAuthenticated) {
    addLabel = 'Log in to add';
    addDisabled = true;
  } else if (alreadyOnTeam) {
    addLabel = 'On team';
    addDisabled = true;
  } else if (teamFull) {
    addLabel = 'Team full';
    addDisabled = true;
  }

  return (
    <div className="pokemon-card">
      <button className="pokemon-card-main" onClick={() => onSelect(pokemon.name)}>
        {pokemon.spriteUrl && <img src={pokemon.spriteUrl} alt={pokemon.name} loading="lazy" />}
        <span className="pokemon-card-name">{pokemon.name}</span>
        <div className="type-badge-row">
          {pokemon.types.map((t) => (
            <span key={t} className={`type-badge type-${t}`}>
              {t}
            </span>
          ))}
        </div>
      </button>
      <button
        className="pokemon-card-add"
        disabled={addDisabled || addToTeam.isPending}
        onClick={() => addToTeam.mutate(pokemon.id)}
      >
        {addLabel}
      </button>
    </div>
  );
}
