use std::collections::HashMap;
use std::fs;
use std::io;

fn main() {
    let mut filename = String::new();
    io::stdin()
        .read_line(&mut filename)
        .expect("Can't input filename");
    let input: String = fs::read_to_string(filename.trim()).expect("Can't read input file");
    let (left, right) = parse_input(input);
    let result = occurrence_sum(left, right);
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

fn occurrence_sum(left: Vec<i64>, right: Vec<i64>) -> i64 {
    let mut num_to_count = HashMap::new();
    let mut result: i64 = 0;
    for l_val in &left {
        num_to_count.entry(l_val).or_insert(0);
    }
    for r_val in &right {
        num_to_count.entry(r_val).and_modify(|count| *count += 1);
    }
    for l_val in &left {
        match num_to_count.get(l_val) {
            Some(count) => {
                result += l_val * count;
            },
            None => {}
        }
    }
    return result;
}
