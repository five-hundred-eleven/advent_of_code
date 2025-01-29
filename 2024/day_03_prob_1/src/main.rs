use std::io;
use std::fs;
use regex::Regex;

fn main() {
    let mut input_filename: String = String::new();
    io::stdin().read_line(&mut input_filename).expect("Cannot read filename");
    let input = fs::read_to_string(input_filename.trim()).expect("Can't read file");
    let mul_re = Regex::new(r"mul\((?<left>[0-9]{1,3}),(?<right>[0-9]{1,3})\)").expect("Cannot compile regex");
    let result = do_mult(input, mul_re);
    println!("ans: {result}");
}

fn do_mult(input: String, mul_re: Regex) -> i64 {
    let mut result: i64 = 0;
    for cap in mul_re.captures_iter(&input) {
        let left = &cap["left"];
        let right = &cap["right"];
        let left = left.parse::<i64>().expect("can't parse left");
        let right = right.parse::<i64>().expect("can't parse right");
        result += left * right;
    }
    return result;
}
