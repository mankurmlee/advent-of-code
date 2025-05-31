use std::{collections::HashMap, env, fs, io};

const MONSTER: &str = "\
..................#.\n\
#....##....##....###\n\
.#..#..#..#..#..#...";

struct Monster {
    body: Vec<(usize, usize)>,
    width: usize,
    height: usize,
}
impl Monster {
    fn new() -> Monster {
        let text: Vec<&str> = MONSTER.lines().collect();
        let width = text[0].len();
        let height = text.len();

        let mut body = Vec::new();
        for y in 0..height {
            for x in 0..width {
                if text[y].as_bytes()[x] == b'#' {
                    body.push((x, y));
                }
            }
        }

        Monster { body, width, height }
    }

    fn rotate(&mut self) {
        let mut new_body = Vec::new();
        for (x, y) in self.body.iter() {
            let new_x = self.height - *y - 1;
            let new_y = *x;
            new_body.push((new_x, new_y));
        }
        self.body = new_body;
        (self.width, self.height) = (self.height, self.width);
    }

    fn flip(&mut self) {
        let mut new_body = Vec::new();
        for (x, y) in self.body.iter() {
            let new_x = self.width - *x - 1;
            new_body.push((new_x, *y));
        }
        self.body = new_body;
    }
}

#[derive(Clone)]
struct Tile {
    id: i32,
    data: [bool; 100],
}
impl Tile {
    fn flip(&mut self) {
        for i in (0..100).step_by(10) {
            self.data[i..i+10].reverse();
        }
    }

    fn rotate(&mut self) {
        let mut new_data = [false; 100];
        for y in 0..10 {
            for x in 0..10 {
                let new_x = 9 - y;
                let new_y = x;
                new_data[new_y * 10 + new_x] = self.data[y * 10 + x];
            }
        }
        self.data = new_data;
    }
}

struct Puzzle {
    tiles: Vec<Tile>,
    grid: Vec<SideTile>,
    side: HashMap<u16, Vec<i32>>,
    links: HashMap<i32, Vec<i32>>,
    width: usize,
}
impl Puzzle {
    fn read_file(filename: &str) -> Result<Self, io::Error> {
        let mut tiles = Vec::new();

        let text = fs::read_to_string(filename)?;

        let mut id = 0;
        let mut data: [bool; 100] = [false; 100];
        let mut i = 0;

        for line in text.lines() {
            if line.trim().is_empty() { continue }
            if line.contains(':') {
                if i > 0 {
                    tiles.push(Tile {id, data});
                    data = [false; 100];
                    i = 0;
                }
                id = line.replace(':', " ")
                    .split_whitespace()
                    .nth(1).unwrap()
                    .parse().unwrap();
            } else {
                for c in line.chars() {
                    data[i] = c == '#';
                    i += 1;
                }
            }
        }
        if i > 0 {
            tiles.push(Tile {id, data});
        }

        // Create lightweight tiles with sides-only for easier manipulation
        // Initially just a list, but will later be sorted into grid-order
        let grid: Vec<SideTile> = tiles.iter()
            .map(SideTile::from)
            .collect();

        // Collect tiles with matching sides into buckets
        let mut side = HashMap::new();
        for t in grid.iter() {
            for i in 0..4 {
                // Get a 'fingerprint' for the side
                let k = t.fingerprint(i);

                // Place tile ID into the bucket for matching
                side.entry(k)
                    .and_modify(|v: &mut Vec<i32>| v.push(t.id))
                    .or_insert(vec![t.id]);
            }
        }

        // Create a lookup for neighbouring tiles
        let mut links = HashMap::new();
        for (_, g) in side.iter() {
            for &i in g.iter() {
                for &j in g.iter() {
                    if i == j { continue }
                    links.entry(i)
                        .and_modify(|v: &mut Vec<i32>| v.push(j))
                        .or_insert(vec![j]);
                }
            }
        }

        let width = (grid.len() as f64).sqrt() as usize;

        Ok(Puzzle { tiles, grid, side, links, width })
    }

    fn corner_ids(&self) -> Vec<i32> {
        self.links.iter()
            .filter(|(_, v)| v.len() == 2)
            .map(|(k, _)| *k)
            .collect()
    }

    fn arrange(&mut self) {
        let n = self.grid.len();
        let w = self.width;
        for i in 0..n {
            // Determine left and up matching
            let mut left_fp = 0;
            let mut top_fp= 0;

            let has_left = (i % w) > 0;
            let mut num_links = 4;
            if has_left {
                left_fp = mirror(self.grid[i-1].sides[1]);
            } else {
                num_links -= 1;
            }

            let has_top = i >= w;
            if has_top {
                top_fp = mirror(self.grid[i-w].sides[2]);
            } else {
                num_links -= 1;
            }

            if (i % w) == w - 1 {
                num_links -= 1;
            }
            if i >= n - w {
                num_links -= 1;
            }

            for j in i..n {
                let tile = &self.grid[j];
                let l = self.links.get(&tile.id).unwrap();
                if l.len() > num_links { continue }
                // Optimisation: Could check links/pool and eliminate candidates
                let mut tile = self.grid[j].clone();

                if !has_left && !has_top {
                    while  self.tile_has_link(&tile, 0)
                        || self.tile_has_link(&tile, 3) {
                        tile.rotate();
                    }
                } else {
                    let fp = if has_left { left_fp } else { top_fp };
                    let n = if has_left { 3 } else { 0 };
                    if !tile.sides.contains(&fp) {
                        tile.flip();
                        if !tile.sides.contains(&fp) {
                            continue;
                        }
                    }
                    while tile.sides[n] != fp {
                        tile.rotate();
                    }
                }

                self.grid[j] = tile;
                self.grid.swap(i, j);
                break;
            }
        }
    }

    fn tile_has_link(&self, tile: &SideTile, n: usize) -> bool {
        let fp = tile.fingerprint(n);
        let v = self.side.get(&fp).unwrap();
        v.len() > 1
    }

    fn create_lut(&self) -> Vec<Tile> {
        let mut lut = HashMap::new();
        for t in self.tiles.iter() {
            lut.insert(t.id, t.clone());
        }

        let mut out = Vec::with_capacity(self.grid.len());
        for s in self.grid.iter() {
            let mut tile = lut.remove(&s.id).unwrap();
            if s.flip {
                tile.flip();
            }
            for _ in 0..s.rot {
                tile.rotate();
            }
            out.push(tile);
        }

        out
    }

    fn create_map(&self) -> Map {
        let lut = self.create_lut();
        let w = self.width;
        let n = w * w;
        let mut grid = Vec::with_capacity(n);
        for x in (0..n).step_by(w) {
            for y in 1..9 {
                for z in 0..w {
                    let t = &lut[x+z];
                    let i = y * 10 + 1;
                    let j = y * 10 + 9;
                    let slice = &t.data[i..j];
                    grid.extend_from_slice(slice);
                }
            }
        }
        let width = self.width * 8;
        Map { width, grid }
    }
}

#[derive(Clone)]
struct SideTile {
    id   : i32,
    sides: [u16; 4],
    rot: i32,
    flip: bool,
}
impl SideTile {
    fn from(tile: &Tile) -> Self {
        let id = tile.id;
        let mut sides = [0; 4];
        for i in 0..10 {
            if tile.data[9 - i] {
                sides[0] |= 1 << i;
            }
            if tile.data[100 - i * 10 - 1] {
                sides[1] |= 1 << i;
            }
            if tile.data[i + 90] {
                sides[2] |= 1 << i;
            }
            if tile.data[10 * i] {
                sides[3] |= 1 << i;
            }
        }
        SideTile {id, sides, rot: 0, flip: false}
    }

    fn fingerprint(&self, side: usize) -> u16 {
        let a = self.sides[side];
        let b = mirror(a);
        a.min(b)
    }

    fn rotate(&mut self) {
        for i in (0..3).rev() {
            self.sides.swap(i, i+1);
        }
        self.rot += 1;
    }

    fn flip(&mut self) {
        self.sides.swap(1, 3);
        for i in 0..4 {
            self.sides[i] = mirror(self.sides[i]);
        }
        self.flip = true;
    }
}

fn mirror(mut n: u16) -> u16 {
    let mut o = 0;
    for _ in 0..10 {
        o = (o << 1) | (n & 1);
        n >>= 1;
    }
    o
}

struct Map {
    width: usize,
    grid: Vec<bool>,
}
impl Map {
    fn print(&self) {
        let n = self.grid.len();
        let w = self.width;
        for i in (0..n).step_by(w) {
            let b = &self.grid[i..i+w];
            let s: String = b.iter()
                .map(|&b| if b {'#'} else {'.'})
                .collect();
            println!("{s} {i}");
        }
    }

    fn find_monsters(&self) -> usize {
        let mut monster = Monster::new();
        for _ in 0..1 {
            for _ in 0..4 {
                let count = self.count_monsters(&monster);
                if count > 0 {
                    return count;
                }
                monster.rotate();
            }
            monster.flip();
        }
        0
    }

    fn count_monsters(&self, monster: &Monster) -> usize {
        let mut count = 0;
        let w = self.width - monster.width + 1;
        let h = self.width - monster.height + 1;
        for y in 0..h {
            for x in 0..w {
                count += 1;
                for (dx, dy) in monster.body.iter() {
                    let i = (y + *dy) * self.width + (x + *dx);
                    if !self.grid[i] {
                        count -= 1;
                        break;
                    }
                }
            }
        }
        count
    }
}

fn main() {
    let args: Vec<String> = env::args().collect();
    let filename = args[1].as_str();

    let mut p = Puzzle::read_file(&filename).unwrap();
    println!("{} tiles loaded.", p.tiles.len());
    p.arrange();

    let prod: u64 = p.corner_ids().iter()
        .map(|&n| n as u64).product();

    let m = p.create_map();
    m.print();

    // Find monsters in the map
    let num_monsters = m.find_monsters();
    println!("{} monsters found", num_monsters);

    let c = m.grid.iter()
        .filter(|&b| *b).count();
    let roughness = c - num_monsters * 15;

    println!("Part 1: {prod}");
    println!("Part 2: {roughness}");
}
