# fractal
A small mandelbrot set renderer in Go

---

## Examples


*Zoom: 1, Center: -0.75 + 0i, Iterations: 100*

<img src="/examples/ex1-zoom-1.png" width="512">

*Zoom: 52.63, Center: -0.722 + 0.246i, Iterations: 250*

<img src="/examples/ex2-zoom-53.png" width="512">

*Zoom: 1.49e6, Center: 0.2929859127507 + 0.6117848324958i, Iterations: 500*

<img src="/examples/ex3-zoom-1.5e6.png" width="512">

*Zoom: 1e11, Center: 0.2549870375144766 - 0.0005679790528465i, Iterations: 1500*

<img src="/examples/ex4-zoom-1e11.png" width="512">

*Zoom: 4e12, Center: -1.99999911758738 + 0i, Iterations: 500*

<img src="/examples/ex5-zoom-4e12.png" width="512">

---

## Sources
Locations: http://www.cuug.ab.ca/dewara/mandelbrot/images.html

Coloring algorithm: http://linas.org/art-gallery/escape/escape.html

Linear mixing: https://github.com/ncruces/go-image/blob/v0.1.0/imageutil/srgb.go

Color model conversion: https://axonflux.com/handy-rgb-to-hsl-and-rgb-to-hsv-color-model-c
