import React, { useEffect } from 'react';
import { useLocation } from 'react-router-dom';
import { HeroSection } from '../components/HeroSection.tsx';
import { Chapter1Section } from '../components/Chapter1Section.tsx';
import { Chapter2Section } from '../components/Chapter2Section.tsx';
import { Chapter3Section } from '../components/Chapter3Section.tsx';
import { Chapter4Section } from '../components/Chapter4Section.tsx';
import { Chapter5Section } from '../components/Chapter5Section.tsx';
import { Chapter6Section } from '../components/Chapter6Section.tsx';
import { Chapter7Section } from '../components/Chapter7Section.tsx';
import { Chapter8Section } from '../components/Chapter8Section.tsx';
import { Chapter9Section } from '../components/Chapter9Section.tsx';
import { Chapter10Section } from '../components/Chapter10Section.tsx';

export const Home: React.FC = () => {
  const { hash } = useLocation();

  useEffect(() => {
    if (hash) {
      const element = document.querySelector(hash);
      if (element) {
        element.scrollIntoView({ behavior: 'smooth' });
      }
    } else {
      window.scrollTo(0, 0);
    }
  }, [hash]);

  return (
    <>
      <HeroSection />
      <Chapter1Section />
      <Chapter2Section />
      <Chapter3Section />
      <Chapter4Section />
      <Chapter5Section />
      <Chapter6Section />
      <Chapter7Section />
      <Chapter8Section />
      <Chapter9Section />
      <Chapter10Section />
    </>
  );
};
