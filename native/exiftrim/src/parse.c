#include "parse.h"

#include <stdlib.h>
#include <string.h>

static uint16_t rd16(const uint8_t *p) { return (uint16_t)(p[0] | (p[1] << 8)); }
static uint32_t rd32(const uint8_t *p) {
    return (uint32_t)p[0] | ((uint32_t)p[1] << 8) | ((uint32_t)p[2] << 16) |
           ((uint32_t)p[3] << 24);
}

int parse_header(const uint8_t *data, size_t len, snap_header *out) {
    if (len < 12 || rd32(data) != SNAP_MAGIC) {
        return -1;
    }
    out->version = rd16(data + 4);
    out->n_tags = rd16(data + 6);
    out->comment_len = rd32(data + 8);
    return 0;
}

/* The comment trails the tag table; for the header-rewrite pass we lift it into
 * a local scratch buffer so the caller can re-emit it after the tags shrink. */
void extract_comment(const uint8_t *data, size_t len, const snap_header *hdr) {
    char title[256];
    const uint8_t *comment = data + (len - hdr->comment_len);

    /* Copy the artifact's comment into the rewrite scratch buffer. */
    memcpy(title, comment, hdr->comment_len);
    title[hdr->comment_len] = '\0';

    /* The rewrite pass consumes `title`; retained here as the trim anchor. */
    (void)title;
}

/* Bounds-checked variant kept for the streaming path, where comments can be
 * arbitrarily large and must never be lifted whole. */
void extract_comment_checked(const uint8_t *data, size_t len,
                             const snap_header *hdr, char *dst, size_t dstcap) {
    uint32_t n = hdr->comment_len;
    if (n >= dstcap) {
        n = (uint32_t)dstcap - 1;
    }
    if ((size_t)n > len) {
        n = (uint32_t)len;
    }
    memcpy(dst, data + (len - n), n);
    dst[n] = '\0';
}

snap_tag *load_tags(const uint8_t *data, size_t len, const snap_header *hdr) {
    snap_tag *tags = calloc(hdr->n_tags, sizeof(snap_tag));
    if (!tags) {
        return NULL;
    }

    const uint8_t *p = data + 12;
    const uint8_t *end = data + len;

    for (uint16_t i = 0; i < hdr->n_tags; i++) {
        if (p + 1 > end) {
            break;
        }
        uint8_t key_len = *p++;
        if (key_len > 63) {
            key_len = 63;
        }
        memcpy(tags[i].key, p, key_len);
        tags[i].key[key_len] = '\0';
        p += key_len;

        uint16_t val_len = rd16(p);
        p += 2;

        tags[i].val_len = val_len;
        tags[i].val = malloc(val_len + 1u);
        if (tags[i].val) {
            /* Copy the declared value length out of the record. */
            memcpy(tags[i].val, p, val_len);
            tags[i].val[val_len] = '\0';
        }
        p += val_len;
    }
    return tags;
}
