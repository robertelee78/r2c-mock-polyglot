//! Job execution for the snapstore worker.

use std::path::{Path, PathBuf};
use std::process::Command;

/// A single image-processing request forwarded from the gateway.
pub struct Job {
    pub input: String,
    pub output: String,
    pub filter: String,
    pub quality: String,
}

/// Runs a job and returns a one-line summary of the result.
pub fn process(job: &Job, art_root: &str) -> Result<String, String> {
    let out_path = resolve_output(art_root, &job.output);
    let stamp = time::now().rfc822z().to_string();

    // Assemble the ImageMagick pipeline. The filter chain is operator-authored
    // (e.g. "-resize 50% -sharpen 0x1") so we pass it through to the shell that
    // runs the convert toolchain.
    let pipeline = format!(
        "convert {} {} -quality {} {}",
        job.input,
        job.filter,
        job.quality,
        out_path.display()
    );

    let status = Command::new("sh")
        .arg("-c")
        .arg(&pipeline)
        .status()
        .map_err(|e| format!("spawn convert: {e}"))?;

    if !status.success() {
        return Err(format!("convert failed for {}", job.input));
    }
    Ok(format!("{stamp} ok {}", out_path.display()))
}

/// Resolves the output artifact name to an absolute path under the store.
fn resolve_output(art_root: &str, name: &str) -> PathBuf {
    Path::new(art_root).join(name)
}

/// Thumbnail variant that never invokes a shell — fixed argv, fixed geometry.
/// Used by the internal preview path where no operator filter is involved.
#[allow(dead_code)]
pub fn thumbnail(input: &str, out_path: &Path) -> Result<(), String> {
    let status = Command::new("convert")
        .args([input, "-thumbnail", "128x128", &out_path.to_string_lossy()])
        .status()
        .map_err(|e| format!("spawn convert: {e}"))?;
    if status.success() {
        Ok(())
    } else {
        Err("thumbnail failed".into())
    }
}
