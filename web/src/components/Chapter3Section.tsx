import React, { useEffect, useRef, useState } from 'react';
import { Terminal } from 'lucide-react';

export const Chapter3Section: React.FC = () => {
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
    <section id="building-a-socket" ref={sectionRef} style={{
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
        backgroundImage: 'url(./images/chapter3.png)',
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
          }}>CHAPTER 03</div>
          
          <h2 className="chapter-heading">
            BUILDING A<br />SOCKET
          </h2>

          <p style={{
            fontFamily: 'var(--font-roboto-mono)',
            fontSize: '14px',
            color: 'var(--color-steel-mid)',
            lineHeight: 1.8,
            marginBottom: '40px'
          }}>
            A socket is the absolute boundary between your software and the operating system. 
            By binding to a port (like 8080), we tell the OS: "Any bytes arriving here belong to us."
            <br/><br/>
            In Go, we don't need heavy abstractions to do this. We just open a literal network listener.
          </p>
        </div>

        {/* Left Column: Code Window */}
        <div className="chapter-column-visual">
          <div style={{
            background: 'rgba(10, 10, 12, 0.8)',
            backdropFilter: 'blur(16px)',
            border: '1px solid rgba(255, 255, 255, 0.1)',
            borderRadius: '8px',
            width: '100%',
            maxWidth: '500px',
            overflow: 'hidden',
            boxShadow: '0 30px 60px rgba(0,0,0,0.6)'
          }}>
            {/* Terminal Header */}
            <div style={{
              background: 'rgba(255, 255, 255, 0.05)',
              padding: '12px 16px',
              display: 'flex',
              alignItems: 'center',
              gap: '8px',
              borderBottom: '1px solid rgba(255, 255, 255, 0.1)'
            }}>
              <Terminal size={14} color="var(--color-steel-mid)" />
              <span style={{ fontFamily: 'var(--font-roboto-mono)', fontSize: '11px', color: 'var(--color-steel-mid)' }}>main.go</span>
            </div>
            
            {/* Code Body */}
            <div style={{ padding: '24px' }}>
              <pre style={{
                fontFamily: 'var(--font-jetbrains-mono)',
                fontSize: '13px',
                lineHeight: 1.6,
                margin: 0,
                color: 'var(--color-paper-white)',
                overflowX: 'auto'
              }}>
                {[
                  { len: 13, delay: 0.5, content: <><span style={{ color: '#c678dd' }}>func</span> <span style={{ color: '#61afef' }}>main</span>() {'{'}</> },
                  { len: 45, delay: 0.89, content: <>{'  '}listener, err := net.<span style={{ color: '#61afef' }}>Listen</span>(<span style={{ color: '#98c379' }}>"tcp"</span>, <span style={{ color: '#98c379' }}>":8080"</span>)</> },
                  { len: 17, delay: 2.24, content: <>{'  '}<span style={{ color: '#c678dd' }}>if</span> err != <span style={{ color: '#d19a66' }}>nil</span> {'{'}</> },
                  { len: 18, delay: 2.75, content: <>{'    '}log.<span style={{ color: '#61afef' }}>Fatal</span>(err)</> },
                  { len: 3,  delay: 3.29, content: <>{'  }'}</> },
                  { len: 24, delay: 3.38, content: <>{'  '}defer listener.<span style={{ color: '#61afef' }}>Close</span>()</> },
                  { len: 42, delay: 4.10, content: <>{'  '}log.<span style={{ color: '#61afef' }}>Println</span>(<span style={{ color: '#98c379' }}>"Server bound to port 8080"</span>)</> },
                  { len: 1,  delay: 5.36, content: <>{'}'}</> },
                ].map((line, i) => (
                  <div key={i} style={{
                    overflow: 'hidden',
                    whiteSpace: 'nowrap',
                    width: 'fit-content',
                    maxWidth: isInView ? `${line.len}ch` : '0ch',
                    transition: isInView ? `max-width ${line.len * 0.03}s steps(${line.len}, end) ${line.delay}s` : 'none',
                    borderRight: isInView ? '2px solid transparent' : 'none', // optional cursor, but transparent is safer
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
