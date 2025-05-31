use std::{collections::HashSet, env, fs, io::{self, BufRead}};

struct Group(Vec<String>);

impl Group {
    fn count_or_yes(&self) -> i32 {
        let mut set = HashSet::new();
        for s in self.0.iter() {
            for c in s.chars() {
                set.insert(c);
            }
        }
        set.len() as i32
    }

    fn count_and_yes(&self) -> i32 {
        let s = match self.0.get(0) {
            Some(s) => s,
            None => return 0,
        };

        let mut set: HashSet<char> = s.chars().collect();
        for s in self.0.iter().skip(1) {
            let common: HashSet<char> = s.chars().filter(
                |c| set.contains(c)
            ).collect();
            set = common;
        }
        set.len() as i32
    }
}

struct Puzzle(Vec<Group>);

impl Puzzle {
    fn read_file(filename: &str) -> Result<Self, io::Error> {
        let file = fs::File::open(filename)?;
        let reader = io::BufReader::new(file);

        let mut groups = Vec::new();
        let mut group = Vec::new();

        for line in reader.lines() {
            let line = line?;
            if line.trim() == "" {
                groups.push(Group(group));
                group = Vec::new();
                continue;
            }
            group.push(line);
        }
        groups.push(Group(group));

        Ok(Puzzle(groups))
    }
}

fn main() {
    let filename = env::args().nth(1).unwrap();
    
    let p = Puzzle::read_file(&filename).unwrap();
    
    let or_yes: i32 = p.0.iter().map(
        |g| g.count_or_yes()
    ).sum();

    println!("Part 1: {or_yes}");

    let and_yes: i32 = p.0.iter().map(
        |g| g.count_and_yes()
    ).sum();

    println!("Part 2: {and_yes}");
}
