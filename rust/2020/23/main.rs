use std::{collections::HashMap, env, fs, io};

// Optimization: Use a vector instead of a hashmap

#[derive(Clone)]
struct Puzzle {
    data: HashMap<usize, usize>,
    current: usize,
}
impl Puzzle {
    fn read_file(filename: &str) -> Result<Self, io::Error> {
        let text = fs::read_to_string(filename)?;

        let text = text.lines().nth(0).unwrap();
        let input: Vec<usize> = text.split("")
            .filter_map(|s| s.parse::<usize>().ok())
            .collect();

        let current = input[0];
        let n = input.len();

        let mut data = HashMap::new();
        for (i, v) in input.iter().enumerate() {
            let next_i = (i + 1) % n;
            data.insert(*v, input[next_i]);
        }

        Ok(Puzzle { data, current })
    }

    fn cup_move(&mut self) {
        let n = self.data.len();
        let m = n - 1;
        let c = self.current;

        // take 3
        let n1 = *self.data.get(&c).unwrap();
        let n2 = *self.data.get(&n1).unwrap();
        let n3 = *self.data.get(&n2).unwrap();
        self.current = *self.data.get(&n3).unwrap();
        self.data.insert(c, self.current);

        // find dest
        let mut dest = (c + m) % n;
        if dest == 0 {
            dest = n;
        }
        while dest == n1 || dest == n2 || dest == n3 {
            dest = (dest + m) % n;
            if dest == 0 {
                dest = n;
            }
        }

        // reinsert
        let n5 = *self.data.get(&dest).unwrap();
        self.data.insert(dest, n1);
        self.data.insert(n3, n5);
    }

    fn check_sum1(&self) -> String {
        let mut l = Vec::new();
        let mut i = 1;
        loop {
            let n = *self.data.get(&i).unwrap();
            if n == 1 {
                return l.join("");
            }
            l.push(format!("{n}"));
            i = n;
        }
    }

    fn million(&self) -> Self {
        let mut out = self.clone();
        let c = out.current;
        let mil = 1_000_000;

        let t = out.data.iter()
            .filter(|(_, &v)| v == c)
            .map(|(&i, _)| i)
            .nth(0).unwrap();

        // Extend the current collection to a million
        out.data.insert(t, 10);
        for i in 10..=mil {
            out.data.insert(i, i+1);
        }

        // Reattach tail
        out.data.insert(mil, c);

        out
    }

    fn check_sum2(&self) -> usize {
        let n1 = *self.data.get(&1).unwrap();
        let n2 = *self.data.get(&n1).unwrap();
        n1 * n2
    }

}

fn main() {
    let args: Vec<String> = env::args().collect();
    let filename = args[1].as_str();

    let mut p = Puzzle::read_file(&filename).unwrap();
    let mut q = p.million();

    for _ in 0..100 {
        p.cup_move();
    }
    let chksum1 = p.check_sum1();

    for _ in 0..10_000_000 {
        q.cup_move();
    }

    let chksum2 = q.check_sum2();

    println!("Part 1: {chksum1}");
    println!("Part 2: {chksum2}");
}
