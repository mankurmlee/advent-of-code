use std::{collections::HashMap, env, fs, io::{self, BufRead}};

struct Puzzle {
    rules: HashMap<String, Vec<(String, usize)>>,
    d_count: HashMap<String, usize>,
}

impl Puzzle {
    fn read_file(filename: &str) -> Result<Self, io::Error> {
        let file = fs::File::open(filename)?;
        let reader = io::BufReader::new(file);

        let mut p = Puzzle {
            rules: HashMap::new(),
            d_count: HashMap::new(),
        };
        for line in reader.lines() {
            p.add_rule(line?);
        }

        Ok(p)
    }

    fn add_rule(&mut self, line: String) {
        let data: Vec<_> = line.split_whitespace().collect();
        let n = data.len();
        if n != 7 && (n < 8 || n % 4 != 0) {
            eprintln!("Bad input: {line}");
            return;
        }
        let mut v = Vec::new();
        if n > 7 {
            for i in (4..=n-1).step_by(4) {
                let count = data[i].parse::<usize>().unwrap();
                let colour = format!("{} {}", data[i+1], data[i+2]);
                v.push((colour, count));
           }
        }
        let k = format!("{} {}", data[0], data[1]);
        self.rules.insert(k, v);
    }

    fn rev_index(&self) -> HashMap<String, Vec<String>> {
        let mut idx = HashMap::new();
        for (k, v) in self.rules.iter() {
            for (c, _) in v.iter() {
                if !idx.contains_key(c) {
                    idx.insert(c.clone(), Vec::new());
                }
                idx.get_mut(c).unwrap().push(k.clone());
            }
        }
        idx
    }

    fn get_ancestors(&self, colour: &str) -> Vec<String> {
        let mut ancs = Vec::new();
        let rev = self.rev_index();

        let mut q = vec![colour];
        while q.len() > 0 {
            let c = q.pop().unwrap();
            let bags = match rev.get(c) {
                Some(bags) => bags,
                None => {
                    // println!("Can't find {c} in reverse index");
                    continue
                },
            };
            for bag in bags {
                let b = bag.as_str();
                if q.contains(&b) {
                    // println!("{b} is already in the queue");
                    continue;
                }
                if ancs.contains(bag) {
                    // println!("{bag} is already in output");
                    continue;
                }
                q.push(b);
                ancs.push(String::from(b));
            }
        }
        ancs
    }

    fn count_descendants(&mut self, colour: &str) -> usize {
        if let Some(&count) = self.d_count.get(colour) {
            println!("cd {colour} => {count} (cached)");
            return count;
        }
        let children = self.rules.get(colour).unwrap().clone();
        let count = children
            .iter()
            .map(|(b, n)| (self.count_descendants(b) + 1) * n)
            .sum();
        self.d_count.insert(String::from(colour), count);
        println!("cd {colour} => {count}");
        count
    }
}

fn main() {
    let filename = env::args().nth(1).unwrap();

    let mut p = Puzzle::read_file(&filename).unwrap();
    let a = p.get_ancestors("shiny gold");
    // println!("Shiny Bag ancestors: {}", a.join(", "));
    println!("Part 1: {}", a.len());

    let n = p.count_descendants("shiny gold");
    println!("Part 2: {n}");
}
