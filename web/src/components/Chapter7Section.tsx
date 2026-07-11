import React, { useEffect, useRef, useState } from 'react';

export const Chapter7Section: React.FC = () => {
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
    <section id="concurrency" ref={sectionRef} style={{
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
        backgroundImage: 'url(/images/chapter7.png)',
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
      <div style={{
        position: 'relative',
        zIndex: 2,
        maxWidth: '1440px',
        margin: '0 auto',
        width: '100%',
        padding: '0 40px',
        display: 'flex',
        gap: '80px',
        opacity: isInView ? 1 : 0,
        transform: `translateY(${isInView ? 0 : '40px'})`,
        transition: 'all 1s cubic-bezier(0.2, 0.8, 0.2, 1)'
      }}>
        {/* Left Column: Text */}
        <div style={{ flex: 1, maxWidth: '500px' }}>
          <div style={{
            fontFamily: 'var(--font-roboto-mono)',
            fontSize: '14px',
            color: 'var(--color-network-cyan)',
            letterSpacing: '0.1em',
            marginBottom: '16px'
          }}>CHAPTER 07</div>
          
          <h2 style={{
            fontFamily: 'var(--font-lambotype)',
            fontSize: '80px',
            lineHeight: 0.9,
            color: 'var(--color-paper-white)',
            marginBottom: '32px',
            textTransform: 'uppercase'
          }}>
            CONCURRENCY
          </h2>

          <p style={{
            fontFamily: 'var(--font-roboto-mono)',
            fontSize: '14px',
            color: 'var(--color-steel-mid)',
            lineHeight: 1.8,
            marginBottom: '40px'
          }}>
            A production server cannot block while handling a single request. 
            When 10,000 users connect simultaneously, the server must handle them all at once.
            <br/><br/>
            In Go, this is handled via Goroutines. By spawning <code style={{ color: 'var(--color-network-cyan)' }}>go handle(conn)</code>, 
            TitanHTTP processes thousands of concurrent requests independently, seamlessly multiplexed by the Go runtime.
          </p>
        </div>

        {/* Right Column: Code Visualization */}
        <div style={{
          flex: 1,
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'flex-end'
        }}>
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
            <div style={{ padding: '24px' }}>
              <pre style={{
                fontFamily: 'var(--font-jetbrains-mono)',
                fontSize: '14px',
                lineHeight: 1.6,
                margin: 0,
                color: 'var(--color-paper-white)'
              }}>
                {[
                  { len: 7,  delay: 0.5, content: <><span style={{ color: '#c678dd' }}>for</span> {'{'}</> },
                  { len: 32, delay: 0.71, content: <>{'  '}conn, err := listener.<span style={{ color: '#61afef' }}>Accept</span>()</> },
                  { len: 17, delay: 1.67, content: <>{'  '}<span style={{ color: '#c678dd' }}>if</span> err != <span style={{ color: '#d19a66' }}>nil</span> {'{'}</> },
                  { len: 12, delay: 2.18, content: <>{'    '}<span style={{ color: '#c678dd' }}>continue</span></> },
                  { len: 3,  delay: 2.54, content: <>{'  }'}</> },
                  { len: 0,  delay: 2.63, content: <>&nbsp;</> },
                  { len: 27, delay: 2.68, content: <>{'  '}<span style={{ 
                    background: 'rgba(0, 217, 255, 0.2)', 
                    padding: '2px 4px', 
                    borderRadius: '2px', 
                    color: 'var(--color-network-cyan)',
                    fontWeight: 'bold',
                    boxShadow: '0 0 10px rgba(0, 217, 255, 0.5)'
                  }}>go handleConnection(conn)</span></> },
                  { len: 1,  delay: 3.49, content: <>{'}'}</> }
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
