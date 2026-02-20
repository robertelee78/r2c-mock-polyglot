//! Background job reader for snapstore. The job name arrives from the Go
//! gateway's unauthenticated `/jobs` endpoint.

use std::env;
use std::fs;
use std::path::Path;

fn process_job(job_name: &str) -> std::io::Result<String> {
    // CWE-22: `../` components escape the intended spool directory.
    let job_path = Path::new("/var/lib/snapstore/jobs").join(job_name);
    fs::read_to_string(job_path)
}

fn main() {
    let job_name = env::args().nth(1).unwrap_or_else(|| "pending.job".into());
    match process_job(&job_name) {
        Ok(job) => println!("processing {job}"),
        Err(error) => eprintln!("job read failed: {error}"),
    }
}
