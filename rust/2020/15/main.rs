use std::{collections::HashMap, env, fs, io};

struct Puzzle {
    data: Vec<usize>,
}

impl Puzzle {
    fn read_file(filename: &str) -> Result<Self, io::Error> {
        let text = fs::read_to_string(filename)?;
        let data = text
            .trim()
            .split(',')
            .map(|x| x.parse().unwrap())
            .collect();
        Ok(Puzzle { data })
    }

    fn get_turn(&self, arg: usize) -> usize {
        if arg <= self.data.len() {
            return self.data[arg-1];
        }
        let mut map = HashMap::new();
        let mut turn = 1;
        for &v in self.data.iter() {
            map.insert(v, turn);
            turn += 1;
        }

        let mut next_v = 0;
        let mut v = next_v;
        while turn < arg {
            next_v = turn - *map.get(&v).unwrap_or(&turn);
            map.insert(v, turn);
            v = next_v;
            turn += 1;
        }

        next_v
    }
}

fn main() {
    let args: Vec<String> = env::args().collect();
    let filename = args[1].as_str();

    let p = Puzzle::read_file(&filename).unwrap();
    let num = p.get_turn(2020);
    println!("Part 1: {num}");

    let num = p.get_turn(30000000);
    println!("Part 2: {num}");
}
