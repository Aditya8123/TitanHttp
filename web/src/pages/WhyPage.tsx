import React, { useEffect } from 'react';
import { Target, ServerCrash, GitMerge, CheckCircle } from 'lucide-react';

export const WhyPage: React.FC = () => {
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
        <div style={{ fontFamily: 'var(--font-roboto-mono)', color: 'var(--color-network-cyan)', fontSize: '14px', letterSpacing: '0.1em', marginBottom: '16px' }}>&gt;_ THE MOTIVATION</div>
        <h1 style={{ fontFamily: 'var(--font-lambotype)', fontSize: 'clamp(40px, 6vw, 80px)', color: 'var(--color-paper-white)', margin: '0 0 24px 0', textTransform: 'uppercase' }}>
          Why build an HTTP Server?
        </h1>
        <p style={{ fontFamily: 'var(--font-suisse-intl)', fontSize: '18px', color: 'var(--color-steel-mid)', maxWidth: '800px', margin: '0 auto', lineHeight: 1.6 }}>
          I built TitanHTTP because I wanted to prove that I don't just know how to use tools—I know how the tools work. Most developers build backend applications using frameworks or high-level standard libraries. While practical for business, it obscures the complex engineering happening underneath.
        </p>
      </div>

      {/* Grid of Challenges */}
      <div>
        <div style={{ fontFamily: 'var(--font-roboto-mono)', color: 'var(--color-paper-white)', fontSize: '20px', letterSpacing: '0.1em', marginBottom: '32px', borderBottom: '1px solid rgba(255,255,255,0.1)', paddingBottom: '16px' }}>
          THE HARDEST TECHNICAL CHALLENGES
        </div>
        
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(300px, 1fr))', gap: '32px' }}>
          
          {/* Card 1 */}
          <div style={{
            background: 'rgba(20, 20, 22, 0.6)',
            border: '1px solid rgba(0, 217, 255, 0.2)',
            borderRadius: '12px',
            padding: '32px',
            display: 'flex',
            flexDirection: 'column',
            gap: '16px',
            boxShadow: '0 10px 30px rgba(0,0,0,0.5)',
            transition: 'transform 0.3s ease',
          }}>
            <div style={{ color: 'var(--color-network-cyan)' }}><GitMerge size={32} /></div>
            <h3 style={{ fontFamily: 'var(--font-roboto-mono)', color: 'var(--color-paper-white)', fontSize: '16px', letterSpacing: '0.05em', margin: 0 }}>The HTTP Parser (Memory)</h3>
            <p style={{ fontFamily: 'var(--font-suisse-intl)', color: 'var(--color-steel-mid)', fontSize: '14px', lineHeight: 1.6, margin: 0 }}>
              When bytes arrive over a TCP socket, they stream in unpredictably. Writing a parser that can read a stream, identify the \r\n\r\n boundary, and extract headers without causing a massive garbage-collection spike required deep knowledge of Go's byte-slice manipulation and sync.Pool.
            </p>
          </div>

          {/* Card 2 */}
          <div style={{
            background: 'rgba(20, 20, 22, 0.6)',
            border: '1px solid rgba(138, 43, 226, 0.2)',
            borderRadius: '12px',
            padding: '32px',
            display: 'flex',
            flexDirection: 'column',
            gap: '16px',
            boxShadow: '0 10px 30px rgba(0,0,0,0.5)',
          }}>
            <div style={{ color: 'var(--color-packet-violet)' }}><Target size={32} /></div>
            <h3 style={{ fontFamily: 'var(--font-roboto-mono)', color: 'var(--color-paper-white)', fontSize: '16px', letterSpacing: '0.05em', margin: 0 }}>Concurrency (Worker Pool)</h3>
            <p style={{ fontFamily: 'var(--font-suisse-intl)', color: 'var(--color-steel-mid)', fontSize: '14px', lineHeight: 1.6, margin: 0 }}>
              Handling 10,000 concurrent requests without crashing is hard. Instead of blindly spawning a new goroutine for every connection, TitanHTTP implements a bounded Worker Pool. Connections are placed in a non-blocking queue and processed by a fixed number of workers, ensuring stable memory.
            </p>
          </div>

          {/* Card 3 */}
          <div style={{
            background: 'rgba(20, 20, 22, 0.6)',
            border: '1px solid rgba(255, 255, 255, 0.1)',
            borderRadius: '12px',
            padding: '32px',
            display: 'flex',
            flexDirection: 'column',
            gap: '16px',
            boxShadow: '0 10px 30px rgba(0,0,0,0.5)',
          }}>
            <div style={{ color: 'var(--color-paper-white)' }}><ServerCrash size={32} /></div>
            <h3 style={{ fontFamily: 'var(--font-roboto-mono)', color: 'var(--color-paper-white)', fontSize: '16px', letterSpacing: '0.05em', margin: 0 }}>Graceful Degradation</h3>
            <p style={{ fontFamily: 'var(--font-suisse-intl)', color: 'var(--color-steel-mid)', fontSize: '14px', lineHeight: 1.6, margin: 0 }}>
              The server is designed to be deterministic. If a client sends malformed headers, it safely drains the socket and returns a 400 Bad Request. If a handler panics, a recovery middleware catches the panic, logs the stack trace, and returns a 500 without killing the listener.
            </p>
          </div>

        </div>
      </div>
      
      {/* Conclusion */}
      <div style={{ display: 'flex', alignItems: 'center', gap: '24px', background: 'rgba(0, 217, 255, 0.05)', border: '1px solid rgba(0, 217, 255, 0.2)', padding: '32px', borderRadius: '12px' }}>
        <CheckCircle size={48} color="var(--color-network-cyan)" />
        <div>
          <h4 style={{ fontFamily: 'var(--font-roboto-mono)', color: 'var(--color-paper-white)', fontSize: '14px', letterSpacing: '0.1em', margin: '0 0 8px 0' }}>THE RESULT</h4>
          <p style={{ fontFamily: 'var(--font-suisse-intl)', color: 'var(--color-steel-mid)', fontSize: '14px', lineHeight: 1.6, margin: 0 }}>
            TitanHTTP is not just a toy project. It is a benchmarked, tested, and structurally sound piece of infrastructure. The codebase is strictly organized, heavily documented, and adheres to idiomatic Go standards.
          </p>
        </div>
      </div>

    </div>
  );
};
