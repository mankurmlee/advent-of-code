use std::{collections::HashMap, env, fs, io::{self, BufRead}};
use regex::Regex;

struct Passport {
    fields: HashMap<String, String>,
}

impl Passport {
    const REQUIRED: [&'static str; 7] = [
        "byr", "iyr", "eyr", "hgt", "hcl", "ecl", "pid",
    ];

    const EYE_COLOURS: [&'static str; 7] = [
        "amb", "blu", "brn", "gry", "grn", "hzl", "oth",
    ];

    fn has_required_fields(&self) -> bool {
        Passport::REQUIRED
            .iter()
            .map(|&k| k.to_string())
            .filter(|k| !self.fields.contains_key(k))
            .count() == 0
    }

    fn has_valid_values(&self) -> bool {
        for (k, v) in &self.fields {
            match k.as_str() {
                "byr" => {
                    let v = v.parse::<i32>().unwrap_or(0);
                    if v < 1920 || v > 2002 {
                        return false;
                    }
                }
                "iyr" => {
                    let v = v.parse::<i32>().unwrap_or(0);
                    if v < 2010 || v > 2020 {
                        return false;
                    }
                }
                "eyr" => {
                    let v = v.parse::<i32>().unwrap_or(0);
                    if v < 2020 || v > 2030 {
                        return false;
                    }
                }
                "hgt" => {
                    let n = v.len();
                    if n < 4 || n > 5 {
                        return false;
                    }
                    match &v[n - 2..] {
                        "cm" => {
                            let cm = v[..n-2]
                                .parse::<i32>()
                                .unwrap_or(0);
                            if cm < 150 || cm > 193 {
                                return false;
                            }
                        },
                        "in" => {
                            let inches = v[..n-2]
                                .parse::<i32>()
                                .unwrap_or(0);
                            if inches < 59 || inches > 76 {
                                return false;
                            }
                        },
                        _ => return false,
                    }
                }
                "hcl" => {
                    let re = Regex::new(r"^#[0-9a-fA-F]{6}$").unwrap();
                    if !re.is_match(v.as_str()) {
                        return false;
                    }
                }
                "ecl" => {
                    if !Passport::EYE_COLOURS.contains(&v.as_str()) {
                        return false;
                    }
                }
                "pid" => {
                    let re = Regex::new(r"^[0-9]{9}$").unwrap();
                    if !re.is_match(v.as_str()) {
                        return false;
                    }
                }
                "cid" => {}
                _ => return false,
            }
        }
        true
    }
}

struct Puzzle {
    passports: Vec<Passport>,
}

impl Puzzle {
    fn read_file(filename: &str) -> Result<Self, io::Error> {
        let file = fs::File::open(filename)?;
        let reader = io::BufReader::new(file);

        let mut passports = Vec::new();
        let mut fields = HashMap::new();

        for line in reader.lines() {
            let line = line?;
            if line.trim() == "" {
                    passports.push(Passport {fields: fields.clone()});
                    fields = HashMap::new();
            } else {
                for f in line.split_whitespace() {
                    let e: Vec<&str> = f.split(":").collect();
                    fields.insert(e[0].to_string(), e[1].to_string());
                }
            }
        }

        passports.push(Passport {fields: fields.clone()});

        Ok(Puzzle { passports })
    }
}

fn main() {
    let filename = env::args().nth(1).unwrap();
    
    let p = Puzzle::read_file(&filename).unwrap();

    let has_fields: Vec<&Passport> = p.passports
        .iter()
        .filter(|&p| p.has_required_fields())
        .collect();

    println!("Part 1: {}", has_fields.len());

    let count = has_fields
        .iter()
        .filter(|&&p| p.has_valid_values())
        .count();

    println!("Part 2: {count}");
}
