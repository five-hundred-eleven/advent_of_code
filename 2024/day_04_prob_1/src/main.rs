use std::env;
use std::fs;

fn main() {
    let args: Vec<String> = env::args().collect();
    if args.len() < 2 {
        println!("need filename");
        return;
    }
    let input = fs::read_to_string(&args[1].trim()).expect("could not read file");
    let result = solve_crossword(input);
    println!("ans: {result}");
}

fn solve_crossword(input: String) -> i64 {
    let mut result: i64 = 0;
    let input_rows: Vec<&str> = input.split('\n').collect();
    let mut rows: Vec<Vec<u8>> = Vec::new();
    for row in &input_rows {
        rows.push(row.as_bytes().to_vec());
    }
    for i in 0..rows.len() {
        for j in 0..rows[i].len() {
            if rows[i][j] == b'X' {
                if i > 3 && search_at_loc(&rows, i, j, -1, 0) {
                    result += 1;
                }
                if i < rows.len() - 4 && search_at_loc(&rows, i, j, 1, 0) {
                    result += 1;
                }
                if j > 3 && search_at_loc(&rows, i, j, 0, -1) {
                    result += 1;
                }
                if j < rows[i].len() - 4 && search_at_loc(&rows, i, j, 0, 1) {
                    result += 1;
                }
                if i > 3 && j > 3 && search_at_loc(&rows, i, j, -1, -1) {
                    result += 1;
                }
                if i < rows.len() - 4 && j > 3 && search_at_loc(&rows, i, j, 1, -1) {
                    result += 1;
                }
                if i > 3 && j < rows[i].len() - 4 && search_at_loc(&rows, i, j, -1, 1) {
                    result += 1;
                }
                if i < rows.len() - 4 && j < rows[i].len() - 4 && search_at_loc(&rows, i, j, 1, 1) {
                    result += 1;
                }
            }
        }
    }
    return result;
}

fn search_at_loc(rows: &Vec<Vec<u8>>, i: usize, j: usize, i_offset: i64, j_offset: i64) -> bool {
    let mut i: i64 = i64::try_from(i).expect("could not convert i to i64");
    let mut j: i64 = i64::try_from(j).expect("could not convert j to i64");
    for letter in [b'X', b'M', b'A', b'S'] {
        if rows[i as usize][j as usize] != letter {
            return false;
        }
        i += i_offset;
        j += j_offset;
    }
    return true;
}
