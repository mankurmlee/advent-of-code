use std::{collections::{HashMap, HashSet}, env, fs, io};

struct Food {
    ingredients: Vec<String>,
    allergens: Vec<String>,
}

struct Puzzle {
    foods: Vec<Food>,
    ingredient_allergens: HashMap<String, String>,
}

impl Puzzle {
    fn read_file(filename: &str) -> Result<Self, io::Error> {
        let text = fs::read_to_string(filename)?;

        let mut foods = Vec::new();
        for line in text.lines() {
            let mut data = line.split(" (contains");
            let ingredients = data.next().unwrap()
                .split_whitespace()
                .map(String::from)
                .collect();
            let allergens = data.next().unwrap()
                .split_whitespace()
                .map(|s| s.chars()
                    .filter(|&b| !",)".contains(b))
                    .collect())
                .collect();
            foods.push(Food {ingredients, allergens});
        }

        let ingredient_allergens = HashMap::new();

        Ok(Puzzle { foods, ingredient_allergens })
    }

    fn deduce_ingredient_allergens(&mut self) {
        let mut allergen_foods = HashMap::new();
        for (i, f) in self.foods.iter().enumerate() {
            for a in f.allergens.iter().cloned() {
                allergen_foods.entry(a)
                    .and_modify(|e: &mut Vec<usize>| e.push(i))
                    .or_insert(vec![i]);
            }
        }
        let mut ingredient_allergens = HashMap::new();
        let mut mapped = HashSet::new();
        loop {
            for (a, f) in allergen_foods.iter() {
                if mapped.contains(a) {continue}
                let ings: Vec<String> = self.intersection(f).into_iter()
                    .filter(|ing| !ingredient_allergens.contains_key(ing))
                    .collect();

                if ings.len() == 1 {
                    let k = ings[0].clone();
                    ingredient_allergens.insert(k, a.clone());
                    mapped.insert(a.clone());
                    break;
                }
            }
            if mapped.len() == allergen_foods.len() {break}
        }
        self.ingredient_allergens = ingredient_allergens;
    }

    fn intersection(&self, f: &[usize]) -> Vec<String> {
        let mut out = Vec::new();
        for ing in self.foods[f[0]].ingredients.iter() {
            let mut found = true;
            for &i in f[1..].iter() {
                if !self.foods[i].ingredients.contains(ing) {
                    found = false;
                    break;
                }
            }
            if found {
                out.push(ing.clone());
            }
        }
        out
    }

    fn count_safe_ingredients(&self) -> usize {
        self.foods.iter()
            .map(|f|
                f.ingredients.iter()
                    .filter(|&ing|
                        !self.ingredient_allergens.contains_key(ing))
                    .count()
            )
            .sum()
    }

    fn get_danger_list(&self) -> String {
        let mut v = self.ingredient_allergens.iter()
            .map(|(k, v)| (v.clone(), k.clone()))
            .collect::<Vec<_>>();
        v.sort();
        v.into_iter()
            .map(|(_, v)| v)
            .collect::<Vec<_>>()
            .join(",")
    }
}

fn main() {
    let args: Vec<String> = env::args().collect();
    let filename = args[1].as_str();

    let mut p = Puzzle::read_file(&filename).unwrap();
    p.deduce_ingredient_allergens();
    dbg!(&p.ingredient_allergens);

    let safe_count = p.count_safe_ingredients();
    println!("Part 1: {safe_count}");

    let bad_list = p.get_danger_list();
    println!("Part 2: {bad_list}");
}
