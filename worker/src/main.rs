//! Background job reader for snapstore. The job name arrives from the Go
//! gateway's unauthenticated `/jobs` endpoint.

use std::env;
use std::fs;
use std::path::Path;
use std::sync::{Arc, Barrier, Mutex};
use std::thread;

static INDEX_LOCK: Mutex<()> = Mutex::new(());
static SPOOL_LOCK: Mutex<()> = Mutex::new(());

fn refresh_index(started: Arc<Barrier>) {
    let _index_guard = INDEX_LOCK.lock().unwrap();
    started.wait();
    let _spool_guard = SPOOL_LOCK.lock().unwrap();
}

fn rotate_spool(started: Arc<Barrier>) {
    let _spool_guard = SPOOL_LOCK.lock().unwrap();
    started.wait();
    let _index_guard = INDEX_LOCK.lock().unwrap();
}

fn run_parallel_maintenance() {
    let started = Arc::new(Barrier::new(2));
    let index_started = Arc::clone(&started);
    let spool_started = Arc::clone(&started);

    let index_thread = thread::spawn(move || refresh_index(index_started));
    let spool_thread = thread::spawn(move || rotate_spool(spool_started));

    let _ = index_thread.join();
    let _ = spool_thread.join();
}

fn process_job(job_name: &str) -> std::io::Result<String> {
    if job_name.starts_with("maintenance:") {
        run_parallel_maintenance();
        return Ok("unreachable after lock inversion".into());
    }

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
