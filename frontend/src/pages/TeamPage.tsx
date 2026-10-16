import { useTeam, useRemoveFromTeam } from '../api/queries';
import { useAuth } from '../context/AuthContext';
import { TeamSlot } from '../components/TeamSlot';

export function TeamPage() {
  const { isAuthenticated } = useAuth();
  const { data, isLoading, isError, error } = useTeam(isAuthenticated);
  const removeFromTeam = useRemoveFromTeam();

  if (isLoading) return <p>Loading your team...</p>;
  if (isError) return <p className="error-text">{(error as Error).message}</p>;
  if (!data) return null;

  return (
    <div className="team-page">
      <h2>
        My Team ({data.size}/{data.max})
      </h2>
      {data.members.length === 0 && <p>You haven't added any Pokemon yet. Go browse the catalog!</p>}
      <div className="team-grid">
        {data.members.map((m) => (
          <TeamSlot
            key={m.pokemonId}
            member={m}
            removing={removeFromTeam.isPending}
            onRemove={(id) => removeFromTeam.mutate(id)}
          />
        ))}
      </div>
    </div>
  );
}
