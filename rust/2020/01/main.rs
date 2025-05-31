use std::{collections::HashSet, env, fs, io::{self, BufRead}};

struct Puzzle {
    entries: HashSet<i32>
}

impl Puzzle {
    fn read_file(filename: &str) -> Result<Self, io::Error> {
        let file = fs::File::open(filename)?;
        let reader = io::BufReader::new(file);

        let mut entries = HashSet::new();

        for line in reader.lines() {
            let line_content = line?;

            if let Ok(num) = line_content.parse::<i32>() {
                entries.insert(num);
            }
        }

        Ok(Puzzle{entries})
    }

    fn sum_to(&self, sum: i32) -> Option<(i32, i32)> {
        for num in &self.entries {
            let other = sum - num;
            if self.entries.contains(&other) {
                return Some((num.clone(), other));
            }
        }
        None
    }

    fn triplet(&self, sum: i32) -> Option<i32> {
        for i in self.entries {
            for j in self.entries {
                if *i == *j {
                    continue;
                }
                let k = sum - i - j;
                if k == *i || k == *j {
                    continue;
                }
                if self.entries.contains(&k) {
                    return Some(*i * *j * k);
                }
            }
        }
        None
    }
}

fn main() {
    let filename = env::args().nth(1).unwrap();

    let p = Puzzle::read_file(&filename).unwrap();

    if let Some((a, b)) = p.sum_to(2020) {
        println!("Part 1: {}", a * b);
    } else {
        eprintln!("Part 1: Match not found!");
    }

    if let Some(triplet) = p.triplet(2020) {
        println!("Part 2: {triplet}");
    } else {
        eprintln!("Part 2: Match not found!");
    }
}
