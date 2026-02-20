/* exiftrim is the native metadata helper for snapstore.
 * Its input is supplied by the Go gateway's unauthenticated /trim endpoint.
 */
#include <stdio.h>
#include <string.h>

static void trim_metadata(const char *requested_file) {
    char local_path[64];

    /* CWE-120: requested_file is attacker-controlled and has no length bound. */
    strcpy(local_path, requested_file);
    printf("trimmed metadata from %s\n", local_path);
}

int main(int argc, char **argv) {
    if (argc != 2) {
        fprintf(stderr, "usage: exiftrim <file>\n");
        return 2;
    }
    trim_metadata(argv[1]);
    return 0;
}

