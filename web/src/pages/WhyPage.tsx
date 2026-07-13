import React, { useEffect } from 'react';
import { Target, ServerCrash, GitMerge, CheckCircle } from 'lucide-react';

export const WhyPage: React.FC = () => {
  useEffect(() => {
    window.scrollTo(0, 0);
  }, []);

  return (
    <>
      <style>{`
        .why-container {
          padding: 160px 40px 80px 40px;
          max-width: 1200px;
          margin: 0 auto;
          min-height: 100vh;
          display: flex;
          flex-direction: column;
          gap: 80px;
        }
        .why-grid-columns {
          display: grid;
          grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
          gap: 48px;
          align-items: start;
        }
        .why-grid-columns-subsystems {
          display: grid;
          grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
          gap: 32px;
        }
        .why-stats-grid {
          display: grid;
          grid-template-columns: 1fr 1fr;
          gap: 16px;
          margin-top: 8px;
          padding-top: 16px;
          border-top: 1px solid rgba(255,255,255,0.08);
        }
        .chapter-roadmap-container::-webkit-scrollbar {
          height: 6px;
        }
        .chapter-roadmap-container::-webkit-scrollbar-track {
          background: rgba(255, 255, 255, 0.02);
        }
        .chapter-roadmap-container::-webkit-scrollbar-thumb {
          background: var(--color-network-cyan);
          border-radius: 4px;
        }

        @media (max-width: 1024px) {
          .why-container {
            padding: 100px 24px 60px 24px !important;
            gap: 48px !important;
          }
          .why-grid-columns {
            grid-template-columns: 1fr !important;
            gap: 32px !important;
          }
          .why-grid-columns-subsystems {
            grid-template-columns: 1fr !important;
            gap: 24px !important;
          }
          .why-container h1 {
            font-size: clamp(36px, 8vw, 72px) !important;
          }
        }
        @media (max-width: 480px) {
          .why-stats-grid {
            grid-template-columns: 1fr !important;
            gap: 12px !important;
          }
        }
      `}</style>
      <div className="why-container">
      {/* Header Section (Full Width) */}
      <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
        <div style={{ 
          fontFamily: 'var(--font-roboto-mono)', 
          color: 'var(--color-network-cyan)', 
          fontSize: '12px', 
          letterSpacing: '0.182em', 
          textTransform: 'uppercase'
        }}>
          &gt;_ SYSTEM DESIGN THESIS
        </div>
        <h1 style={{ 
          fontFamily: 'var(--font-lambotype)', 
          fontSize: 'clamp(44px, 7vw, 100px)', 
          color: 'var(--color-paper-white)', 
          margin: 0, 
          textTransform: 'uppercase',
          letterSpacing: '0.023em',
          lineHeight: 0.95
        }}>
          The Case for<br/>Raw Infrastructure
        </h1>
      </div>

      {/* Intro & Telemetry Section (Two Columns) */}
      <div className="why-grid-columns">
        {/* Left Column: Thesis Copy */}
        <div style={{ display: 'flex', flexDirection: 'column', gap: '24px' }}>
          <p style={{ 
            fontFamily: 'var(--font-suisse-intl)', 
            fontSize: '17px', 
            color: 'var(--color-steel-mid)', 
            lineHeight: 1.6,
            margin: 0,
            fontWeight: 400
          }}>
            Modern web application development is built upon extensive layers of high-level abstractions. While third-party frameworks and standard libraries accelerate development velocity, they isolate developers from the underlying mechanics of network sockets, operating systems, and kernel interactions. 
          </p>
          <p style={{ 
            fontFamily: 'var(--font-suisse-intl)', 
            fontSize: '17px', 
            color: 'var(--color-steel-mid)', 
            lineHeight: 1.6,
            margin: 0,
            fontWeight: 400
          }}>
            TitanHTTP was engineered to break through these layers. By building a production-inspired HTTP/1.1 engine from the raw TCP socket level up, this project implements core network protocols, memory boundary management, and concurrent synchronization directly. It stands as a demonstration of software architecture built from first principles.
          </p>
        </div>

        {/* Right Column: Telemetry Comparison Console */}
        <div style={{
          background: 'var(--color-frosted-graphite)',
          backdropFilter: 'blur(16px)',
          WebkitBackdropFilter: 'blur(16px)',
          border: '1px solid rgba(255, 255, 255, 0.05)',
          borderRadius: 'var(--radius-glass)',
          padding: '32px',
          boxShadow: '0 20px 50px rgba(0,0,0,0.7)',
          display: 'flex',
          flexDirection: 'column',
          gap: '24px'
        }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', borderBottom: '1px solid rgba(255,255,255,0.08)', paddingBottom: '16px' }}>
            <div style={{ fontFamily: 'var(--font-roboto-mono)', fontSize: '11px', color: 'var(--color-paper-white)', letterSpacing: '0.1em' }}>COMPARATIVE TELEMETRY</div>
            <div style={{ fontFamily: 'var(--font-roboto-mono)', fontSize: '9px', color: 'var(--color-network-cyan)', background: 'rgba(0, 217, 255, 0.1)', padding: '4px 8px', borderRadius: '2px' }}>v1.0.0 STABLE</div>
          </div>
          
          <div style={{ display: 'flex', flexDirection: 'column', gap: '20px' }}>
            {/* Metric 1 */}
            <div>
              <div style={{ display: 'flex', justifyContent: 'space-between', fontFamily: 'var(--font-suisse-intl)', fontSize: '13px', marginBottom: '8px' }}>
                <span style={{ color: 'var(--color-paper-white)', fontWeight: 500 }}>Max Raw Throughput (Ping)</span>
                <span style={{ color: 'var(--color-network-cyan)', fontFamily: 'var(--font-roboto-mono)', fontSize: '12px' }}>+41.7% Winner</span>
              </div>
              <div style={{ display: 'grid', gridTemplateColumns: '1fr auto', gap: '16px', alignItems: 'center' }}>
                <div style={{ height: '8px', background: 'rgba(255,255,255,0.05)', borderRadius: '4px', overflow: 'hidden', position: 'relative' }}>
                  <div style={{ position: 'absolute', height: '100%', left: 0, top: 0, width: '100%', background: 'var(--color-network-cyan)', borderRadius: '4px' }} />
                  <div style={{ position: 'absolute', height: '100%', left: 0, top: 0, width: '76%', background: 'rgba(255,255,255,0.3)', borderRadius: '4px' }} />
                </div>
                <div style={{ fontFamily: 'var(--font-roboto-mono)', fontSize: '12px', color: 'var(--color-paper-white)' }}>
                  123.4K <span style={{ color: 'var(--color-steel-mid)', fontSize: '10px' }}>vs 94.7K req/s</span>
                </div>
              </div>
            </div>

            {/* Metric 2 */}
            <div>
              <div style={{ display: 'flex', justifyContent: 'space-between', fontFamily: 'var(--font-suisse-intl)', fontSize: '13px', marginBottom: '8px' }}>
                <span style={{ color: 'var(--color-paper-white)', fontWeight: 500 }}>Large Payload (10KB)</span>
                <span style={{ color: 'var(--color-network-cyan)', fontFamily: 'var(--font-roboto-mono)', fontSize: '12px' }}>+115% Winner</span>
              </div>
              <div style={{ display: 'grid', gridTemplateColumns: '1fr auto', gap: '16px', alignItems: 'center' }}>
                <div style={{ height: '8px', background: 'rgba(255,255,255,0.05)', borderRadius: '4px', overflow: 'hidden', position: 'relative' }}>
                  <div style={{ position: 'absolute', height: '100%', left: 0, top: 0, width: '100%', background: 'var(--color-network-cyan)', borderRadius: '4px' }} />
                  <div style={{ position: 'absolute', height: '100%', left: 0, top: 0, width: '47%', background: 'rgba(255,255,255,0.3)', borderRadius: '4px' }} />
                </div>
                <div style={{ fontFamily: 'var(--font-roboto-mono)', fontSize: '12px', color: 'var(--color-paper-white)' }}>
                  74.2K <span style={{ color: 'var(--color-steel-mid)', fontSize: '10px' }}>vs 34.9K req/s</span>
                </div>
              </div>
            </div>

            {/* Metric 3 */}
            <div>
              <div style={{ display: 'flex', justifyContent: 'space-between', fontFamily: 'var(--font-suisse-intl)', fontSize: '13px', marginBottom: '8px' }}>
                <span style={{ color: 'var(--color-paper-white)', fontWeight: 500 }}>P99 Latency (Max Load)</span>
                <span style={{ color: 'var(--color-packet-violet)', fontFamily: 'var(--font-roboto-mono)', fontSize: '12px' }}>-56.9% Lower</span>
              </div>
              <div style={{ display: 'grid', gridTemplateColumns: '1fr auto', gap: '16px', alignItems: 'center' }}>
                <div style={{ height: '8px', background: 'rgba(255,255,255,0.05)', borderRadius: '4px', overflow: 'hidden', position: 'relative' }}>
                  <div style={{ position: 'absolute', height: '100%', left: 0, top: 0, width: '43%', background: 'var(--color-packet-violet)', borderRadius: '4px' }} />
                  <div style={{ position: 'absolute', height: '100%', left: 0, top: 0, width: '100%', background: 'rgba(255,255,255,0.1)', borderRadius: '4px' }} />
                </div>
                <div style={{ fontFamily: 'var(--font-roboto-mono)', fontSize: '12px', color: 'var(--color-paper-white)' }}>
                  4.76ms <span style={{ color: 'var(--color-steel-mid)', fontSize: '10px' }}>vs 11.06ms</span>
                </div>
              </div>
            </div>
          </div>

          <div className="why-stats-grid">
            <div style={{ textAlign: 'center', background: 'rgba(255,255,255,0.02)', padding: '12px', borderRadius: '8px' }}>
              <div style={{ fontFamily: 'var(--font-roboto-mono)', fontSize: '10px', color: 'var(--color-steel-mid)', letterSpacing: '0.05em' }}>RFC COMPLIANCE</div>
              <div style={{ fontFamily: 'var(--font-roboto-mono)', fontSize: '18px', color: 'var(--color-paper-white)', marginTop: '4px' }}>20 / 20 <span style={{ fontSize: '10px', color: 'var(--color-network-cyan)' }}>PASS</span></div>
            </div>
            <div style={{ textAlign: 'center', background: 'rgba(255,255,255,0.02)', padding: '12px', borderRadius: '8px' }}>
              <div style={{ fontFamily: 'var(--font-roboto-mono)', fontSize: '10px', color: 'var(--color-steel-mid)', letterSpacing: '0.05em' }}>SECURITY ATTACKS</div>
              <div style={{ fontFamily: 'var(--font-roboto-mono)', fontSize: '18px', color: 'var(--color-paper-white)', marginTop: '4px' }}>8 / 8 <span style={{ fontSize: '10px', color: 'var(--color-network-cyan)' }}>BLOCKED</span></div>
            </div>
          </div>
        </div>
      </div>

      {/* Chapter Pipeline */}
      <div>
        <div style={{ 
          fontFamily: 'var(--font-roboto-mono)', 
          color: 'var(--color-paper-white)', 
          fontSize: '14px', 
          letterSpacing: '0.182em', 
          marginBottom: '32px', 
          borderBottom: '1px solid rgba(255,255,255,0.1)', 
          paddingBottom: '16px',
          textTransform: 'uppercase'
        }}>
          THE GUIDED NARRATIVE ROADMAP
        </div>
        
        <div style={{
          background: 'rgba(255,255,255,0.02)',
          border: '1px solid rgba(255,255,255,0.05)',
          borderRadius: 'var(--radius-glass)',
          padding: '40px 24px',
          display: 'flex',
          flexDirection: 'column',
          gap: '32px',
          overflowX: 'auto'
        }}>
          <div style={{ 
            display: 'flex', 
            justifyContent: 'space-between', 
            position: 'relative',
            minWidth: '900px',
            padding: '0 20px'
          }}>
            {/* Connecting line */}
            <div style={{ 
              position: 'absolute', 
              top: '15px', 
              left: '40px', 
              right: '40px', 
              height: '1px', 
              background: 'linear-gradient(to right, var(--color-network-cyan), var(--color-packet-violet), var(--color-steel-mid))',
              zIndex: 0
            }} />
            
            {[
              { num: 'C1', title: 'Why HTTP', label: 'HTTP Roots' },
              { num: 'C2', title: 'TCP socket', label: 'Transport' },
              { num: 'C3', title: 'Socket Bind', label: 'Syscalls' },
              { num: 'C4', title: 'Read Bytes', label: 'Buffer' },
              { num: 'C5', title: 'Parse Req', label: 'Parsing' },
              { num: 'C6', title: 'Radix Route', label: 'Routing' },
              { num: 'C7', title: 'Worker Pool', label: 'Concurrency' },
              { num: 'C8', title: 'Compression', label: 'Production' },
              { num: 'C9', title: 'Benchmarking', label: 'Performance' },
              { num: 'C10', title: 'Source Code', label: 'Orchestrator' }
            ].map((ch, idx) => (
              <div 
                key={idx} 
                style={{ 
                  display: 'flex', 
                  flexDirection: 'column', 
                  alignItems: 'center', 
                  gap: '12px',
                  zIndex: 1,
                  width: '80px',
                  cursor: 'pointer'
                }}
                onMouseEnter={(e) => {
                  const node = e.currentTarget.querySelector('.node-dot') as HTMLDivElement;
                  if (node) {
                    node.style.borderColor = 'var(--color-network-cyan)';
                    node.style.boxShadow = '0 0 12px var(--color-network-cyan)';
                  }
                }}
                onMouseLeave={(e) => {
                  const node = e.currentTarget.querySelector('.node-dot') as HTMLDivElement;
                  if (node) {
                    node.style.borderColor = idx < 9 ? 'var(--color-network-cyan)' : 'var(--color-steel-mid)';
                    node.style.boxShadow = 'none';
                  }
                }}
              >
                <div 
                  className="node-dot"
                  style={{
                    width: '30px',
                    height: '30px',
                    borderRadius: '50%',
                    background: 'var(--color-void-black)',
                    border: `2px solid ${idx < 9 ? 'var(--color-network-cyan)' : 'var(--color-steel-mid)'}`,
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    fontFamily: 'var(--font-roboto-mono)',
                    fontSize: '9px',
                    color: 'var(--color-paper-white)',
                    transition: 'all var(--motion-fast) var(--ease-mechanical)'
                  }}
                >
                  {ch.num}
                </div>
                <div style={{ textAlign: 'center' }}>
                  <div style={{ fontFamily: 'var(--font-roboto-mono)', fontSize: '10px', color: 'var(--color-paper-white)', fontWeight: 500 }}>{ch.title}</div>
                  <div style={{ fontFamily: 'var(--font-suisse-intl)', fontSize: '9px', color: 'var(--color-steel-mid)', marginTop: '2px' }}>{ch.label}</div>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* Grid of Challenges */}
      <div>
        <div style={{ 
          fontFamily: 'var(--font-roboto-mono)', 
          color: 'var(--color-paper-white)', 
          fontSize: '14px', 
          letterSpacing: '0.182em', 
          marginBottom: '32px', 
          borderBottom: '1px solid rgba(255,255,255,0.1)', 
          paddingBottom: '16px',
          textTransform: 'uppercase'
        }}>
          CRITICAL SUBSYSTEM ENGINEERING
        </div>
        
        <div className="why-grid-columns-subsystems">
          
          {/* Card 1 */}
          <div 
            style={{
              background: 'var(--color-frosted-graphite)',
              backdropFilter: 'blur(16px)',
              WebkitBackdropFilter: 'blur(16px)',
              border: '1px solid rgba(255, 255, 255, 0.05)',
              borderRadius: 'var(--radius-glass)',
              padding: '40px',
              display: 'flex',
              flexDirection: 'column',
              gap: '24px',
              boxShadow: '0 10px 40px rgba(0,0,0,0.6)',
              transition: 'all var(--motion-fast) var(--ease-mechanical)',
              cursor: 'default'
            }}
            onMouseEnter={(e) => {
              e.currentTarget.style.borderColor = 'rgba(0, 217, 255, 0.3)';
              e.currentTarget.style.transform = 'translateY(-4px)';
              e.currentTarget.style.boxShadow = '0 15px 45px rgba(0, 217, 255, 0.05)';
            }}
            onMouseLeave={(e) => {
              e.currentTarget.style.borderColor = 'rgba(255, 255, 255, 0.05)';
              e.currentTarget.style.transform = 'translateY(0)';
              e.currentTarget.style.boxShadow = '0 10px 40px rgba(0,0,0,0.6)';
            }}
          >
            <div style={{ color: 'var(--color-network-cyan)' }}><GitMerge size={32} /></div>
            <div>
              <h3 style={{ 
                fontFamily: 'var(--font-roboto-mono)', 
                color: 'var(--color-paper-white)', 
                fontSize: '15px', 
                letterSpacing: '0.05em', 
                margin: '0 0 12px 0',
                textTransform: 'uppercase'
              }}>
                Zero-Allocation Request Parsing
              </h3>
              
              {/* Inline SVG showing Sawtooth (spiky heap allocations) vs. Flat Line (zero-allocation) */}
              <div style={{ height: '60px', width: '100%', margin: '16px 0', borderBottom: '1px solid rgba(255,255,255,0.05)', position: 'relative' }}>
                <svg width="100%" height="100%" viewBox="0 0 300 60" preserveAspectRatio="none">
                  <path 
                    d="M 0 50 L 30 10 L 30 50 L 60 10 L 60 50 L 90 10 L 90 50 L 120 10 L 120 50 L 150 10 L 150 50 L 180 10 L 180 50 L 210 10 L 210 50 L 240 10 L 240 50 L 270 10 L 270 50 L 300 10" 
                    fill="none" 
                    stroke="rgba(255, 255, 255, 0.15)" 
                    strokeWidth="1.5" 
                    strokeDasharray="3,3"
                  />
                  <line 
                    x1="0" y1="45" x2="300" y2="45" 
                    stroke="var(--color-network-cyan)" 
                    strokeWidth="2.5" 
                  />
                </svg>
                <div style={{ display: 'flex', justifyContent: 'space-between', fontFamily: 'var(--font-roboto-mono)', fontSize: '8px', color: 'var(--color-steel-mid)', marginTop: '4px' }}>
                  <span>HEAP ALLOCATIONS (SAWTOOTH SPIKES)</span>
                  <span style={{ color: 'var(--color-network-cyan)' }}>TITANHTTP FLAT BOUND</span>
                </div>
              </div>

              <p style={{ 
                fontFamily: 'var(--font-suisse-intl)', 
                color: 'var(--color-steel-mid)', 
                fontSize: '14px', 
                lineHeight: 1.6, 
                margin: 0 
              }}>
                Parsing unstructured, stream-oriented TCP byte chunks presents strict performance boundaries. Standard parsing patterns produce high allocation volume, inducing garbage collector sweeps and introducing latency variance. TitanHTTP resolves this by utilizing an optimized state-machine parser combined with a recycling <code style={{ fontFamily: 'var(--font-jetbrains-mono)', color: 'var(--color-network-cyan)', fontSize: '13px' }}>sync.Pool</code> layer. By indexing header slices lazily and preserving buffers across connection lifecycles, memory fragmentation is effectively negated.
              </p>
            </div>
          </div>

          {/* Card 2 */}
          <div 
            style={{
              background: 'var(--color-frosted-graphite)',
              backdropFilter: 'blur(16px)',
              WebkitBackdropFilter: 'blur(16px)',
              border: '1px solid rgba(255, 255, 255, 0.05)',
              borderRadius: 'var(--radius-glass)',
              padding: '40px',
              display: 'flex',
              flexDirection: 'column',
              gap: '24px',
              boxShadow: '0 10px 40px rgba(0,0,0,0.6)',
              transition: 'all var(--motion-fast) var(--ease-mechanical)',
              cursor: 'default'
            }}
            onMouseEnter={(e) => {
              e.currentTarget.style.borderColor = 'rgba(132, 125, 255, 0.4)';
              e.currentTarget.style.transform = 'translateY(-4px)';
              e.currentTarget.style.boxShadow = '0 15px 45px rgba(132, 125, 255, 0.05)';
            }}
            onMouseLeave={(e) => {
              e.currentTarget.style.borderColor = 'rgba(255, 255, 255, 0.05)';
              e.currentTarget.style.transform = 'translateY(0)';
              e.currentTarget.style.boxShadow = '0 10px 40px rgba(0,0,0,0.6)';
            }}
          >
            <div style={{ color: 'var(--color-packet-violet)' }}><Target size={32} /></div>
            <div>
              <h3 style={{ 
                fontFamily: 'var(--font-roboto-mono)', 
                color: 'var(--color-paper-white)', 
                fontSize: '15px', 
                letterSpacing: '0.05em', 
                margin: '0 0 12px 0',
                textTransform: 'uppercase'
              }}>
                Bounded Concurrency & Shedding
              </h3>

              {/* Inline SVG showing job queue pipeline and load shedding diversion */}
              <div style={{ height: '60px', width: '100%', margin: '16px 0', borderBottom: '1px solid rgba(255,255,255,0.05)', position: 'relative' }}>
                <svg width="100%" height="100%" viewBox="0 0 300 60" fill="none">
                  <rect x="70" y="15" width="130" height="30" rx="4" stroke="rgba(255,255,255,0.15)" strokeWidth="1.5" />
                  <text x="135" y="33" fill="var(--color-steel-mid)" fontSize="8" fontFamily="var(--font-roboto-mono)" textAnchor="middle">JOB QUEUE [100]</text>
                  <circle cx="90" cy="30" r="4" fill="var(--color-packet-violet)" />
                  <circle cx="105" cy="30" r="4" fill="var(--color-packet-violet)" />
                  <circle cx="120" cy="30" r="4" fill="var(--color-packet-violet)" />
                  <path d="M 10 30 L 60 30" stroke="var(--color-network-cyan)" strokeWidth="2" strokeDasharray="4,4" />
                  <circle cx="30" cy="30" r="3" fill="var(--color-network-cyan)" />
                  <path d="M 50 30 C 50 15, 60 5, 80 5" stroke="#ff4a4a" strokeWidth="1.5" strokeDasharray="3,3" />
                  <circle cx="80" cy="5" r="3" fill="#ff4a4a" />
                  <text x="90" y="8" fill="#ff4a4a" fontSize="7" fontFamily="var(--font-roboto-mono)">503 SHED</text>
                </svg>
                <div style={{ display: 'flex', justifyContent: 'space-between', fontFamily: 'var(--font-roboto-mono)', fontSize: '8px', color: 'var(--color-steel-mid)', marginTop: '4px' }}>
                  <span>WORKER CAP EXHAUSTION</span>
                  <span style={{ color: '#ff4a4a' }}>AUTO LOAD SHEDDING ACTIVE</span>
                </div>
              </div>

              <p style={{ 
                fontFamily: 'var(--font-suisse-intl)', 
                color: 'var(--color-steel-mid)', 
                fontSize: '14px', 
                lineHeight: 1.6, 
                margin: 0 
              }}>
                Unbounded execution allocation poses risk under high concurrency workloads. Spawning new goroutines dynamically for thousands of active connections exposes servers to socket starvation and memory limits. TitanHTTP implements a fixed Worker Pool pattern to bound scheduler overhead. Connections enter a thread-safe task channel; when concurrency thresholds are saturated, the system initiates deterministic load-shedding, serving immediate 503 limits to guarantee baseline integrity.
              </p>
            </div>
          </div>

          {/* Card 3 */}
          <div 
            style={{
              background: 'var(--color-frosted-graphite)',
              backdropFilter: 'blur(16px)',
              WebkitBackdropFilter: 'blur(16px)',
              border: '1px solid rgba(255, 255, 255, 0.05)',
              borderRadius: 'var(--radius-glass)',
              padding: '40px',
              display: 'flex',
              flexDirection: 'column',
              gap: '24px',
              boxShadow: '0 10px 40px rgba(0,0,0,0.6)',
              transition: 'all var(--motion-fast) var(--ease-mechanical)',
              cursor: 'default'
            }}
            onMouseEnter={(e) => {
              e.currentTarget.style.borderColor = 'rgba(255, 255, 255, 0.2)';
              e.currentTarget.style.transform = 'translateY(-4px)';
              e.currentTarget.style.boxShadow = '0 15px 45px rgba(255, 255, 255, 0.05)';
            }}
            onMouseLeave={(e) => {
              e.currentTarget.style.borderColor = 'rgba(255, 255, 255, 0.05)';
              e.currentTarget.style.transform = 'translateY(0)';
              e.currentTarget.style.boxShadow = '0 10px 40px rgba(0,0,0,0.6)';
            }}
          >
            <div style={{ color: 'var(--color-paper-white)' }}><ServerCrash size={32} /></div>
            <div>
              <h3 style={{ 
                fontFamily: 'var(--font-roboto-mono)', 
                color: 'var(--color-paper-white)', 
                fontSize: '15px', 
                letterSpacing: '0.05em', 
                margin: '0 0 12px 0',
                textTransform: 'uppercase'
              }}>
                Fault Isolation & Graceful Recovery
              </h3>

              {/* Inline SVG showing request smuggling filter and path traversal blocks */}
              <div style={{ height: '60px', width: '100%', margin: '16px 0', borderBottom: '1px solid rgba(255,255,255,0.05)', position: 'relative' }}>
                <svg width="100%" height="100%" viewBox="0 0 300 60" fill="none">
                  <path d="M 150 10 L 190 18 L 190 38 L 150 48 L 110 38 L 110 18 Z" fill="rgba(255, 255, 255, 0.02)" stroke="rgba(255, 255, 255, 0.15)" strokeWidth="1.5" />
                  <path d="M 20 30 L 110 30" stroke="var(--color-network-cyan)" strokeWidth="2" />
                  <path d="M 190 30 L 280 30" stroke="var(--color-network-cyan)" strokeWidth="2" />
                  <path d="M 40 12 L 120 22" stroke="#ff4a4a" strokeWidth="1.5" strokeDasharray="3,3" />
                  <circle cx="120" cy="22" r="3" fill="#ff4a4a" />
                  <text x="50" y="10" fill="#ff4a4a" fontSize="7" fontFamily="var(--font-roboto-mono)">CL-TE SMUGGLE</text>
                  <path d="M 40 48 L 120 38" stroke="#ff4a4a" strokeWidth="1.5" strokeDasharray="3,3" />
                  <circle cx="120" cy="38" r="3" fill="#ff4a4a" />
                  <text x="50" y="55" fill="#ff4a4a" fontSize="7" fontFamily="var(--font-roboto-mono)">SLOWLORIS SHIELD</text>
                </svg>
                <div style={{ display: 'flex', justifyContent: 'space-between', fontFamily: 'var(--font-roboto-mono)', fontSize: '8px', color: 'var(--color-steel-mid)', marginTop: '4px' }}>
                  <span>RFC STANDARDS GATEWAY</span>
                  <span style={{ color: 'var(--color-network-cyan)' }}>8 / 8 PROACTIVE DEFENSES</span>
                </div>
              </div>

              <p style={{ 
                fontFamily: 'var(--font-suisse-intl)', 
                color: 'var(--color-steel-mid)', 
                fontSize: '14px', 
                lineHeight: 1.6, 
                margin: 0 
              }}>
                Reliable backends must exhibit graceful degradation under anomalous states. TitanHTTP treats connection states as independent domains. If socket negotiation or stream syntax fails, the connection buffer is safely drained and recycled. If a handler crashes, an isolated recovery middleware intercepts the panic via Go's deferred recovery structures, recording stack trace telemetry and returning a 500 response while keeping the main listener functional.
              </p>
            </div>
          </div>

        </div>
      </div>
      
      {/* Conclusion / Verification */}
      <div style={{ 
        display: 'flex', 
        alignItems: 'center', 
        gap: '32px', 
        background: 'rgba(0, 217, 255, 0.03)', 
        border: '1px solid rgba(0, 217, 255, 0.1)', 
        padding: '40px', 
        borderRadius: 'var(--radius-glass)',
        backdropFilter: 'blur(10px)',
        WebkitBackdropFilter: 'blur(10px)'
      }}>
        <CheckCircle size={48} color="var(--color-network-cyan)" style={{ flexShrink: 0 }} />
        <div>
          <h4 style={{ 
            fontFamily: 'var(--font-roboto-mono)', 
            color: 'var(--color-paper-white)', 
            fontSize: '12px', 
            letterSpacing: '0.182em', 
            margin: '0 0 8px 0',
            textTransform: 'uppercase'
          }}>
            VERIFICATION & DETERMINISTIC BENCHMARKS
          </h4>
          <p style={{ 
            fontFamily: 'var(--font-suisse-intl)', 
            color: 'var(--color-steel-mid)', 
            fontSize: '14px', 
            lineHeight: 1.6, 
            margin: 0 
          }}>
            TitanHTTP is not a simplified mock implementation. It is a benchmarked, profile-optimized, and spec-compliant HTTP engine. Through intensive execution testing, regression sweeps against the Go native library, and profiling sweeps via CPU and heap pprof tooling, we guarantee a stable, deterministic, and structurally mature codebase designed for modern performance engineering.
          </p>
        </div>
      </div>

      {/* Recruiter Quote Block */}
      <div style={{
        background: 'rgba(255, 255, 255, 0.02)',
        border: '1px solid rgba(255,255,255,0.06)',
        borderRadius: 'var(--radius-glass)',
        padding: '64px 40px',
        textAlign: 'center',
        position: 'relative',
        overflow: 'hidden'
      }}>
        <div style={{
          position: 'absolute',
          top: '-20px',
          left: '20px',
          fontSize: '120px',
          fontFamily: 'var(--font-lambotype)',
          color: 'rgba(255,255,255,0.02)',
          lineHeight: 1,
          userSelect: 'none'
        }}>&ldquo;</div>
        <p style={{
          fontFamily: 'var(--font-suisse-intl)',
          fontSize: '20px',
          color: 'var(--color-paper-white)',
          lineHeight: 1.6,
          maxWidth: '850px',
          margin: '0 auto 16px auto',
          fontStyle: 'italic',
          fontWeight: 300
        }}>
          &ldquo;Would a senior backend engineer enjoy reviewing this code, and would a recruiter trust the developer who built it?&rdquo;
        </p>
        <div style={{
          fontFamily: 'var(--font-roboto-mono)',
          fontSize: '11px',
          color: 'var(--color-network-cyan)',
          letterSpacing: '0.15em',
          textTransform: 'uppercase'
        }}>
          — The TitanHTTP Engineering Rule
        </div>
      </div>

    </div>
    </>
  );
};
