use std::{collections::HashMap, env, fs, io};

struct Puzzle {
    data: Vec<String>,
}

impl Puzzle {
    fn read_file(filename: &str) -> Result<Self, io::Error> {
        let text = fs::read_to_string(filename)?;
        let mut data = Vec::new();
        for s in text.lines() {
            data.push(s.to_string());
        }
        Ok(Puzzle { data })
    }

    fn part_one(&self) -> u64 {
        let mut or_mask = 0;
        let mut and_mask = !or_mask;
        let mut mem: HashMap<u64, u64> = HashMap::new();
        for line in self.data.iter() {
            let line = line.replace(|c| "[]".contains(c), " ");
            let data = line.split_whitespace().collect::<Vec<_>>();
            match data[0] {
                "mask" => (or_mask, and_mask) = mask_split(data[2]),
                "mem" => {
                    if let Ok(addr) = data[1].parse() {
                        if let Ok(value) = data[3].parse::<u64>() {
                            mem.insert(addr, value & and_mask | or_mask);
                        }
                    }
                },
                _ => eprintln!("Bad input!"),
            }
        }
        mem.into_iter().map(|(_, v)| v).sum()
    }

    fn part_two(&self) -> u64 {
        let mut mem: HashMap<u64, u64> = HashMap::new();
        let mut mask = String::new();
        for line in self.data.iter() {
            let line = line.replace(|c| {
                    "[]".contains(c)
                }, " ");
            let data = line.split_whitespace().collect::<Vec<_>>();
            match data[0] {
                "mask" => mask = String::from(data[2]),
                "mem" => {
                    let n = data[1].parse::<u64>().unwrap();
                    let addr = format!("{:036b}", n);
                    let addrs = gen_addrs(
                            addr.as_bytes(),
                            mask.as_bytes(),
                        )
                        .into_iter()
                        .filter_map(|s| {
                            u64::from_str_radix(&s, 2).ok()
                        })
                        .collect::<Vec<_>>();

                    let val = data[3].parse().unwrap();
                    for a in addrs {
                        mem.insert(a, val);
                    }
                },
                _ => eprintln!("Bad input!"),
            }
        }
        mem.into_iter().map(|(_, v)| v).sum()
    }
}

fn gen_addrs(addr: &[u8], mask: &[u8]) -> Vec<String> {
    let mut out = Vec::new();
    if addr.len() == 0 {
        out.push(String::new());
        return out;
    }
    let mut result = String::with_capacity(36);
    for (i, m) in mask.iter().enumerate() {
        match m {
            b'0' => result.push(addr[i] as char),
            b'1' => result.push('1'),
            _ => {
                let j = i + 1;
                let combis = gen_addrs(&addr[j..], &mask[j..]);
                for bit in ['0', '1'] {
                    for c in combis.iter() {
                        let mut r = String::with_capacity(36);
                        r.push_str(&result);
                        r.push(bit);
                        r.push_str(c);
                        out.push(r);
                    }
                }
                return out;
            },
        }
    }
    out.push(result);
    out
}

fn mask_split(mask: &str) -> (u64, u64) {
    let mut or_mask = 0;
    let mut and_mask = !or_mask;
    for (i, v) in mask.as_bytes().into_iter().enumerate() {
        let n = 1 << (35 - i);
        match v {
            b'1' => or_mask |= n,
            b'0' => and_mask &= !n,
            _ => {},
        }
    }
    (or_mask, and_mask)
}

fn main() {
    let args: Vec<String> = env::args().collect();
    let filename = args[1].as_str();

    let p = Puzzle::read_file(&filename).unwrap();

    println!("Part 1: {}", p.part_one());
    println!("Part 2: {}", p.part_two());
}
