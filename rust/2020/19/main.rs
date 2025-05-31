use std::{collections::HashMap, env, fs, io};

use regex::Regex;

#[derive(Debug)]
enum Rule {
    Sequence(Vec<i32>),
    Choice(Vec<i32>, Vec<i32>),
    Char(char),
}

#[derive(Debug)]
struct Puzzle {
    rules: HashMap<i32, Rule>,
    messages: Vec<String>,
}

impl Puzzle {
    fn read_file(filename: &str) -> Result<Self, io::Error> {
        let mut rules = HashMap::new();
        let mut messages = Vec::new();

        let text = fs::read_to_string(filename)?;

        let mut block = 0;
        for line in text.lines() {
            if line.trim().is_empty() {
                block += 1;
                continue;
            }
            if block == 0 {
                let (id, rule) = Puzzle::parse_rule(line);
                rules.insert(id, rule);
            } else {
                messages.push(String::from(line));
            }
        }

        Ok(Puzzle { rules, messages })
    }

    fn parse_rule(line: &str) -> (i32, Rule) {
        let data: Vec<&str> = line.split(':').collect();

        let id = data[0].trim();
        let rule = data[1].trim();

        let id = id.parse().unwrap();

        // Check for Char type
        if rule.contains('\"') {
            let c = rule.trim()
                .trim_matches('\"')
                .as_bytes();
            return (id, Rule::Char(c[0] as char));
        }

        // Check for Choice type
        if rule.contains('|') {
            let data: Vec<&str> = rule.split('|').collect();
            let lhs = Puzzle::parse_ints(data[0]);
            let rhs = Puzzle::parse_ints(data[1]);

            return (id, Rule::Choice(lhs, rhs));
        }

        // Check for Sequence type
        let seq = Puzzle::parse_ints(rule);
        (id, Rule::Sequence(seq))
    }

    fn parse_ints(data: &str) -> Vec<i32> {
        data.split_whitespace()
            .map(|i| i.parse().unwrap())
            .collect()
    }

    fn get_pattern(&self, id: i32) -> String {
        match self.rules.get(&id).unwrap() {
            Rule::Sequence(s) => {

                if id == 11 {
                    let lhs = self.get_pattern(42);
                    let rhs = self.get_pattern(31);
                    return format!("({lhs}{rhs}|{lhs}{{2}}{rhs}{{2}}|{lhs}{{3}}{rhs}{{3}}|{lhs}{{4}}{rhs}{{4}}|{lhs}{{5}}{rhs}{{5}})");
                }

                let patt = self.seq_eval(s);

                // Part 2
                if id == 8 {
                    return patt + "+";
                }
                patt
            },
            Rule::Choice(lhs, rhs) => {
                let lhs = self.seq_eval(lhs);
                let rhs = self.seq_eval(rhs);
                format!("({lhs}|{rhs})")
            },
            Rule::Char(c) => c.to_string(),
        }
    }

    fn seq_eval(&self, s: &Vec<i32>) -> String {
        s.iter()
            .map(|id| self.get_pattern(*id))
            .collect::<Vec<_>>()
            .join("")
    }
}

fn main() {
    let args: Vec<String> = env::args().collect();
    let filename = args[1].as_str();

    let p = Puzzle::read_file(&filename).unwrap();

    let patt = format!("^{}$", p.get_pattern(0));
    println!("Pattern: {patt}");

    let regex = Regex::new(&patt).unwrap();

    // for m in p.messages.iter() {
    //     if regex.is_match(m) {
    //         println!("{m} matches");
    //     } else {
    //         println!("{m} does not match");
    //     }
    // }

    let num_match = p.messages.iter()
        .filter(|&m| regex.is_match(m))
        .count();
    println!("Num matches: {num_match}");
}
