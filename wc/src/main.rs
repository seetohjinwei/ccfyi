use clap::Parser;
use std::{fs::File, io::Read};

/// wc
#[derive(Parser, Debug)]
#[command(about)]
pub struct Args {
    /// Line count
    #[arg(short, default_value_t = false)]
    lines: bool,

    /// Word count
    #[arg(short, default_value_t = false)]
    words: bool,

    /// Character count
    #[arg(short, default_value_t = false)]
    chars: bool,

    /// File path
    #[arg(num_args(1))]
    path: String,
}

fn wc(file: &File, args: Args) {
    let mut chars = 0i32;
    let mut lines = 0i32;
    let mut words = 0i32; // TODO: fix words

    // A word is a non-zero-length sequence of characters delimited by whitespace.
    let mut in_word = true;

    for byte in file.bytes() {
        let byte = byte.unwrap();

        chars += 1;
        if byte.is_ascii_whitespace() {
            if in_word {
                in_word = false;
                words += 1;
            }
        } else {
            in_word = true;
        }
        if byte == '\n' as u8 {
            lines += 1;
        }
    }

    if in_word {
        words += 1;
    }

    print!("\t");

    if args.lines {
        print!("{}\t", lines);
    }
    if args.words {
        print!("{}\t", words);
    }
    if args.chars {
        print!("{}\t", chars);
    }

    print!("{}", args.path);

    println!("");
}

fn main() {
    let mut args = Args::parse();

    if !args.chars && !args.lines && !args.words {
        args.chars = true;
        args.lines = true;
        args.words = true;
    }

    let f = File::open(args.path.clone());
    match f {
        Err(e) => panic!("file not found: {}", e),
        Ok(f) => wc(&f, args),
    }
}
