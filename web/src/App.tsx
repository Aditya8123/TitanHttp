import { useEffect } from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import Lenis from 'lenis';
import { NavBar } from './components/NavBar.tsx';
import { Home } from './pages/Home.tsx';
import { WhyPage } from './pages/WhyPage.tsx';
import { ArchitecturePage } from './pages/ArchitecturePage.tsx';
import { DecisionsPage } from './pages/DecisionsPage.tsx';

function App() {
  useEffect(() => {
    const lenis = new Lenis({
      duration: 1.2,
      easing: (t) => Math.min(1, 1.001 - Math.pow(2, -10 * t)),
      orientation: 'vertical',
      gestureOrientation: 'vertical',
      smoothWheel: true,
      wheelMultiplier: 1,
      touchMultiplier: 2,
    });

    function raf(time: number) {
      lenis.raf(time);
      requestAnimationFrame(raf);
    }

    requestAnimationFrame(raf);

    return () => lenis.destroy();
  }, []);

  return (
    <BrowserRouter basename={import.meta.env.BASE_URL}>
      <div className="app-container" style={{ backgroundColor: 'var(--color-void-black)', color: 'var(--color-paper-white)', minHeight: '100vh', overflowX: 'clip' }}>
        <NavBar />
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/why" element={<WhyPage />} />
          <Route path="/architecture" element={<ArchitecturePage />} />
          <Route path="/decisions" element={<DecisionsPage />} />
        </Routes>
      </div>
    </BrowserRouter>
  );
}

export default App;
