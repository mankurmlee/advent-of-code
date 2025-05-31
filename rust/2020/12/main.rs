use std::{env, fs, io::{self, BufRead}};

struct Instruction {
    cmd: u8,
    arg: i32,
}

struct Route {
    data: Vec<Instruction>,
}

impl Route {
    fn read_file(filename: &str) -> Result<Self, io::Error> {
        let file = fs::File::open(filename)?;
        let reader = io::BufReader::new(file);

        let mut data = Vec::new();

        for line in reader.lines() {
            let line = line?;
            let cmd = line.as_bytes()[0];
            let arg = if let Ok(v) = line[1..].parse::<i32>()
                { v } else { continue };
            data.push(Instruction{ cmd, arg });
        }

        Ok(Route { data })
    }
}

struct Ship {
    x: i32,
    y: i32,
    fc: i32,
    dx: i32,
    dy: i32,
}

impl Ship {
    fn new() -> Self {
        Ship {
            x : 0,
            y : 0,
            fc: 90,
            dx: 10,
            dy: 1,
        }
    }

    fn sail(&mut self, r: &Route) -> i32 {
        for i in r.data.iter() {
            self.move_ship(i);
        }
        self.x.abs() + self.y.abs()
    }

    fn use_waypoint(&mut self, r: &Route) -> i32 {
        for i in r.data.iter() {
            self.move_wp(i);
        }
        self.x.abs() + self.y.abs()
    }

    fn move_ship(&mut self, ins: &Instruction) {
        let scalar = ins.arg;
        match ins.cmd {
            b'N' => self.y += scalar,
            b'E' => self.x += scalar,
            b'S' => self.y -= scalar,
            b'W' => self.x -= scalar,
            b'L' => {
                self.fc -= scalar;
                if self.fc < 0 {
                    self.fc += 360;
                }
            },
            b'R' => {
                self.fc += scalar;
                self.fc %= 360;
            },
            b'F' => {
                match self.fc {
                    0   => self.y += scalar,
                    90  => self.x += scalar,
                    180 => self.y -= scalar,
                    270 => self.x -= scalar,
                    f => {
                        println!("Unexpected facing: {f}");
                        return;
                    }
                }
            },
            v => {
                println!("Unexpected instruction: {v}");
                return;
            },
        }
    }

    fn move_wp(&mut self, ins: &Instruction) {
        let scalar = ins.arg;
        match ins.cmd {
            b'N' => self.dy += scalar,
            b'E' => self.dx += scalar,
            b'S' => self.dy -= scalar,
            b'W' => self.dx -= scalar,
            b'L' => {
                self.wp_rotate(360 - scalar);
            },
            b'R' => {
                self.wp_rotate(scalar);
            },
            b'F' => {
                self.x += self.dx * scalar;
                self.y += self.dy * scalar;
            },
            v => {
                println!("Unexpected instruction: {v}");
                return;
            },
        }
    }

    fn wp_rotate(&mut self, scalar: i32) {
        let times = scalar / 90;
        for _ in 0..times {
            (self.dx, self.dy) = (self.dy, -self.dx)
        }
    }

}

fn main() {
    let args: Vec<String> = env::args().collect();
    let filename = args[1].as_str();

    let r = Route::read_file(&filename).unwrap();

    let mut f = Ship::new();
    println!("Part 1: {}", f.sail(&r));

    let mut g = Ship::new();
    println!("Part 2: {}", g.use_waypoint(&r));
}
