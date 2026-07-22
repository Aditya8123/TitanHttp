import React, { useEffect, useRef, useState } from 'react';
import { GitMerge } from 'lucide-react';
import chapter6Img from '../assets/images/chapter6.png';

const ROUTES = [
  { method: 'GET', path: '/api/users', handler: 'handleUsers()' },
  { method: 'POST', path: '/api/posts', handler: 'handlePosts()' },
  { method: 'GET', path: '/', handler: 'handleIndex()' }
];

export const Chapter6Section: React.FC = () => {
  const sectionRef = useRef<HTMLElement>(null);
  const [isInView, setIsInView] = useState(false);
  const [activeIndex, setActiveIndex] = useState(0);

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

  useEffect(() => {
    if (!isInView) return;
    const interval = setInterval(() => {
      setActiveIndex(prev => (prev + 1) % ROUTES.length);
    }, 4000);
    return () => clearInterval(interval);
  }, [isInView]);

  return (
    <section id="routing" ref={sectionRef} style={{
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
        backgroundImage: `url(${chapter6Img})`,
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
          }}>CHAPTER 06</div>
          
          <h2 className="chapter-heading">
            ROUTING
          </h2>

          <p style={{
            fontFamily: 'var(--font-roboto-mono)',
            fontSize: '14px',
            color: 'var(--color-steel-mid)',
            lineHeight: 1.8,
            marginBottom: '40px'
          }}>
            Once a request is parsed, the server must decide what to do with it. 
            A router is a specialized switchboard, matching the HTTP Method and Path to specific 
            functions known as Handlers.
            <br/><br/>
            In TitanHTTP, routing is deterministic and fast, mapping <code style={{ color: 'var(--color-paper-white)' }}>GET /users</code> directly to the logic that fetches user data.
          </p>
        </div>

        {/* Left Column: Routing Visualization */}
        <div className="chapter-column-visual">
          <div key={activeIndex} style={{
            background: 'rgba(10, 10, 12, 0.8)',
            border: '1px solid rgba(255, 255, 255, 0.1)',
            borderRadius: '16px',
            padding: '40px',
            width: '100%',
            maxWidth: '500px',
            boxShadow: '0 20px 40px rgba(0,0,0,0.5)',
            backdropFilter: 'blur(16px)',
            '--scan-height': `${48 + activeIndex * 71}px`
          } as React.CSSProperties}>
            <style>{`
              @keyframes reqDrop {
                0% { transform: translateY(-20px); opacity: 0; }
                10%, 80% { transform: translateY(0); opacity: 1; }
                90%, 100% { transform: translateY(10px); opacity: 0; }
              }
              @keyframes trunkScan {
                0%, 10% { height: 0%; opacity: 0; }
                15%, 80% { height: var(--scan-height); opacity: 1; }
                90%, 100% { height: var(--scan-height); opacity: 0; }
              }
              @keyframes branchFill {
                0%, 15% { width: 0%; opacity: 0; }
                20%, 80% { width: 100%; opacity: 1; }
                90%, 100% { width: 100%; opacity: 0; }
              }
              @keyframes handlerPulse {
                0%, 20% { 
                  background: rgba(255,255,255,0.02); 
                  border-color: rgba(255,255,255,0.1); 
                  color: var(--color-steel-mid); 
                  box-shadow: none; 
                  transform: translateX(0); 
                }
                25%, 75% { 
                  background: rgba(0,217,255,0.1); 
                  border-color: var(--color-network-cyan); 
                  color: var(--color-paper-white); 
                  box-shadow: inset 0 0 20px rgba(0,217,255,0.1), 0 0 15px rgba(0,217,255,0.2); 
                  transform: translateX(10px); 
                }
                85%, 100% { 
                  background: rgba(255,255,255,0.02); 
                  border-color: rgba(255,255,255,0.1); 
                  color: var(--color-steel-mid); 
                  box-shadow: none; 
                  transform: translateX(0); 
                }
              }
            `}</style>

            {/* Request Badge */}
            <div style={{ 
              display: 'flex', alignItems: 'center', gap: '16px', marginBottom: '32px',
              animation: 'reqDrop 4s infinite ease-out'
            }}>
               <GitMerge color="var(--color-network-cyan)" size={24} />
               <div style={{
                  background: 'rgba(0, 217, 255, 0.1)',
                  border: '1px solid var(--color-network-cyan)',
                  padding: '8px 16px',
                  borderRadius: '6px',
                  fontFamily: 'var(--font-jetbrains-mono)',
                  color: '#fff',
                  boxShadow: '0 0 15px rgba(0, 217, 255, 0.2)',
                  overflowX: 'auto',
                  maxWidth: 'calc(100% - 40px)'
               }}>
                 <span style={{ color: 'var(--color-network-cyan)', marginRight: '8px' }}>{ROUTES[activeIndex].method}</span>
                 {ROUTES[activeIndex].path}
               </div>
            </div>

            {/* Handlers List */}
            <div style={{ display: 'flex', flexDirection: 'column', gap: '20px', position: 'relative', paddingLeft: '32px' }}>
               
               {/* Vertical glowing spine */}
               <div style={{
                 position: 'absolute',
                 left: 0,
                 top: '-16px', /* Start slightly above to connect with the icon area */
                 width: '2px',
                 height: 'calc(100% - 20px)',
                 background: 'rgba(255,255,255,0.1)'
               }}>
                 {/* The scanner */}
                 <div style={{
                   width: '100%',
                   background: 'var(--color-network-cyan)',
                   boxShadow: '0 0 10px var(--color-network-cyan)',
                   animation: 'trunkScan 4s infinite linear'
                 }}/>
               </div>

               {ROUTES.map((route, idx) => {
                 const isActive = idx === activeIndex;
                 return (
                   <div key={idx} style={{ position: 'relative', opacity: isActive ? 1 : 0.4 }}>
                      {/* Branch Line */}
                      <div style={{
                        position: 'absolute',
                        left: '-32px',
                        top: '50%',
                        width: '32px',
                        height: '2px',
                        background: 'rgba(255,255,255,0.1)'
                      }}>
                        {/* Glowing Branch Fill */}
                        {isActive && (
                          <div style={{
                            height: '100%',
                            background: 'var(--color-network-cyan)',
                            boxShadow: '0 0 10px var(--color-network-cyan)',
                            animation: 'branchFill 4s infinite ease-in-out'
                          }}/>
                        )}
                      </div>
                      
                      {/* Handler Card */}
                      <div style={{
                         padding: '16px',
                         border: '1px solid rgba(255,255,255,0.1)',
                         borderRadius: '8px',
                         background: isActive ? 'rgba(255,255,255,0.02)' : 'transparent',
                         display: 'flex',
                         justifyContent: 'space-between',
                         fontFamily: 'var(--font-jetbrains-mono)',
                         fontSize: '13px',
                         color: 'var(--color-steel-mid)',
                         animation: isActive ? 'handlerPulse 4s infinite cubic-bezier(0.4, 0, 0.2, 1)' : 'none',
                         gap: '8px',
                         overflowX: 'auto'
                      }}>
                         <span style={{ color: isActive ? 'inherit' : 'var(--color-steel-mid)' }}>{route.path}</span>
                         <span>{route.handler}</span>
                      </div>
                   </div>
                 );
               })}
            </div>

          </div>
        </div>
      </div>
    </section>
  );
};
