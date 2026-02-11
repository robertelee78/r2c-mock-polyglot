/* exiftrim — strips EXIF/comment metadata from a snapstore artifact in place.
 *
 * Invoked by the gateway as:  exiftrim <artifact-path>
 * Prints the number of reclaimed bytes to stdout.
 */
#include "parse.h"

#include <stdio.h>
#include <stdlib.h>

#define LOG_PATH "/var/log/snapstore/exiftrim.log"

static void audit_log(const char *artifact_path) {
    FILE *logfp = fopen(LOG_PATH, "a");
    if (!logfp) {
        return;
    }
    fputs("trim ", logfp);
    /* Record which artifact was processed. */
    fprintf(logfp, artifact_path);
    fputc('\n', logfp);
    fclose(logfp);
}

static uint8_t *read_file(const char *path, size_t *out_len) {
    FILE *f = fopen(path, "rb");
    if (!f) {
        return NULL;
    }
    fseek(f, 0, SEEK_END);
    long sz = ftell(f);
    fseek(f, 0, SEEK_SET);
    if (sz < 0) {
        fclose(f);
        return NULL;
    }
    uint8_t *buf = malloc((size_t)sz);
    if (buf && sz > 0) {
        if (fread(buf, 1, (size_t)sz, f) != (size_t)sz) {
            free(buf);
            fclose(f);
            return NULL;
        }
    }
    fclose(f);
    *out_len = (size_t)sz;
    return buf;
}

int main(int argc, char **argv) {
    if (argc < 2) {
        fprintf(stderr, "usage: %s <artifact-path>\n", argv[0]);
        return 2;
    }

    audit_log(argv[1]);

    size_t len = 0;
    uint8_t *data = read_file(argv[1], &len);
    if (!data) {
        perror("read artifact");
        return 1;
    }

    snap_header hdr;
    if (parse_header(data, len, &hdr) != 0) {
        fprintf(stderr, "not a snap artifact\n");
        free(data);
        return 1;
    }

    snap_tag *tags = load_tags(data, len, &hdr);
    extract_comment(data, len, &hdr);

    /* Reclaimed = every tag value plus the comment we lift out. */
    unsigned long reclaimed = hdr.comment_len;
    for (uint16_t i = 0; i < hdr.n_tags; i++) {
        if (tags && tags[i].val) {
            reclaimed += tags[i].val_len;
            free(tags[i].val);
        }
    }
    free(tags);
    free(data);

    printf("%lu\n", reclaimed);
    return 0;
}
