import React from 'react';

interface TerminalCardProps {
  title?: string;
  glowColor?: 'cyan' | 'violet' | 'none';
  children: React.ReactNode;
  style?: React.CSSProperties;
  className?: string;
}

export const TerminalCard: React.FC<TerminalCardProps> = ({
  title,
  glowColor = 'none',
  children,
  style,
  className
}) => {
  const getGlowColorHex = () => {
    if (glowColor === 'cyan') return 'var(--color-network-cyan)';
    if (glowColor === 'violet') return 'var(--color-packet-violet)';
    return 'var(--color-paper-white)';
  };

  return (
    <div 
      className={className}
      style={{
        background: 'rgba(10, 10, 12, 0.8)',
        backdropFilter: 'blur(12px)',
        WebkitBackdropFilter: 'blur(12px)',
        border: '1px solid rgba(255, 255, 255, 0.05)',
        borderLeft: glowColor !== 'none' ? `2px solid ${getGlowColorHex()}` : undefined,
        padding: '24px',
        width: '100%',
        boxShadow: '0 20px 40px rgba(0,0,0,0.5)',
        ...style
      }}
    >
      {title && (
        <div style={{ 
          fontFamily: 'var(--font-roboto-mono)', 
          fontSize: '10px', 
          color: getGlowColorHex(), 
          letterSpacing: '0.1em', 
          marginBottom: '16px' 
        }}>
          {title}
        </div>
      )}
      <pre style={{ 
        margin: 0, 
        fontFamily: 'var(--font-jetbrains-mono)', 
        fontSize: '12px', 
        color: 'var(--color-steel-mid)', 
        lineHeight: 1.6,
        whiteSpace: 'pre-wrap',
        wordBreak: 'break-word'
      }}>
        {children}
      </pre>
    </div>
  );
};
