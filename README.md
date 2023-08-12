# fractals

**fractals** is a customizable renderer for the Mandelbrot set written in Go. It uses Go's **goroutines** to achieve high performance.

### 🚀 Featured in [Golang Weekly #464](https://golangweekly.com/issues/464) 🚀 

## Usage
```sh
git clone https://github.com/joweich/fractals.git
cd fractals
go build 
./fractals -h  # to see list of available customizations
./fractals -height 1000 -width 1000 # fractals.exe for Windows systems
```

## Examples
<table>
  <tr>
    <td>
      <img src="/examples/ex1-zoom-1.png" width="350">
    </td>
    <td>
      <img src="/examples/ex2-zoom-53.png" width="350">
    </td>
  </tr>
  <tr>
    <td>
      <img src="/examples/ex5-zoom-4e12.png" width="350">
    </td>
    <td>
      <img src="/examples/ex4-zoom-1e11.png" width="350">
    </td>
  </tr>
</table>

## About the Algorithm
### The Math in a Nutshell
The Mandelbrot set is defined as the set of complex numbers $z_0$ for which the series 

$$z_{n+1} = z²_n + z_0$$

is bounded for all $n ≥ 0$. In other words, $z_0$ is part of the Mandelbrot set if $z_n$ does not approach infinity. This is equivalent to the  magnitude $|z_n| ≤ 2$ for all $n ≥ 0$.

### But how is this visualized in a colorful image?
The image is interpreted as complex plane, i.e. the horizontal axis being the real part and the vertical axis representing the complex part of $z_0$. 

The colors are determined by the so-called **naïve escape time algorithm**. It's as simple as that: A pixel is painted in a predefined color (often black) if it's in the set and will have a color if it's not. The color is determined by the number of iterations $n$ needed for $z_n$ to exceed $|z_n| = 2$. This $n$ is the escape time, and $|z_n| ≥ 2$ is the escape condition. In our implementation, this is done via the _hue_ parameter in the [HSL color model](https://en.wikipedia.org/wiki/HSL_and_HSV).

### And how does it leverage Goroutines?
Each row of the image is added as a job to a [channel](https://go.dev/doc/effective_go#channels). These jobs are distributed using [goroutines](https://go.dev/doc/effective_go#goroutines) (lightweight threads managed by the Go runtime) that are spun off by consuming from the channel until it's empty.

## Advanced Rendering Features
* Linear color mixing ([source](https://github.com/ncruces/go-image/blob/v0.1.0/imageutil/srgb.go))
* Anti-aliasing by random sampling ([source](https://www.fractalus.com/info/antialias.htm))
* _Normative iteration count_ to smooth stair-step function ([math behind](http://linas.org/art-gallery/escape/escape.html))
