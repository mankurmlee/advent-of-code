use std::{collections::HashMap, env, fs, io};

struct Range {
    lo: i32,
    hi: i32,
}
impl Range {
    fn contains(&self, n: i32) -> bool {
        n >= self.lo && n <= self.hi
    }
}

struct Pair {
    a: Range,
    b: Range,
}
impl Pair {
    fn num_is_valid(&self, n: i32) -> bool {
        self.a.contains(n) || self.b.contains(n)
    }

    fn col_is_valid(&self, col: &[i32]) -> bool {
        for &n in col.iter() {
            if !self.num_is_valid(n) {
                return false;
            }
        }
        true
    }
}

struct Puzzle {
    rules: HashMap<String,Pair>,
    my_ticket: Vec<i32>,
    nearby: Vec<Vec<i32>>,
}
impl Puzzle {
    fn read_file(filename: &str) -> Result<Self, io::Error> {
        let text = fs::read_to_string(filename)?;
        let mut lines = text.lines();

        let mut rules = HashMap::new();
        let mut my_ticket = Vec::new();
        let mut nearby = Vec::new();

        loop {
            let line = if
                let Some(v) = lines.next() { v }
                else { break };
            if line.trim() == "" { break }
            let kv: Vec<&str> = line
                .split(':')
                .collect();
            let replaced = kv[1].replace('-', " ");
            let data: Vec<i32> = replaced
                .split_whitespace()
                .filter_map(|s| s.parse().ok())
                .collect();

            let pair = Pair {
                a: Range { lo: data[0], hi: data[1] },
                b: Range { lo: data[2], hi: data[3] },
            };
            rules.insert(String::from(kv[0]), pair);
        }

        lines.next();
        if let Some(line) = lines.next() {
            my_ticket = line
                .split(',')
                .filter_map(|n| n.parse().ok())
                .collect();
        }

        lines.next();
        lines.next();
        loop {
            let line = if
                let Some(v) = lines.next() { v }
                else { break };
            if line.trim() == "" { break }
            let ticket = line
                .split(',')
                .filter_map(|n| n.parse().ok())
                .collect();
            nearby.push(ticket);
        }

        Ok(Puzzle { rules, my_ticket, nearby })
    }

    fn get_error_rate(&self) -> i32 {
        self.nearby
            .iter()
            .flat_map(|v| v.iter())
            .filter(|&n| !self.num_is_valid(*n))
            .sum()
    }

    fn num_is_valid(&self, n: i32) -> bool {
        for (_, v) in self.rules.iter() {
            if v.num_is_valid(n) {
                return true;
            }
        }
        false
    }

    fn valid_tickets(&self) -> Vec<Vec<i32>> {
        let mut valid = Vec::new();
        valid.push(self.my_ticket.clone());
        for t in self.nearby.iter() {
            let bad = t.iter().filter(|&n| !self.num_is_valid(*n)).count();
            if bad == 0 {
                valid.push(t.clone());
            }
        }
        valid
    }

    fn get_rules_that_fit_column(&self, col: &[i32]) -> Vec<String> {
        let mut fit = Vec::new();
        for (k, v) in self.rules.iter() {
            if v.col_is_valid(col) {
                fit.push(k.clone());
            }
        }
        fit
    }

    fn get_col_rules(&self) -> HashMap<usize, Vec<String>> {
        let mut matches = HashMap::new();
        let cols = transpose(self.valid_tickets());
        for (i, col) in cols.iter().enumerate() {
            let fit = self.get_rules_that_fit_column(col);
            matches.insert(i, fit);
        }
        matches
    }

    fn map_ticket(&self) -> HashMap<String, i32> {
        let mut myticket = HashMap::new();
        let colrules = self.get_col_rules();
        let n = self.rules.len();
        while myticket.len() < n {
            let mut cycle = HashMap::new();
            for (k, v) in colrules.iter() {
                let mut rules = Vec::new();
                for s in v.iter() {
                    if myticket.contains_key(s) { continue }
                    rules.push(s.clone());
                }
                if rules.len() != 1 { continue }
                let rule = rules[0].clone();
                cycle.insert(rule, self.my_ticket[*k]);
            }
            myticket.extend(cycle);
        }
        myticket
    }
}

fn main() {
    let args: Vec<String> = env::args().collect();
    let filename = args[1].as_str();

    let p = Puzzle::read_file(&filename).unwrap();

    let error_rate = p.get_error_rate();
    println!("Part 1: {}", error_rate);

    let map = p.map_ticket();
    let prod = get_product_departure(&map);
    println!("Part 2: {prod}");
}

fn get_product_departure(map: &HashMap<String, i32>) -> u64 {
    let mut prod = 1;
    for (k, &v) in map.iter() {
        if !k.starts_with("departure") { continue }
        prod *= v as u64;
    }
    prod
}

fn transpose(tickets: Vec<Vec<i32>>) -> Vec<Vec<i32>> {
    let num_tickets = tickets.len();
    let num_fields = tickets[0].len();

    let mut out = vec![vec![0; num_tickets]; num_fields];

    for j in 0..num_tickets {
        for i in 0..num_fields {
            out[i][j] = tickets[j][i];
        }
    }

    out
}
