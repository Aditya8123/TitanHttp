import React, { useEffect, useRef, useState } from 'react';
import { ArrowRight } from 'lucide-react';
import { TerminalCard } from './ui/TerminalCard.tsx';

export const Chapter1Section: React.FC = () => {
  const sectionRef = useRef<HTMLElement>(null);
  const [scrollY, setScrollY] = useState(0);
  const [isInView, setIsInView] = useState(false);

  useEffect(() => {
    const handleScroll = () => setScrollY(window.scrollY);
    window.addEventListener('scroll', handleScroll);
    return () => window.removeEventListener('scroll', handleScroll);
  }, []);

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
    <section id="the-request" ref={sectionRef} style={{
      position: 'relative',
      width: '100%',
      minHeight: '120vh',
      display: 'flex',
      alignItems: 'center',
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
        backgroundImage: `url(${import.meta.env.BASE_URL}images/chapter1.png)`,
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
        background: 'linear-gradient(to right, rgba(0,0,0,0.7) 0%, rgba(0,0,0,0.2) 100%)',
        zIndex: 1
      }} />

      {/* Content Container */}
      <div className="chapter-container" style={{
        position: 'relative',
        zIndex: 2,
        opacity: isInView ? 1 : 0,
        transform: `translateY(${isInView ? 0 : '40px'})`
      }}>
        {/* Left Column: Text */}
        <div className="chapter-column-text">
          <div style={{
            fontFamily: 'var(--font-roboto-mono)',
            fontSize: '14px',
            color: 'var(--color-network-cyan)',
            letterSpacing: '0.1em',
            marginBottom: '16px'
          }}>CHAPTER 01</div>
          
          <h2 className="chapter-heading">
            WHY HTTP<br/>EXISTS
          </h2>

          <p style={{
            fontFamily: 'var(--font-roboto-mono)',
            fontSize: '14px',
            color: 'var(--color-steel-mid)',
            lineHeight: 1.8,
            marginBottom: '40px'
          }}>
            In 1989, networks were chaotic. Every server spoke a different language. 
            HTTP was born not as a complex protocol, but as a simple, human-readable text exchange 
            to request documents over the wire.
            <br/><br/>
            Before we build a server, we must understand the raw bytes.
          </p>

          <button style={{
            background: 'transparent',
            color: 'var(--color-paper-white)',
            border: '1px solid rgba(255, 255, 255, 0.2)',
            padding: '16px 32px',
            fontFamily: 'var(--font-roboto-mono)',
            fontSize: '12px',
            letterSpacing: '0.1em',
            borderRadius: '4px',
            display: 'inline-flex',
            alignItems: 'center',
            gap: '12px',
            cursor: 'pointer',
            transition: 'all 0.3s ease'
          }}
          onMouseEnter={(e) => {
            e.currentTarget.style.borderColor = 'var(--color-network-cyan)';
            e.currentTarget.style.color = 'var(--color-network-cyan)';
          }}
          onMouseLeave={(e) => {
            e.currentTarget.style.borderColor = 'rgba(255, 255, 255, 0.2)';
            e.currentTarget.style.color = 'var(--color-paper-white)';
          }}>
            EXPLORE THE PROTOCOL <ArrowRight size={16} />
          </button>
        </div>

        {/* Right Column: Terminal Card */}
        <div className="chapter-column-visual" style={{
          transform: `translateY(${(scrollY - 500) * -0.1}px)` // Opposing parallax
        }}>
          <TerminalCard title="RAW REQUEST EXAMPLE" glowColor="cyan" style={{ maxWidth: '400px', width: '100%' }}>
            GET / HTTP/1.1{'\n'}
            Host: localhost:8080{'\n'}
            User-Agent: curl/7.81.0{'\n'}
            Accept: */*{'\n'}
          </TerminalCard>
        </div>
      </div>
    </section>
  );
};
