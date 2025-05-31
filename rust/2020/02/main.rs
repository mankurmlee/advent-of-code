use std::{env, fs, io::{self, BufRead}};

struct Record {
    range: (usize, usize),
    letter: char,
    password: String,
}

impl Record {
    fn count_valid(&self) -> bool {
        let c = self.password.chars()
            .filter(|&c| c == self.letter)
            .count();

        c >= self.range.0 && c <= self.range.1
    }

    fn index_valid(&self) -> bool {
        let i = self.range.0 - 1;
        let j = self.range.1 - 1;
        let a = self.password.chars().nth(i).unwrap() == self.letter;
        let b = self.password.chars().nth(j).unwrap() == self.letter;
        (a || b) && !(a && b)
    }
}

struct Puzzle {
    records: Vec<Record>,
}

impl Puzzle {
    fn read_file(filename: &str) -> Result<Self, io::Error> {
        let file = fs::File::open(filename)?;
        let reader = io::BufReader::new(file);

        let records = reader.lines().map(|line| {
            let data = line.unwrap()
                .replace("-", " ")
                .replace(":", " ");

            let data: Vec<&str> = data
                .split_whitespace()
                .collect();
            
            let range = (
                data[0].parse::<usize>().unwrap(),
                data[1].parse::<usize>().unwrap(),
            );

            let letter = data[2].chars().nth(0).unwrap();

            let password = data[3].to_string();

            Record { range, letter, password }
        }).collect();

        Ok(Puzzle{ records })
    }
}

fn main() {
    let filename = env::args().nth(1).unwrap();
    
    let p = Puzzle::read_file(&filename).unwrap();

    let c = p.records.iter()
        .filter(|&r| r.count_valid())
        .count();
    
    println!("Part 1: {c}");
    
    let c = p.records.iter()
        .filter(|&r| r.index_valid())
        .count();
    
    println!("Part 2: {c}");
}
