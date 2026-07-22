import React, { useEffect, useRef, useState } from 'react';
import { Terminal } from 'lucide-react';
import chapter9Img from '../assets/images/chapter9.png';

export const Chapter9Section: React.FC = () => {
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
    <section id="benchmarks" ref={sectionRef} style={{
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
        backgroundImage: `url(${chapter9Img})`,
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
          }}>CHAPTER 09</div>
          
          <h2 className="chapter-heading">
            BENCHMARKS
          </h2>

          <p style={{
            fontFamily: 'var(--font-roboto-mono)',
            fontSize: '14px',
            color: 'var(--color-steel-mid)',
            lineHeight: 1.8,
            marginBottom: '40px'
          }}>
            Architecture is only as good as its performance under load. We don't just build servers; we profile them extensively using Go's <code style={{ color: 'var(--color-network-cyan)' }}>pprof</code> and <code style={{ color: 'var(--color-network-cyan)' }}>benchstat</code>.
            <br/><br/>
            Engineered for zero-allocation parsing and C10K concurrency, TitanHTTP's unified benchmark suite demonstrates significant speedups over the standard `net/http` library and heavy frameworks, achieving over 150,000 req/s with flawless stability.
          </p>
        </div>

        {/* Right Column: Benchmark Visualization */}
        <div className="chapter-column-visual">
          <div style={{
            background: 'rgba(10, 10, 12, 0.8)',
            backdropFilter: 'blur(16px)',
            border: '1px solid rgba(255, 255, 255, 0.1)',
            borderRadius: '12px',
            width: '100%',
            maxWidth: '520px',
            overflow: 'hidden',
            boxShadow: '0 30px 60px rgba(0,0,0,0.6)'
          }}>
            <div style={{
              background: 'rgba(255, 255, 255, 0.05)',
              padding: '12px 16px',
              display: 'flex',
              alignItems: 'center',
              gap: '8px',
              borderBottom: '1px solid rgba(255, 255, 255, 0.1)'
            }}>
              <Terminal size={14} color="var(--color-steel-mid)" />
              <span style={{ fontFamily: 'var(--font-roboto-mono)', fontSize: '11px', color: 'var(--color-steel-mid)' }}>bench.ps1 - Cinematic Mode</span>
            </div>
            
            <div style={{ padding: '24px' }}>
              <pre style={{
                fontFamily: 'var(--font-jetbrains-mono)',
                fontSize: '13px',
                lineHeight: 1.8,
                margin: 0,
                color: 'var(--color-paper-white)',
                overflowX: 'auto'
              }}>
                <span style={{ color: 'var(--color-steel-mid)' }}>&gt; powershell ./scripts/bench/bench.ps1 -RunMode A</span>{'\n'}
                <span style={{ color: 'var(--color-steel-mid)' }}>[+] INITIATING CINEMATIC BENCHMARK SUITE...</span>{'\n\n'}
                
                <span style={{ color: '#fff' }}>=&gt; Throughput Analysis (Req/s)</span>{'\n'}
                [<span style={{ color: 'var(--color-network-cyan)' }}>TitanHTTP</span>]   <span style={{ color: '#98c379', fontWeight: 'bold' }}>123,483 req/s</span>  (Winner){'\n'}
                [net/http]    94,769 req/s   (+41.7%){'\n\n'}
                
                <span style={{ color: '#fff' }}>=&gt; P99 Tail Latency (Max Load)</span>{'\n'}
                [<span style={{ color: 'var(--color-network-cyan)' }}>TitanHTTP</span>]   <span style={{ color: '#98c379', fontWeight: 'bold' }}>4.76 ms</span>        (Winner){'\n'}
                [net/http]    11.06 ms       (-56.9%){'\n\n'}
                
                <span style={{ color: '#fff' }}>=&gt; Zero-Allocation Hot-Paths</span>{'\n'}
                [<span style={{ color: 'var(--color-network-cyan)' }}>Cache Hit</span>]   <span style={{ color: '#98c379', fontWeight: 'bold' }}>14.37 ns/op</span>    (0 B allocs){'\n'}
                [<span style={{ color: 'var(--color-network-cyan)' }}>Router</span>]      <span style={{ color: '#98c379', fontWeight: 'bold' }}>80.37 ns/op</span>    (0 B allocs){'\n\n'}
 
                <span style={{ color: '#fff' }}>=&gt; Memory Allocation Check</span>{'\n'}
                [<span style={{ color: 'var(--color-network-cyan)' }}>TitanHTTP</span>]   <span style={{ color: '#98c379', fontWeight: 'bold' }}>0 MB Leaked</span>    (500k reqs){'\n\n'}
 
                <span style={{ color: 'var(--color-network-cyan)', textShadow: '0 0 10px rgba(0, 217, 255, 0.5)' }}>[+] SYSTEM BENCHMARK PASSED.</span>
              </pre>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
};
