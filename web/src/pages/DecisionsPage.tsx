import React, { useEffect } from 'react';
import { ShieldCheck, GitBranch, ServerCog, Zap } from 'lucide-react';

export const DecisionsPage: React.FC = () => {
  useEffect(() => {
    window.scrollTo(0, 0);
  }, []);

  return (
    <div style={{
      padding: '160px 40px 80px 40px',
      maxWidth: '1200px',
      margin: '0 auto',
      minHeight: '100vh',
      display: 'flex',
      flexDirection: 'column',
      gap: '80px'
    }}>
      {/* Header */}
      <div style={{ textAlign: 'center' }}>
        <div style={{ fontFamily: 'var(--font-roboto-mono)', color: 'var(--color-network-cyan)', fontSize: '14px', letterSpacing: '0.1em', marginBottom: '16px' }}>&gt;_ ARCHITECTURAL DECISION RECORDS</div>
        <h1 style={{ fontFamily: 'var(--font-lambotype)', fontSize: 'clamp(40px, 6vw, 80px)', color: 'var(--color-paper-white)', margin: '0 0 24px 0', textTransform: 'uppercase' }}>
          Decisions & Trade-offs
        </h1>
        <p style={{ fontFamily: 'var(--font-suisse-intl)', fontSize: '18px', color: 'var(--color-steel-mid)', maxWidth: '800px', margin: '0 auto', lineHeight: 1.6 }}>
          We document the context, trade-offs, and reasoning behind significant design choices. This ensures that the <em>why</em> is preserved alongside the <em>how</em>.
        </p>
      </div>

      {/* Grid of ADRs */}
      <div>
        <div style={{ display: 'flex', flexDirection: 'column', gap: '32px' }}>
          
          {/* ADR 001 */}
          <div style={{
            background: 'rgba(20, 20, 22, 0.4)',
            border: '1px solid rgba(0, 217, 255, 0.2)',
            borderRadius: '12px',
            padding: '40px',
            display: 'flex',
            gap: '32px',
            boxShadow: '0 10px 30px rgba(0,0,0,0.5)',
            alignItems: 'flex-start'
          }}>
            <div style={{ color: 'var(--color-network-cyan)', background: 'rgba(0, 217, 255, 0.1)', padding: '16px', borderRadius: '12px' }}><ShieldCheck size={40} /></div>
            <div style={{ flex: 1 }}>
              <div style={{ fontFamily: 'var(--font-roboto-mono)', color: 'var(--color-network-cyan)', fontSize: '12px', letterSpacing: '0.1em', marginBottom: '8px' }}>ADR 001</div>
              <h3 style={{ fontFamily: 'var(--font-roboto-mono)', color: 'var(--color-paper-white)', fontSize: '20px', margin: '0 0 16px 0' }}>Building From Scratch in Pure Go</h3>
              <p style={{ fontFamily: 'var(--font-suisse-intl)', color: 'var(--color-steel-mid)', fontSize: '15px', lineHeight: 1.6, margin: '0 0 16px 0' }}>
                We rejected the standard `net/http` package. Instead, we use the `net` package to interact directly with TCP sockets. We manually implement parsing, routing, and connection management.
              </p>
              <div style={{ display: 'flex', gap: '16px', fontFamily: 'var(--font-roboto-mono)', fontSize: '12px' }}>
                <span style={{ color: '#4ade80' }}>+ Deep Educational Value</span>
                <span style={{ color: '#f87171' }}>- Increased Complexity</span>
              </div>
            </div>
          </div>

          {/* ADR 003 */}
          <div style={{
            background: 'rgba(20, 20, 22, 0.4)',
            border: '1px solid rgba(138, 43, 226, 0.2)',
            borderRadius: '12px',
            padding: '40px',
            display: 'flex',
            gap: '32px',
            boxShadow: '0 10px 30px rgba(0,0,0,0.5)',
            alignItems: 'flex-start'
          }}>
            <div style={{ color: 'var(--color-packet-violet)', background: 'rgba(138, 43, 226, 0.1)', padding: '16px', borderRadius: '12px' }}><GitBranch size={40} /></div>
            <div style={{ flex: 1 }}>
              <div style={{ fontFamily: 'var(--font-roboto-mono)', color: 'var(--color-packet-violet)', fontSize: '12px', letterSpacing: '0.1em', marginBottom: '8px' }}>ADR 003</div>
              <h3 style={{ fontFamily: 'var(--font-roboto-mono)', color: 'var(--color-paper-white)', fontSize: '20px', margin: '0 0 16px 0' }}>Radix Tree for Dynamic Routing</h3>
              <p style={{ fontFamily: 'var(--font-suisse-intl)', color: 'var(--color-steel-mid)', fontSize: '15px', lineHeight: 1.6, margin: '0 0 16px 0' }}>
                Replaced O(1) Hash Maps with an optimized Radix Tree to support dynamic path parameters (`/users/:id`) and wildcards without falling back to slow regular expressions.
              </p>
              <div style={{ display: 'flex', gap: '16px', fontFamily: 'var(--font-roboto-mono)', fontSize: '12px' }}>
                <span style={{ color: '#4ade80' }}>+ O(k) Pattern Matching</span>
                <span style={{ color: '#f87171' }}>- Harder to Debug</span>
              </div>
            </div>
          </div>

          {/* ADR 005 */}
          <div style={{
            background: 'rgba(20, 20, 22, 0.4)',
            border: '1px solid rgba(255, 255, 255, 0.1)',
            borderRadius: '12px',
            padding: '40px',
            display: 'flex',
            gap: '32px',
            boxShadow: '0 10px 30px rgba(0,0,0,0.5)',
            alignItems: 'flex-start'
          }}>
            <div style={{ color: 'var(--color-paper-white)', background: 'rgba(255, 255, 255, 0.05)', padding: '16px', borderRadius: '12px' }}><ServerCog size={40} /></div>
            <div style={{ flex: 1 }}>
              <div style={{ fontFamily: 'var(--font-roboto-mono)', color: 'var(--color-paper-white)', fontSize: '12px', letterSpacing: '0.1em', marginBottom: '8px' }}>ADR 005</div>
              <h3 style={{ fontFamily: 'var(--font-roboto-mono)', color: 'var(--color-paper-white)', fontSize: '20px', margin: '0 0 16px 0' }}>Bounded Worker Pool</h3>
              <p style={{ fontFamily: 'var(--font-suisse-intl)', color: 'var(--color-steel-mid)', fontSize: '15px', lineHeight: 1.6, margin: '0 0 16px 0' }}>
                Instead of unbounded goroutines, we spawn a fixed number of workers. If the queue is full, the server aggressively applies Load Shedding (`503 Service Unavailable`) to prevent catastrophic crashes.
              </p>
              <div style={{ display: 'flex', gap: '16px', fontFamily: 'var(--font-roboto-mono)', fontSize: '12px' }}>
                <span style={{ color: '#4ade80' }}>+ Predictable Memory</span>
                <span style={{ color: '#f87171' }}>- Drops Requests during Bursts</span>
              </div>
            </div>
          </div>

          {/* ADR 008 */}
          <div style={{
            background: 'rgba(20, 20, 22, 0.4)',
            border: '1px solid rgba(0, 217, 255, 0.2)',
            borderRadius: '12px',
            padding: '40px',
            display: 'flex',
            gap: '32px',
            boxShadow: '0 10px 30px rgba(0,0,0,0.5)',
            alignItems: 'flex-start'
          }}>
            <div style={{ color: 'var(--color-network-cyan)', background: 'rgba(0, 217, 255, 0.1)', padding: '16px', borderRadius: '12px' }}><Zap size={40} /></div>
            <div style={{ flex: 1 }}>
              <div style={{ fontFamily: 'var(--font-roboto-mono)', color: 'var(--color-network-cyan)', fontSize: '12px', letterSpacing: '0.1em', marginBottom: '8px' }}>ADR 008</div>
              <h3 style={{ fontFamily: 'var(--font-roboto-mono)', color: 'var(--color-paper-white)', fontSize: '20px', margin: '0 0 16px 0' }}>Lazy Evaluation Rate Limiting</h3>
              <p style={{ fontFamily: 'var(--font-suisse-intl)', color: 'var(--color-steel-mid)', fontSize: '15px', lineHeight: 1.6, margin: '0 0 16px 0' }}>
                Instead of locking the Token Bucket map every second with a background Ticker, tokens are calculated and refilled mathematically only when a request arrives: `(time.Now() - lastRefill) * rate`.
              </p>
              <div style={{ display: 'flex', gap: '16px', fontFamily: 'var(--font-roboto-mono)', fontSize: '12px' }}>
                <span style={{ color: '#4ade80' }}>+ O(1) CPU Usage</span>
                <span style={{ color: '#f87171' }}>- Complex Math Edge Cases</span>
              </div>
            </div>
          </div>

          {/* ADR 010 */}
          <div style={{
            background: 'rgba(20, 20, 22, 0.4)',
            border: '1px solid rgba(138, 43, 226, 0.2)',
            borderRadius: '12px',
            padding: '40px',
            display: 'flex',
            gap: '32px',
            boxShadow: '0 10px 30px rgba(0,0,0,0.5)',
            alignItems: 'flex-start'
          }}>
            <div style={{ color: 'var(--color-packet-violet)', background: 'rgba(138, 43, 226, 0.1)', padding: '16px', borderRadius: '12px' }}><GitBranch size={40} /></div>
            <div style={{ flex: 1 }}>
              <div style={{ fontFamily: 'var(--font-roboto-mono)', color: 'var(--color-packet-violet)', fontSize: '12px', letterSpacing: '0.1em', marginBottom: '8px' }}>ADR 010</div>
              <h3 style={{ fontFamily: 'var(--font-roboto-mono)', color: 'var(--color-paper-white)', fontSize: '20px', margin: '0 0 16px 0' }}>Native Documentation Routing</h3>
              <p style={{ fontFamily: 'var(--font-suisse-intl)', color: 'var(--color-steel-mid)', fontSize: '15px', lineHeight: 1.6, margin: '0 0 16px 0' }}>
                Instead of external GitHub links, we integrated react-router-dom to build native documentation pages. This extracts core concepts into premium, glassmorphic React UI cards.
              </p>
              <div style={{ display: 'flex', gap: '16px', fontFamily: 'var(--font-roboto-mono)', fontSize: '12px' }}>
                <span style={{ color: '#4ade80' }}>+ Unbroken User Immersion</span>
                <span style={{ color: '#f87171' }}>- Duplicates Markdown Content</span>
              </div>
            </div>
          </div>

          {/* ADR 011 */}
          <div style={{
            background: 'rgba(20, 20, 22, 0.4)',
            border: '1px solid rgba(255, 255, 255, 0.1)',
            borderRadius: '12px',
            padding: '40px',
            display: 'flex',
            gap: '32px',
            boxShadow: '0 10px 30px rgba(0,0,0,0.5)',
            alignItems: 'flex-start'
          }}>
            <div style={{ color: 'var(--color-paper-white)', background: 'rgba(255, 255, 255, 0.05)', padding: '16px', borderRadius: '12px' }}><ServerCog size={40} /></div>
            <div style={{ flex: 1 }}>
              <div style={{ fontFamily: 'var(--font-roboto-mono)', color: 'var(--color-paper-white)', fontSize: '12px', letterSpacing: '0.1em', marginBottom: '8px' }}>ADR 011</div>
              <h3 style={{ fontFamily: 'var(--font-roboto-mono)', color: 'var(--color-paper-white)', fontSize: '20px', margin: '0 0 16px 0' }}>Strict Bespoke Design System</h3>
              <p style={{ fontFamily: 'var(--font-suisse-intl)', color: 'var(--color-steel-mid)', fontSize: '15px', lineHeight: 1.6, margin: '0 0 16px 0' }}>
                We rejected Tailwind CSS and external UI libraries. Instead, we established a strict bespoke design system using pure Vanilla CSS governed by explicit color tokens and typography.
              </p>
              <div style={{ display: 'flex', gap: '16px', fontFamily: 'var(--font-roboto-mono)', fontSize: '12px' }}>
                <span style={{ color: '#4ade80' }}>+ Pixel-Perfect Control</span>
                <span style={{ color: '#f87171' }}>- Verbose Raw CSS</span>
              </div>
            </div>
          </div>

        </div>
      </div>
    </div>
  );
};
