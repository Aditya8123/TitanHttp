import React from 'react';
import { useParams } from 'react-router-dom';

export const Chapter: React.FC = () => {
  const { id } = useParams<{ id: string }>();

  return (
    <div style={{
      display: 'flex',
      flexDirection: 'column',
      justifyContent: 'center',
      alignItems: 'center',
      height: '100%',
      color: 'var(--color-paper-white)'
    }}>
      <h1 style={{ fontFamily: 'var(--font-lambotype)', fontSize: '80px', margin: 0 }}>CHAPTER {id}</h1>
      <p style={{ fontFamily: 'var(--font-suisse-intl)', fontSize: '18px', color: 'var(--color-steel-mid)' }}>
        Content for this chapter is coming soon.
      </p>
    </div>
  );
};
