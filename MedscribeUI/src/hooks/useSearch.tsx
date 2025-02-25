import { useState, useEffect } from 'react';
import { search } from 'fast-fuzzy';

function useSearch<T extends { id: string }>(
  list: T[],
  keySelector: (obj: T) => string,
  debounce = 0,
): [T[], string, (query: string) => void] {
  const [query, setQuery] = useState('');
  const [filteredResults, setFilteredResults] = useState(list);

  useEffect(() => {
    const handler = setTimeout(() => {
      if (query.trim() === '') {
        setFilteredResults(list);
      } else {
        setFilteredResults(
          search(query, list, {
            keySelector: keySelector,
          }),
        );
      }
    }, debounce);

    return () => {
      clearTimeout(handler);
    };
  }, [query, list, debounce, keySelector]);

  return [filteredResults, query, setQuery];
}

export default useSearch;
