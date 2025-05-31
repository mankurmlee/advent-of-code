use std::{collections::HashSet, env, fs, io};

struct Puzzle {
    data: Vec<String>,
    flipped: HashSet<(i32, i32)>,
    lo: (i32, i32),
    hi: (i32, i32),
}

impl Puzzle {
    const ADJACENT: [(i32, i32); 6] = [
        (0, 1), (1, 1),
        (-1, 0), (1, 0),
        (-1, -1), (0, -1),
    ];

    fn read_file(filename: &str) -> Result<Self, io::Error> {
        let text = fs::read_to_string(filename)?;

        let data = text.lines().map(String::from).collect();
        let flipped = HashSet::new();

        let lo = (0, 0);
        let hi = (0, 0);

        Ok(Puzzle { data, flipped, lo, hi })
    }

    fn flip(&mut self, c: &(i32, i32)) -> bool {
        if self.flipped.remove(&c) { return false }
        self.flipped.insert(c.clone());
        true
    }

    fn count_flipped(&mut self) -> usize {
        let to_flip: Vec<(i32, i32)> = self.data.iter()
            .map(|txt| get_coords(txt))
            .collect();
        let mut lo = self.lo;
        let mut hi = self.hi;
        for c in to_flip {
            if self.flip(&c) {
                if c.0 < lo.0 {
                    lo.0 = c.0;
                }
                if c.1 < lo.1 {
                    lo.1 = c.1;
                }
                if c.0 > hi.0 {
                    hi.0 = c.0;
                }
                if c.1 > hi.1 {
                    hi.1 = c.1;
                }
            }
        }
        self.lo = lo;
        self.hi = hi;
        self.flipped.len()
    }

    fn advance_day(&mut self) -> usize {
        let mut out = HashSet::new();
        let lo_x = self.lo.0 - 1;
        let lo_y = self.lo.1 - 1;
        let hi_x = self.hi.0 + 1;
        let hi_y = self.hi.1 + 1;
        let mut lo = self.lo;
        let mut hi = self.hi;
        for y in lo_y..=hi_y {
            for x in lo_x..=hi_x {
                if self.adjust(x, y) {
                    if x < lo.0 {
                        lo.0 = x;
                    }
                    if y < lo.1 {
                        lo.1 = y;
                    }
                    if x > hi.0 {
                        hi.0 = x;
                    }
                    if y > hi.1 {
                        hi.1 = y;
                    }
                    out.insert((x, y));
                }
            }
        }
        self.flipped = out;
        self.lo = lo;
        self.hi = hi;
        self.flipped.len()
    }

    fn adjust(&mut self, x: i32, y: i32) -> bool {
        let n = Puzzle::ADJACENT.iter()
            .filter(|(dx, dy)| {
                let nx = x + dx;
                let ny = y + dy;
                self.flipped.contains(&(nx, ny))
            })
            .count();
        if self.flipped.contains(&(x, y)) {
            return n == 1 || n == 2;
        }
        n == 2
    }
}

fn main() {
    let args: Vec<String> = env::args().collect();
    let filename = args[1].as_str();

    let mut p = Puzzle::read_file(&filename).unwrap();

    let n = p.count_flipped();
    println!("Part 1: {n}");

    let mut n = 0;
    for _ in 0..100 {
        n = p.advance_day();
    }
    println!("Part 2: {n}");
}

fn get_coords(text: &str) -> (i32, i32) {
    let n = text.len();
    let mut i = 0;
    let mut x = 0;
    let mut y = 0;
    while i < n {
        let snip = &text[i..];
        if snip.starts_with("e") {
            x += 1;
            y += 1;
            i += 1;
        } else if snip.starts_with("se") {
            x += 1;
            i += 2;
        } else if snip.starts_with("sw") {
            y -= 1;
            i += 2;
        } else if snip.starts_with("w") {
            x -= 1;
            y -= 1;
            i += 1;
        } else if snip.starts_with("nw") {
            x -= 1;
            i += 2;
        } else if snip.starts_with("ne") {
            y += 1;
            i += 2;
        }
    }
    // println!("{text} => ({x}, {y})");
    (x, y)
}
