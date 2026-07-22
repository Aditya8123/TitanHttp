import React, { useEffect, useRef, useState } from 'react';
import { Terminal, Code } from 'lucide-react';
import chapter10Img from '../assets/images/chapter10.jpeg';

export const Chapter10Section: React.FC = () => {
  const sectionRef = useRef<HTMLElement>(null);
  const [isInView, setIsInView] = useState(false);

  useEffect(() => {
    const observer = new IntersectionObserver(
      ([entry]) => {
        setIsInView(entry.isIntersecting);
      },
      { threshold: 0.3 }
    );
    if (sectionRef.current) observer.observe(sectionRef.current);
    return () => observer.disconnect();
  }, []);

  return (
    <section id="alive" ref={sectionRef} style={{
      position: 'relative',
      width: '100%',
      minHeight: '100vh',
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center',
      padding: '100px 0',
      overflow: 'hidden'
    }}>
      {/* Background Image Parallax */}
      <div style={{
        position: 'absolute',
        top: 0,
        left: 0,
        width: '100%',
        height: '100%',
        backgroundImage: `url(${chapter10Img})`,
        backgroundSize: 'cover',
        backgroundPosition: 'center',
        backgroundAttachment: 'fixed',
        zIndex: 0
      }} />

      <div style={{
        position: 'absolute',
        top: 0,
        left: 0,
        width: '100%',
        height: '100%',
        background: 'radial-gradient(circle at center, transparent 0%, rgba(0,0,0,0.4) 100%)',
        zIndex: 1
      }} />

      <div style={{
        position: 'absolute',
        top: 0,
        left: 0,
        width: '100%',
        height: '100%',
        background: 'rgba(0,0,0,0.75)',
        zIndex: 1
      }} />

      <style>{`
        .c10-buttons-container {
          display: flex;
          gap: 24px;
          align-items: center;
        }
        .c10-footer {
          position: absolute;
          bottom: 0;
          left: 0;
          width: 100%;
          padding: 24px 40px;
          display: flex;
          justify-content: space-between;
          align-items: center;
          border-top: 1px solid rgba(255, 255, 255, 0.05);
          background: rgba(0, 0, 0, 0.5);
          backdrop-filter: blur(8px);
          z-index: 10;
          font-family: var(--font-roboto-mono);
          font-size: 12px;
          color: var(--color-steel-mid);
        }
        @media (max-width: 640px) {
          .c10-buttons-container {
            flex-direction: column !important;
            gap: 16px !important;
            width: 100%;
            padding: 0 20px;
          }
          .c10-buttons-container a {
            width: 100% !important;
            justify-content: center !important;
          }
          .c10-footer {
            flex-direction: column !important;
            gap: 12px !important;
            padding: 16px 20px !important;
            text-align: center !important;
            position: relative !important;
            margin-top: 40px;
          }
        }
      `}</style>

      {/* Content Container */}
      <div style={{
        position: 'relative',
        zIndex: 2,
        maxWidth: '800px',
        margin: '0 auto',
        width: '100%',
        padding: '0 24px',
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        textAlign: 'center',
        opacity: isInView ? 1 : 0,
        transform: `scale(${isInView ? 1 : 0.95})`,
        transition: 'all 1s cubic-bezier(0.2, 0.8, 0.2, 1)'
      }}>
        <div style={{
          fontFamily: 'var(--font-roboto-mono)',
          fontSize: '14px',
          color: 'var(--color-network-cyan)',
          letterSpacing: '0.1em',
          marginBottom: '16px'
        }}>CHAPTER 10</div>

        <h2 className="chapter-heading" style={{
          textShadow: '0 0 40px rgba(0, 217, 255, 0.3)'
        }}>
          SOURCE CODE
        </h2>

        <p style={{
          fontFamily: 'var(--font-roboto-mono)',
          fontSize: '16px',
          color: 'var(--color-steel-mid)',
          lineHeight: 1.8,
          marginBottom: '40px',
          maxWidth: '600px'
        }}>
          A cinematic orchestration of raw infrastructure. TitanHTTP is built completely from scratch in Go. No standard library crutches. No shortcuts. Just pure networking, memory optimization, and zero-allocation routing.
        </p>

        {/* The Quick Start Terminal */}
        <div style={{
          background: 'rgba(10, 10, 12, 0.8)',
          border: '1px solid rgba(255, 255, 255, 0.1)',
          borderRadius: '12px',
          padding: '24px',
          width: '100%',
          maxWidth: '600px',
          backdropFilter: 'blur(16px)',
          textAlign: 'left',
          marginBottom: '40px',
          boxShadow: '0 20px 40px rgba(0,0,0,0.5)'
        }}>
          <div style={{ display: 'flex', gap: '8px', marginBottom: '16px' }}>
            <div style={{ width: '12px', height: '12px', borderRadius: '50%', background: '#ff5f56' }} />
            <div style={{ width: '12px', height: '12px', borderRadius: '50%', background: '#ffbd2e' }} />
            <div style={{ width: '12px', height: '12px', borderRadius: '50%', background: '#27c93f' }} />
          </div>
          <pre style={{
            fontFamily: 'var(--font-jetbrains-mono)',
            fontSize: '14px',
            color: 'var(--color-paper-white)',
            margin: 0,
            lineHeight: 1.8,
            overflowX: 'auto'
          }}>
            <span style={{ color: 'var(--color-steel-mid)' }}># Clone the repository</span>{'\n'}
            <span style={{ color: 'var(--color-network-cyan)' }}>git</span> clone https://github.com/Aditya8123/TitanHttp.git{'\n\n'}
            <span style={{ color: 'var(--color-steel-mid)' }}># Enter the datacenter</span>{'\n'}
            <span style={{ color: 'var(--color-network-cyan)' }}>cd</span> TitanHttp{'\n\n'}
            <span style={{ color: 'var(--color-steel-mid)' }}># Boot the server</span>{'\n'}
            <span style={{ color: 'var(--color-network-cyan)' }}>go</span> run cmd/titanhttp/main.go
          </pre>
        </div>

        <div className="c10-buttons-container">
          <a href="https://github.com/Aditya8123/TitanHttp" target="_blank" rel="noreferrer" style={{
            background: 'var(--color-paper-white)',
            color: '#000',
            border: 'none',
            padding: '16px 32px',
            fontFamily: 'var(--font-roboto-mono)',
            fontSize: '14px',
            fontWeight: 'bold',
            letterSpacing: '0.1em',
            borderRadius: '4px',
            display: 'inline-flex',
            alignItems: 'center',
            gap: '12px',
            cursor: 'pointer',
            boxShadow: '0 0 20px rgba(255, 255, 255, 0.2)',
            transition: 'transform 0.2s ease',
            textDecoration: 'none'
          }}
            onMouseEnter={(e) => e.currentTarget.style.transform = 'translateY(-2px)'}
            onMouseLeave={(e) => e.currentTarget.style.transform = 'translateY(0)'}
          >
            <Code size={18} /> VIEW ON GITHUB
          </a>

          <a href="https://github.com/Aditya8123/TitanHttp/blob/main/README.md" target="_blank" rel="noreferrer" style={{
            background: 'rgba(255, 255, 255, 0.05)',
            color: 'var(--color-paper-white)',
            border: '1px solid rgba(255, 255, 255, 0.2)',
            padding: '16px 32px',
            fontFamily: 'var(--font-roboto-mono)',
            fontSize: '14px',
            letterSpacing: '0.1em',
            borderRadius: '4px',
            display: 'inline-flex',
            alignItems: 'center',
            gap: '12px',
            cursor: 'pointer',
            transition: 'all 0.2s ease',
            textDecoration: 'none'
          }}
            onMouseEnter={(e) => {
              e.currentTarget.style.background = 'rgba(255, 255, 255, 0.1)';
              e.currentTarget.style.transform = 'translateY(-2px)';
            }}
            onMouseLeave={(e) => {
              e.currentTarget.style.background = 'rgba(255, 255, 255, 0.05)';
              e.currentTarget.style.transform = 'translateY(0)';
            }}
          >
            <Terminal size={18} /> READ DOCUMENTATION
          </a>
        </div>
      </div>

      {/* Footer */}
      <div className="c10-footer">
        <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
          <span style={{ color: 'var(--color-network-cyan)' }}>TITAN</span>HTTP © {new Date().getFullYear()}
        </div>
        <div style={{ display: 'flex', gap: '24px' }}>
          <a href="https://github.com/Aditya8123/TitanHttp" target="_blank" rel="noreferrer" style={{ color: 'inherit', textDecoration: 'none', transition: 'color 0.2s' }} onMouseEnter={(e) => e.currentTarget.style.color = 'var(--color-network-cyan)'} onMouseLeave={(e) => e.currentTarget.style.color = 'var(--color-steel-mid)'}>GitHub</a>
          <a href="#" style={{ color: 'inherit', textDecoration: 'none', transition: 'color 0.2s' }} onMouseEnter={(e) => e.currentTarget.style.color = 'var(--color-network-cyan)'} onMouseLeave={(e) => e.currentTarget.style.color = 'var(--color-steel-mid)'}>Documentation</a>
        </div>
      </div>
    </section>
  );
};
