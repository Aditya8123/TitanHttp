import React, { useEffect, useRef, useState } from 'react';

export const Chapter5Section: React.FC = () => {
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
    <section id="parsing" ref={sectionRef} style={{
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
        backgroundImage: `url(${import.meta.env.BASE_URL}images/chapter5.png)`,
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
          }}>CHAPTER 05</div>
          
          <h2 className="chapter-heading">
            PARSING<br />REQUESTS
          </h2>

          <p style={{
            fontFamily: 'var(--font-roboto-mono)',
            fontSize: '14px',
            color: 'var(--color-steel-mid)',
            lineHeight: 1.8,
            marginBottom: '40px'
          }}>
            The network gives us raw bytes. It's our job to bring order to the chaos. 
            A parser reads the byte stream line by line, isolating the Method, Path, Protocol, and Headers.
            <br/><br/>
            We construct a strictly typed Go `Request` struct, transforming a chaotic string into a structured data object.
          </p>
        </div>

        {/* Right Column: Parsing Visualization */}
        <div className="chapter-column-visual">
          <div style={{
            background: 'rgba(20, 20, 20, 0.8)',
            border: '1px solid rgba(255, 255, 255, 0.1)',
            borderRadius: '8px',
            padding: '40px',
            maxWidth: '500px',
            width: '100%',
            display: 'flex',
            flexDirection: 'column',
            gap: '24px',
            boxShadow: '0 20px 40px rgba(0,0,0,0.5)'
          }}>
            {/* Raw String */}
            <div style={{ borderBottom: '1px solid rgba(255,255,255,0.1)', paddingBottom: '16px' }}>
              <div style={{ fontFamily: 'var(--font-roboto-mono)', fontSize: '10px', color: 'var(--color-steel-mid)', marginBottom: '8px' }}>RAW BYTE STREAM</div>
              <div style={{ fontFamily: 'var(--font-jetbrains-mono)', fontSize: '14px', color: '#e06c75', overflowX: 'auto' }}>
                "GET /api/v1/users HTTP/1.1\r\n"
              </div>
            </div>

            {/* Down Arrow */}
            <div style={{ 
              textAlign: 'center', 
              color: 'var(--color-network-cyan)',
              animation: 'pulse 1.5s infinite ease-in-out'
            }}>↓</div>

            {/* Parsed Struct */}
            <div>
              <div style={{ fontFamily: 'var(--font-roboto-mono)', fontSize: '10px', color: 'var(--color-steel-mid)', marginBottom: '8px' }}>HTTP.REQUEST STRUCT</div>
              <pre style={{ 
                fontFamily: 'var(--font-jetbrains-mono)', 
                fontSize: '13px', 
                color: 'var(--color-paper-white)', 
                margin: 0, 
                lineHeight: 1.6,
                overflowX: 'auto'
              }}>
                {[
                  { len: 20, delay: 0.5, content: <>type Request struct {'{'}</> },
                  { len: 17, delay: 1.1, content: <>{'  '}Method: <span style={{ color: '#98c379' }}>"GET"</span>,</> },
                  { len: 24, delay: 1.61, content: <>{'  '}Path:   <span style={{ color: '#98c379' }}>"/api/v1/users"</span>,</> },
                  { len: 20, delay: 2.33, content: <>{'  '}Proto:  <span style={{ color: '#98c379' }}>"HTTP/1.1"</span>,</> },
                  { len: 1,  delay: 2.93, content: <>{'}'}</> },
                ].map((line, i) => (
                  <div key={i} style={{
                    overflow: 'hidden',
                    whiteSpace: 'nowrap',
                    width: 'fit-content',
                    maxWidth: isInView ? `${line.len}ch` : '0ch',
                    transition: isInView ? `max-width ${line.len * 0.03}s steps(${line.len}, end) ${line.delay}s` : 'none',
                  }}>
                    {line.content}
                  </div>
                ))}
              </pre>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
};
