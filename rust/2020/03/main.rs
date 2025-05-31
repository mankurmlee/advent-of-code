use std::{env, fs, io::{self, BufRead}};

struct Grid {
    width: usize,
    height: usize,
    data: Vec<char>,
}

impl Grid {
    fn load(filename: &str) -> Result<Grid, io::Error> {
        let file = fs::File::open(filename)?;
        let reader = io::BufReader::new(file);

        let mut width = 0;
        let mut height = 0;
        let mut data = Vec::new();

        for line in reader.lines() {
            let line = line?;
            width = width.max(line.chars().count());
            height += 1;
            data.extend(line.chars());
        }

        Ok(Grid { width, height, data })
    }

    fn slope_hits(&self, slope: &(i32, i32)) -> i32 {
        let w = self.width as i32;
        let h = self.height as i32;

        let mut x = 0;
        let mut y = 0;
        let mut hits = 0;

        while y < h {
            let i = (y * w + x) as usize;
            if self.data[i] == '#' {
                hits += 1;
            }
            x = (x + slope.0) % w;
            y += slope.1;
        }

        hits
    }
}

fn main() {
    let filename = env::args().nth(1).unwrap();

    let p = Grid::load(&filename).unwrap();

    let slopes = vec![(1, 1), (3, 1), (5, 1), (7, 1), (1, 2)];

    println!("Part 1: {}", p.slope_hits(&slopes[1]));

    let product = slopes
        .iter()
        .map(|s| p.slope_hits(s))
        .fold(1, |acc, sum| acc * sum);

    println!("Part 2: {product}");
}
