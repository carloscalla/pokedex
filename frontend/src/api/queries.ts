import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { request } from './client';
import type { PaginatedResponse, PokemonSummary, PokemonDetail, TeamResponse } from './types';

export interface PokemonListParams {
  page: number;
  search: string;
  type: string;
  sort: string;
}

function buildListUrl(params: PokemonListParams): string {
  const qs = new URLSearchParams();
  qs.set('page', String(params.page));
  if (params.search) qs.set('search', params.search);
  if (params.type) qs.set('type', params.type);
  if (params.sort) qs.set('sort', params.sort);
  return `/api/pokemon?${qs.toString()}`;
}

export function usePokemonList(params: PokemonListParams) {
  return useQuery({
    // The query key includes every filter param, so each distinct
    // page/search/type/sort combination caches independently.
    queryKey: ['pokemon', params],
    queryFn: () => request<PaginatedResponse<PokemonSummary>>(buildListUrl(params)),
    placeholderData: (previous) => previous,
  });
}

export function usePokemonDetail(name: string | null) {
  return useQuery({
    queryKey: ['pokemon', 'detail', name],
    queryFn: () => request<PokemonDetail>(`/api/pokemon/${encodeURIComponent(name ?? '')}`),
    enabled: !!name,
  });
}

export function useTeam(enabled: boolean) {
  return useQuery({
    queryKey: ['team'],
    queryFn: () => request<TeamResponse>('/api/team'),
    enabled,
  });
}

export function useAddToTeam() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (pokemonId: number) => request<TeamResponse>(`/api/team/${pokemonId}`, { method: 'POST' }),
    // The backend returns the updated team on a successful mutation, so we
    // write it straight into the cache instead of invalidating + refetching.
    onSuccess: (data) => queryClient.setQueryData(['team'], data),
  });
}

export function useRemoveFromTeam() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (pokemonId: number) => request<TeamResponse>(`/api/team/${pokemonId}`, { method: 'DELETE' }),
    onSuccess: (data) => queryClient.setQueryData(['team'], data),
  });
}
