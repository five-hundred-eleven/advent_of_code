use regex::Regex;
use std::env;
use std::fs;

fn main() {
    let args: Vec<String> = env::args().collect();
    if args.len() < 2 {
        println!("needed filename! {args:#?}");
        return;
    }
    let input = fs::read_to_string(&args[1].trim()).expect("Can't read file");
    let mul_re = Regex::new(r"(?<enable>do\(\))|(?<disable>don't\(\))|(?<mul>mul\((?<left>[0-9]{1,3}),(?<right>[0-9]{1,3})\))").expect("Cannot compile regex");
    let result = do_mult(input, mul_re);
    println!("ans: {result}");
}

fn do_mult(input: String, mul_re: Regex) -> i64 {
    let mut result: i64 = 0;
    let mut is_enabled: bool = true;
    for cap in mul_re.captures_iter(&input) {
        if cap.name("enable").is_some() {
            is_enabled = true;
        } else if cap.name("disable").is_some() {
            is_enabled = false;
        } else if is_enabled && cap.name("mul").is_some() {
            let left = &cap["left"];
            let right = &cap["right"];
            let left = left.parse::<i64>().expect("can't parse left");
            let right = right.parse::<i64>().expect("can't parse right");
            result += left * right;
        }
    }
    return result;
}
