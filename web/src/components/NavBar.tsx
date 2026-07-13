import React, { useEffect, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { Server, ChevronDown, Menu, X } from 'lucide-react';

const GithubIcon = ({ size = 16, color = "currentColor" }) => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    width={size}
    height={size}
    viewBox="0 0 24 24"
    fill="none"
    stroke={color}
    strokeWidth="2"
    strokeLinecap="round"
    strokeLinejoin="round"
  >
    <path d="M15 22v-4a4.8 4.8 0 0 0-1-3.2c3-.3 6-1.5 6-6.5a4.6 4.6 0 0 0-1.3-3.2 4.2 4.2 0 0 0-.1-3.2s-1.1-.3-3.5 1.3a12.3 12.3 0 0 0-6.2 0C6.5 2.8 5.4 3.1 5.4 3.1a4.2 4.2 0 0 0-.1 3.2A4.6 4.6 0 0 0 4 9.5c0 5 3 6.2 6 6.5a4.8 4.8 0 0 0-1 3.2v4" />
  </svg>
);

const journeyChapters = [
  { name: '01 WHY HTTP EXISTS', link: '/#the-request' },
  { name: '02 TCP', link: '/#tcp' },
  { name: '03 BUILDING A SOCKET', link: '/#building-a-socket' },
  { name: '04 READING BYTES', link: '/#the-listener-loop' },
  { name: '05 PARSING REQUESTS', link: '/#parsing' },
  { name: '06 ROUTING', link: '/#routing' },
  { name: '07 CONCURRENCY', link: '/#concurrency' },
  { name: '08 PRODUCTION FEATURES', link: '/#production' },
  { name: '09 BENCHMARKS', link: '/#benchmarks' },
  { name: '10 SOURCE CODE', link: '/#alive' }
];

export const NavBar: React.FC = () => {
  const [scrolled, setScrolled] = useState(false);
  const [journeyOpen, setJourneyOpen] = useState(false);
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
  const [mobileJourneyOpen, setMobileJourneyOpen] = useState(false);
  const navigate = useNavigate();

  useEffect(() => {
    const handleScroll = () => {
      setScrolled(window.scrollY > 50);
    };
    window.addEventListener('scroll', handleScroll);
    return () => window.removeEventListener('scroll', handleScroll);
  }, []);

  // Prevent scrolling when mobile menu is open
  useEffect(() => {
    if (mobileMenuOpen) {
      document.body.style.overflow = 'hidden';
    } else {
      document.body.style.overflow = '';
    }
    return () => {
      document.body.style.overflow = '';
    };
  }, [mobileMenuOpen]);

  return (
    <>
      <style>{`
        .premium-nav-link {
          position: relative;
          color: var(--color-paper-white);
          text-decoration: none;
          font-family: var(--font-roboto-mono);
          font-size: 11px;
          letter-spacing: 0.1em;
          text-transform: uppercase;
          opacity: 0.7;
          transition: all 0.3s ease;
        }
        .premium-nav-link:hover {
          opacity: 1;
          color: var(--color-network-cyan);
          text-shadow: 0 0 10px rgba(0, 217, 255, 0.5);
        }
        .premium-nav-link::after {
          content: '';
          position: absolute;
          width: 0;
          height: 2px;
          bottom: -6px;
          left: 50%;
          background: var(--color-network-cyan);
          box-shadow: 0 0 10px var(--color-network-cyan);
          transition: all 0.3s ease;
          transform: translateX(-50%);
          border-radius: 2px;
        }
        .premium-nav-link:hover::after {
          width: 100%;
        }
        
        .dropdown-link {
          color: var(--color-steel-mid);
          text-decoration: none;
          padding: 10px 20px;
          font-family: var(--font-roboto-mono);
          font-size: 10px;
          letter-spacing: 0.1em;
          transition: all 0.2s ease;
          display: block;
        }
        .dropdown-link:hover {
          background: rgba(0, 217, 255, 0.05);
          color: var(--color-network-cyan);
          padding-left: 24px;
        }

        /* Mobile specific styles */
        .mobile-toggle {
          display: none;
        }

        .mobile-drawer {
          position: fixed;
          top: 0;
          left: 0;
          width: 100vw;
          height: 100vh;
          background: rgba(10, 10, 12, 0.98);
          backdrop-filter: blur(24px);
          -webkit-backdrop-filter: blur(24px);
          z-index: 95;
          display: flex;
          flex-direction: column;
          padding: 100px 24px 60px 24px;
          transform: translateY(-100%);
          transition: transform 0.4s cubic-bezier(0.4, 0, 0.2, 1);
          overflow-y: scroll;
          -webkit-overflow-scrolling: touch;
        }

        .mobile-drawer.open {
          transform: translateY(0);
        }

        .mobile-drawer-link {
          font-family: var(--font-roboto-mono);
          color: var(--color-paper-white);
          font-size: 16px;
          letter-spacing: 0.1em;
          text-transform: uppercase;
          text-decoration: none;
          padding: 16px 0;
          border-bottom: 1px solid rgba(255, 255, 255, 0.05);
          transition: color 0.3s ease;
        }

        .mobile-drawer-link:hover {
          color: var(--color-network-cyan);
        }

        .mobile-journey-dropdown {
          display: flex;
          flex-direction: column;
          max-height: 0;
          overflow: hidden;
          transition: max-height 0.4s ease-out;
          background: rgba(255,255,255,0.02);
          border-radius: 4px;
        }

        .mobile-journey-dropdown.open {
          max-height: 600px;
          padding: 8px 0;
          margin-top: 8px;
          overflow-y: auto;
        }

        .mobile-dropdown-link {
          color: var(--color-steel-mid);
          text-decoration: none;
          padding: 12px 16px;
          font-family: var(--font-roboto-mono);
          font-size: 11px;
          letter-spacing: 0.1em;
          display: block;
          border-bottom: 1px solid rgba(255, 255, 255, 0.02);
        }

        .mobile-dropdown-link:hover {
          color: var(--color-network-cyan);
          background: rgba(0, 217, 255, 0.05);
        }

        @media (max-width: 1024px) {
          .nav-desktop-links {
            display: none !important;
          }
          .nav-desktop-cta {
            display: none !important;
          }
          .mobile-toggle {
            display: flex !important;
            align-items: center;
            justify-content: center;
            background: transparent;
            border: none;
            color: var(--color-paper-white);
            cursor: pointer;
            z-index: 101;
            padding: 8px;
            border-radius: 4px;
            transition: background 0.3s;
          }
          .mobile-toggle:hover {
            background: rgba(255, 255, 255, 0.05);
          }
        }
      `}</style>

      <nav className="nav-container" style={{
        position: 'fixed',
        top: '0',
        left: '0',
        width: '100%',
        zIndex: 100,
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'space-between',
        padding: scrolled ? '16px 24px' : '24px 24px',
        background: (scrolled || mobileMenuOpen) ? 'rgba(10, 10, 12, 0.95)' : 'transparent',
        backdropFilter: (scrolled || mobileMenuOpen) ? 'blur(24px) saturate(150%)' : 'none',
        WebkitBackdropFilter: (scrolled || mobileMenuOpen) ? 'blur(24px) saturate(150%)' : 'none',
        borderBottom: (scrolled || mobileMenuOpen) ? '1px solid rgba(255, 255, 255, 0.05)' : '1px solid transparent',
        transition: 'all 0.5s cubic-bezier(0.4, 0, 0.2, 1)'
      }}>
        {/* Left: Brand */}
        <div 
          style={{ display: 'flex', alignItems: 'center', gap: '16px', cursor: 'pointer', zIndex: 101 }} 
          onClick={() => {
            setMobileMenuOpen(false);
            navigate('/');
            window.scrollTo({top: 0, behavior: 'smooth'});
          }}
        >
          <div style={{
            width: '36px',
            height: '36px',
            borderRadius: '8px',
            background: 'rgba(0, 217, 255, 0.05)',
            border: '1px solid rgba(0, 217, 255, 0.3)',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            boxShadow: '0 0 15px rgba(0, 217, 255, 0.15)',
            transition: 'all 0.3s ease'
          }}>
            <Server size={18} color="var(--color-network-cyan)" />
          </div>
          <span style={{ 
            fontFamily: 'var(--font-lambotype)', 
            fontSize: '20px', 
            letterSpacing: '2px', 
            textTransform: 'uppercase',
            color: 'var(--color-paper-white)',
            textShadow: scrolled ? '0 2px 10px rgba(0,0,0,0.5)' : 'none'
          }}>
            TitanHTTP
          </span>
        </div>

        {/* Center: Links (Desktop) */}
        <div className="nav-desktop-links" style={{ 
          display: 'flex', 
          gap: '32px', 
          alignItems: 'center',
          background: scrolled ? 'transparent' : 'rgba(255, 255, 255, 0.03)',
          padding: scrolled ? '0' : '12px 32px',
          borderRadius: '100px',
          border: scrolled ? 'none' : '1px solid rgba(255, 255, 255, 0.05)',
          transition: 'all 0.5s ease'
        }}>
          <Link to="/why" className="premium-nav-link">WHY</Link>
          <Link to="/architecture" className="premium-nav-link">ARCHITECTURE</Link>
          <Link to="/decisions" className="premium-nav-link">DECISIONS</Link>
          
          <div 
            style={{ position: 'relative' }} 
            onMouseEnter={() => setJourneyOpen(true)} 
            onMouseLeave={() => setJourneyOpen(false)}
          >
            <span className="premium-nav-link" style={{ display: 'flex', alignItems: 'center', gap: '4px', cursor: 'pointer' }}>
              JOURNEY <ChevronDown size={12} style={{ transform: journeyOpen ? 'rotate(180deg)' : 'none', transition: 'transform 0.2s ease' }} />
            </span>
            
            {/* Dropdown Menu */}
            <div style={{
              position: 'absolute',
              top: '100%',
              left: '50%',
              transform: `translateX(-50%) translateY(${journeyOpen ? '0' : '10px'})`,
              opacity: journeyOpen ? 1 : 0,
              visibility: journeyOpen ? 'visible' : 'hidden',
              background: 'rgba(10, 10, 12, 0.95)',
              backdropFilter: 'blur(24px)',
              border: '1px solid rgba(255, 255, 255, 0.1)',
              borderRadius: '8px',
              padding: '12px 0',
              marginTop: '24px',
              width: '240px',
              boxShadow: '0 20px 40px rgba(0,0,0,0.5)',
              transition: 'all 0.3s cubic-bezier(0.4, 0, 0.2, 1)',
              pointerEvents: journeyOpen ? 'auto' : 'none'
            }}>
              {/* Invisible bridge to keep hover active between link and dropdown */}
              <div style={{ position: 'absolute', top: '-24px', left: 0, width: '100%', height: '24px' }} />
              
              <div style={{ padding: '0 20px', marginBottom: '8px', fontSize: '9px', fontFamily: 'var(--font-roboto-mono)', color: 'var(--color-network-cyan)', letterSpacing: '0.1em' }}>
                SELECT CHAPTER
              </div>
              
              {journeyChapters.map(ch => (
                <Link key={ch.name} to={ch.link} className="dropdown-link" onClick={() => setJourneyOpen(false)}>
                  {ch.name}
                </Link>
              ))}
            </div>
          </div>
        </div>

        {/* Right: CTA (Desktop) */}
        <div className="nav-desktop-cta">
          <a href="/#alive" style={{ textDecoration: 'none' }}>
            <button style={{
              background: 'var(--color-paper-white)',
              color: '#000',
              border: 'none',
              padding: '12px 28px',
              borderRadius: '100px',
              fontFamily: 'var(--font-roboto-mono)',
              fontSize: '11px',
              fontWeight: 'bold',
              letterSpacing: '0.1em',
              cursor: 'pointer',
              display: 'flex',
              alignItems: 'center',
              gap: '8px',
              boxShadow: '0 4px 14px rgba(255, 255, 255, 0.25)',
              transition: 'all 0.3s cubic-bezier(0.4, 0, 0.2, 1)'
            }}
            onMouseEnter={(e) => {
              e.currentTarget.style.transform = 'translateY(-2px)';
              e.currentTarget.style.boxShadow = '0 6px 20px rgba(255, 255, 255, 0.4)';
            }}
            onMouseLeave={(e) => {
              e.currentTarget.style.transform = 'translateY(0)';
              e.currentTarget.style.boxShadow = '0 4px 14px rgba(255, 255, 255, 0.25)';
            }}
            >
              <GithubIcon size={16} /> SOURCE CODE
            </button>
          </a>
        </div>

        {/* Mobile Hamburger Toggle Button */}
        <button 
          className="mobile-toggle"
          onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
          aria-label="Toggle navigation menu"
        >
          {mobileMenuOpen ? <X size={24} /> : <Menu size={24} />}
        </button>
      </nav>

      {/* Mobile Navigation Drawer */}
      <div className={`mobile-drawer ${mobileMenuOpen ? 'open' : ''}`}>
        <Link to="/why" className="mobile-drawer-link" onClick={() => setMobileMenuOpen(false)}>WHY</Link>
        <Link to="/architecture" className="mobile-drawer-link" onClick={() => setMobileMenuOpen(false)}>ARCHITECTURE</Link>
        <Link to="/decisions" className="mobile-drawer-link" onClick={() => setMobileMenuOpen(false)}>DECISIONS</Link>
        
        {/* Mobile Journey Dropdown */}
        <div style={{ display: 'flex', flexDirection: 'column', borderBottom: '1px solid rgba(255, 255, 255, 0.05)' }}>
          <button 
            className="mobile-drawer-link" 
            style={{ 
              background: 'transparent', 
              border: 'none', 
              textAlign: 'left', 
              width: '100%', 
              display: 'flex', 
              alignItems: 'center', 
              justifyContent: 'space-between',
              cursor: 'pointer',
              outline: 'none'
            }}
            onClick={() => setMobileJourneyOpen(!mobileJourneyOpen)}
          >
            <span>JOURNEY</span>
            <ChevronDown size={18} style={{ transform: mobileJourneyOpen ? 'rotate(180deg)' : 'none', transition: 'transform 0.2s ease' }} />
          </button>
          
          <div className={`mobile-journey-dropdown ${mobileJourneyOpen ? 'open' : ''}`}>
            {journeyChapters.map(ch => (
              <Link 
                key={ch.name} 
                to={ch.link} 
                className="mobile-dropdown-link" 
                onClick={() => {
                  setMobileMenuOpen(false);
                  setMobileJourneyOpen(false);
                }}
              >
                {ch.name}
              </Link>
            ))}
          </div>
        </div>

        {/* CTA (Mobile) */}
        <div style={{ marginTop: 'auto', paddingTop: '40px', paddingBottom: '40px' }}>
          <a href="/#alive" style={{ textDecoration: 'none' }} onClick={() => setMobileMenuOpen(false)}>
            <button style={{
              background: 'var(--color-paper-white)',
              color: '#000',
              border: 'none',
              padding: '16px 28px',
              borderRadius: '100px',
              fontFamily: 'var(--font-roboto-mono)',
              fontSize: '12px',
              fontWeight: 'bold',
              letterSpacing: '0.15em',
              cursor: 'pointer',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              gap: '12px',
              width: '100%',
              boxShadow: '0 4px 14px rgba(255, 255, 255, 0.25)',
            }}>
              <GithubIcon size={18} /> SOURCE CODE
            </button>
          </a>
        </div>
      </div>
    </>
  );
};
