import React, { useEffect, useRef, useState } from 'react';
import chapter8Img from '../assets/images/chapter8.png';

export const Chapter8Section: React.FC = () => {
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
    <section id="production" ref={sectionRef} style={{
      position: 'relative',
      width: '100%',
      minHeight: '120vh',
      display: 'flex',
      alignItems: 'center',
      padding: '100px 0',
      overflow: 'hidden'
    }}>
      {/* Background Image */}
      <div style={{
        position: 'absolute',
        top: 0,
        left: 0,
        width: '100%',
        height: '100%',
        backgroundImage: `url(${chapter8Img})`,
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
        background: 'linear-gradient(to left, rgba(0,0,0,0.7) 0%, rgba(0,0,0,0.2) 100%)',
        zIndex: 1
      }} />

      {/* Content Container */}
      <div className="chapter-container reverse" style={{
        position: 'relative',
        zIndex: 2,
        opacity: isInView ? 1 : 0,
        transform: `translateY(${isInView ? 0 : '40px'})`
      }}>
        {/* Right Column: Text */}
        <div className="chapter-column-text">
          <div style={{
            fontFamily: 'var(--font-roboto-mono)',
            fontSize: '14px',
            color: 'var(--color-network-cyan)',
            letterSpacing: '0.1em',
            marginBottom: '16px'
          }}>CHAPTER 08</div>
          
          <h2 className="chapter-heading">
            PRODUCTION<br/>FEATURES
          </h2>

          <p style={{
            fontFamily: 'var(--font-roboto-mono)',
            fontSize: '14px',
            color: 'var(--color-steel-mid)',
            lineHeight: 1.8,
            marginBottom: '40px'
          }}>
            A raw HTTP server is easily overwhelmed in the real world. Production features and security hardening prepare the server for malicious traffic and edge cases.
            <br/><br/>
            TitanHTTP actively defends against Slowloris, Slow POST, header floods, and CL-TE Request Smuggling exploits. Through rigorous validation, it passes 100% of RFC Compliance tests for Method validation, Keep-Alive, and Range requests.
          </p>
        </div>

        {/* Left Column: Security Matrix Visualization */}
        <div className="chapter-column-visual">
          <div style={{
            background: 'rgba(10, 10, 12, 0.8)',
            border: '1px solid rgba(255, 255, 255, 0.1)',
            borderRadius: '16px',
            padding: '32px',
            width: '100%',
            maxWidth: '500px',
            boxShadow: '0 20px 40px rgba(0,0,0,0.5)',
            backdropFilter: 'blur(16px)',
            fontFamily: 'var(--font-jetbrains-mono)',
            position: 'relative',
            overflow: 'hidden'
          }}>
            <style>{`
              @keyframes scanline {
                0% { transform: translateY(-100px); }
                100% { transform: translateY(400px); }
              }
              @keyframes textGlow {
                0%, 100% { text-shadow: 0 0 5px rgba(0, 217, 255, 0.3); opacity: 0.8; }
                50% { text-shadow: 0 0 20px rgba(0, 217, 255, 0.8), 0 0 10px rgba(0, 217, 255, 0.5); opacity: 1; }
              }
            `}</style>

            {/* Scanline */}
            <div style={{
              position: 'absolute', top: 0, left: 0, right: 0, height: '2px',
              background: 'var(--color-network-cyan)',
              opacity: 0.5,
              boxShadow: '0 0 20px var(--color-network-cyan), 0 0 10px var(--color-network-cyan)',
              animation: 'scanline 4s linear infinite'
            }} />

            <div style={{ 
              color: 'var(--color-paper-white)', 
              marginBottom: '24px', 
              fontSize: '14px', 
              borderBottom: '1px solid rgba(255,255,255,0.1)', 
              paddingBottom: '12px',
              letterSpacing: '0.05em',
              overflowX: 'auto',
              whiteSpace: 'nowrap'
            }}>
              &gt; INITIATING SECURITY VALIDATION...
            </div>

            <div style={{ display: 'flex', flexDirection: 'column', gap: '20px', fontSize: '13px' }}>
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', gap: '8px' }}>
                <span style={{ color: 'var(--color-steel-mid)' }}>RFC Compliance Suite</span>
                <span style={{ color: 'var(--color-network-cyan)', animation: 'textGlow 2s infinite' }}>[20/20 PASS]</span>
              </div>
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', gap: '8px' }}>
                <span style={{ color: 'var(--color-steel-mid)' }}>CL-TE Request Smuggling</span>
                <span style={{ color: 'var(--color-network-cyan)', animation: 'textGlow 2s infinite 0.5s' }}>[BLOCKED]</span>
              </div>
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', gap: '8px' }}>
                <span style={{ color: 'var(--color-steel-mid)' }}>Slowloris & Slow POST</span>
                <span style={{ color: 'var(--color-network-cyan)', animation: 'textGlow 2s infinite 1.0s' }}>[DEFENDED]</span>
              </div>
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', gap: '8px' }}>
                <span style={{ color: 'var(--color-steel-mid)' }}>Header Flood Guard</span>
                <span style={{ color: 'var(--color-network-cyan)', animation: 'textGlow 2s infinite 1.5s' }}>[ACTIVE]</span>
              </div>
            </div>

            <div style={{ 
              marginTop: '24px', 
              paddingTop: '16px', 
              borderTop: '1px solid rgba(255,255,255,0.1)', 
              color: 'var(--color-network-cyan)', 
              fontSize: '12px', 
              textAlign: 'right',
              letterSpacing: '0.1em'
            }}>
              SYSTEM SECURE.
            </div>
          </div>
        </div>
      </div>
    </section>
  );
};
