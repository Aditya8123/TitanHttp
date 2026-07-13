import React from 'react';

interface GlassPanelProps {
  children: React.ReactNode;
  style?: React.CSSProperties;
  className?: string;
  glowColor?: 'cyan' | 'violet' | 'none';
}

export const GlassPanel: React.FC<GlassPanelProps> = ({ 
  children, 
  style, 
  className,
  glowColor = 'none' 
}) => {
  const getGlowStyle = () => {
    if (glowColor === 'cyan') return { borderLeft: '2px solid var(--color-network-cyan)' };
    if (glowColor === 'violet') return { borderLeft: '2px solid var(--color-packet-violet)' };
    return { borderTop: '1px solid rgba(255, 255, 255, 0.2)' };
  };

  return (
    <div 
      className={className}
      style={{
        background: 'var(--color-frosted-graphite)',
        backdropFilter: 'blur(16px)',
        WebkitBackdropFilter: 'blur(16px)',
        border: '1px solid rgba(255, 255, 255, 0.05)',
        borderRadius: 'var(--radius-glass)',
        padding: '32px',
        boxShadow: '0 20px 40px rgba(0,0,0,0.5)',
        ...getGlowStyle(),
        ...style
      }}
    >
      {children}
    </div>
  );
};
