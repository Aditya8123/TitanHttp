import React, { useEffect } from 'react';
import { Layers, Cpu, Network, ArrowRight, User, Users, FileText, GitMerge, Code, Send, Server, ShieldAlert, Database } from 'lucide-react';

const FlowNode = ({ icon, title, desc }: { icon: React.ReactNode, title: string, desc: string }) => (
  <div style={{ 
    background: 'rgba(20, 20, 22, 0.6)', 
    border: '1px solid rgba(0, 217, 255, 0.2)', 
    borderRadius: '8px', 
    padding: '16px 24px', 
    display: 'flex', 
    flexDirection: 'column', 
    alignItems: 'center', 
    gap: '12px',
    minWidth: '160px',
    boxShadow: '0 4px 20px rgba(0,0,0,0.3)',
    textAlign: 'center',
    width: '100%',
    maxWidth: '220px'
  }}>
    <div style={{ color: 'var(--color-network-cyan)' }}>{icon}</div>
    <div>
      <div style={{ fontFamily: 'var(--font-roboto-mono)', fontSize: '11px', color: 'var(--color-paper-white)', letterSpacing: '0.05em' }}>{title}</div>
      <div style={{ fontFamily: 'var(--font-suisse-intl)', fontSize: '10px', color: 'var(--color-steel-mid)', marginTop: '4px' }}>{desc}</div>
    </div>
  </div>
);

export const ArchitecturePage: React.FC = () => {
  useEffect(() => {
    window.scrollTo(0, 0);
  }, []);

  return (
    <>
      <style>{`
        .arch-container {
          padding: 160px 40px 80px 40px;
          max-width: 1200px;
          margin: 0 auto;
          min-height: 100vh;
          display: flex;
          flex-direction: column;
          gap: 80px;
        }
        .arch-flow-wrapper {
          display: flex;
          flex-wrap: wrap;
          gap: 16px;
          align-items: center;
          justify-content: center;
          padding: 40px;
          background: rgba(0, 217, 255, 0.02);
          border: 1px dashed rgba(0, 217, 255, 0.1);
          border-radius: 16px;
        }
        .arch-components-grid {
          display: grid;
          grid-template-columns: repeat(auto-fit, minmax(360px, 1fr));
          gap: 32px;
        }
        .arch-subsystems-grid {
          display: grid;
          grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
          gap: 32px;
        }
        
        @media (max-width: 1024px) {
          .arch-container {
            padding: 100px 24px 60px 24px !important;
            gap: 48px !important;
          }
          .arch-flow-wrapper {
            flex-direction: column !important;
            padding: 24px 16px !important;
          }
          .arch-components-grid {
            grid-template-columns: 1fr !important;
            gap: 24px !important;
          }
          .arch-subsystems-grid {
            grid-template-columns: 1fr !important;
            gap: 24px !important;
          }
          .arch-components-grid > div {
            padding: 24px !important;
          }
        }
      `}</style>
      <div className="arch-container">
        {/* Header */}
        <div style={{ textAlign: 'center' }}>
          <div style={{ fontFamily: 'var(--font-roboto-mono)', color: 'var(--color-network-cyan)', fontSize: '14px', letterSpacing: '0.1em', marginBottom: '16px' }}>&gt;_ SYSTEM BLUEPRINT</div>
          <h1 style={{ fontFamily: 'var(--font-lambotype)', fontSize: 'clamp(40px, 6vw, 80px)', color: 'var(--color-paper-white)', margin: '0 0 24px 0', textTransform: 'uppercase' }}>
            Architecture
          </h1>
          <p style={{ fontFamily: 'var(--font-suisse-intl)', fontSize: '18px', color: 'var(--color-steel-mid)', maxWidth: '800px', margin: '0 auto', lineHeight: 1.6 }}>
            TitanHTTP is designed to process HTTP requests efficiently, cleanly, and concurrently. It is modeled as a strict pipeline, transforming raw byte streams into structured responses.
          </p>
        </div>

        {/* The Request Lifecycle (Pipeline) */}
        <div>
          <div style={{ fontFamily: 'var(--font-roboto-mono)', color: 'var(--color-paper-white)', fontSize: '20px', letterSpacing: '0.1em', marginBottom: '32px', borderBottom: '1px solid rgba(255,255,255,0.1)', paddingBottom: '16px', textAlign: 'center' }}>
            THE REQUEST LIFECYCLE
          </div>
          
          <div className="arch-flow-wrapper">
            <FlowNode icon={<User size={24} />} title="CLIENT" desc="TCP Connection" />
            <ArrowRight size={20} className="flow-arrow-icon" color="var(--color-steel-mid)" />
            
            <FlowNode icon={<Network size={24} />} title="TCP LISTENER" desc="Accept Connection" />
            <ArrowRight size={20} className="flow-arrow-icon" color="var(--color-steel-mid)" />
            
            <FlowNode icon={<Users size={24} />} title="WORKER POOL" desc="Raw Bytes" />
            <ArrowRight size={20} className="flow-arrow-icon" color="var(--color-steel-mid)" />
            
            <FlowNode icon={<FileText size={24} />} title="HTTP PARSER" desc="Valid Request" />
            <ArrowRight size={20} className="flow-arrow-icon" color="var(--color-steel-mid)" />
            
            <FlowNode icon={<GitMerge size={24} />} title="ROUTER" desc="Match Route" />
            <ArrowRight size={20} className="flow-arrow-icon" color="var(--color-steel-mid)" />
            
            <FlowNode icon={<Layers size={24} />} title="MIDDLEWARE" desc="Pipeline" />
            <ArrowRight size={20} className="flow-arrow-icon" color="var(--color-steel-mid)" />
            
            <FlowNode icon={<Code size={24} />} title="HANDLER" desc="Generate Response" />
            <ArrowRight size={20} className="flow-arrow-icon" color="var(--color-steel-mid)" />
            
            <FlowNode icon={<Send size={24} />} title="RESPONSE WRITER" desc="Bytes over Socket" />
          </div>
        </div>

        {/* Core Components Grid */}
        <div>
          <div style={{ fontFamily: 'var(--font-roboto-mono)', color: 'var(--color-paper-white)', fontSize: '20px', letterSpacing: '0.1em', marginBottom: '32px', borderBottom: '1px solid rgba(255,255,255,0.1)', paddingBottom: '16px' }}>
            CORE COMPONENTS
          </div>
          
          <div className="arch-components-grid">
            
            <div style={{ background: 'rgba(20, 20, 22, 0.4)', border: '1px solid rgba(255, 255, 255, 0.05)', borderRadius: '8px', padding: '32px', display: 'flex', gap: '24px', transition: 'all 0.3s ease' }} onMouseEnter={(e) => e.currentTarget.style.borderColor = 'rgba(0, 217, 255, 0.3)'} onMouseLeave={(e) => e.currentTarget.style.borderColor = 'rgba(255, 255, 255, 0.05)'}>
              <div style={{ color: 'var(--color-network-cyan)', flexShrink: 0 }}><Network size={32} /></div>
              <div>
                <h3 style={{ fontFamily: 'var(--font-roboto-mono)', color: 'var(--color-paper-white)', fontSize: '16px', margin: '0 0 12px 0', letterSpacing: '0.05em' }}>1. TCP Listener & Manager</h3>
                <p style={{ fontFamily: 'var(--font-suisse-intl)', color: 'var(--color-steel-mid)', fontSize: '14px', lineHeight: 1.6, margin: 0 }}>
                  The foundation. Opens a raw socket and listens for incoming connections. Accepts connections, manages timeouts, and cleanly closes sockets. Each accepted connection is immediately handed off to the worker pool to ensure the main listener thread never blocks.
                </p>
              </div>
            </div>

            <div style={{ background: 'rgba(20, 20, 22, 0.4)', border: '1px solid rgba(255, 255, 255, 0.05)', borderRadius: '8px', padding: '32px', display: 'flex', gap: '24px', transition: 'all 0.3s ease' }} onMouseEnter={(e) => e.currentTarget.style.borderColor = 'rgba(138, 43, 226, 0.3)'} onMouseLeave={(e) => e.currentTarget.style.borderColor = 'rgba(255, 255, 255, 0.05)'}>
              <div style={{ color: 'var(--color-packet-violet)', flexShrink: 0 }}><Layers size={32} /></div>
              <div>
                <h3 style={{ fontFamily: 'var(--font-roboto-mono)', color: 'var(--color-paper-white)', fontSize: '16px', margin: '0 0 12px 0', letterSpacing: '0.05em' }}>2. HTTP Parser</h3>
                <p style={{ fontFamily: 'var(--font-suisse-intl)', color: 'var(--color-steel-mid)', fontSize: '14px', lineHeight: 1.6, margin: 0 }}>
                  Reads raw bytes from the TCP socket and translates them into structured Request objects. Must be strictly memory efficient, handling CRLF endings and Content-Length/Transfer-Encoding without allocating unnecessary strings.
                </p>
              </div>
            </div>

            <div style={{ background: 'rgba(20, 20, 22, 0.4)', border: '1px solid rgba(255, 255, 255, 0.05)', borderRadius: '8px', padding: '32px', display: 'flex', gap: '24px', transition: 'all 0.3s ease' }} onMouseEnter={(e) => e.currentTarget.style.borderColor = 'rgba(255, 255, 255, 0.3)'} onMouseLeave={(e) => e.currentTarget.style.borderColor = 'rgba(255, 255, 255, 0.05)'}>
              <div style={{ color: 'var(--color-paper-white)', flexShrink: 0 }}><Cpu size={32} /></div>
              <div>
                <h3 style={{ fontFamily: 'var(--font-roboto-mono)', color: 'var(--color-paper-white)', fontSize: '16px', margin: '0 0 12px 0', letterSpacing: '0.05em' }}>3. Radix Router & Middleware</h3>
                <p style={{ fontFamily: 'var(--font-suisse-intl)', color: 'var(--color-steel-mid)', fontSize: '14px', lineHeight: 1.6, margin: 0 }}>
                  The traffic controller. Uses a highly optimized Radix Tree per HTTP Method enabling O(k) pattern matching without slow regex. Middleware forms a chain of functions executing before and after the main handler.
                </p>
              </div>
            </div>

            <div style={{ background: 'rgba(20, 20, 22, 0.4)', border: '1px solid rgba(255, 255, 255, 0.05)', borderRadius: '8px', padding: '32px', display: 'flex', gap: '24px', transition: 'all 0.3s ease' }} onMouseEnter={(e) => e.currentTarget.style.borderColor = 'rgba(0, 217, 255, 0.3)'} onMouseLeave={(e) => e.currentTarget.style.borderColor = 'rgba(255, 255, 255, 0.05)'}>
              <div style={{ color: 'var(--color-network-cyan)', flexShrink: 0 }}><Users size={32} /></div>
              <div>
                <h3 style={{ fontFamily: 'var(--font-roboto-mono)', color: 'var(--color-paper-white)', fontSize: '16px', margin: '0 0 12px 0', letterSpacing: '0.05em' }}>4. Worker Pool & Sync</h3>
                <p style={{ fontFamily: 'var(--font-suisse-intl)', color: 'var(--color-steel-mid)', fontSize: '14px', lineHeight: 1.6, margin: 0 }}>
                  The concurrency engine. A fixed number of goroutines block on a job channel. If the queue is full, the server aggressively applies Load Shedding (503 Service Unavailable). Uses lock-free sync/atomic for metrics.
                </p>
              </div>
            </div>

          </div>
        </div>

        {/* Advanced Subsystems Grid */}
        <div>
          <div style={{ fontFamily: 'var(--font-roboto-mono)', color: 'var(--color-paper-white)', fontSize: '20px', letterSpacing: '0.1em', marginBottom: '32px', borderBottom: '1px solid rgba(255,255,255,0.1)', paddingBottom: '16px' }}>
            ADVANCED BACKEND SUBSYSTEMS
          </div>
          
          <div className="arch-subsystems-grid">
            
            <div style={{ background: 'var(--color-void-black)', border: '1px solid rgba(0, 217, 255, 0.2)', borderRadius: '8px', padding: '32px', boxShadow: '0 10px 30px rgba(0,217,255,0.05)' }}>
              <Server size={32} color="var(--color-network-cyan)" style={{ marginBottom: '16px' }} />
              <h3 style={{ fontFamily: 'var(--font-roboto-mono)', color: 'var(--color-paper-white)', fontSize: '16px', margin: '0 0 12px 0', letterSpacing: '0.05em' }}>Reverse Proxy & Load Balancer</h3>
              <p style={{ fontFamily: 'var(--font-suisse-intl)', color: 'var(--color-steel-mid)', fontSize: '14px', lineHeight: 1.6, margin: 0 }}>
                Acts as a gateway, forwarding incoming requests to external servers using net.DialTimeout. Wraps a pool of backends utilizing an atomic Round Robin selection algorithm with Active Health Checks.
              </p>
            </div>

            <div style={{ background: 'var(--color-void-black)', border: '1px solid rgba(138, 43, 226, 0.2)', borderRadius: '8px', padding: '32px', boxShadow: '0 10px 30px rgba(138,43,226,0.05)' }}>
              <ShieldAlert size={32} color="var(--color-packet-violet)" style={{ marginBottom: '16px' }} />
              <h3 style={{ fontFamily: 'var(--font-roboto-mono)', color: 'var(--color-paper-white)', fontSize: '16px', margin: '0 0 12px 0', letterSpacing: '0.05em' }}>Rate Limiter</h3>
              <p style={{ fontFamily: 'var(--font-suisse-intl)', color: 'var(--color-steel-mid)', fontSize: '14px', lineHeight: 1.6, margin: 0 }}>
                A dynamic defense mechanism. Implements both an ultra-efficient TokenBucket (using lazy-evaluation refills) and a strict-boundary SlidingWindow, supported by autonomous background sweepers to purge stale IPs.
              </p>
            </div>

            <div style={{ background: 'var(--color-void-black)', border: '1px solid rgba(255, 255, 255, 0.2)', borderRadius: '8px', padding: '32px', boxShadow: '0 10px 30px rgba(255,255,255,0.05)' }}>
              <Database size={32} color="var(--color-paper-white)" style={{ marginBottom: '16px' }} />
              <h3 style={{ fontFamily: 'var(--font-roboto-mono)', color: 'var(--color-paper-white)', fontSize: '16px', margin: '0 0 12px 0', letterSpacing: '0.05em' }}>Caching Layer</h3>
              <p style={{ fontFamily: 'var(--font-suisse-intl)', color: 'var(--color-steel-mid)', fontSize: '14px', lineHeight: 1.6, margin: 0 }}>
                A MemoryCache designed to bypass handler execution for frequently accessed resources. Protected by sync.RWMutex and fully parses Cache-Control directives with automatic TTL expiration via a sweeper.
              </p>
            </div>

          </div>
        </div>

      </div>
    </>
  );
};
