import React, { useEffect, useRef } from 'react';
import gsap from 'gsap';
import {
  Code, Zap, Shield, BarChart, BookOpen,
  ArrowRight, Globe, Layers, TerminalSquare,
  Download, FileText, GitMerge, Users, Lock, BarChart2
} from 'lucide-react';
import { GlassPanel } from './ui/GlassPanel.tsx';
import { TerminalCard } from './ui/TerminalCard.tsx';

const GoLogo = ({ width = 40, height = 15, color = "var(--color-network-cyan)" }) => (
  <svg width={width} height={height} viewBox="0 0 100 38" fill="none" xmlns="http://www.w3.org/2000/svg">
    <path d="M41.7 13.9C40.6 13.9 39.7 14.3 39 15.1C38.3 15.9 37.9 17 37.9 18.3C37.9 19.6 38.3 20.7 39 21.5C39.7 22.3 40.6 22.7 41.7 22.7C42.8 22.7 43.7 22.3 44.4 21.5C45.1 20.7 45.4 19.6 45.4 18.3C45.4 17 45.1 15.9 44.4 15.1C43.7 14.3 42.8 13.9 41.7 13.9ZM41.7 26.6C39.4 26.6 37.5 25.8 36 24.3C34.5 22.8 33.7 20.8 33.7 18.3C33.7 15.8 34.5 13.8 36 12.3C37.5 10.8 39.4 10 41.7 10C44 10 45.9 10.8 47.4 12.3C48.9 13.8 49.6 15.8 49.6 18.3C49.6 20.8 48.9 22.8 47.4 24.3C45.9 25.8 44 26.6 41.7 26.6ZM15.4 14.1C14.3 14.1 13.4 14.5 12.6 15.3C11.9 16.1 11.5 17.2 11.5 18.5C11.5 19.8 11.9 20.9 12.6 21.7C13.4 22.5 14.3 22.9 15.4 22.9C16.3 22.9 17.1 22.7 17.7 22.2V18.1H15V15H21.5V23.7C20.6 24.8 19.4 25.6 18 26.1C16.6 26.6 15.2 26.8 13.7 26.8C11 26.8 8.8 25.9 7 24.2C5.2 22.5 4.3 20.2 4.3 17.4C4.3 14.6 5.2 12.3 7 10.6C8.8 8.9 11 8 13.7 8C16.1 8 18.2 8.7 19.9 10L17.5 13C16.8 12.4 15.9 12.1 14.8 12.1C13.8 12.1 12.9 12.5 12.2 13.3C11.6 14.1 11.3 15.1 11.3 16.4C11.3 17.7 11.6 18.7 12.2 19.5C12.9 20.3 13.8 20.7 14.8 20.7C15.5 20.7 16 20.5 16.5 20.2V17.3H14.8V14.1H15.4Z" fill={color} />
  </svg>
);

const chapters = [
  { num: '01', title: 'WHY HTTP EXISTS', icon: Globe, link: '#the-request' },
  { num: '02', title: 'TCP', icon: Layers, link: '#tcp' },
  { num: '03', title: 'BUILDING A SOCKET', icon: TerminalSquare, link: '#building-a-socket' },
  { num: '04', title: 'READING BYTES', icon: Download, link: '#the-listener-loop' },
  { num: '05', title: 'PARSING REQUESTS', icon: FileText, link: '#parsing' },
  { num: '06', title: 'ROUTING', icon: GitMerge, link: '#routing' },
  { num: '07', title: 'CONCURRENCY', icon: Users, link: '#concurrency' },
  { num: '08', title: 'PRODUCTION FEATURES', icon: Lock, link: '#production' },
  { num: '09', title: 'BENCHMARKS', icon: BarChart2, link: '#benchmarks' },
  { num: '10', title: 'SOURCE CODE', icon: Code, link: '#alive' }
];

const FeatureCard = ({ icon, title, desc1, desc2 }: { icon: React.ReactNode, title: string, desc1: string, desc2: string }) => (
  <div style={{ display: 'flex', gap: '16px', alignItems: 'flex-start', flex: 1, minWidth: '180px' }}>
    <div style={{ color: 'var(--color-network-cyan)' }}>{icon}</div>
    <div>
      <div style={{ fontFamily: 'var(--font-roboto-mono)', fontSize: '11px', color: 'var(--color-paper-white)', letterSpacing: '0.05em', marginBottom: '8px' }}>{title}</div>
      <div style={{ fontFamily: 'var(--font-suisse-intl)', fontSize: '11px', color: 'var(--color-steel-mid)', lineHeight: 1.6 }}>{desc1}<br />{desc2}</div>
    </div>
  </div>
);

export const HeroSection: React.FC = () => {
  const sectionRef = useRef<HTMLElement>(null);
  const heroTextRef = useRef<HTMLHeadingElement>(null);
  const cardsRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (heroTextRef.current) {
      gsap.fromTo(heroTextRef.current,
        { y: 40, opacity: 0 },
        { y: 0, opacity: 1, duration: 1.6, ease: "power3.out", delay: 0.2 }
      );
    }
    
    if (cardsRef.current) {
      const cards = Array.from(cardsRef.current.children);
      gsap.fromTo(cards,
        { x: 30, opacity: 0 },
        { x: 0, opacity: 1, duration: 1.2, stagger: 0.15, ease: "power3.out", delay: 0.6 }
      );
      
      gsap.to(cards, {
        y: -8,
        duration: 3,
        yoyo: true,
        repeat: -1,
        ease: "sine.inOut",
        stagger: 0.4
      });
    }
  }, []);

  return (
    <>
      <style>{`
        .responsive-hero-section {
          position: relative;
          width: 100%;
          min-height: 100vh;
          padding: 120px 40px 40px 40px;
          overflow: hidden;
          display: flex;
          flex-direction: column;
          justify-content: space-between;
        }
        .responsive-hero-content {
          position: relative;
          z-index: 2;
          maxWidth: 1600px;
          margin: 0 auto;
          width: 100%;
          display: flex;
          gap: 40px;
          flex: 1;
        }
        .responsive-hero-left {
          flex: 1.5;
          display: flex;
          flex-direction: column;
          justify-content: center;
        }
        .responsive-hero-right {
          flex: 1;
          position: relative;
          display: flex;
          flex-direction: column;
          gap: 24px;
          padding-top: 80px;
          align-items: flex-end;
        }
        .responsive-hero-features {
          display: flex;
          gap: 24px;
          margin-top: 60px;
          padding-top: 32px;
          border-top: 1px solid rgba(255, 255, 255, 0.1);
        }
        .responsive-hero-journey-grid {
          display: grid;
          grid-template-columns: repeat(10, 1fr);
          gap: 8px;
          margin-bottom: 24px;
        }
        .responsive-hero-journey-footer {
          display: flex;
          gap: 24px;
          align-items: center;
          font-family: var(--font-roboto-mono);
          font-size: 10px;
          color: var(--color-steel-mid);
          letter-spacing: 0.1em;
        }

        @media (max-width: 1024px) {
          .responsive-hero-section {
            padding: 100px 24px 40px 24px !important;
            min-height: auto !important;
          }
          .responsive-hero-content {
            flex-direction: column !important;
            gap: 48px !important;
          }
          .responsive-hero-left {
            align-items: center;
            text-align: center;
          }
          .responsive-hero-left h1 {
            font-size: clamp(40px, 9vw, 84px) !important;
            line-height: 0.95 !important;
            text-align: center;
          }
          .responsive-hero-left > div:first-of-type {
            margin-bottom: 16px !important;
          }
          .responsive-hero-left > div:nth-of-type(2) {
            border-left: none !important;
            border-top: 2px solid var(--color-network-cyan);
            padding-left: 0 !important;
            padding-top: 20px;
            margin-bottom: 32px !important;
            max-width: 100% !important;
          }
          .responsive-hero-features {
            flex-wrap: wrap !important;
            gap: 24px !important;
            margin-top: 40px !important;
            justify-content: center;
          }
          .responsive-hero-right {
            align-items: center !important;
            padding-top: 0 !important;
          }
          .responsive-hero-right > div {
            max-width: 420px !important;
            width: 100% !important;
            transform: none !important;
          }
          .responsive-hero-journey-grid {
            grid-template-columns: repeat(5, 1fr) !important;
            gap: 8px !important;
          }
          .responsive-hero-journey-footer {
            flex-wrap: wrap !important;
            gap: 16px !important;
            justify-content: center;
          }
        }

        @media (max-width: 640px) {
          .responsive-hero-journey-grid {
            grid-template-columns: repeat(2, 1fr) !important;
          }
          .responsive-hero-journey-grid a {
            height: 90px !important;
          }
          .responsive-hero-journey-header {
            flex-direction: column !important;
            align-items: center !important;
            gap: 8px !important;
            text-align: center !important;
          }
          .responsive-hero-journey-header div {
            text-align: center !important;
          }
        }
      `}</style>
      <section id="hero" ref={sectionRef} className="responsive-hero-section">
        {/* Background Image */}
        <div style={{
          position: 'absolute',
          top: 0,
          left: 0,
          width: '100%',
          height: '100%',
          backgroundImage: `url(${import.meta.env.BASE_URL}images/landing_page_image.png)`,
          backgroundSize: 'cover',
          backgroundPosition: 'center',
          backgroundRepeat: 'no-repeat',
          backgroundAttachment: 'fixed',
          zIndex: 0
        }} />

        {/* Extreme Dark Overlay */}
        <div style={{
          position: 'absolute',
          top: 0,
          left: 0,
          width: '100%',
          height: '100%',
          background: 'radial-gradient(circle at 70% 50%, rgba(0,0,0,0.5) 0%, rgba(0,0,0,0.95) 80%)',
          zIndex: 1
        }} />

        {/* Content Top */}
        <div className="responsive-hero-content">
          {/* Left Column: Dense Typography & Matrix */}
          <div className="responsive-hero-left">

            <div style={{ fontFamily: 'var(--font-roboto-mono)', fontSize: '12px', letterSpacing: '0.1em', marginBottom: '24px' }}>
              <span style={{ color: 'var(--color-network-cyan)' }}>&gt;_ TITANHTTP</span>
              <span style={{ color: 'var(--color-steel-mid)', marginLeft: '16px' }}>BUILT FROM FIRST PRINCIPLES</span>
            </div>

            <h1 ref={heroTextRef} style={{
              fontFamily: 'var(--font-lambotype)',
              fontSize: 'clamp(50px, 8vw, 120px)',
              lineHeight: 0.85,
              margin: '0 0 32px 0',
              color: 'var(--color-paper-white)',
              textTransform: 'uppercase',
              letterSpacing: '-0.02em',
              textShadow: '0 0 60px rgba(0, 217, 255, 0.15)'
            }}>
              HTTP<span style={{ color: 'var(--color-network-cyan)' }}>.</span><br />
              BUILT RAW<span style={{ color: 'var(--color-network-cyan)' }}>.</span>
            </h1>

            <div style={{
              borderLeft: '2px solid var(--color-network-cyan)',
              paddingLeft: '24px',
              marginBottom: '40px',
              maxWidth: '600px'
            }}>
              <h3 style={{
                fontFamily: 'var(--font-roboto-mono)',
                fontSize: '14px',
                color: 'var(--color-paper-white)',
                letterSpacing: '0.05em',
                margin: '0 0 12px 0'
              }}>
                A PRODUCTION-INSPIRED HTTP SERVER<br />
                BUILT FROM SCRATCH IN <span style={{ color: 'var(--color-network-cyan)' }}>GO</span>.
              </h3>
              <p style={{
                fontFamily: 'var(--font-suisse-intl)',
                fontSize: '14px',
                color: 'var(--color-steel-mid)',
                lineHeight: 1.6,
                margin: 0
              }}>
                From TCP sockets to concurrent request handling — every byte, every step, built to understand, not just to ship.
              </p>
            </div>

            <div>
              <a href="#the-request" style={{ textDecoration: 'none' }}>
                <button style={{
                  background: 'var(--color-paper-white)',
                  color: '#000',
                  border: 'none',
                  padding: '14px 28px',
                  borderRadius: '4px',
                  fontFamily: 'var(--font-roboto-mono)',
                  fontSize: '12px',
                  fontWeight: 'bold',
                  letterSpacing: '0.1em',
                  cursor: 'pointer',
                  display: 'inline-flex',
                  alignItems: 'center',
                  gap: '12px',
                  boxShadow: '0 0 20px rgba(255, 255, 255, 0.15)',
                  transition: 'transform 0.2s ease'
                }}
                  onMouseEnter={(e) => e.currentTarget.style.transform = 'translateY(-2px)'}
                  onMouseLeave={(e) => e.currentTarget.style.transform = 'translateY(0)'}
                >
                  EXPLORE THE JOURNEY <ArrowRight size={16} />
                </button>
              </a>
            </div>

            {/* Feature Matrix */}
            <div className="responsive-hero-features">
              <FeatureCard icon={<Code size={24} />} title="100% FROM SCRATCH" desc1="No frameworks." desc2="No shortcuts." />
              <FeatureCard icon={<Zap size={24} />} title="FAST" desc1="Concurrent by design." desc2="Goroutines + non-blocking I/O." />
              <FeatureCard icon={<Shield size={24} />} title="RELIABLE" desc1="HTTP/1.1 compliant." desc2="Persistent connections." />
              <FeatureCard icon={<BarChart size={24} />} title="MEASURED" desc1="Benchmarked against" desc2="industry standard servers." />
              <FeatureCard icon={<BookOpen size={24} />} title="EDUCATIONAL" desc1="Deep documentation." desc2="Every concept visually explained." />
            </div>
          </div>

          {/* Right Column: Widgets */}
          <div ref={cardsRef} className="responsive-hero-right">
            
            <TerminalCard title="REQUEST" glowColor="cyan" style={{ maxWidth: '320px' }}>
              <span style={{ color: 'var(--color-paper-white)' }}>GET</span> / HTTP/1.1{'\n'}
              Host: localhost{'\n'}
              User-Agent: TitanHTTP{'\n\n'}
              <span style={{ color: 'var(--color-network-cyan)' }}>234 BYTES</span>
            </TerminalCard>

            <TerminalCard title="RESPONSE" glowColor="violet" style={{ maxWidth: '320px', transform: 'translateX(-40px)' }}>
              HTTP/1.1 <span style={{ color: 'var(--color-paper-white)' }}>200 OK</span>{'\n'}
              Content-Type: text/html{'\n'}
              Content-Length: 1024{'\n\n'}
              <span style={{ color: 'var(--color-packet-violet)' }}>1024 BYTES</span>
            </TerminalCard>

            <GlassPanel glowColor="none" style={{ maxWidth: '400px', marginTop: 'auto', padding: '24px' }}>
              <div style={{ fontFamily: 'var(--font-roboto-mono)', fontSize: '10px', color: 'var(--color-paper-white)', letterSpacing: '0.1em', marginBottom: '16px' }}>MAX WORKER POOL</div>
              <div style={{ display: 'flex', gap: '16px', alignItems: 'center', justifyContent: 'space-between', width: '100%' }}>
                <div style={{ display: 'flex', flexDirection: 'column' }}>
                  <span style={{ fontFamily: 'var(--font-lambotype)', fontSize: '32px', color: 'var(--color-paper-white)', lineHeight: 1 }}>10,000</span>
                  <span style={{ fontFamily: 'var(--font-roboto-mono)', fontSize: '9px', color: 'var(--color-network-cyan)', marginTop: '6px', letterSpacing: '0.05em' }}>BOUNDED GOROUTINES</span>
                </div>
                <div style={{ display: 'flex', gap: '4px' }}>
                  {[...Array(16)].map((_, i) => (
                    <div key={i} style={{
                      width: '6px',
                      height: '20px',
                      background: i < 13 ? 'var(--color-network-cyan)' : 'rgba(255, 255, 255, 0.1)',
                      boxShadow: i < 13 ? '0 0 10px rgba(0, 217, 255, 0.4)' : 'none',
                      opacity: 0.8
                    }} />
                  ))}
                </div>
              </div>
            </GlassPanel>
          </div>
        </div>

        {/* The Journey Navigator & Tech Stack */}
        <div style={{
          position: 'relative',
          zIndex: 2,
          maxWidth: '1600px',
          margin: '0 auto',
          width: '100%',
          marginTop: '60px'
        }}>
          <div className="responsive-hero-journey-header" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-end', marginBottom: '16px' }}>
            <div style={{ fontFamily: 'var(--font-roboto-mono)', fontSize: '14px', letterSpacing: '0.1em', color: 'var(--color-paper-white)' }}>
              THE JOURNEY <span style={{ color: 'var(--color-steel-mid)', fontSize: '10px', marginLeft: '16px' }}>10 CHAPTERS. ONE PROTOCOL. ENDLESS LEARNING.</span>
            </div>
            <div style={{ fontFamily: 'var(--font-roboto-mono)', fontSize: '10px', color: 'var(--color-steel-mid)', letterSpacing: '0.1em' }}>
              SCROLL TO EXPLORE ↓
            </div>
          </div>

          <div className="responsive-hero-journey-grid">
            {chapters.map((ch, idx) => {
              const Icon = ch.icon;
              return (
                <a key={idx} href={ch.link} style={{
                  textDecoration: 'none',
                  border: '1px solid rgba(255, 255, 255, 0.1)',
                  background: 'rgba(20, 20, 22, 0.5)',
                  padding: '16px 12px',
                  display: 'flex',
                  flexDirection: 'column',
                  alignItems: 'center',
                  textAlign: 'center',
                  gap: '12px',
                  height: '100px',
                  justifyContent: 'center',
                  transition: 'all 0.2s ease',
                  cursor: 'pointer'
                }}
                  onMouseEnter={(e) => {
                    e.currentTarget.style.background = 'rgba(0, 217, 255, 0.05)';
                    e.currentTarget.style.borderColor = 'rgba(0, 217, 255, 0.3)';
                  }}
                  onMouseLeave={(e) => {
                    e.currentTarget.style.background = 'rgba(20, 20, 22, 0.5)';
                    e.currentTarget.style.borderColor = 'rgba(255, 255, 255, 0.1)';
                  }}
                >
                  <div style={{ alignSelf: 'flex-start', fontFamily: 'var(--font-roboto-mono)', fontSize: '10px', color: 'var(--color-steel-mid)' }}>{ch.num}</div>
                  <Icon size={20} color={idx === 0 || idx === 9 ? "var(--color-network-cyan)" : "var(--color-paper-white)"} />
                  <div style={{ fontFamily: 'var(--font-roboto-mono)', fontSize: '9px', letterSpacing: '0.05em', color: 'var(--color-steel-mid)' }}>{ch.title}</div>
                </a>
              );
            })}
          </div>

          <div className="responsive-hero-journey-footer">
            <div>BUILT WITH</div>
            <GoLogo width={60} height={48} color="var(--color-network-cyan)" />
            <div style={{ flex: 1, textOverflow: 'ellipsis', overflow: 'hidden', whiteSpace: 'nowrap' }}>| NET | TCP | GOROUTINES | EPOLL (IOCP ON WINDOWS) | ZERO DEPENDENCIES</div>
          </div>
        </div>
      </section>
    </>
  );
};
