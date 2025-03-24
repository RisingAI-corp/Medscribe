import { useEffect, useRef, useState } from 'react';
import { useDebouncedCallback } from 'use-debounce';

interface UseDebouncedNameChangeProps {
  name: string;
  onChange: (value: string) => void;
  handleUpdateName: (value: string) => void;
}

export function useDebouncedNameChange({
  name,
  onChange,
  handleUpdateName,
}: UseDebouncedNameChangeProps) {
  const nameRef = useRef<HTMLInputElement>(null);
  const [nameValue, setNameValue] = useState(name);

  const debouncedNameChange = useDebouncedCallback((value: string) => {
    if (value.trim() !== '' && value !== name) {
      onChange(value);
      handleUpdateName(value);
    }
  }, 700);

  useEffect(() => {
    setNameValue(name);
  }, [name]);

  useEffect(() => {
    return () => {
      debouncedNameChange.cancel();
    };
  }, [debouncedNameChange]);

  return {
    nameRef,
    nameValue,
    setNameValue,
    debouncedNameChange,
  };
}
