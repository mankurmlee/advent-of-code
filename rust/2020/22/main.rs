use std::{collections::HashSet, env, fs, io};

#[derive(Debug, Clone)]
struct Puzzle {
    decks: Vec<Vec<usize>>,
    memo: HashSet<String>,
}

impl Puzzle {
    fn read_file(filename: &str) -> Result<Self, io::Error> {
        let text = fs::read_to_string(filename)?;

        let mut decks = Vec::new();
        let mut deck = Vec::new();
        for line in text.lines() {
            if line.trim().is_empty() {
                decks.push(deck);
                deck = Vec::new();
                continue
            }
            if line.starts_with("Player") { continue }
            let card = line.parse().unwrap();
            deck.push(card);
        }
        decks.push(deck);

        let memo = HashSet::new();

        Ok(Puzzle { decks, memo })
    }

    fn play_round(&mut self) -> bool {
        if self.decks[0].len() == 0 || self.decks[1].len() == 0 {
            return false;
        }

        let a = self.decks[0].remove(0);
        let b = self.decks[1].remove(0);

        if a > b {
            self.decks[0].push(a);
            self.decks[0].push(b);
        } else {
            self.decks[1].push(b);
            self.decks[1].push(a);
        }
        true
    }

    fn play_game(&mut self) {
        while self.play_round() {}
    }

    fn get_score(&self) -> usize {
        let i = if self.decks[0].len() > 0 { 0 } else { 1 };
        self.decks[i].iter().rev().enumerate()
            .map(|(i, &v)| (i+1) * v)
            .sum()
    }

    fn play_recursive_game(&mut self) -> u8 {
        //println!("Starting game: {:?} {:?}", self.decks[0], self.decks[1]);
        loop {
            let s = self.serialise();
            if self.memo.contains(&s) { return 0 }
            if self.decks[0].len() == 0 { return 1 }
            if self.decks[1].len() == 0 { return 0 }
            self.memo.insert(s);
            self.play_recursive_round();
        }
    }

    fn serialise(&self) -> String {
        self.decks.iter()
            .map(|d| {
                d.iter()
                    .map(|n| format!("{n}"))
                    .collect::<Vec<_>>()
                    .join(",")
            })
            .collect::<Vec<_>>()
            .join(";")
    }

    fn play_recursive_round(&mut self) {
        let a = self.decks[0].remove(0);
        let b = self.decks[1].remove(0);

        let mut winner = 0;
        if a <= self.decks[0].len() && b <= self.decks[1].len() {
            let mut g = self.create_sub_game(a, b);
            winner = g.play_recursive_game();
        } else if a < b {
            winner = 1;
        }

        if winner == 0 {
            self.decks[0].push(a);
            self.decks[0].push(b);
        } else {
            self.decks[1].push(b);
            self.decks[1].push(a);
        }
    }

    fn create_sub_game(&self, a: usize, b: usize) -> Puzzle {
        let mut decks = Vec::new();
        let deck = self.decks[0][..a].to_vec();
        decks.push(deck);
        let deck = self.decks[1][..b].to_vec();
        decks.push(deck);
        let memo = HashSet::new();
        Puzzle { decks, memo }
    }
}

fn main() {
    let args: Vec<String> = env::args().collect();
    let filename = args[1].as_str();

    let mut p = Puzzle::read_file(&filename).unwrap();
    let mut q = p.clone();

    p.play_game();

    let score = p.get_score();
    println!("Part 1: {score}");

    q.play_recursive_game();

    let score = q.get_score();
    println!("Part 2: {score}");
}
