use std::fs;

fn main() {
    let input: String = fs::read_to_string("input.txt")
        .expect("Can't read input file");
    let (left, right) = parse_input(input.trim());
    let result = ordered_difference(left, right);
    println!("result: {result}");
}

// input: is given ownership, further references in calling function are not allowed
// return values: gives ownership, calling function will own them
fn parse_input(input: String) -> (Vec<i64>, Vec<i64>) {
    let mut left: Vec<i64> = Vec::new();
    let mut right: Vec<i64> = Vec::new();
    for line in input.split('\n') {
        let line: &str = line.trim();
        let line_parts: Vec<&str> = line.split_ascii_whitespace().collect();
        if line_parts.len() != 2 {
            println!("Line does not have 2 parts: {line}");
            continue;
        }
        let l_result = line_parts[0].parse::<i64>();
        if l_result.is_err() {
            println!("Line left part is not int: {line}");
            continue;
        }
        let r_result = line_parts[1].parse::<i64>();
        if r_result.is_err() {
            println!("Line right part is not int: {line}");
            continue;
        }
        left.push(l_result.unwrap());
        right.push(r_result.unwrap());
    }
    return (left, right);
}

fn ordered_difference(mut left: Vec<i64>, mut right: Vec<i64>) -> i64 {
    left.sort();
    right.sort();
    let mut result: i64 = 0;
    for (l_val, r_val) in left.iter().zip(right.iter()) {
        if l_val < r_val {
            result += r_val - l_val;
        } else {
            result += l_val - r_val;
        }
    }
    return result;
}
