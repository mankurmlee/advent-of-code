use std::{env, fs, io};

struct Puzzle {
    data: Vec<String>,
}

enum NumExp {
    Exp(Expr),
    Num(i64),
}

struct Expr {
    left: String,
    op: u8,
    right: String,
}

impl Puzzle {
    fn read_file(filename: &str) -> Result<Self, io::Error> {
        let text = fs::read_to_string(filename)?;

        let data = text.lines().map(String::from).collect();

        Ok(Puzzle { data })
    }

    fn sum(&self, adv: bool) -> i64 {
        self.data.iter()
            .map(|s| eval(s, adv))
            .sum()
    }
}

fn eval(s: &str, adv: bool) -> i64 {
    let rev: String = s.trim()
        .chars()
        .rev()
        .collect();
    let n = real_eval(&rev, adv);
    println!("{s} becomes {n}");
    n
}

fn real_eval(s: &str, adv: bool) -> i64 {
    let expr = match parse_expr(s, adv) {
        NumExp::Exp(e) => e,
        NumExp::Num(n) => return n,
    };

    let left  = real_eval(&expr.left,  adv);
    let right = real_eval(&expr.right, adv);

    match expr.op {
        b'*' => left * right,
        _    => left + right,
    }
}

fn parse_expr(s: &str, adv: bool) -> NumExp {
    if let Ok(n) = s.parse() {
        return NumExp::Num(n);
    }

    let mut s = &s[..];

    let i = loop {
        if adv {
            let i = find_op(&s, adv);
            if i > 0 {
                break i;
            }
        }
        let i = find_op(&s, false);
        if i > 0 {
            break i;
        }
        let n = s.len();
        s = &s[1..n-1];
    };

    let left = s[..i].trim().to_string();
    let op = s.bytes().nth(i).unwrap();
    let right = s[i+1..].trim().to_string();
    NumExp::Exp(Expr { left, op, right })
}

fn find_op(s: &str, mul_only: bool) -> usize {
    let mut nest = 0;
    for (i, c) in s.char_indices() {
        match c {
            ')' => nest += 1,
            '(' => nest -= 1,
            '+' => if nest == 0 && !mul_only { return i },
            '*' => if nest == 0 { return i },
            _ => {},
        }
    }
    0
}

fn main() {
    let args: Vec<String> = env::args().collect();
    let filename = args[1].as_str();

    let p = Puzzle::read_file(&filename).unwrap();
    println!("Part 1: {}", p.sum(false));
    println!("Part 2: {}", p.sum(true));
}
