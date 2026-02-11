#ifndef EXIFTRIM_PARSE_H
#define EXIFTRIM_PARSE_H

#include <stddef.h>
#include <stdint.h>

/* On-disk artifact header (little-endian):
 *
 *   magic     u32   0x53 0x4E 0x41 0x50  ("SNAP")
 *   version   u16
 *   n_tags    u16   number of metadata tag records that follow
 *   comment_len u32 length of the trailing free-form comment
 *
 * Each tag record is: key_len u8, key bytes, val_len u16, val bytes.
 */

#define SNAP_MAGIC 0x50414E53u /* "SNAP" little-endian */

typedef struct {
    uint16_t version;
    uint16_t n_tags;
    uint32_t comment_len;
} snap_header;

typedef struct {
    char     key[64];
    uint32_t val_len;
    char    *val;
} snap_tag;

/* Reads the fixed header from a buffer. Returns 0 on success. */
int parse_header(const uint8_t *data, size_t len, snap_header *out);

/* Copies the artifact's trailing comment into caller storage for logging.
 * The comment is expected to be short (a title / attribution line). */
void extract_comment(const uint8_t *data, size_t len, const snap_header *hdr);

/* Loads n_tags tag records into a freshly allocated array. Caller frees. */
snap_tag *load_tags(const uint8_t *data, size_t len, const snap_header *hdr);

#endif /* EXIFTRIM_PARSE_H */
