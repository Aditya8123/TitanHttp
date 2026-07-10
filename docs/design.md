# design.md

# TitanHTTP — Style Reference

> A cinematic staging of raw infrastructure. Visualizing the invisible.

**Theme:** dark / immersive / educational

TitanHTTP presents backend infrastructure as a cinematic, interactive engineering experience. It combines premium editorial documentation, immersive storytelling, and real-time system visualization into a single cohesive design language.

The interface treats backend architecture like a supercar unveiling: an aggressive, cinematic pure-black canvas where data packets move like light through a datacenter. Display typography is industrial and massive (120px) for high-impact hero storytelling, seamlessly transitioning into editorial grace and premium whitespace for documentation.

## Design Goals

Every page and interaction within TitanHTTP must satisfy at least one of these goals:

* Teach a backend engineering concept.
* Visualize invisible systems.
* Demonstrate engineering quality.
* Encourage exploration.
* Build recruiter confidence.
* Reward curiosity.
* Feel premium without unnecessary decoration.

## Design Principles

TitanHTTP is not just a UI; it is an educational journey. The interface should feel **calm, confident, technical, premium, educational, minimal,** and **intentional.**

* **Invisible systems become visible.**
* **Motion always has meaning.**
* **The interface never competes with the content.**
* **Backend engineering deserves beautiful presentation.**
* **Documentation is a first-class experience.**
* **Complexity is revealed progressively.**
* **Data is more important than decoration.**
* **Consistency beats novelty.**

## The Recruiter Journey

A first-time visitor should experience a clear path of discovery:

* **30 Seconds:** They should understand exactly what TitanHTTP is.
* **2 Minutes:** They should understand the high-level architecture.
* **5 Minutes:** They should understand *why* it was built and the engineering challenges solved.
* **10 Minutes:** They should trust the engineer who built it.

## Emotion Timeline

The user experience is mapped to a strict emotional progression:

Landing
↓
Curiosity
↓
Understanding
↓
Exploration
↓
Confidence
↓
Admiration
↓
"I want to read the code."

## Experience Chapters

The website is a guided narrative rather than a collection of pages:

* **Chapter 1:** Why HTTP Exists
* **Chapter 2:** TCP
* **Chapter 3:** Building a Socket
* **Chapter 4:** Reading Bytes
* **Chapter 5:** Parsing Requests
* **Chapter 6:** Routing
* **Chapter 7:** Concurrency
* **Chapter 8:** Production Features
* **Chapter 9:** Benchmarks
* **Chapter 10:** Source Code

## Tokens — Colors

| Name | Value | Token | Role |
| --- | --- | --- | --- |
| Void Black | `#000000` | `--color-void-black` | Universal page canvas, cinematic hero background, and absolute dark. The depth of the datacenter. |
| Obsidian | `#0f1011` | `--color-obsidian` | Base document background for editorial documentation sections. |
| Carbon | `#090a0b` | `--color-carbon` | Elevated container surfaces for code blocks and terminal windows. |
| Frosted Graphite | `rgba(46,46,46,0.3)` | `--color-frosted-graphite` | Translucent substrate for floating metrics and dashboard cards. Used exclusively with background-blur for high-end glassmorphism. |
| Bone | `#f5f5f7` | `--color-bone` | Hairline borders, dividers, and secondary text on light documentation surfaces. |
| Paper White | `#ffffff` | `--color-paper-white` | Primary display text, editorial copy, and filled CTA pills. |
| Steel Mid | `#7d7d7d` | `--color-steel-mid` | Monospace UI labels, inactive nodes, and muted chrome. |
| Network Cyan | `#00d9ff` | `--color-network-cyan` | The single functional accent. Used exclusively for traveling packets, active architecture nodes, graphs, and the request path. Never used for buttons. |
| Packet Violet | `#847dff` | `--color-packet-violet` | Secondary data accent for contrasting request flows, multiplexed connections, or concurrent worker visualization. |

## Tokens — Typography

### LamboType — Hero and cinematic chapter titles. Renders strictly in UPPERCASE with a consistent 0.0230em letter-spacing. Display sizes hit 120px with tight 0.92 line-height to create monumental, engineered text blocks. · `--font-lambotype`

* **Substitute:** Barlow Condensed or Bebas Neue
* **Weights:** 400
* **Sizes:** 54px, 80px, 120px
* **Line height:** 0.92, 1.00
* **Role:** Cinematic openers and massive chapter markers ("TCP IS RELIABLE").

### Lyon Display — Editorial headers for the documentation and premium SaaS sections. The weight 300 whispers authority and slows the user down to read technical specs with the reverence of a definitive engineering paper. · `--font-lyon-display`

* **Substitute:** GT Sectra Display or Playfair Display
* **Weights:** 300, 400
* **Sizes:** 38px, 80px
* **Line height:** 0.90, 1.00
* **Role:** Section openers in the documentation and dashboard views.

### Suisse Int'l — Sans-serif workhorse for body, subheadings, and UI copy. Geometric humanist proportions provide high readability against deep black surfaces without competing with the display fonts. · `--font-suisse-intl`

* **Substitute:** Inter or Söhne
* **Weights:** 300, 400
* **Sizes:** 14px, 16px, 18px
* **Line height:** 1.50, 1.60
* **Role:** Main documentation body copy and descriptive storytelling.

### JetBrains Mono — Authentic engineering familiarity for all code blocks, syntax highlighting, and raw HTTP request payloads. · `--font-jetbrains-mono`

* **Weights:** 400, 500
* **Sizes:** 14px, 16px
* **Line height:** 1.60
* **Role:** Code blocks, JSON payloads, and raw technical outputs.

### Roboto Mono — Monospace for floating metrics, stamped UI labels, and dashboard data. Wide-tracked (0.182em) for micro-labels to feel pressed into the glass. · `--font-roboto-mono`

* **Weights:** 400, 500
* **Sizes:** 10px, 12px, 16px, 24px
* **Letter spacing:** 0.021em (data), 0.182em (labels)
* **Role:** Dashboard latency numbers, memory usage, CPU, worker stats, and chart axes.

## Tokens — Spacing & Shapes

**Base unit:** 8px
**Density:** Cinematic / Generous

### Spacing Scale

| Name | Value | Token |
| --- | --- | --- |
| 8 | 8px | `--spacing-8` |
| 16 | 16px | `--spacing-16` |
| 24 | 24px | `--spacing-24` |
| 40 | 40px | `--spacing-40` |
| 80 | 80px | `--spacing-80` |
| 160 | 160px | `--spacing-160` |

### Border Radius

| Element | Value |
| --- | --- |
| terminal-windows | 0px |
| nav-items | 2px |
| glass-panels | 16px |
| cta-buttons | 1440px |

## Tokens — Motion & Sound

Motion is the core educational tool of TitanHTTP. Nothing simply appears; everything is built, routed, or parsed.

### Motion Philosophy

* Motion should reveal architecture.
* Never distract.
* Never loop forever; animations stop after teaching.
* Progress always moves forward.
* Packets never teleport.
* Data never fades randomly.
* State changes should be explainable.
* **No elastic easing.** (Backend engineering is deterministic, not bouncy).
* **3D as Functional Depth:** 3D elements are used to separate layers of architecture (e.g., physical infrastructure vs. logical routing) or to demonstrate state transitions, never just for visual flair.

### 3D Transitions & Animations (WebGL / Three.js)

* **Volumetric Parallax:** Scroll-linked WebGL canvases provide true spatial depth. Foreground UI (glassmorphism) floats above deep, volumetric point-cloud representations of data structures.
* **Cinematic Camera Choreography:** Camera movements (Z-axis plunging, slow isometric panning) must feel like a heavy, precision drone moving through a server farm. Use easing functions that mimic physical mass, avoiding lightweight CSS "spring" physics.
* **Exploded Architecture (Z-Space):** Deconstruct HTTP requests by exploding headers, method, and body along the Z-axis, allowing the user's camera to fly *through* the parsed data.
* **Shader-Driven State Morphing:** Transition states using custom GLSL shaders (e.g., noise-driven particle dispersal when a connection drops, or a light-sweep when a packet is successfully routed).

### Durations

| Name | Value | Token | Role |
| --- | --- | --- | --- |
| fast | 250ms | `--motion-250` | Hover states, micro-interactions, button transitions |
| base | 400ms | `--motion-400` | Component expansions, terminal typing speed per line |
| slow | 800ms | `--motion-800` | Packet flow across a section, architecture diagram highlighting |
| cinematic | 1600ms | `--motion-1600` | Full section unveils, hero text fade-ins, 3D camera pans |

### Sound

* No audio.
* No autoplay sound.
* No gimmicks.
* The experience should feel calm and focused.

## Accessibility

Demonstrating engineering maturity requires building for everyone.

* Minimum AA contrast on all text surfaces.
* Full keyboard navigation.
* Reduced motion support (animations respect `prefers-reduced-motion`).
* Screen reader labels on all interactive nodes.
* Every visualization has a textual explanation.

## Experience Hierarchy (Experience → Section → Component)

TitanHTTP is organized not by isolated components, but by narrative flows.

### Example: The Request Journey Experience

**↓ Section: TCP Handshake**

* **Component:** Timeline Canvas (Pure Black, 100vw).
* **Component:** SYN / ACK Packets (Network Cyan, moving 800ms).
* **Component:** Socket Node (Glowing interaction point).

**↓ Section: The Parser**

* **Component:** Full-bleed dark terminal (0px radius).
* **Component:** Typewriter Animation (JetBrains Mono).
* **Component:** Header Expansion UI (Highlights headers vs body in real-time).

**↓ Section: Production Telemetry**

* **Component:** Frosted Glass Dashboard (16px radius, heavy blur).
* **Component:** Live Metric Nodes (Roboto Mono, counting up).
* **Component:** Data Visualization Charts.

## Components

### Cinematic Hero Stage (3D)

Pure Black HTML canvas overlaying a deep WebGL scene. A slow-moving, volumetric particle system simulates fiber-optic data flow in Z-space, vanishing into infinite darkness. Subtle depth-of-field (bokeh) blurs the deepest nodes. No navbar initially. A single sentence drops in, stark and massive in LamboType 120px:

**THE INTERNET**
**STARTS**
**WITH A REQUEST**

A glowing 3D Network Cyan sphere (the packet) emerges from the depth of field, piercing the 2D plane of the screen, and begins a scroll-linked plunge downward into the architecture (`Browser ↓ TCP ↓ Socket ↓ TitanHTTP`).

### Floating Metrics Dashboard

High-end, premium data visualization. Built exclusively with frosted glass (`rgba(46,46,46,0.3)`) and a strong background-blur. This allows the dark, cinematic server-rack background to bleed softly through the panels, grounding the precise Roboto Mono telemetry data in physical infrastructure.

### Primary CTA Pill

Paper White (`#ffffff`) fill, Pure Black text, 1440px radius (full pill). Trailing arrow glyph. Sits starkly against the black void as the only true click action.

### Ghost Nav Button

Transparent background, 2px border radius, Paper White text, Roboto Mono 12px with 0.182em tracking. 1px Bone border on hover. Reads as quiet system instrumentation.

## Data Visualization

Graphs must feel like high-performance engineering tooling.

* **Latency & Response Times:** Rendered as jagged, precise area charts. No soft bezier curves.
* **Connections & Workers:** Visualized as glowing nodes or dots on a scatter-plot timeline.
* **Throughput (Memory/CPU):** Live-updating bar charts or sparklines.
* **Palette:** Network Cyan for primary metric, Packet Violet for secondary/comparative metric. Bone and Steel Mid for axes and gridlines.

## Illustrations & Storytelling Rules

### Storytelling Rules

1. Every screen answers one question.
2. Every animation teaches one concept.
3. Every interaction reveals architecture.
4. Nothing moves without purpose.
5. Every graph tells a story.
6. Every metric has context.
7. Every section ends with curiosity.

### Image & 3D Style (Premium WebGL)

Imagery is not static; it is rendered, volumetric, and alive.

**Only use:**

* **Dark Mode WebGL:** Real-time rendered scenes using Three.js or custom shaders. The background is `#000000`. Lighting should be low-key, utilizing rim lights and emissive materials (specifically Network Cyan and Packet Violet) to define geometry.
* **Point Clouds & Particle Systems:** Represent bytes, streams, and raw data using millions of particles. Use GLSL shaders to flow them along bezier splines mimicking copper traces.
* **Volumetric Grids:** Infinite, fading wireframe floors to ground isometric architectures, providing scale and perspective.
* **Glass Shaders:** 3D nodes should use physical-based rendering (PBR) glass materials with high index of refraction (IOR) to distort the light passing behind them.
* Monochromatic, hardware-accelerated architecture diagrams.

**Strictly Prohibited:**

* No stock photos.
* No low-poly or "cartoony" 3D models.
* No flat vector illustrations of people (e.g., "Corporate Memphis").
* No smiling developers or office environments.
* No generic coding/laptop illustrations.

## Do's and Don'ts

### Do

* **Do** treat the page like a continuous film strip: users follow a single packet through chapters.
* **Do** limit glassmorphism strictly to where there is *live information* (Dashboards, metrics, floating controls, navigation). Never use glass for articles or documentation.
* **Do** use Network Cyan (`#00D9FF`) *only* for data: traveling packets, active architecture nodes, and graphs.
* **Do** format documentation with generous whitespace and Lyon Display serifs for editorial readability.
* **Do** strictly control motion: linear and smooth, never bouncy or elastic.

### Don't

* **Don't** use lifestyle or finance-specific language (e.g., Budget, Money). Use engineering terms (Requests, Packets, Worker Pool).
* **Don't** use atmospheric sky/cloud photography. Replace entirely with datacenters and fiber optics.
* **Don't** use phone mockups. Showcase the product via architecture diagrams, terminal output, and benchmark graphs.
* **Don't** use Network Cyan for buttons or marketing CTAs. Use stark white pills to maintain gallery-grade restraint.

## Agent Prompt Guide

### Quick Reference

* **Page Canvas:** #000000 (Void Black)
* **Doc Background:** #0f1011 (Obsidian)
* **Glass Panel:** rgba(46,46,46,0.3) with backdrop-blur
* **Primary Text:** #ffffff (Paper White)
* **Data Highlight:** #00d9ff (Network Cyan)
* **Primary Action:** #ffffff (filled pill)

### Example Component Prompts

1. **WebGL Architecture Node:** Initialize a Three.js scene on a pure black background. Create an array of translucent dark glass cubes representing the worker pool. Inject an emissive `#00d9ff` point light (the packet) that travels along a defined 3D spline from the router node to a worker node over 1600ms. Apply a post-processing bloom pass to make the cyan light streak.
2. **Dashboard Metric Panel:** Create a CSS glassmorphic card layered *above* the WebGL canvas. Background `rgba(46,46,46,0.3)`, backdrop filter blur 16px, border radius 16px. Top left label "TCP CONNECTIONS" in Roboto Mono 10px, uppercase, tracking 0.182em, color `#9f9fa0`. Center value "1,024" in Roboto Mono 24px, color `#ffffff`.
3. **Scroll-Driven Exploded View:** On scroll, transition a 3D server rack model into an exploded view along the Y-axis. Fade in HTML labels (`DOM` overlay) that track the 3D coordinates of the exposed motherboard components, connecting to them via 1px solid `#7d7d7d` SVG lines.

## Quick Start

### CSS Custom Properties

```css
:root {
  /* Colors */
  --color-void-black: #000000;
  --color-obsidian: #0f1011;
  --color-carbon: #090a0b;
  --color-frosted-graphite: rgba(46, 46, 46, 0.3);
  --color-bone: #f5f5f7;
  --color-paper-white: #ffffff;
  --color-steel-mid: #7d7d7d;
  --color-network-cyan: #00d9ff;
  --color-packet-violet: #847dff;

  /* Typography */
  --font-lambotype: 'LamboType', 'Barlow Condensed', sans-serif;
  --font-lyon-display: 'Lyon Display', 'Playfair Display', serif;
  --font-suisse-intl: 'Suisse Int'l', 'Inter', sans-serif;
  --font-jetbrains-mono: 'JetBrains Mono', monospace;
  --font-roboto-mono: 'Roboto Mono', monospace;

  /* Motion */
  --motion-fast: 250ms;
  --motion-base: 400ms;
  --motion-slow: 800ms;
  --motion-cinematic: 1600ms;
  --ease-mechanical: cubic-bezier(0.4, 0.0, 0.2, 1);
  --ease-linear: linear;

  /* Spacing & Layout */
  --spacing-8: 8px;
  --spacing-16: 16px;
  --spacing-24: 24px;
  --spacing-40: 40px;
  --spacing-80: 80px;
  --spacing-160: 160px;
  
  /* Border Radius */
  --radius-none: 0px;
  --radius-sm: 2px;
  --radius-glass: 16px;
  --radius-pill: 1440px;
}

```

### Tailwind v4

```css
@theme {
  /* Colors */
  --color-void-black: #000000;
  --color-obsidian: #0f1011;
  --color-carbon: #090a0b;
  --color-frosted-graphite: rgba(46, 46, 46, 0.3);
  --color-bone: #f5f5f7;
  --color-paper-white: #ffffff;
  --color-steel-mid: #7d7d7d;
  --color-network-cyan: #00d9ff;
  --color-packet-violet: #847dff;

  /* Typography */
  --font-lambotype: 'LamboType', 'Barlow Condensed', sans-serif;
  --font-lyon-display: 'Lyon Display', 'Playfair Display', serif;
  --font-suisse-intl: 'Suisse Int'l', 'Inter', sans-serif;
  --font-jetbrains-mono: 'JetBrains Mono', monospace;
  --font-roboto-mono: 'Roboto Mono', monospace;

  /* Motion */
  --animate-fast: 250ms;
  --animate-base: 400ms;
  --animate-slow: 800ms;
  --animate-cinematic: 1600ms;

  /* Spacing */
  --spacing-8: 8px;
  --spacing-16: 16px;
  --spacing-24: 24px;
  --spacing-40: 40px;
  --spacing-80: 80px;
  --spacing-160: 160px;

  /* Border Radius */
  --radius-none: 0px;
  --radius-sm: 2px;
  --radius-glass: 16px;
  --radius-pill: 1440px;
}

```