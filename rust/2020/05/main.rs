use std::{collections::HashSet, env, fs, io::{self, BufRead}, num::ParseIntError};

struct Seat(String);

impl Seat {
    fn seat_id(&self) -> Result<i32, ParseIntError> {
        let bin = &self.0
            .replace("F","0")
            .replace("B","1")
            .replace("L","0")
            .replace("R","1");
        i32::from_str_radix(&bin, 2)
    }
}

struct Puzzle (Vec<Seat>);

impl Puzzle {
    fn read_file(filename: &str) -> Result<Self, io::Error> {
        let file = fs::File::open(filename)?;
        let reader = io::BufReader::new(file);
        let data: Vec<_> = reader
            .lines()
            .map(|l| Seat(l.unwrap()))
            .collect();
        Ok(Puzzle(data))
    }
}

fn main() {
    let filename = env::args().nth(1).unwrap();
    
    let p = Puzzle::read_file(&filename).unwrap();

    for s in &p.0 {
        println!("{} => {}", s.0, s.seat_id().unwrap());
    }

    let seat_ids: HashSet<_> = p.0
        .iter()
        .map(|s| s.seat_id().unwrap())
        .collect();
    
    println!("Part 1: {}", seat_ids.iter().max().unwrap());

    for id in 9..=1015 {
        if seat_ids.contains(&id) {
            continue;
        }
        let a = id - 1;
        if !seat_ids.contains(&a) {
            continue;
        }
        let b = id + 1;
        if !seat_ids.contains(&b) {
            continue;
        }
        println!("Part 2: {id}");
    }
}
