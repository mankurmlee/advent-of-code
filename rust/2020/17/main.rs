use std::{collections::HashSet, env, fs, io};

#[derive(PartialEq, Eq, Hash, Clone)]
struct Point ( i32, i32, i32, i32 );
impl Point {
    fn add(&self, arg: (i32, i32, i32, i32)) -> Point {
        Point( self.0 + arg.0, self.1 + arg.1, self.2 + arg.2, self.3 + arg.3 )
    }
}

#[derive(Clone)]
struct Cube {
    lo: Point,
    hi: Point,
}

#[derive(Clone)]
struct Puzzle<'a> {
    alive: HashSet<Point>,
    bounds: Cube,
    adjacent: &'a Vec<(i32, i32, i32, i32)>,
}

impl <'a> Puzzle<'a> {
    fn read_file(filename: &str, adjacent: &'a Vec<(i32, i32, i32, i32)>) -> Result<Self, io::Error> {
        let text = fs::read_to_string(filename)?;

        let data = text.lines()
            .map(|s| s.as_bytes().to_vec())
            .collect::<Vec<_>>();

        let h = data.len();
        let w = data[0].len();

        let mut alive = HashSet::new();
        for y in 0..h {
            for x in 0..w {
                if data[y][x] == b'#' {
                    alive.insert(Point ( x as i32, y as i32, 0, 0 ));
                }
            }
        }

        let bounds = Cube {
            lo: Point(0, 0, 0, 0),
            hi: Point(w as i32 - 1, h as i32 - 1, 0, 0),
        };

        Ok(Puzzle { alive, bounds, adjacent })
    }

    fn simulate(&self) -> Self {
        let mut alive = HashSet::new();
        let bounds = Cube {
            lo: self.bounds.lo.add((-1, -1, -1, -1)),
            hi: self.bounds.hi.add((1, 1, 1, 1)),
        };

        for w in bounds.lo.3..=bounds.hi.3 {
            for z in bounds.lo.2..=bounds.hi.2 {
                for y in bounds.lo.1..=bounds.hi.1 {
                    for x in bounds.lo.0..=bounds.hi.0 {
                        let p = Point(x, y, z, w);
                        let n = self.count_neighbours(&p);
                        if n == 3 || n == 2 && self.alive.contains(&p) {
                            alive.insert(p);
                        }
                    }
                }
            }
        }

        Self { alive, bounds, adjacent: self.adjacent }
    }

    fn count_neighbours(&self, p: &Point) -> usize {
        self.adjacent.iter()
            .filter(|&&o| {
                let n = p.add(o);
                self.alive.contains(&n)
            })
            .count()
    }

}

fn main() {
    let args: Vec<String> = env::args().collect();
    let filename = args[1].as_str();

    let adjacent = generate_adjacent();

    let p = Puzzle::read_file(&filename, &adjacent).unwrap();

    let mut q = p.clone();
    println!("Starting with {} cubes alive", p.alive.len());
    for i in 0..6 {
        q = q.simulate();
        println!("After {} cycles, {} cubes are alive", i+1, q.alive.len());
    }
}

fn generate_adjacent() -> Vec<(i32, i32, i32, i32)> {
    let mut out = Vec::new();
    for w in -1..=1 {
        for z in -1..=1 {
            for y in -1..=1 {
                for x in -1..=1 {
                    if x == 0 && y == 0 && z == 0 && w == 0 { continue }
                    let p = (x, y, z, w);
                    out.push(p);
                }
            }
        }
    }
    out
}
