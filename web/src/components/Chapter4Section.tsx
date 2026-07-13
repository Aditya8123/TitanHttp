import React, { useEffect, useRef, useState } from 'react';

export const Chapter4Section: React.FC = () => {
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
    <section id="the-listener-loop" ref={sectionRef} style={{
      position: 'relative',
      width: '100%',
      minHeight: '120vh',
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
        backgroundImage: 'url(/images/chapter4.png)',
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
        background: 'linear-gradient(to right, rgba(0,0,0,0.7) 0%, rgba(0,0,0,0.2) 100%)', // Lighter gradient overlay
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
          }}>CHAPTER 04</div>

          <h2 className="chapter-heading">
            READING BYTES
          </h2>
          <p style={{
            fontFamily: 'var(--font-roboto-mono)',
            fontSize: '14px',
            color: 'var(--color-steel-mid)',
            lineHeight: 1.8
          }}>
            Once a connection is established, the server must read the raw incoming bytes.
            An HTTP request is just a stream of ASCII characters ending in <code style={{ color: 'var(--color-network-cyan)' }}>\r\n\r\n</code>.
            TitanHTTP continuously reads from the socket into a fixed-size buffer.
          </p>
        </div>

        <style>{`
          @keyframes byteFlow {
            0% { left: -100px; opacity: 0; transform: scale(0.8); }
            10% { opacity: 1; transform: scale(1); }
            90% { opacity: 1; transform: scale(1); }
            100% { left: 100%; opacity: 0; transform: scale(0.5); }
          }
          @keyframes bufferPulse {
            0%, 100% { box-shadow: 0 0 30px rgba(0, 217, 255, 0.1); border-color: rgba(255, 255, 255, 0.1); }
            50% { box-shadow: 0 0 60px rgba(0, 217, 255, 0.4); border-color: rgba(0, 217, 255, 0.8); }
          }
          @keyframes dataStreamGlow {
            0%, 100% { opacity: 0.3; }
            50% { opacity: 1; }
          }
          @keyframes hexFall {
            0% { transform: translateY(-50px); opacity: 0; }
            20% { opacity: 1; }
            100% { transform: translateY(150px); opacity: 0; }
          }
          
          @media (max-width: 640px) {
            .responsive-byte-stream-wrapper {
              flex-direction: column !important;
              gap: 24px !important;
              width: 100%;
            }
            .responsive-byte-stream-pipe {
              width: 100% !important;
              height: 50px !important;
              border-right: 1px dashed rgba(0, 217, 255, 0.3) !important;
              border-bottom: none !important;
            }
            .responsive-byte-stream-buffer {
              width: 180px !important;
              height: 180px !important;
            }
          }
        `}</style>

        {/* Right Column: Cinematic Byte Stream Animation */}
        <div className="chapter-column-visual">
          <div className="responsive-byte-stream-wrapper" style={{
            width: '100%',
            maxWidth: '600px',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center'
          }}>
            {/* Socket Pipe (Left) */}
            <div className="responsive-byte-stream-pipe" style={{
              flex: 1,
              height: '60px',
              background: 'linear-gradient(90deg, rgba(0,0,0,0) 0%, rgba(0, 217, 255, 0.05) 100%)',
              borderTop: '1px dashed rgba(0, 217, 255, 0.3)',
              borderBottom: '1px dashed rgba(0, 217, 255, 0.3)',
              position: 'relative',
              display: 'flex',
              alignItems: 'center',
              overflow: 'hidden',
              borderRight: 'none'
            }}>
              {/* Core Laser */}
              <div style={{
                position: 'absolute',
                top: '50%',
                left: 0,
                width: '100%',
                height: '2px',
                background: 'linear-gradient(90deg, transparent, var(--color-network-cyan))',
                transform: 'translateY(-50%)',
                animation: 'dataStreamGlow 2s infinite linear'
              }} />

              {/* Traveling Data Packets */}
              {['GET', '/api', 'HTTP'].map((text, i) => (
                <div key={i} style={{
                  position: 'absolute',
                  left: 0,
                  fontFamily: 'var(--font-jetbrains-mono)',
                  fontSize: '12px',
                  fontWeight: 'bold',
                  background: 'var(--color-network-cyan)',
                  color: '#000',
                  padding: '4px 8px',
                  borderRadius: '4px',
                  boxShadow: '0 0 15px var(--color-network-cyan)',
                  animation: `byteFlow 2.4s infinite cubic-bezier(0.4, 0, 0.2, 1)`,
                  animationDelay: `${i * 0.8}s`,
                  zIndex: 2
                }}>
                  {text}
                </div>
              ))}
            </div>

            {/* Buffer Box (Right) */}
            <div className="responsive-byte-stream-buffer" style={{
              width: '240px',
              height: '240px',
              background: 'rgba(10, 10, 12, 0.8)',
              border: '2px solid rgba(255, 255, 255, 0.1)',
              borderRadius: '16px',
              display: 'flex',
              flexDirection: 'column',
              alignItems: 'center',
              justifyContent: 'center',
              backdropFilter: 'blur(16px)',
              position: 'relative',
              overflow: 'hidden',
              animation: 'bufferPulse 2.4s infinite ease-in-out',
              animationDelay: '0.4s' // Offset to match packet arrival
            }}>
              {/* Falling Hex Codes inside Buffer */}
              {[0, 1, 2, 3].map(i => (
                <div key={`hex-${i}`} style={{
                  position: 'absolute',
                  top: 0,
                  left: `${20 + i * 20}%`,
                  fontFamily: 'var(--font-jetbrains-mono)',
                  fontSize: '10px',
                  color: 'rgba(0, 217, 255, 0.5)',
                  animation: 'hexFall 2s infinite linear',
                  animationDelay: `${i * 0.4}s`
                }}>
                  0x{Math.floor(Math.random() * 256).toString(16).padStart(2, '0').toUpperCase()}
                </div>
              ))}

              {/* Scroll-based fill level */}
              <div style={{
                position: 'absolute',
                bottom: 0, left: 0, width: '100%',
                height: `${Math.min(100, Math.max(0, (scrollY - 1500) * 0.15))}%`,
                background: 'linear-gradient(to top, rgba(0, 217, 255, 0.2), rgba(0, 217, 255, 0.05))',
                transition: 'height 0.1s ease',
                borderTop: '1px solid var(--color-network-cyan)',
                boxShadow: '0 -10px 20px rgba(0, 217, 255, 0.1)'
              }} />

              <div style={{
                fontFamily: 'var(--font-roboto-mono)',
                fontSize: '12px',
                color: 'var(--color-steel-mid)',
                marginBottom: '8px',
                zIndex: 2,
                background: 'rgba(0,0,0,0.5)',
                padding: '4px 8px',
                borderRadius: '4px'
              }}>
                [4096]byte
              </div>
              <div style={{
                fontFamily: 'var(--font-jetbrains-mono)',
                fontSize: '18px',
                color: 'var(--color-paper-white)',
                fontWeight: 'bold',
                zIndex: 2,
                letterSpacing: '0.1em'
              }}>
              </div>
            </div>
          </div>
        </div>

      </div>
    </section>
  );
};
