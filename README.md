#  Performative ASCII — Real-Time Generative Visual Art Engine

<p align="center">
  <img src="https://img.shields.io/badge/language-Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" />
  <img src="https://img.shields.io/badge/platform-Terminal-black?style=for-the-badge&logo=windowsterminal&logoColor=white" />
  <img src="https://img.shields.io/badge/type-Generative%20Art-blueviolet?style=for-the-badge" />
  <img src="https://img.shields.io/badge/render-ASCII%20%2B%20ANSI-green?style=for-the-badge" />
</p>

> **Seni performatif digital yang hidup di dalam terminal.** Sebuah mesin visual generatif real-time berbasis ASCII art dengan 13 layer visual yang dapat dimainkan secara interaktif — dirancang untuk pertunjukan seni digital, live coding, dan eksplorasi estetika komputasional.

---

## Tentang Proyek

**Performative ASCII** adalah sebuah aplikasi terminal yang mengubah konsol teks biasa menjadi kanvas seni generatif yang dinamis. Ditulis sepenuhnya dalam **Go**, proyek ini merender animasi prosedural secara real-time (~25 FPS) menggunakan karakter ASCII dan escape sequence ANSI untuk warna.

Proyek ini terinspirasi dari seni performatif, ritual, dan estetika glitch — menggabungkan matematika (gelombang sinus, noise procedural, persamaan spiral) dengan simbol-simbol budaya.

---
*Salah satu layer: Topeng Sakral — ASCII art topeng tradisional dengan efek breathing dan glowing eyes.*

---

## Layer Visual (13 Sequences)

| # | Nama Layer | Deskripsi | Warna |
|---|-----------|-----------|-------|
| 1 | **Portal Looping** | Cincin konsentris yang berputar dengan efek spiral dan center glow | 🔵 Cyan |
| 2 | **Ball Matrix** | Partikel bola yang memantul dengan trail dan physics simulation | 🟡 Yellow |
| 3 | **Topeng Sakral** | ASCII art topeng tradisional dengan breathing effect dan glowing eyes | 🔴 Red |
| 4 | **Matrix Data Typist** | Hujan karakter ala Matrix dengan huruf Latin, angka, dan Katakana | 🟢 Green |
| 5 | **Tubuh Yang Mengingat** | Pola gelombang organik — visualisasi memori tubuh | 🟠 Orange |
| 6 | **3D Mountain 360** | Terrain 3D wireframe dengan perspektif, atmospheric haze, dan panning camera | 🔵 Electric Blue |
| 7 | **Barcode Glitch** | Barcode vertikal dengan efek scan line, displacement, dan data corruption | ⚪ White |
| 8 | **Wave Interference** | Interferensi multi-gelombang sinusoidal | 🔵 Blue |
| 9 | **Hyper Spiral** | Spiral multi-arm hipnotik dengan fungsi tanh | 🟣 Magenta |
| 10 | **Void Abyss** | Lubang hitam dengan accretion disk dan ripple | 🟡 Gold |
| 11 | **Glitch Storm** | Block corruption, scanline, dan static noise | 🟣 Purple |
| 12 | **Chaos Entropy** | Fractal Brownian Motion (fBm) noise berlapis | ⚫ Dark Gray |
| 13 | **Particle Wave Field** | Medan partikel yang bergelombang mengikuti pola wave terrain | ⚪ Bright White |

---

## 🎮 Kontrol

### Navigasi Layer
| Tombol | Fungsi |
|--------|--------|
| `1` - `9` | Pilih layer 1–9 |
| `0` | Pilih layer 10 |
| `-` | Pilih layer 11 |
| `=` | Pilih layer 12 |
| `+` | Pilih layer 13 |

### Kontrol Layer
| Tombol | Fungsi |
|--------|--------|
| `Space` / `A` | Toggle layer ON/OFF (dengan fade transition) |
| `X` | Solo mode — aktifkan hanya layer terpilih |
| `Enter` | Toggle Auto Mode (siklus otomatis semua layer) |
| `N` | Skip ke sequence berikutnya (auto mode) |

### Parameter Tuning
| Tombol | Fungsi |
|--------|--------|
| `↑` / `↓` | Navigasi parameter (Speed, Density, Brightness, Scale, Chaos) |
| `←` / `→` | Kurangi / Tambah nilai parameter |

### Global
| Tombol | Fungsi |
|--------|--------|
| `[` / `]` | Kurangi / Tambah global brightness |
| `C` | Toggle color mode ON/OFF |
| `F` / `S` | Percepat / Perlambat transisi antar layer |
| `H` / `G` | Perpanjang / Persingkat hold duration (auto mode) |
| `R` | Reset partikel (Ball Matrix & Matrix Typist) |
| `Q` / `Esc` | Keluar dari program |

---

## 🛠️ Instalasi & Penggunaan

### Prasyarat
- **Go** 1.18+ terinstal
- Terminal yang mendukung **ANSI escape codes** dan **Unicode** (disarankan: kitty, alacritty, Windows Terminal, iTerm2, atau terminal Linux modern)
- Ukuran terminal disarankan: minimal **120×30** karakter

### Clone & Run

```bash
# Clone repository
git clone https://github.com/username/performative-ascii.git
cd performative-ascii

# Install dependencies
go mod tidy

# Jalankan
go run main.go

---

┌─────────────────────────────────────────────────┐
│                   main loop                      │
│              (~25 FPS / 40ms tick)               │
├─────────────┬───────────────────┬───────────────┤
│  Input       │   Render Engine   │  Transition   │
│  Handler     │                   │  System       │
│  (goroutine) │  13 Layer Renders │  (Alpha Fade) │
├──────────────┼───────────────────┼───────────────┤
│  keyboard    │  Per-pixel math:  │  Auto/Manual  │
│  events      │  • Trigonometric  │  Mode with    │
│              │  • Noise (fBm)    │  smooth easing│
│              │  • Physics sim    │               │
│              │  • Ray marching   │               │
│              │  • ASCII mapping  │               │
└──────────────┴───────────────────┴───────────────┘

---

Teknik Rendering yang Digunakan
Procedural Noise: Hash-based 2D noise dengan smooth interpolation (cubic smoothstep)
Fractal Brownian Motion (fBm): Multi-octave noise untuk terrain generation
Ray Marching: Iterative depth solving untuk 3D mountain perspective
Domain Warping: Distorsi koordinat untuk efek organik
Easing Functions: Cubic dan quadratic easing untuk transisi halus
Physics Simulation: Bouncing balls dengan velocity, random perturbation, dan boundary collision
Layered Compositing: Alpha blending multi-layer dengan dominant color selection




