import type { TeamMember } from '../api/types';

interface TeamSlotProps {
  member: TeamMember;
  onRemove: (pokemonId: number) => void;
  removing: boolean;
}

export function TeamSlot({ member, onRemove, removing }: TeamSlotProps) {
  return (
    <div className="team-slot">
      {member.spriteUrl && <img src={member.spriteUrl} alt={member.name} />}
      <span className="pokemon-card-name">{member.name}</span>
      <div className="type-badge-row">
        {member.types.map((t) => (
          <span key={t} className={`type-badge type-${t}`}>
            {t}
          </span>
        ))}
      </div>
      <button disabled={removing} onClick={() => onRemove(member.pokemonId)}>
        Remove
      </button>
    </div>
  );
}
