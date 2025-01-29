use std::env;
use std::fs;
use std::collections::HashMap;

fn main() {
    let args: Vec<String> = env::args().collect();
    if args.len() < 2 {
        println!("needed filename arg");
        return;
    }
    let input = fs::read_to_string(&args[1]).expect("need filename");
    let ans = num_antinodes(&input);
    println!("ans: {ans}");
}

fn num_antinodes(input: &String) -> i64 {
    let mut freq_to_coords: HashMap<u8, Vec<[i64; 2]>> = HashMap::new();
    let rows: Vec<&str> = input.trim().split('\n').collect();
    let height = rows.len() as i64;
    let width = rows[0].len() as i64;
    let mut visited: Vec<Vec<bool>> = Vec::new();
    for i in 0..rows.len() {
        visited.push(Vec::new());
        let row = rows[i].as_bytes();
        for j in 0..row.len() {
            visited[i].push(false);
            if row[j] != b'.' {
                freq_to_coords.entry(row[j]).or_insert(Vec::new()).push([i as i64, j as i64]);
            }
        }
    }
    let mut result: i64 = 0;
    for (k, v) in freq_to_coords.iter() {
        let v_len = v.len();
        println!("k: {k}, v: {v_len}");
        for i in 0..v_len {
            for j in (i+1)..v_len {
                get_antinodes(&v[i], &v[j], width, height, &mut visited, &mut result);
                get_antinodes(&v[j], &v[i], width, height, &mut visited, &mut result);
            }
        }
    }
    return result;
}

fn get_antinodes(p1: &[i64; 2], p2: &[i64; 2], width: i64, height: i64, visited: &mut Vec<Vec<bool>>, result: &mut i64) {
    let dx = p1[0] - p2[0];
    let dy = p1[1] - p2[1];
    if dx == 0 && dy == 0 {
        return;
    }
    let mut xx: i64 = p1[0];
    let mut yy: i64 = p1[1];
    while xx >= 0 && xx < height && yy >= 0 && yy < width {
        let xx_usize = xx as usize;
        let yy_usize = yy as usize;
        if !visited[xx_usize][yy_usize] {
            *result += 1;
            visited[xx_usize][yy_usize] = true;
        }
        xx += dx;
        yy += dy;
    }
}
