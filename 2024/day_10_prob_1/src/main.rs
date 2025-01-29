use std::env;
use std::fs;

fn main() {
    let args: Vec<String> = env::args().collect();
    if args.len() < 2 {
        println!("need filename");
        return;
    }
    let input = fs::read_to_string(&args[1]).expect("could not read file");
    let grid = parse_input(input);
    let mut result = 0;
    let width = grid[0].len();
    let height = grid.len();
    let mut visited: Vec<Vec<bool>> = Vec::new();
    for row in &grid {
        let mut visited_row = Vec::new();
        for _ in row {
            visited_row.push(false);
        }
        visited.push(visited_row);
    }
    for i in 0..height {
        for j in 0..width {
            if grid[i][j] == b'0' {
                println!("found trailhead: {i} {j}");
                result += get_num_trails(&grid, &mut visited, i, j);
            }
        }
    }
    println!("ans: {result}");
}

fn parse_input(input: String) -> Vec<Vec<u8>> {
    let mut result: Vec<Vec<u8>> = Vec::new();
    for row in input.split('\n') {
        if row.len() == 0 {
            continue;
        }
        result.push(row.as_bytes().to_vec());
    }
    return result;
}

fn get_num_trails(grid: &Vec<Vec<u8>>, zero_visited: &mut Vec<Vec<bool>>, row: usize, col: usize) -> i64 {
    let mut visited: Vec<Vec<bool>> = Vec::new();
    for row in grid {
        let mut visited_row = Vec::new();
        for _ in row {
            visited_row.push(false);
        }
        visited.push(visited_row);
    }
    let mut current: Vec<[usize; 2]> = Vec::new();
    current.push([row, col]);
    let mut result = 0;
    let width = grid[0].len();
    let height = grid.len();
    while current.len() > 0 {
        let coords = current.pop().expect("could not pop");
        if visited[coords[0]][coords[1]] {
            continue;
        }
        visited[coords[0]][coords[1]] = true;
        let tile = grid[coords[0]][coords[1]];
        if tile == b'9'
        {
            result += 1;
            continue;
        }
        if coords[0] > 0 && is_valid(tile, grid[coords[0] - 1][coords[1]]) {
            current.push([coords[0] - 1, coords[1]]);
        }
        if coords[0] < height - 1 && is_valid(tile, grid[coords[0] + 1][coords[1]]) {
            current.push([coords[0] + 1, coords[1]]);
        }
        if coords[1] > 0 && is_valid(tile, grid[coords[0]][coords[1] - 1]) {
            current.push([coords[0], coords[1] - 1]);
        }
        if coords[1] < width - 1 && is_valid(tile, grid[coords[0]][coords[1] + 1]) {
            current.push([coords[0], coords[1] + 1]);
        }
    }
    return result;
}

fn is_valid(g1: u8, g2: u8) -> bool {
    if g1 - 1 == g2 || g1 + 1 == g2 {
        return true;
    }
    return false;
}
