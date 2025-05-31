use std::{collections::HashMap, env, fs, io::{self, BufRead}};

struct Puzzle {
    data: Vec<usize>,
    cache: HashMap<usize, usize>,
}

impl Puzzle {
    fn read_file(filename: &str) -> Result<Self, io::Error> {
        let file = fs::File::open(filename)?;
        let reader = io::BufReader::new(file);

        let mut data = Vec::new();

        for line in reader.lines() {
            data.push(line?.parse::<usize>().unwrap());
        }

        data.sort();

        Ok(Puzzle { 
            data,
            cache: HashMap::new(),
        })
    }
    
    fn count_jolts_by_size(&self) -> HashMap<usize, usize> {
        let mut joltmap = HashMap::new();
        joltmap.insert(3, 1);

        let mut jolts = 0;
        for r in self.data.iter() {
            let size = r - &jolts;
            *joltmap.entry(size).or_insert(0) += 1;
            jolts = *r;
        }
        joltmap
    }
    
    fn count_arrangements(&mut self, joltage: usize) -> usize {
        if let Some(v) = self.cache.get(&joltage) {
            return *v;
        }

        let adapters = self.data
            .iter()
            .filter(|&v| *v > joltage && *v <= joltage + 3)
            .cloned()
            .collect::<Vec<_>>();

        let arr = if adapters.len() == 0 { 1 } else {
            adapters
                .iter()
                .map(|&v| self.count_arrangements(v))
                .sum::<usize>()
        };

        self.cache.insert(joltage, arr);
        arr
    }
}

fn main() {
    let args: Vec<String> = env::args().collect();
    let filename = args[1].as_str();

    let mut p = Puzzle::read_file(&filename).unwrap();

    let j = p.count_jolts_by_size();
    let a = j.get(&1).unwrap_or(&0);
    let b = j.get(&3).unwrap_or(&0);
    println!("Part 1: {}", a * b);

    let count = p.count_arrangements(0);
    println!("Part 2: {count}");
}
