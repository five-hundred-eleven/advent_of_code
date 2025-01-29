use std::io;
use std::fs;

fn main() {
    let mut input_filename = String::new();
    io::stdin()
        .read_line(&mut input_filename)
        .expect("Can't get filename");
    let input = fs::read_to_string(input_filename.trim()).expect("Can't read file");
    let levels = parse_input(input);
    let ans = num_valid_levels(levels);
    println!("ans: {ans}");
}

fn parse_input(input: String) -> Vec<Vec<i64>> {
    let mut result = Vec::new();
    for line in input.split('\n') {
        let line = line.trim();
        let line_parts: Vec<&str> = line.split_ascii_whitespace().collect();
        if line_parts.len() == 0 {
            continue;
        }
        let mut level = Vec::new();
        for part in line_parts {
            let part_result = part.parse::<i64>();
            if part_result.is_err() {
                println!("could not parse: {part}");
                continue
            }
            level.push(part_result.unwrap());
        }
        result.push(level);
    }
    return result;
}

fn num_valid_levels(levels: Vec<Vec<i64>>) -> i64 {
    let mut result: i64 = 0;
    for level in levels {
        if level.len() < 2 {
            println!("level has 1 or fewer elements");
            continue;
        }
        let mut num_asc = 0;
        let mut num_desc = 0;
        for l in &level {
            if *l < level[0] {
                num_desc += 1;
            }
            if *l > level[0] {
                num_asc += 1;
            }
        }
        let asc = num_asc > num_desc;
        if is_level_valid(&level, asc, true) {
            result += 1;
        }
    }
    return result;
}

fn is_level_valid(level: &Vec<i64>, asc: bool, remove: bool) -> bool {
    println!("{level:#?} {asc} {remove}");
    for i in 1..level.len() {
        let diff = level[i] - level[i - 1];
        if (asc && (diff < 1 || diff > 3)) || (!asc && (diff > -1 || diff < -3)) {
            if !remove {
                return false;
            }
            let mut cp = level.to_vec();
            cp.remove(i);
            let result = is_level_valid(&cp, asc, false);
            if result {
                println!("true: {i} removed");
                return true;
            }
            let mut cp = level.to_vec();
            cp.remove(i-1);
            let result = is_level_valid(&cp, asc, false);
            println!("{result:#?}: {i} removed");
            return result;
        }
    }
    println!("true");
    return true;
}
