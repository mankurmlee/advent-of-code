use std::{collections::HashSet, env, fs, io::{self, BufRead}};

#[derive(Clone)]
struct Core {
    ax: i32,
    ip: usize,
}

#[derive(Clone)]
struct Instruction {
    op: String,
    arg: i32
}

impl Instruction {
    fn from(line: &str) -> Self {
        let data: Vec<&str> = line.split_whitespace().collect();
        let op = data.get(0).unwrap().to_string();
        let arg = data
            .get(1)
            .unwrap()
            .parse::<i32>()
            .unwrap();
        Instruction { op, arg }
    }
}

#[derive(Clone)]
struct Puzzle {
    core: Core,
    prog: Vec<Instruction>,
    been: HashSet<usize>,
}

impl Puzzle {
    fn read_file(filename: &str) -> Result<Self, io::Error> {
        let file = fs::File::open(filename)?;
        let reader = io::BufReader::new(file);

        let mut prog = Vec::new();

        for line in reader.lines() {
            let line = line?;
            let ins = Instruction::from(&line);
            prog.push(ins);
        }

        let mut been = HashSet::new();
        been.insert(0);

        Ok(Puzzle{
            core: Core {
                ax: 0,
                ip: 0,
            },
            prog,
            been,
        })
    }

    fn run(&mut self) {
        loop {
            self.step();
            let ip = self.core.ip;
            if self.been.contains(&ip) {
                return;
            }
            self.been.insert(ip);
        }
    }

    fn step(&mut self) {
        let ins = match self.prog.get(self.core.ip) {
            Some(ins) => ins,
            None => return,
        };
        match ins.op.as_str() {
            "acc" => {
                self.core.ax += ins.arg;
                self.core.ip += 1;
            },
            "jmp" => {
                self.core.ip = self.core.ip.wrapping_add(ins.arg as usize);
            },
            "nop" => {
                self.core.ip += 1;
            },
            _ => {
                eprintln!("Unexpected instruction {}", ins.op);
            },
        }
    }

    fn reset(&mut self) {
        self.core.ax = 0;
        self.core.ip = 0;
        self.been.clear();
    }

    fn toggle(&mut self, i: usize) -> bool {
        match self.prog[i].op.as_str() {
            "jmp" => self.prog[i].op = String::from("nop"),
            "nop" => self.prog[i].op = String::from("jmp"),
            _ => return false,
        };
        true
    }
}

fn main() {
    let filename = env::args().nth(1).unwrap();

    let mut p = Puzzle::read_file(&filename).unwrap();
    p.run();
    println!("Part 1: {}", p.core.ax);

    p.reset();
    let n = p.prog.len();
    for i in 0..=n-1 {
        let mut q = p.clone();
        if !q.toggle(i) {
            continue;
        }
        q.run();
        if q.core.ip == n {
            println!("Corruption found on line {}", i+1);
            println!("Part 2: {}", q.core.ax);
            break;
        }
    }
}
