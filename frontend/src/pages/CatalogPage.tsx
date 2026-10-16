import { useState } from 'react';
import { usePokemonList } from '../api/queries';
import { useDebouncedValue } from '../hooks/useDebouncedValue';
import { SearchBar } from '../components/SearchBar';
import { TypeFilterSelect } from '../components/TypeFilterSelect';
import { SortSelect } from '../components/SortSelect';
import { Pagination } from '../components/Pagination';
import { PokemonCard } from '../components/PokemonCard';
import { PokemonDetailModal } from '../components/PokemonDetailModal';

export function CatalogPage() {
  const [searchInput, setSearchInput] = useState('');
  const [type, setType] = useState('');
  const [sort, setSort] = useState('');
  const [page, setPage] = useState(1);
  const [selectedName, setSelectedName] = useState<string | null>(null);

  const search = useDebouncedValue(searchInput, 300);

  const { data, isLoading, isError, error, isFetching } = usePokemonList({ page, search, type, sort });

  function updateSearch(value: string) {
    setSearchInput(value);
    setPage(1);
  }
  function updateType(value: string) {
    setType(value);
    setPage(1);
  }
  function updateSort(value: string) {
    setSort(value);
    setPage(1);
  }

  return (
    <div className="catalog-page">
      <div className="catalog-controls">
        <SearchBar value={searchInput} onChange={updateSearch} />
        <TypeFilterSelect value={type} onChange={updateType} />
        <SortSelect value={sort} onChange={updateSort} />
      </div>

      {isLoading && <p>Loading Pokemon...</p>}
      {isError && <p className="error-text">{(error as Error).message}</p>}

      {data && data.items.length === 0 && (
        <p>No Pokemon match your search{search ? ` for "${search}"` : ''}.</p>
      )}

      {data && data.items.length > 0 && (
        <>
          <div className={`pokemon-grid ${isFetching ? 'is-fetching' : ''}`}>
            {data.items.map((p) => (
              <PokemonCard key={p.id} pokemon={p} onSelect={setSelectedName} />
            ))}
          </div>
          <Pagination page={data.page} totalPages={data.totalPages} onChange={setPage} />
        </>
      )}

      {selectedName && <PokemonDetailModal name={selectedName} onClose={() => setSelectedName(null)} />}
    </div>
  );
}
