use std::{env, fs, io::{self, BufRead}, fmt};

#[derive(Clone)]
struct Puzzle {
    width: usize,
    height: usize,
    data: Vec<u8>,
}

impl Puzzle {
    const ADJACENT: [(i32, i32); 8] = [
        (-1, -1), ( 0, -1), ( 1, -1),
        (-1,  0),           ( 1,  0),
        (-1,  1), ( 0,  1), ( 1,  1),
    ];

    fn read_file(filename: &str) -> Result<Self, io::Error> {
        let file = fs::File::open(filename)?;
        let reader = io::BufReader::new(file);

        let mut width = 0;
        let mut height = 0;
        let mut data = Vec::new();

        for line in reader.lines() {
            data.extend(line?.into_bytes());
            if width == 0 {
                width = data.len();
            }
            height += 1;
        }

        Ok(Puzzle { width, height, data })
    }

    fn shuffle(&mut self, los: bool) -> bool {
        let mut changed = false;
        let mut buffer = self.data.clone();
        let w = self.width;
        let h = self.height;
        let tolerance = if los { 5 } else { 4 };
        for y in 0..h {
            for x in 0..w {
                let n = if los {
                    self.count_nearest(x as i32, y as i32)
                } else {
                    self.count_neighbours(x as i32, y as i32)
                };
                let i = y * w + x;
                match self.data[i] {
                    b'L' => {
                        if n == 0 {
                            changed = true;
                            buffer[i] = b'#';
                        }
                    },
                    b'#' => {
                        if n >= tolerance {
                            changed = true;
                            buffer[i] = b'L';
                        }
                    },
                    _ => {},
                }
            }
        }
        if changed {
            self.data = buffer;
        }
        changed
    }

    fn count_neighbours(&self, x: i32, y: i32) -> usize {
        let mut n = 0;
        let w = self.width as i32;
        let h = self.height as i32;

        for (dx, dy) in Puzzle::ADJACENT {
            let nx = x + dx;
            let ny = y + dy;
            if nx < 0 || nx >= w || ny < 0 || ny >= h { continue }
            let i = ny * w + nx;
            if self.data[i as usize] == b'#' {
                n += 1;
            }
        }

        n
    }
    fn count_nearest(&self, x: i32, y: i32) -> usize {
        let mut n = 0;
        let w = self.width as i32;
        let h = self.height as i32;

        for (dx, dy) in Puzzle::ADJACENT {
            let mut nx = x + dx;
            let mut ny = y + dy;
            while nx >= 0 && nx < w && ny >= 0 && ny < h {
                let i = ny * w + nx;
                let found = match self.data[i as usize] {
                    b'#' => {
                        n += 1;
                        true
                    },
                    b'L' => true,
                    _ => {
                        nx += dx;
                        ny += dy;
                        false
                    },
                };
                if found { break }
            }
        }

        n
    }

    fn occupied(&self) -> usize {
        self.data
            .iter()
            .filter(|&b| *b == b'#')
            .count()
    }
}

fn main() {
    let args: Vec<String> = env::args().collect();
    let filename = args[1].as_str();

    let mut p = Puzzle::read_file(&filename).unwrap();
    let mut q = p.clone();

    while p.shuffle(false) {}
    let o = p.occupied();
    println!("Part 1: {o}");

    while q.shuffle(true) {}
    let o = q.occupied();
    println!("Part 2: {o}");
}

impl fmt::Display for Puzzle {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        for v in self.data.chunks(self.width) {
            let chunk = String::from_utf8_lossy(v);
            writeln!(f, "{chunk}")?
        }
        Ok(())
    }
}
