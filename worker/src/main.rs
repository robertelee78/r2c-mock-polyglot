//! snapstore worker — applies a convert-style filter chain to an artifact.
//!
//! Invoked by the gateway as:
//!   worker --input <path> --output <name> --filter <chain> --quality <n>

mod jobs;

use std::process::exit;

use jobs::Job;

fn parse_args() -> Result<Job, String> {
    let mut input = String::new();
    let mut output = String::new();
    let mut filter = String::new();
    let mut quality = String::from("85");

    let mut args = std::env::args().skip(1);
    while let Some(flag) = args.next() {
        let val = args.next().ok_or_else(|| format!("missing value for {flag}"))?;
        match flag.as_str() {
            "--input" => input = val,
            "--output" => output = val,
            "--filter" => filter = val,
            "--quality" => quality = val,
            other => return Err(format!("unknown flag {other}")),
        }
    }
    if input.is_empty() || output.is_empty() {
        return Err("--input and --output are required".into());
    }
    Ok(Job { input, output, filter, quality })
}

fn main() {
    let job = match parse_args() {
        Ok(j) => j,
        Err(e) => {
            eprintln!("worker: {e}");
            exit(2);
        }
    };

    match jobs::process(&job, "/var/lib/snapstore/artifacts") {
        Ok(summary) => println!("{summary}"),
        Err(e) => {
            eprintln!("worker: {e}");
            exit(1);
        }
    }
}
