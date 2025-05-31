use std::{env, fs};

#[derive(Debug)]
struct Equation {
    m: usize,
    c: usize,
}

struct Puzzle {
    start: i32,
    buses: Vec<String>,
}

impl Puzzle {
    fn load(filename: &str) -> Self {
        let mut start = 0;
        let mut buses = Vec::new();

        if let Ok(text) = fs::read_to_string(filename) {
            let data = text.lines().collect::<Vec<&str>>();
            if data.len() == 2 {
                start = data[0].parse::<i32>()
                    .expect("start time parse error");
                buses = data[1]
                    .split(',')
                    .map(|s| s.to_string())
                    .collect();
            }
        }

        Self { start, buses }
    }

    fn earliest_bustime(&self) -> i32 {
        let start = self.start;
        let mut eid = 0;
        let mut ewait = 0;
        let buses = self.buses
            .iter()
            .filter_map(|s| s.parse::<i32>().ok());
        for id in buses {
            let wait = id - (start % id);
            if eid == 0 || ewait > wait {
                eid = id;
                ewait = wait;
            }
        }
        eid * ewait
    }

    fn get_equations(&self) -> Vec<Equation> {
        let mut eqs = Vec::new();
        for (i, s) in self.buses.iter().enumerate() {
            if let Ok(v) = s.parse() {
                let c = if i == 0 { 0 } else {
                    (0 - i as i32) % v as i32 + v as i32
                };
                eqs.push(Equation { m: v, c: c as usize });
            }
        }
        eqs
    }
}

fn find_solution(a: &Equation, b: &Equation) -> Equation {
    // println!("Finding solution to {:?} and {:?}", a, b);
    let m = a.m * b.m;
    let mut c = a.c;
    while c % b.m != b.c {
        c += a.m;
    }
    Equation {m, c}
}

fn main() {
    let args: Vec<String> = env::args().collect();
    let filename = args[1].as_str();

    let p = Puzzle::load(&filename);
    let t = p.earliest_bustime();
    println!("Part 1: {t}");

    let eqs = p.get_equations();

    let res = eqs.iter()
        .fold(Equation {m: 1, c: 0},|acc, x| {
            find_solution(&acc, x)
        });

    println!("Product: {}", res.m);
    println!("Part 2: {}", res.c);
}
