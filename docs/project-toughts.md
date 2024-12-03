# Intro

Hi everyone!
This file was created by [@gucio321](https://github.com/gucio321) and it was not a part of
the original prject. I just wanted to share my thoughts about this project, what am I planning to do here
(or rather what should be done here if I had more time - I don't believe I'll be able to do it myself).

# HellSpawner & AbyssEngine

My idea is to abandon the initial idea of AbyssEngine.
In my opinion It'd be better to do a tuned, step-by-step migration from current hard-coded OpenDiablo2 to the
desired state where the engine is fully configurable by HellSpawner.

# Roadmap

In this section I'm going to discuss overall what should be done on hellspawner and od2 in a general steps.

- [ ] Clean up HellSpawner a bit (this project died when there was giu v0.4.x (or something) and now it is giu v0.11.0 - we did much on giu-side since that).
      Especially things like hardcoded IDs and so on could be removed from this project (this will remove hundreds of code lines)
- [ ] Clean up OpenDiablo2. This project also was left long time ago so ebiten had several releases since that.
- [ ] Create a project decoder and extend the current Project type.
- [ ] Rewrite things from od2 to the new project decoder.
