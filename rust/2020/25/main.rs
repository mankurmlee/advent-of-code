use std::{env, fs, io};

struct Puzzle {
    keys: Vec<u64>,
}

impl Puzzle {
    fn read_file(filename: &str) -> Result<Self, io::Error> {
        let text = fs::read_to_string(filename)?;

        let keys = text.lines()
            .map(|s| s.parse().unwrap())
            .collect();

        Ok(Puzzle { keys })
    }
}

fn main() {
    let args: Vec<String> = env::args().collect();
    let filename = args[1].as_str();

    let p = Puzzle::read_file(&filename).unwrap();

    let key1 = p.keys[0];
    let key2 = p.keys[1];

    println!("Trying to find loop value for {key1}");
    let l1 = find_loop(key1, 7);
    println!("Loop is {l1}!");

    println!("Transforming {key2} with {l1} loops...");
    let k =  transform(key2, l1);
    println!("Key is {k}");
}

fn transform(subject: u64, loop_size: i32) -> u64 {
    let mut s = 1;
    for _ in 0..loop_size {
        s *= subject;
        s %= 20201227;
    }
    s
}

fn find_loop(key: u64, subject: u64) -> i32 {
    let mut loops = 0;
    let mut s = 1;
    while s != key {
        s *= subject;
        s %= 20201227;
        loops += 1;
    }
    loops
}
