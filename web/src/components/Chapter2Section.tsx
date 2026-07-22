import React, { useEffect, useRef, useState } from 'react';
import { Monitor, Server } from 'lucide-react';
import { GlassPanel } from './ui/GlassPanel.tsx';

export const Chapter2Section: React.FC = () => {
  const sectionRef = useRef<HTMLElement>(null);

  const [isInView, setIsInView] = useState(false);
  const [progress, setProgress] = useState(0);

  useEffect(() => {
    const handleScroll = () => {
      if (sectionRef.current) {
        const rect = sectionRef.current.getBoundingClientRect();
        const windowHeight = window.innerHeight;
        // Calculate progress from 0 to 1 as the section scrolls through the viewport
        if (rect.top < windowHeight && rect.bottom > 0) {
          const totalDistance = windowHeight + rect.height;
          const currentDistance = windowHeight - rect.top;
          setProgress(Math.max(0, Math.min(1, currentDistance / totalDistance)));
        }
      }
    };
    window.addEventListener('scroll', handleScroll);
    return () => window.removeEventListener('scroll', handleScroll);
  }, []);

  useEffect(() => {
    const observer = new IntersectionObserver(
      ([entry]) => {
        setIsInView(entry.isIntersecting);
      },
      { threshold: 0.1 }
    );
    if (sectionRef.current) observer.observe(sectionRef.current);
    return () => observer.disconnect();
  }, []);

  // Map scroll progress to packet movement (0 to 1 for each phase)
  // With minHeight: 400vh, the sticky phase is between progress 0.2 and 0.8
  // Phase 1 (SYN): 0.25 to 0.40
  const synProgress = Math.min(1, Math.max(0, (progress - 0.25) / 0.15));
  // Phase 2 (SYN-ACK): 0.45 to 0.60
  const synAckProgress = Math.min(1, Math.max(0, (progress - 0.45) / 0.15));
  // Phase 3 (ACK): 0.65 to 0.80
  const ackProgress = Math.min(1, Math.max(0, (progress - 0.65) / 0.15));

  const clientGlow = synAckProgress >= 1 ? '0 0 30px rgba(0, 217, 255, 0.4)' : 'none';
  const serverGlow = synProgress >= 1 ? '0 0 30px rgba(0, 217, 255, 0.4)' : 'none';
  const clientColor = synAckProgress >= 1 ? 'var(--color-network-cyan)' : 'var(--color-steel-mid)';
  const serverColor = synProgress >= 1 ? 'var(--color-network-cyan)' : 'var(--color-steel-mid)';

  return (
    <section id="tcp" ref={sectionRef} style={{
      position: 'relative',
      width: '100%',
      minHeight: '400vh' // Extra height to allow scrolling the animation while sticky
    }}>
      {/* Sticky Wrapper */}
      <div style={{
        position: 'sticky',
        top: 0,
        width: '100%',
        height: '100vh',
        display: 'flex',
        flexDirection: 'column',
        justifyContent: 'center',
        paddingTop: '80px',
        overflow: 'hidden'
      }}>
        {/* Background Image Parallax */}
        <div style={{
          position: 'absolute',
          top: 0,
          left: 0,
          width: '100%',
          height: '100%',
          backgroundImage: `url(${import.meta.env.BASE_URL}images/chapter2.png)`,
          backgroundSize: 'cover',
          backgroundPosition: 'center',
          zIndex: 0
        }} />

        <div style={{
          position: 'absolute',
          top: 0,
          left: 0,
          width: '100%',
          height: '100%',
          background: 'rgba(0,0,0,0.6)', // Reduced overlay to show more background
          zIndex: 1
        }} />

        {/* Content Container */}
        <div style={{
          position: 'relative',
          zIndex: 2,
          maxWidth: '1440px',
          margin: '0 auto',
          width: '100%',
          padding: '0 24px',
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
          opacity: isInView ? 1 : 0,
          transform: `translateY(${isInView ? 0 : '40px'})`,
          transition: 'opacity 1s ease, transform 1s ease'
        }}>
        
        <div style={{ textAlign: 'center', marginBottom: '40px' }}>
          <div style={{
            fontFamily: 'var(--font-roboto-mono)',
            fontSize: '14px',
            color: 'var(--color-network-cyan)',
            letterSpacing: '0.1em',
            marginBottom: '16px'
          }}>CHAPTER 02</div>
          
          <h2 className="chapter-heading">
            TCP
          </h2>
          <p style={{
            fontFamily: 'var(--font-roboto-mono)',
            fontSize: '14px',
            color: 'var(--color-steel-mid)',
            maxWidth: '600px',
            margin: '0 auto',
            lineHeight: 1.8
          }}>
            Before HTTP can send a single byte of data, a reliable connection must be established. 
            TCP guarantees delivery through a precise, 3-step dance between the client and the server.
          </p>
        </div>

        {/* CSS/SVG Diagram for TCP Handshake */}
        <GlassPanel className="tcp-diagram-panel">
          
          {/* Client Icon */}
          <div className="tcp-node-box" style={{
            color: clientColor,
            textShadow: clientGlow
          }}>
            <div className="tcp-node-icon-wrapper" style={{
              border: `1px solid ${synAckProgress >= 1 ? 'rgba(0, 217, 255, 0.5)' : 'rgba(255,255,255,0.1)'}`,
              boxShadow: clientGlow,
            }}>
              <Monitor size={48} strokeWidth={1.5} />
            </div>
            <div style={{ fontFamily: 'var(--font-roboto-mono)', fontSize: '14px', fontWeight: 'bold', letterSpacing: '0.1em' }}>CLIENT</div>
          </div>

          {/* Connection Tracks */}
          <div className="tcp-tracks-container">
            
            {/* SYN Track */}
            <div style={{ position: 'relative', height: '2px', background: 'rgba(255,255,255,0.05)' }}>
              {/* Path highlight behind packet */}
              <div style={{
                position: 'absolute',
                top: 0, left: 0, height: '100%',
                width: `${synProgress * 100}%`,
                background: 'linear-gradient(90deg, transparent, rgba(0, 217, 255, 0.5))'
              }}/>
              {/* Moving Packet */}
              <div style={{
                position: 'absolute',
                top: '50%',
                left: `${synProgress * 100}%`,
                transform: 'translate(-50%, -50%)',
                width: '12px',
                height: '12px',
                background: 'var(--color-network-cyan)',
                borderRadius: '50%',
                boxShadow: '0 0 15px var(--color-network-cyan), 0 0 30px var(--color-network-cyan)',
                opacity: synProgress > 0 && synProgress < 1 ? 1 : (synProgress === 1 ? 0 : 0),
                transition: 'opacity 0.2s'
              }} />
              <div style={{
                position: 'absolute',
                top: '-24px',
                left: '50%',
                transform: 'translateX(-50%)',
                fontFamily: 'var(--font-roboto-mono)',
                fontSize: '12px',
                color: synProgress > 0 ? 'var(--color-network-cyan)' : 'rgba(255,255,255,0.2)',
                textShadow: synProgress > 0 ? '0 0 10px rgba(0, 217, 255, 0.5)' : 'none',
                transition: 'color 0.3s'
              }}>SYN</div>
            </div>

            {/* SYN-ACK Track */}
            <div style={{ position: 'relative', height: '2px', background: 'rgba(255,255,255,0.05)' }}>
              {/* Path highlight behind packet */}
              <div style={{
                position: 'absolute',
                top: 0, right: 0, height: '100%',
                width: `${synAckProgress * 100}%`,
                background: 'linear-gradient(-90deg, transparent, rgba(0, 217, 255, 0.5))'
              }}/>
              {/* Moving Packet (Right to left) */}
              <div style={{
                position: 'absolute',
                top: '50%',
                right: `${synAckProgress * 100}%`,
                transform: 'translate(50%, -50%)',
                width: '12px',
                height: '12px',
                background: 'var(--color-network-cyan)',
                borderRadius: '50%',
                boxShadow: '0 0 15px var(--color-network-cyan), 0 0 30px var(--color-network-cyan)',
                opacity: synAckProgress > 0 && synAckProgress < 1 ? 1 : (synAckProgress === 1 ? 0 : 0),
                transition: 'opacity 0.2s'
              }} />
              <div style={{
                position: 'absolute',
                top: '-24px',
                left: '50%',
                transform: 'translateX(-50%)',
                fontFamily: 'var(--font-roboto-mono)',
                fontSize: '12px',
                color: synAckProgress > 0 ? 'var(--color-network-cyan)' : 'rgba(255,255,255,0.2)',
                textShadow: synAckProgress > 0 ? '0 0 10px rgba(0, 217, 255, 0.5)' : 'none',
                transition: 'color 0.3s'
              }}>SYN-ACK</div>
            </div>

            {/* ACK Track */}
            <div style={{ position: 'relative', height: '2px', background: 'rgba(255,255,255,0.05)' }}>
              {/* Path highlight behind packet */}
              <div style={{
                position: 'absolute',
                top: 0, left: 0, height: '100%',
                width: `${ackProgress * 100}%`,
                background: 'linear-gradient(90deg, transparent, rgba(0, 217, 255, 0.5))'
              }}/>
              {/* Moving Packet */}
              <div style={{
                position: 'absolute',
                top: '50%',
                left: `${ackProgress * 100}%`,
                transform: 'translate(-50%, -50%)',
                width: '12px',
                height: '12px',
                background: 'var(--color-network-cyan)',
                borderRadius: '50%',
                boxShadow: '0 0 15px var(--color-network-cyan), 0 0 30px var(--color-network-cyan)',
                opacity: ackProgress > 0 && ackProgress < 1 ? 1 : (ackProgress === 1 ? 0 : 0),
                transition: 'opacity 0.2s'
              }} />
              <div style={{
                position: 'absolute',
                top: '-24px',
                left: '50%',
                transform: 'translateX(-50%)',
                fontFamily: 'var(--font-roboto-mono)',
                fontSize: '12px',
                color: ackProgress > 0 ? 'var(--color-network-cyan)' : 'rgba(255,255,255,0.2)',
                textShadow: ackProgress > 0 ? '0 0 10px rgba(0, 217, 255, 0.5)' : 'none',
                transition: 'color 0.3s'
              }}>ACK</div>
            </div>
            
            {/* Connecting Vertical Line indicating established connection */}
            <div style={{
              position: 'absolute',
              top: '50%',
              left: '50%',
              transform: 'translate(-50%, -50%)',
              height: '100%',
              width: '1px',
              background: 'var(--color-network-cyan)',
              opacity: ackProgress >= 1 ? 0.3 : 0,
              boxShadow: '0 0 20px var(--color-network-cyan)',
              transition: 'opacity 1s ease'
            }} />

          </div>

          {/* Server Icon */}
          <div className="tcp-node-box" style={{
            color: serverColor,
            textShadow: serverGlow
          }}>
            <div className="tcp-node-icon-wrapper" style={{
              border: `1px solid ${synProgress >= 1 ? 'rgba(0, 217, 255, 0.5)' : 'rgba(255,255,255,0.1)'}`,
              boxShadow: serverGlow,
            }}>
              <Server size={48} strokeWidth={1.5} />
            </div>
            <div style={{ fontFamily: 'var(--font-roboto-mono)', fontSize: '14px', fontWeight: 'bold', letterSpacing: '0.1em' }}>SERVER</div>
          </div>

        </GlassPanel>

      </div>
      </div>
    </section>
  );
};

