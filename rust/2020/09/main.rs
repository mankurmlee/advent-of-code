use std::{collections::HashSet, env, fs, io::{self, BufRead}};

struct Puzzle {
    data: Vec<usize>,
}

impl Puzzle {
    fn read_file(filename: &str) -> Result<Self, io::Error> {
        let file = fs::File::open(filename)?;
        let reader = io::BufReader::new(file);

        let mut data = Vec::new();

        for line in reader.lines() {
            data.push(line?.parse::<usize>().unwrap());
        }

        Ok(Puzzle { data })
    }
    
    fn find_error(&self, samples: usize) -> Option<usize> {
        let n = self.data.len();

        let mut pre: HashSet<usize> = self.data
            .iter()
            .take(samples)
            .cloned()
            .collect();

        for i in samples..=n-1 {
            let num = self.data[i];
            let mut bad = true;
            for &other in pre.iter() {
                if other > num { continue }
                let diff = num - other;
                if other != diff && pre.contains(&diff) {
                    bad = false;
                    break;
                }
            }
            if bad {
                return Some(num);
            }
            pre.remove(&self.data[i - samples]);
            pre.insert(num);
        }

        None
    }
    
    fn find_weakness(&self, bad: usize) -> Option<usize> {
        let (i, j) = match self.find_indices(bad) {
            Some(x) => x,
            None => return None,
        };
        let slice = &self.data[i..j+1];
        let min = *slice.iter().min().unwrap();
        let max = *slice.iter().max().unwrap();
        Some(min + max)
    }
    
    fn find_indices(&self, bad: usize) -> Option<(usize, usize)> {
        let n = self.data.len();
        
        for j in (1..=n-1).rev() {
            let mut total = self.data[j];
            if total > bad { continue }
            
            for i in (0..=j-1).rev() {
                total += self.data[i];
                if total == bad { 
                    println!("Found {i} and {j}");
                    return Some((i, j));
                }
                if total > bad { break }
            }
        }
        None
    }
}

fn main() {
    let args: Vec<String> = env::args().collect();
    if args.len() != 3 { return }
    let filename = args[1].as_str();
    let samples = args[2].parse::<usize>().unwrap();
    let p = Puzzle::read_file(filename).unwrap();
    
    let bad = p.find_error(samples).unwrap();
    println!("Part 1: {bad}");

    let weakness = p.find_weakness(bad).unwrap();
    println!("Part 2: {weakness}");
}
