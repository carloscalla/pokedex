export interface PokemonSummary {
  id: number;
  name: string;
  spriteUrl: string | null;
  types: string[];
}

export interface PokemonDetail {
  id: number;
  name: string;
  spriteUrl: string | null;
  height: number;
  weight: number;
  types: string[];
  stats: Record<string, number>;
  abilities: string[];
}

export interface PaginatedResponse<T> {
  items: T[];
  page: number;
  pageSize: number;
  total: number;
  totalPages: number;
}

export interface AuthUser {
  id: number;
  username: string;
}

export interface AuthResponse {
  token: string;
  user: AuthUser;
}

export interface TeamMember {
  pokemonId: number;
  name: string;
  spriteUrl: string | null;
  types: string[];
  position: number;
}

export interface TeamResponse {
  members: TeamMember[];
  size: number;
  max: number;
}

export interface ApiErrorBody {
  error: string;
}
